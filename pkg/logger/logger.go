package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"time"
)

var Logger *zap.SugaredLogger

// InitLogger 初始化日志
func InitLogger(level string, filePath string, maxSize, maxBackups, maxAge int) error {
	// ✅ 确保日志目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// 解析日志级别
	var logLevel zapcore.Level
	if err := logLevel.UnmarshalText([]byte(level)); err != nil {
		logLevel = zapcore.DebugLevel
	}

	// 配置日志轮转（增强版：支持日期 + 大小轮转）
	lumberJackLogger := &lumberjack.Logger{
		Filename:   getLogFilePath(filePath),
		MaxSize:    maxSize,    // 使用配置参数
		MaxBackups: maxBackups, // 使用配置参数
		Compress:   true,       // 压缩旧文件
		LocalTime:  true,       // 使用本地时间
	}

	// 配置编码器（生产环境 JSON 格式）
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.LevelKey = "level"
	encoderConfig.MessageKey = "message"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // ✅ 彩色级别

	// 创建多个输出核心
	var cores []zapcore.Core

	// 1. 文件输出（JSON 格式）
	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(lumberJackLogger),
		logLevel,
	)
	cores = append(cores, fileCore)

	// 2. 控制台输出（彩色格式，仅开发环境）
	if os.Getenv("GIN_MODE") != "release" {
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		consoleCore := zapcore.NewCore(
			consoleEncoder,
			zapcore.AddSync(os.Stdout),
			logLevel,
		)
		cores = append(cores, consoleCore)
	}

	// 合并核心
	combinedCore := zapcore.NewTee(cores...)

	// 创建 logger
	logger := zap.New(combinedCore, zap.AddCaller(), zap.AddCallerSkip(1))
	Logger = logger.Sugar()

	Logger.Infof("Logger initialized with level: %s, file: %s", level, filePath)
	return nil
}

// getLogFilePath 生成带日期的日志文件路径
func getLogFilePath(filePath string) string {
	dir := filepath.Dir(filePath)
	filename := filepath.Base(filePath)
	ext := filepath.Ext(filename)
	name := filename[:len(filename)-len(ext)]

	// 生成格式：logs/app-YYYY-MM-DD.log
	today := time.Now().Format("2006-01-02")
	return filepath.Join(dir, fmt.Sprintf("%s-%s%s", name, today, ext))
}

// Debug 调试日志
func Debug(args ...interface{}) {
	Logger.Debug(args...)
}

// Info 信息日志
func Info(args ...interface{}) {
	Logger.Info(args...)
}

// Warn 警告日志
func Warn(args ...interface{}) {
	Logger.Warn(args...)
}

// Error 错误日志
func Error(args ...interface{}) {
	Logger.Error(args...)
}

// Debugf 格式化调试日志
func Debugf(format string, args ...interface{}) {
	Logger.Debugf(format, args...)
}

// Infof 格式化信息日志
func Infof(format string, args ...interface{}) {
	Logger.Infof(format, args...)
}

// Warnf 格式化警告日志
func Warnf(format string, args ...interface{}) {
	Logger.Warnf(format, args...)
}

// Errorf 格式化错误日志
func Errorf(format string, args ...interface{}) {
	Logger.Errorf(format, args...)
}

// Fatal 致命错误日志
// ⚠️ 警告：此方法会立即退出程序，不会执行 defer 中的清理操作
// 建议使用 Errorf + os.Exit 的组合代替
func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}

// Fatalf 格式化致命错误日志
// ⚠️ 警告：此方法会立即退出程序，不会执行 defer 中的清理操作
// 建议使用 Errorf + os.Exit 的组合代替
func Fatalf(format string, args ...interface{}) {
	Logger.Fatalf(format, args...)
}

// Sync 同步日志
func Sync() error {
	return Logger.Sync()
}

// Close 关闭日志（程序退出前调用）
func Close() error {
	if Logger != nil {
		return Logger.Sync()
	}
	return nil
}
