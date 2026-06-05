package errors

import "fmt"

// Error level definition
type ErrorLevel int

const (
	LevelInfo ErrorLevel = iota
	LevelWarn
	LevelError
	LevelFatal
)

// Error code definition
const (
	Success           = 0
	InternalError     = 1001
	InvalidParam      = 1002
	Unauthorized      = 1003
	Forbidden         = 1004
	NotFound          = 1005
	DuplicateEntry    = 1006
	RecordNotFound    = 1007
	PasswordIncorrect = 1008
	TokenExpired      = 1009
	TokenInvalid      = 1010
	PermissionDenied  = 1011
	FileUploadFailed  = 1012
	FileNotFound      = 1013
	DatabaseError     = 1014
	CacheError        = 1015

	// User related errors (2001-2099)
	UserNotFoundCode        = 2001
	OldPasswordWrongCode    = 2002
	PasswordHashFailedCode  = 2003
	SessionUpdateFailedCode = 2004

	// Auth related errors (2101-2199)
	UsernameExistsCode      = 2101
	EmailExistsCode         = 2102
	UserDisabledCode        = 2103
	InvalidCredentialsCode  = 2104
	UpdateSessionFailedCode = 2105
	GenerateTokenFailedCode = 2106
	InvalidTokenCode        = 2107
	LogoutFailedCode        = 2108
	RegistrationFailedCode  = 2109

	// Article related errors (2201-2299)
	ArticleNotFoundCode  = 2201
	ArticleNotAuthorCode = 2202
	CategoryNotFoundCode = 2203
	CategoryExistsCode   = 2204
	CategoryInUseCode    = 2205
	TagNotFoundCode      = 2206
	TagNameExistsCode    = 2207

	// Comment related errors (2301-2399)
	CommentNotAuthorCode      = 2301
	ParentCommentNotFoundCode = 2302
	CommentNotFoundCode       = 2303

	// Like related errors (2401-2499)
	AlreadyLikedCode = 2401
	NotLikedYetCode  = 2402

	// File upload related errors (2501-2599)
	InvalidFileTypeCode = 2501
	FileTooLargeCode    = 2502
	FileOpenFailedCode  = 2503
	FileSaveFailedCode  = 2504
)

// 错误消息映射
var messages = map[int]string{
	Success:           "success",
	InternalError:     "Internal server error",
	InvalidParam:      "Invalid parameter",
	Unauthorized:      "Unauthorized access",
	Forbidden:         "Forbidden",
	NotFound:          "Resource not found",
	DuplicateEntry:    "Duplicate entry",
	RecordNotFound:    "Record not found",
	PasswordIncorrect: "Password incorrect",
	TokenExpired:      "Token expired",
	TokenInvalid:      "Token invalid",
	PermissionDenied:  "Permission denied",
	FileUploadFailed:  "File upload failed",
	FileNotFound:      "File not found",
	DatabaseError:     "Database error",
	CacheError:        "Cache error",

	// User related errors
	UserNotFoundCode:        "User not found",
	OldPasswordWrongCode:    "Old password is incorrect",
	PasswordHashFailedCode:  "Password hashing failed",
	SessionUpdateFailedCode: "Failed to update session",

	// Auth related errors
	UsernameExistsCode:      "Username already exists",
	EmailExistsCode:         "Email already registered",
	UserDisabledCode:        "User is disabled",
	InvalidCredentialsCode:  "Invalid username or password",
	UpdateSessionFailedCode: "Failed to update login session",
	GenerateTokenFailedCode: "Failed to generate token",
	InvalidTokenCode:        "Invalid token",
	LogoutFailedCode:        "Failed to logout",
	RegistrationFailedCode:  "Registration failed, please try again later",

	// Article related errors
	ArticleNotFoundCode:  "Article not found",
	ArticleNotAuthorCode: "Only author can modify this article",
	CategoryNotFoundCode: "Category not found",
	CategoryExistsCode:   "Category name already exists",
	CategoryInUseCode:    "Category is in use by articles",
	TagNotFoundCode:      "Tag not found",
	TagNameExistsCode:    "Tag name already exists",

	// Comment related errors
	CommentNotAuthorCode:      "Only author can modify this comment",
	ParentCommentNotFoundCode: "Parent comment not found",
	CommentNotFoundCode:       "Comment not found",

	// Like related errors
	AlreadyLikedCode: "Already liked",
	NotLikedYetCode:  "Not liked yet",

	// File upload related errors
	InvalidFileTypeCode: "Invalid file type",
	FileTooLargeCode:    "File too large",
	FileOpenFailedCode:  "Failed to open file",
	FileSaveFailedCode:  "Failed to save file",
}

