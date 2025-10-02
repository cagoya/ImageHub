package global

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogFileConfig struct {
	MaxSize  int      `mapstructure:"maxSize"`
	Backups  int      `mapstructure:"backups"`
	Compress bool     `mapstructure:"compress"`
	Output   []string `mapstructure:"output"`
}

type ZapConfig struct {
	Prefix     string        `mapstructure:"prefix"`
	TimeFormat string        `mapstructure:"timeFormat"`
	Level      string        `mapstructure:"level"`
	Encoding   string        `mapstructure:"encoding"`
	Caller     bool          `mapstructure:"caller"`
	Stacktrace bool          `mapstructure:"stackTrace"`
	Writer     string        `mapstructure:"writer"`
	LogFile    LogFileConfig `mapstructure:"logFile"`
}

var logger *zap.Logger

func InitZap() error {
	var zapConfig ZapConfig
	if err := viper.UnmarshalKey("zap", &zapConfig); err != nil {
		return fmt.Errorf("解析zap配置失败：%v", err)
	}

	// 设置日志级别
	level, err := zapcore.ParseLevel(zapConfig.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}

	// 设置编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout(zapConfig.TimeFormat),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 设置编码器
	var encoder zapcore.Encoder
	if zapConfig.Encoding == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 设置输出目标
	var cores []zapcore.Core

	// 控制台输出
	if zapConfig.Writer == "console" || zapConfig.Writer == "both" {
		consoleCore := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level)
		cores = append(cores, consoleCore)
	}

	// 文件输出
	if zapConfig.Writer == "file" || zapConfig.Writer == "both" {
		for _, output := range zapConfig.LogFile.Output {
			// 确保日志目录存在
			dir := filepath.Dir(output)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("创建日志目录失败：%v", err)
			}

			// 设置日志轮转
			fileWriter := &lumberjack.Logger{
				Filename:   output,
				MaxSize:    zapConfig.LogFile.MaxSize,
				MaxBackups: zapConfig.LogFile.Backups,
				MaxAge:     30, // 保留30天
				Compress:   zapConfig.LogFile.Compress,
			}

			fileCore := zapcore.NewCore(encoder, zapcore.AddSync(fileWriter), level)
			cores = append(cores, fileCore)
		}
	}

	// 合并所有输出
	core := zapcore.NewTee(cores...)

	// 创建logger
	var options []zap.Option
	if zapConfig.Caller {
		options = append(options, zap.AddCaller())
	}
	if zapConfig.Stacktrace {
		options = append(options, zap.AddStacktrace(zapcore.ErrorLevel))
	}

	logger = zap.New(core, options...)
	return nil
}

func GetLogger() *zap.Logger {
	return logger
}
