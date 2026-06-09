package service

import (
	"blog/internal/model/dto"
	"blog/internal/model/entity"
	"blog/internal/repository"
	"blog/pkg/errors"
	"blog/pkg/logger"
	"blog/pkg/utils"
	"blog/pkg/validator"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthService struct {
	userRepo    *repository.UserRepository
	jwtSecret   string
	expireHours int
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string, expireHours int) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		jwtSecret:   jwtSecret,
		expireHours: expireHours,
	}
}

// Register 用户注册（返回 DTO 而非 Entity，避免敏感字段泄露）
func (s *AuthService) Register(req dto.RegisterRequest, avatarPath string) (*dto.UserResponse, error) {
	// 1. 验证密码强度
	if err := validator.ValidatePassword(req.Password); err != nil {
		return nil, err
	}

	// 2. 检查用户名是否存在
	if s.userRepo.ExistsByUsername(req.Username) {
		logger.Warnf("User registration failed: username already exists, username=%s", req.Username)
		return nil, errors.New(errors.UsernameExistsCode)
	}

	// 3. 检查邮箱是否存在
	if s.userRepo.ExistsByEmail(req.Email) {
		logger.Warnf("User registration failed: email already registered, email=%s", req.Email)
		return nil, errors.New(errors.EmailExistsCode)
	}

	// 4. 密码加密
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		logger.Errorf("User registration failed: password hashing failed, username=%s, error=%v", req.Username, err)
		return nil, errors.New(errors.PasswordHashFailedCode)
	}

	// 5. 创建用户实体
	user := &entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Nickname: req.Username,
		Avatar:   avatarPath,
		Role:     "user",
		Status:   1,
	}

	// 6. 保存到数据库（带并发保护）
	if err := s.userRepo.Create(user); err != nil {
		// 检查是否唯一索引冲突（并发场景）
		errStr := err.Error()
		if strings.Contains(errStr, "username") ||
			strings.Contains(errStr, "uk_username") ||
			strings.Contains(errStr, "Duplicate") {
			logger.Warnf("User registration failed: username conflict (possible concurrent registration), username=%s, error=%v", req.Username, err)
			return nil, errors.New(errors.UsernameExistsCode)
		}
		if strings.Contains(errStr, "email") ||
			strings.Contains(errStr, "uk_email") {
			logger.Warnf("User registration failed: email conflict, email=%s, error=%v", req.Email, err)
			return nil, errors.New(errors.EmailExistsCode)
		}
		logger.Errorf("User registration failed: database error, username=%s, email=%s, error=%v", req.Username, req.Email, err)
		return nil, errors.New(errors.RegistrationFailedCode)
	}

	// 7. 转换为 DTO 返回（过滤敏感字段）
	return &dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Bio:      user.Bio,
		Role:     user.Role,
	}, nil
}

// Login 用户登录
func (s *AuthService) Login(req dto.LoginRequest) (string, error) {
	// 1. 查询用户
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		logger.Warnf("User login failed: user not found, username=%s", req.Username)
		return "", errors.New(errors.InvalidCredentialsCode)
	}

	// 2. 验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		logger.Warnf("User login failed: incorrect password, username=%s", req.Username)
		return "", errors.New(errors.InvalidCredentialsCode)
	}

	// 3. 检查用户状态
	if user.Status != 1 {
		logger.Warnf("User login failed: user disabled, username=%s, status=%d", req.Username, user.Status)
		return "", errors.New(errors.UserDisabledCode)
	}

	// 4. 生成新的 login_session（带时间戳并发控制）
	now := time.Now()
	loginSession := fmt.Sprintf("%s_%d", uuid.New().String(), now.UnixNano())

	// 重试机制处理并发冲突
	var updateErr error
	for i := 0; i < 3; i++ {
		if err := s.userRepo.UpdateLoginSession(user.ID, loginSession, now); err != nil {
			updateErr = err
			// 如果是并发冲突（错误码 2105），短暂等待后重试
			if errors.IsUpdateSessionFailed(err) {
				logger.Warnf("Concurrent login detected, retrying... userID=%d, username=%s, attempt=%d", user.ID, req.Username, i+1)
				time.Sleep(50 * time.Millisecond)
				continue
			}
			// 其他错误直接返回
			logger.Errorf("User login failed: update session failed, userID=%d, username=%s, error=%v", user.ID, req.Username, err)
			return "", errors.New(errors.UpdateSessionFailedCode)
		}
		// 更新成功，退出循环
		updateErr = nil
		break
	}

	if updateErr != nil {
		logger.Warnf("User login failed: failed to update session after retries, userID=%d, username=%s", user.ID, req.Username)
		return "", errors.New(errors.UpdateSessionFailedCode)
	}

	// 5. 生成 JWT Token（包含 session）
	token, err := s.GenerateToken(user, loginSession)
	if err != nil {
		logger.Errorf("User login failed: generate token failed, userID=%d, username=%s, error=%v", user.ID, req.Username, err)
		return "", errors.New(errors.GenerateTokenFailedCode)
	}

	logger.Debugf("User login successful: userID=%d, username=%s", user.ID, req.Username)
	return token, nil
}