// Error level mapping
var errorLevels = map[int]ErrorLevel{
	Success:           LevelInfo,
	InternalError:     LevelError,
	InvalidParam:      LevelWarn,
	Unauthorized:      LevelWarn,
	Forbidden:         LevelWarn,
	NotFound:          LevelWarn,
	DuplicateEntry:    LevelWarn,
	RecordNotFound:    LevelWarn,
	PasswordIncorrect: LevelWarn,
	TokenExpired:      LevelWarn,
	TokenInvalid:      LevelWarn,
	PermissionDenied:  LevelWarn,
	FileUploadFailed:  LevelError,
	FileNotFound:      LevelWarn,
	DatabaseError:     LevelError,
	CacheError:        LevelError,

	// User related errors
	UserNotFoundCode:        LevelWarn,
	OldPasswordWrongCode:    LevelWarn,
	PasswordHashFailedCode:  LevelError,
	SessionUpdateFailedCode: LevelError,

	// Auth related errors
	UsernameExistsCode:      LevelWarn,
	EmailExistsCode:         LevelWarn,
	UserDisabledCode:        LevelWarn,
	InvalidCredentialsCode:  LevelWarn,
	UpdateSessionFailedCode: LevelError,
	GenerateTokenFailedCode: LevelError,
	InvalidTokenCode:        LevelWarn,
	LogoutFailedCode:        LevelError,
	RegistrationFailedCode:  LevelError,

	// Article related errors
	ArticleNotFoundCode:  LevelWarn,
	ArticleNotAuthorCode: LevelWarn,
	CategoryNotFoundCode: LevelWarn,
	CategoryExistsCode:   LevelWarn,
	CategoryInUseCode:    LevelWarn,
	TagNotFoundCode:      LevelWarn,
	TagNameExistsCode:    LevelWarn,

	// Comment related errors
	CommentNotAuthorCode:      LevelWarn,
	ParentCommentNotFoundCode: LevelWarn,
	CommentNotFoundCode:       LevelWarn,

	// Like related errors
	AlreadyLikedCode: LevelWarn,
	NotLikedYetCode:  LevelWarn,

	// File upload related errors
	InvalidFileTypeCode: LevelWarn,
	FileTooLargeCode:    LevelWarn,
	FileOpenFailedCode:  LevelError,
	FileSaveFailedCode:  LevelError,
}

