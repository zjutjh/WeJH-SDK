package zapHelper

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// InfoConfig 用于定义日志配置的结构体
type InfoConfig struct {
	DisableStacktrace bool   // 是否禁用堆栈跟踪
	StacktraceLevel   string // 堆栈跟踪日志级别
	ConsoleLevel      string // 控制台日志级别
	Name              string // 日志名称
	Writer            string // 日志输出方式 "console" "file" 或 "all"
	LoggerDir         string // 日志文件目录
	LogMaxSize        int    // 日志文件最大大小（单位：MB）
	LogMaxAge         int    // 日志文件最大保存天数
	LogCompress       bool   // 是否压缩日志
}

// loggerLevelMap 映射日志级别字符串到 zapcore.Level
var loggerLevelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

const (
	// WriterConsole 表示控制台输出
	WriterConsole = "console"

	// WriterFile 表示文件输出
	WriterFile = "file"

	// WriterAll 表示同时输出到控制台和文件
	WriterAll = "all"
)

// GetDefaultConfig 返回默认的日志配置
func GetDefaultConfig() *InfoConfig {
	return &InfoConfig{
		ConsoleLevel:      "info",
		DisableStacktrace: false,
		LogCompress:       false,
		LogMaxAge:         7,
		LogMaxSize:        10,
		LoggerDir:         "./logs",
		Name:              "default",
		StacktraceLevel:   "error",
		Writer:            WriterAll,
	}
}

// Init 初始化 Zap 日志记录器
func Init(cfg *InfoConfig) (*zap.Logger, error) {
	cfg.LoggerDir = strings.TrimRight(cfg.LoggerDir, "/ ") // 去除尾部斜杠和空格

	// 创建日志目录
	if err := os.MkdirAll(cfg.LoggerDir, 0750); err != nil {
		return nil, err
	}
	encoder := createEncoder()

	var cores []zapcore.Core
	options := []zap.Option{zap.Fields(zap.String("serviceName", cfg.Name))}

	// 根据配置选择输出方式
	cores = append(cores, createLogCores(cfg, encoder)...)
	combinedCore := zapcore.NewTee(cores...)

	// 添加堆栈跟踪
	if !cfg.DisableStacktrace {
		options = append(options, zap.AddStacktrace(zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= GetZapLevel(cfg.StacktraceLevel)
		})))
	}

	logger := zap.New(combinedCore, options...) // 创建新的日志记录器
	return logger, nil
}

// GetZapLevel 返回日志级别
func GetZapLevel(l string) zapcore.Level {
	level, exist := loggerLevelMap[strings.ToLower(l)]
	if !exist {
		return zapcore.InfoLevel // 默认返回 Info 级别
	}
	return level
}

// getFileCore 返回一个把所有级别日志输出到文件的核心
func getFileCore(encoder zapcore.Encoder, cfg *InfoConfig) zapcore.Core {
	allWriter := &lumberjack.Logger{ // lumberjack 接管日志滚动
		Filename: cfg.LoggerDir + "/" + cfg.Name + ".log",
		MaxSize:  cfg.LogMaxSize,
		MaxAge:   cfg.LogMaxAge,
		Compress: cfg.LogCompress,
	}
	allLevel := zap.LevelEnablerFunc(func(_ zapcore.Level) bool {
		return true // 记录所有级别
	})
	return zapcore.NewCore(encoder, zapcore.AddSync(allWriter), allLevel)
}

// createEncoder 创建日志编码器
func createEncoder() zapcore.Encoder {
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	return zapcore.NewJSONEncoder(encoderCfg) // JSON 编码器
}

// createLogCores 创建日志核心
func createLogCores(cfg *InfoConfig, encoder zapcore.Encoder) []zapcore.Core {
	var cores []zapcore.Core
	if cfg.Writer == WriterConsole || cfg.Writer == WriterAll {
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), GetZapLevel(cfg.ConsoleLevel)))
	}
	if cfg.Writer == WriterFile || cfg.Writer == WriterAll {
		cores = append(cores, getFileCore(encoder, cfg))
	}
	return cores
}
