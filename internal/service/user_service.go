package service

import (
	"blog/internal/cache"
	"blog/internal/model/dto"
	"blog/internal/model/entity"
	"blog/internal/repository"
	"blog/pkg/errors"
	"blog/pkg/logger"
	"blog/pkg/utils"
	"blog/pkg/validator"
	"time"
)

type UserService struct {
	userRepo     *repository.UserRepository
	sessionCache *cache.SessionCache
}

func NewUserService(userRepo *repository.UserRepository, sessionCache *cache.SessionCache) *UserService {
	return &UserService{
		userRepo:     userRepo,
		sessionCache: sessionCache,
	}
}

// convertToPublicUserResponse 转换为公开用户响应
func (s *UserService) convertToPublicUserResponse(user *entity.User) *dto.PublicUserResponse {
	return &dto.PublicUserResponse{
		ID:       user.ID,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Bio:      user.Bio,
	}
}

// convertToUserResponse 转换为用户响应（根据场景决定是否包含 Email）
func (s *UserService) convertToUserResponse(user *entity.User, includeEmail bool) *dto.UserResponse {
	resp := &dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Bio:      user.Bio,
		Role:     user.Role,
	}

	if includeEmail {
		resp.Email = user.Email
	}

	return resp
}

// convertToAdminUserResponse 转换为管理员用户响应（包含完整信息）
func (s *UserService) convertToAdminUserResponse(user *entity.User) *dto.AdminUserResponse {
	return &dto.AdminUserResponse{
		ID:            user.ID,
		Username:      user.Username,
		Email:         user.Email,
		Nickname:      user.Nickname,
		Avatar:        user.Avatar,
		Bio:           user.Bio,
		Role:          user.Role,
		Status:        user.Status,
		LastLoginTime: user.LastLoginTime,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}
}

// GetPublicUserInfo 获取公开用户信息（无需登录）
func (s *UserService) GetPublicUserInfo(targetUserID uint) (*dto.PublicUserResponse, error) {
	user, err := s.userRepo.GetByID(targetUserID)
	if err != nil {
		logger.Warnf("Get public user info failed: user not found, userID=%d", targetUserID)
		return nil, errors.New(errors.UserNotFoundCode)
	}

	return s.convertToPublicUserResponse(user), nil
}

// GetUserInfoForSelf 查看自己的用户信息（包含 Email）
func (s *UserService) GetUserInfoForSelf(userID uint) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		logger.Warnf("Get self user info failed: user not found, userID=%d", userID)
		return nil, errors.New(errors.UserNotFoundCode)
	}

	return s.convertToUserResponse(user, true), nil
}

// GetUserInfoForOther 查看他人的用户信息（不包含 Email）
func (s *UserService) GetUserInfoForOther(targetUserID uint) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(targetUserID)
	if err != nil {
		logger.Warnf("Get other user info failed: user not found, userID=%d", targetUserID)
		return nil, errors.New(errors.UserNotFoundCode)
	}

	return s.convertToUserResponse(user, false), nil
}

// GetUserInfoForAdmin 管理员查看用户完整信息
func (s *UserService) GetUserInfoForAdmin(targetUserID uint) (*dto.AdminUserResponse, error) {
	user, err := s.userRepo.GetByID(targetUserID)
	if err != nil {
		logger.Warnf("Get admin user info failed: user not found, userID=%d", targetUserID)
		return nil, errors.New(errors.UserNotFoundCode)
	}

	return s.convertToAdminUserResponse(user), nil
}

// GetPublicUserList 获取公开用户列表（普通用户）
func (s *UserService) GetPublicUserList(req dto.UserListRequest) (*dto.PublicUserListResponse, error) {
	users, total, err := s.userRepo.GetList(req.Page, req.PageSize, req.Role, req.Status, req.Keyword)
	if err != nil {
		logger.Errorf("Get public user list failed: error=%v", err)
		return nil, errors.New(errors.DatabaseError)
	}

	var list []dto.PublicUserResponse
	for _, user := range users {
		list = append(list, *s.convertToPublicUserResponse(&user))
	}

	return &dto.PublicUserListResponse{
		Total: total,
		List:  list,
	}, nil
}

// GetUserListForAdmin 管理员获取用户列表（包含完整信息）
func (s *UserService) GetUserListForAdmin(req dto.UserListRequest) (*dto.AdminUserListResponse, error) {
	users, total, err := s.userRepo.GetList(req.Page, req.PageSize, req.Role, req.Status, req.Keyword)
	if err != nil {
		logger.Errorf("Get admin user list failed: error=%v", err)
		return nil, errors.New(errors.DatabaseError)
	}

	var list []dto.AdminUserResponse
	for _, user := range users {
		list = append(list, *s.convertToAdminUserResponse(&user))
	}

	return &dto.AdminUserListResponse{
		Total: total,
		List:  list,
	}, nil
}