// AppError 应用错误
type AppError struct {
	Code    int        `json:"code"`
	Message string     `json:"message"`
	Level   ErrorLevel `json:"level"`
	Err     error      `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// New create new error
func New(code int) *AppError {
	level, exists := errorLevels[code]
	if !exists {
		level = LevelError
	}
	return &AppError{
		Code:    code,
		Message: messages[code],
		Level:   level,
	}
}

// NewWithMessage create error with custom message
func NewWithMessage(code int, message string) *AppError {
	level, exists := errorLevels[code]
	if !exists {
		level = LevelError
	}
	return &AppError{
		Code:    code,
		Message: message,
		Level:   level,
	}
}

// NewWithErr create error with original error
func NewWithErr(code int, err error) *AppError {
	level, exists := errorLevels[code]
	if !exists {
		level = LevelError
	}
	return &AppError{
		Code:    code,
		Message: messages[code],
		Level:   level,
		Err:     err,
	}
}

// NewWithLevel create error with custom level
func NewWithLevel(code int, level ErrorLevel) *AppError {
	return &AppError{
		Code:    code,
		Message: messages[code],
		Level:   level,
	}
}

// Wrap 包装错误
func (e *AppError) Wrap(err error) *AppError {
	e.Err = err
	return e
}

// WithMessage 设置错误消息
func (e *AppError) WithMessage(message string) *AppError {
	e.Message = message
	return e
}

// IsNotFound 判断是否未找到错误
func IsNotFound(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == NotFound || appErr.Code == RecordNotFound
	}
	return false
}

// IsUnauthorized 判断是否未授权错误
func IsUnauthorized(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == Unauthorized
	}
	return false
}

// IsForbidden 判断是否禁止访问错误
func IsForbidden(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == Forbidden
	}
	return false
}

// IsInvalidParam 判断是否参数错误
func IsInvalidParam(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == InvalidParam
	}
	return false
}

// IsUserNotFound determine if user not found
func IsUserNotFound(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == UserNotFoundCode
	}
	return false
}

// IsPermissionDenied determine if permission denied
func IsPermissionDenied(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == PermissionDenied
	}
	return false
}

// IsOldPasswordWrong determine if old password is wrong
func IsOldPasswordWrong(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == OldPasswordWrongCode
	}
	return false
}

// IsUsernameExists determine if username already exists
func IsUsernameExists(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == UsernameExistsCode
	}
	return false
}

// IsEmailExists determine if email already registered
func IsEmailExists(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == EmailExistsCode
	}
	return false
}

// IsPasswordHashFailed determine if password hash failed
func IsPasswordHashFailed(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == PasswordHashFailedCode
	}
	return false
}

// IsRegistrationFailed determine if registration failed
func IsRegistrationFailed(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == RegistrationFailedCode
	}
	return false
}

// IsUserDisabled determine if user is disabled
func IsUserDisabled(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == UserDisabledCode
	}
	return false
}

// IsInvalidCredentials determine if credentials are invalid
func IsInvalidCredentials(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == InvalidCredentialsCode
	}
	return false
}

// IsUpdateSessionFailed determine if update session failed
func IsUpdateSessionFailed(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == UpdateSessionFailedCode
	}
	return false
}

// IsGenerateTokenFailed determine if generate token failed
func IsGenerateTokenFailed(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == GenerateTokenFailedCode
	}
	return false
}

// IsArticleNotFound determine if article not found
func IsArticleNotFound(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == ArticleNotFoundCode || appErr.Code == ArticleNotAuthorCode
	}
	return false
}

// IsCategoryNotFound determine if category not found
func IsCategoryNotFound(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == CategoryNotFoundCode
	}
	return false
}

// IsCategoryExists determine if category name already exists
func IsCategoryExists(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == CategoryExistsCode
	}
	return false
}

// IsCategoryInUse determine if category is in use by articles
func IsCategoryInUse(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == CategoryInUseCode
	}
	return false
}

// IsCommentNotFound determine if comment not found
func IsCommentNotFound(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == CommentNotFoundCode || appErr.Code == CommentNotAuthorCode
	}
	return false
}

// IsTagNotFound determine if tag not found
func IsTagNotFound(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == TagNotFoundCode
	}
	return false
}

// IsTagNameExists determine if tag name already exists
func IsTagNameExists(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == TagNameExistsCode
	}
	return false
}

// IsAlreadyLiked determine if already liked
func IsAlreadyLiked(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == AlreadyLikedCode
	}
	return false
}

// IsNotLikedYet determine if not liked yet
func IsNotLikedYet(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == NotLikedYetCode
	}
	return false
}

// IsInvalidFileType determine if file type is invalid
func IsInvalidFileType(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == InvalidFileTypeCode
	}
	return false
}

// IsFileTooLarge determine if file is too large
func IsFileTooLarge(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == FileTooLargeCode
	}
	return false
}

// Formatf format error message
func Formatf(code int, format string, args ...interface{}) *AppError {
	level, exists := errorLevels[code]
	if !exists {
		level = LevelError
	}
	return &AppError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
		Level:   level,
	}
}

// HandleValidationError 统一处理参数验证错误
// 将 Gin 的验证错误转换为友好的用户消息
func HandleValidationError(err error) string {
	// 详细错误会记录到日志（由 Handler 层负责日志）
	// 这里只返回统一的友好消息
	return "参数验证失败，请检查输入格式"
}

// GetErrorCode 根据错误类型获取错误码
func GetErrorCode(err error) int {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code
	}
	return InternalError
}

// GetErrorMessage 根据错误码获取友好的错误消息
func GetErrorMessage(code int) string {
	if msg, exists := messages[code]; exists {
		return msg
	}
	return "操作失败"
}
