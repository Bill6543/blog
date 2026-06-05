package repository

import (
	"blog/internal/model/entity"
	"blog/pkg/errors"
	"time"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create 创建用户
func (r *UserRepository) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

// GetByID 根据 ID 查询用户
func (r *UserRepository) GetByID(id uint) (*entity.User, error) {
	var user entity.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.UserNotFoundCode)
		}
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名查询
func (r *UserRepository) GetByUsername(username string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.UserNotFoundCode)
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱查询
func (r *UserRepository) GetByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.UserNotFoundCode)
		}
		return nil, err
	}
	return &user, nil
}

// Update 更新用户
func (r *UserRepository) Update(user *entity.User) error {
	return r.db.Save(user).Error
}

// Delete 删除用户
func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&entity.User{}, id).Error
}

// ExistsByUsername 检查用户名是否存在
func (r *UserRepository) ExistsByUsername(username string) bool {
	var count int64
	r.db.Model(&entity.User{}).Where("username = ?", username).Count(&count)
	return count > 0
}

// ExistsByEmail 检查邮箱是否存在
func (r *UserRepository) ExistsByEmail(email string) bool {
	var count int64
	r.db.Model(&entity.User{}).Where("email = ?", email).Count(&count)
	return count > 0
}

// UpdateLoginSession 更新用户的 login_session（带时间戳并发控制）
// - lastLoginTime 为零值时：强制更新（用于登出、修改密码）
// - lastLoginTime 为非零值时：乐观锁更新（用于登录，防止并发冲突）
func (r *UserRepository) UpdateLoginSession(userID uint, loginSession string, lastLoginTime time.Time) error {
	var result *gorm.DB

	if lastLoginTime.IsZero() {
		// 零值时间表示强制更新（登出、修改密码场景）
		// 不需要时间戳控制，直接根据 ID 更新
		result = r.db.Model(&entity.User{}).
			Where("id = ?", userID).
			Updates(map[string]interface{}{
				"login_session":   loginSession,
				"last_login_time": lastLoginTime,
			})
	} else {
		// 非零值时间表示正常登录，需要并发控制
		// 只有当新的 lastLoginTime 大于数据库中的值时才更新
		result = r.db.Model(&entity.User{}).
			Where("id = ? AND (last_login_time IS NULL OR last_login_time < ?)", userID, lastLoginTime).
			Updates(map[string]interface{}{
				"login_session":   loginSession,
				"last_login_time": lastLoginTime,
			})
	}

	if result.Error != nil {
		return result.Error
	}

	// 如果没有更新任何行，说明存在并发冲突（有其他登录请求先完成）
	if result.RowsAffected == 0 {
		return errors.New(errors.UpdateSessionFailedCode)
	}

	return nil
}

// GetList 获取用户列表（分页、筛选）
func (r *UserRepository) GetList(page, pageSize int, role string, status *int, keyword string) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64

	query := r.db.Model(&entity.User{})

	// 角色筛选
	if role != "" {
		query = query.Where("role = ?", role)
	}

	// 状态筛选
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// 关键词搜索（用户名、邮箱、昵称）
	if keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ? OR nickname LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页（按创建时间倒序）
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// UpdateRole 更新用户角色
func (r *UserRepository) UpdateRole(userID uint, role string) error {
	return r.db.Model(&entity.User{}).
		Where("id = ?", userID).
		Update("role", role).Error
}

// UpdateStatus 更新用户状态
func (r *UserRepository) UpdateStatus(userID uint, status int) error {
	return r.db.Model(&entity.User{}).
		Where("id = ?", userID).
		Update("status", status).Error
}

// SoftDelete 软删除用户（更新 status 为 -1）
func (r *UserRepository) SoftDelete(userID uint) error {
	return r.db.Model(&entity.User{}).
		Where("id = ?", userID).
		Update("status", -1).Error
}