// GenerateToken 生成 JWT Token
func (s *AuthService) GenerateToken(user *entity.User, loginSession string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":       user.ID,
		"username":      user.Username,
		"role":          user.Role,
		"login_session": loginSession,
		"exp":           time.Now().Add(time.Hour * time.Duration(s.expireHours)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// ParseToken 解析 JWT Token
func (s *AuthService) ParseToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}

	return nil, errors.New(errors.InvalidTokenCode)
}

// Logout 用户登出
func (s *AuthService) Logout(userID uint) error {
	// 清空 login_session，使所有该用户的 token 失效
	// 登出不需要时间戳控制，传入零值时间
	if err := s.userRepo.UpdateLoginSession(userID, "", time.Time{}); err != nil {
		logger.Errorf("User logout failed: update session failed, userID=%d, error=%v", userID, err)
		return errors.New(errors.LogoutFailedCode)
	}

	logger.Debugf("User logout successfully: userID=%d", userID)
	return nil
}

// ForgotPassword 忘记密码 - 生成重置令牌
func (s *AuthService) ForgotPassword(email string) (string, error) {
	// 1. 查询用户
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		logger.Warnf("Forgot password failed: email not found, email=%s", email)
		return "", errors.New(errors.UserNotFoundCode)
	}

	// 2. 生成重置令牌
	resetToken := uuid.New().String()
	expiresAt := time.Now().Add(1 * time.Hour) // 1小时后过期

	// 3. 保存令牌到数据库
	if err := s.userRepo.UpdateResetToken(user.ID, resetToken, expiresAt); err != nil {
		logger.Errorf("Forgot password failed: update reset token failed, email=%s, error=%v", email, err)
		return "", errors.New(errors.UpdateFailedCode)
	}

	// 4. 记录日志（用于调试和审计）
	logger.Infof("Password reset requested: email=%s, token=%s, expires_at=%s",
		email, resetToken, expiresAt.Format("2006-01-02 15:04:05"))

	logger.Debugf("Forgot password successful: email=%s", email)
	return resetToken, nil
}

// ResetPassword 重置密码
func (s *AuthService) ResetPassword(token, newPassword string) error {
	// 1. 验证密码强度
	if err := validator.ValidatePassword(newPassword); err != nil {
		return err
	}

	// 2. 根据令牌查找用户
	user, err := s.userRepo.GetByResetToken(token)
	if err != nil {
		logger.Warnf("Reset password failed: invalid token, token=%s", token)
		return errors.New(errors.InvalidTokenCode)
	}

	// 3. 检查令牌是否过期
	if user.ResetTokenExpiresAt == nil || time.Now().After(*user.ResetTokenExpiresAt) {
		logger.Warnf("Reset password failed: token expired, token=%s", token)
		return errors.New(errors.TokenExpiredCode)
	}

	// 4. 加密新密码
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		logger.Errorf("Reset password failed: password hashing failed, userID=%d, error=%v", user.ID, err)
		return errors.New(errors.PasswordHashFailedCode)
	}

	// 5. 更新密码
	if err := s.userRepo.UpdatePassword(user.ID, hashedPassword); err != nil {
		logger.Errorf("Reset password failed: update password failed, userID=%d, error=%v", user.ID, err)
		return errors.New(errors.UpdateFailedCode)
	}

	// 6. 清空重置令牌
	if err := s.userRepo.ClearResetToken(user.ID); err != nil {
		logger.Errorf("Reset password failed: clear reset token failed, userID=%d, error=%v", user.ID, err)
		return errors.New(errors.UpdateFailedCode)
	}

	// 7. 清空登录会话（强制用户重新登录）
	if err := s.userRepo.UpdateLoginSession(user.ID, "", time.Time{}); err != nil {
		logger.Warnf("Reset password: failed to clear login session, userID=%d, error=%v", user.ID, err)
	}

	logger.Infof("Reset password successful: userID=%d, username=%s", user.ID, user.Username)
	return nil
}
