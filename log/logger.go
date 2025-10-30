package logger

import "go.uber.org/zap"

type Logger struct {
	*zap.Logger
}

// NewLogger 创建一个新的 Logger 实例（默认生产模式，输出到 stdout）
func NewLogger() *Logger {
	// 生产环境配置：高性能 JSON 格式
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"} // 输出到控制台，K8s 可直接捕获
	config.ErrorOutputPaths = []string{"stderr"}

	// 可选：如果希望日志中带颜色（开发调试用），可以使用 DevelopmentConfig()
	// config := zap.NewDevelopmentConfig()

	zapLogger, err := config.Build()
	if err != nil {
		panic("failed to initialize zap logger: " + err.Error())
	}

	return &Logger{zapLogger}
}

// SyncLogger 在程序退出前调用，确保日志全部刷出（比如在 defer 中调用）
func (l *Logger) SyncLogger() {
	_ = l.Sync() // ignore error
}

func (l *Logger) InfoWith(msg string, fields ...zap.Field) {
	l.Info(msg, fields...)
}

func (l *Logger) ErrorWith(msg string, fields ...zap.Field) {
	l.Error(msg, fields...)
}

func (l *Logger) DebugWith(msg string, fields ...zap.Field) {
	l.Debug(msg, fields...)
}

func (l *Logger) WarnWith(msg string, fields ...zap.Field) {
	l.Warn(msg, fields...)
}