// UpdateUserInfo 更新用户信息（支持角色判断：管理员可改角色/状态，普通用户只能改资料）
func (s *UserService) UpdateUserInfo(targetUserID uint, req dto.AdminUpdateUserRequest, viewerUserID uint, viewerRole string) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(targetUserID)
	if err != nil {
		logger.Warnf("Update user info failed: user not found, userID=%d", targetUserID)
		return nil, errors.New(errors.UserNotFoundCode)
	}

	// 权限校验：普通用户只能修改自己，管理员可以修改任何人
	if viewerRole != "admin" && targetUserID != viewerUserID {
		logger.Warnf("Update user info failed: permission denied, target=%d, viewer=%d", targetUserID, viewerUserID)
		return nil, errors.New(errors.PermissionDenied)
	}

	// 更新字段
	if req.Nickname != nil {
		user.Nickname = *req.Nickname
	}
	if req.Avatar != nil {
		user.Avatar = *req.Avatar
	}
	if req.Bio != nil {
		user.Bio = *req.Bio
	}

	// 管理员专属字段
	if viewerRole == "admin" {
		if req.Role != nil {
			user.Role = *req.Role
		}
		if req.Status != nil {
			// 禁用用户时清空 session
			if *req.Status == 0 {
				if err := s.userRepo.UpdateLoginSession(user.ID, "", time.Time{}); err != nil {
					logger.Errorf("Failed to clear session when disabling user: userID=%d, error=%v", user.ID, err)
					return nil, errors.New(errors.SessionUpdateFailedCode)
				}
				// 清除会话缓存
				if s.sessionCache != nil {
					_ = s.sessionCache.Del(user.ID)
				}
			}
			user.Status = *req.Status
		}
	}

	if err := s.userRepo.Update(user); err != nil {
		logger.Errorf("Update user info failed: database error, userID=%d, error=%v", targetUserID, err)
		return nil, errors.New(errors.DatabaseError)
	}

	return s.convertToUserResponse(user, true), nil
}

// DeleteUser 删除用户（仅管理员）
func (s *UserService) DeleteUser(targetUserID uint, operatorRole string) error {
	if operatorRole != "admin" {
		logger.Warnf("Delete user failed: permission denied, operator=%s", operatorRole)
		return errors.New(errors.PermissionDenied)
	}

	user, err := s.userRepo.GetByID(targetUserID)
	if err != nil {
		logger.Warnf("Delete user failed: user not found, userID=%d", targetUserID)
		return errors.New(errors.UserNotFoundCode)
	}

	// 软删除
	if err := s.userRepo.SoftDelete(user.ID); err != nil {
		logger.Errorf("Delete user failed: database error, userID=%d, error=%v", targetUserID, err)
		return errors.New(errors.DatabaseError)
	}

	// 清空登录会话，使被删除用户的 token 立即失效（修复删除后仍可访问的问题）
	if err := s.userRepo.UpdateLoginSession(user.ID, "", time.Time{}); err != nil {
		logger.Warnf("Delete user: failed to clear login session, userID=%d, error=%v", user.ID, err)
	}
	if s.sessionCache != nil {
		_ = s.sessionCache.Del(user.ID)
	}

	return nil
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(userID uint, req dto.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		logger.Warnf("Change password failed: user not found, userID=%d", userID)
		return errors.New(errors.UserNotFoundCode)
	}

	// 验证旧密码
	if !utils.CheckPassword(req.OldPassword, user.Password) {
		logger.Warnf("Change password failed: incorrect old password, userID=%d", userID)
		return errors.New(errors.OldPasswordWrongCode)
	}

	// 验证新密码强度
	if err := validator.ValidatePassword(req.NewPassword); err != nil {
		logger.Warnf("Change password failed: weak password, userID=%d", userID)
		return err
	}

	// 加密新密码
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		logger.Errorf("Change password failed: hash password failed, userID=%d, error=%v", userID, err)
		return errors.New(errors.PasswordHashFailedCode)
	}

	// 清除 login_session（让所有已登录的设备下线）
	// 修改密码不需要时间戳控制，传入零值时间
	if err := s.userRepo.UpdateLoginSession(userID, "", time.Time{}); err != nil {
		logger.Errorf("Change password failed: update session failed, userID=%d, error=%v", userID, err)
		return errors.New(errors.SessionUpdateFailedCode)
	}

	// 清除会话缓存
	if s.sessionCache != nil {
		_ = s.sessionCache.Del(userID)
	}

	// 更新密码
	user.Password = hashedPassword
	if err := s.userRepo.Update(user); err != nil {
		logger.Errorf("Change password failed: database error, userID=%d, error=%v", userID, err)
		return errors.New(errors.DatabaseError)
	}

	return nil
}
