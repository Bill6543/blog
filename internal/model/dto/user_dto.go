package dto

import "time"

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `form:"username" binding:"required,min=3,max=50"`
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required,min=6,max=50"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserResponse 用户响应（已登录用户可见，仅包含基本信息）
type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
	Role     string `json:"role"`
}

// PublicUserResponse 公开场景的用户响应（无需登录）
type PublicUserResponse struct {
	ID       uint   `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
}

// PublicUserListResponse 公开场景的用户列表响应
type PublicUserListResponse struct {
	Total int64                `json:"total"`
	List  []PublicUserResponse `json:"list"`
}

// UpdateUserRequest 更新用户请求（使用指针字段区分不更新和清空）
type UpdateUserRequest struct {
	Nickname *string `json:"nickname,omitempty" binding:"omitempty,min=3,max=50"`
	Avatar   *string `json:"avatar,omitempty" binding:"omitempty,max=255"`
	Bio      *string `json:"bio,omitempty" binding:"omitempty,max=500"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=6,max=50"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=50"`
}

// 用户列表查询请求
type UserListRequest struct {
	Page     int    `json:"page" form:"page" binding:"min=1"`
	PageSize int    `json:"page_size" form:"page_size" binding:"min=1,max=100"`
	Role     string `json:"role" form:"role"`
	Status   *int   `json:"status" form:"status"`
	Keyword  string `json:"keyword" form:"keyword"`
}

// 用户列表响应
type UserListResponse struct {
	Total int64          `json:"total"`
	List  []UserResponse `json:"list"`
}

// 更新用户信息请求（管理员专用，支持修改角色和状态）
type AdminUpdateUserRequest struct {
	Nickname *string `json:"nickname,omitempty" binding:"omitempty,min=3,max=50"`
	Avatar   *string `json:"avatar,omitempty" binding:"omitempty,max=255"`
	Bio      *string `json:"bio,omitempty" binding:"omitempty,max=500"`
	Role     *string `json:"role,omitempty" binding:"omitempty,oneof=admin user"`
	Status   *int    `json:"status,omitempty" binding:"omitempty,oneof=0 1"`
}

// AdminUserResponse 管理员专用的用户响应（包含完整信息）
type AdminUserResponse struct {
	ID            uint       `json:"id"`
	Username      string     `json:"username"`
	Email         string     `json:"email"`
	Nickname      string     `json:"nickname"`
	Avatar        string     `json:"avatar"`
	Bio           string     `json:"bio"`
	Role          string     `json:"role"`
	Status        int        `json:"status"`
	LastLoginTime *time.Time `json:"last_login_time"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// AdminUserListResponse 管理员专用的用户列表响应
type AdminUserListResponse struct {
	Total int64               `json:"total"`
	List  []AdminUserResponse `json:"list"`
}

// ForgotPasswordRequest 忘记密码请求
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=50"`
}
