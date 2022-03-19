package servicelog

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	logger *zap.SugaredLogger
}

//LogFileConfig : log file options
type LogFileConfig struct {
	FileAddress string
	MaxSize     int
	MaxBackups  int
	MaxAge      int
}

//New :initial looger object
func New(level string, logFileConfig *LogFileConfig, stdout bool) Logger {
	srvLogger := Logger{}
	var loglevel zapcore.Level
	//set loglevel
	switch level {
	case "info":
		loglevel = zapcore.InfoLevel
	case "debug":
		loglevel = zapcore.DebugLevel
	case "error":
		loglevel = zapcore.ErrorLevel
	case "warn":
		loglevel = zapcore.WarnLevel
	case "fatal":
		loglevel = zapcore.FatalLevel
	default:
		loglevel = zapcore.DebugLevel
	}
	// -------- log endocing config ---------//
	encodConfig := zapcore.NewJSONEncoder(
		zapcore.EncoderConfig{
			MessageKey:  "message",
			TimeKey:     "timestamp",
			EncodeTime:  zapcore.ISO8601TimeEncoder,
			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,
		})

	zapCoreOpt := []zapcore.Core{}
	if stdout {
		consoleErrors := zapcore.Lock(os.Stderr)
		zapCoreOpt = append(zapCoreOpt, zapcore.NewCore(encodConfig, consoleErrors, zap.NewAtomicLevelAt(loglevel)))
	}
	if logFileConfig != nil {
		fileLimberJackOutput := zapcore.AddSync(&lumberjack.Logger{
			Filename:   logFileConfig.FileAddress,
			MaxSize:    logFileConfig.MaxSize, // megabytes
			MaxBackups: logFileConfig.MaxBackups,
			MaxAge:     logFileConfig.MaxAge, // days
		})
		zapCoreOpt = append(zapCoreOpt, zapcore.NewCore(encodConfig, fileLimberJackOutput, zap.NewAtomicLevelAt(loglevel)))
	}
	core := zapcore.NewTee(zapCoreOpt...)
	zlogger := zap.New(core)
	sugar := zlogger.Sugar()
	srvLogger.logger = sugar
	return srvLogger
}

//***********************************************************************//

func (lg Logger) Infolog(msg string) {
	lg.logger.Infow(msg)
}
func (lg Logger) Warnlog(msg string) {
	lg.logger.Warnw(msg)
}
func (lg Logger) Debuglog(msg string) {
	lg.logger.Debugw(msg)
}
func (lg Logger) ErrorLog(msg string) {
	lg.logger.Errorw(msg)
}
func (lg Logger) FatalLog(msg string) {
	lg.logger.Fatalw(msg)
}
func (lg Logger) AccessLog(Url string, Method string, RemoteAddr string) {
	lg.logger.Infow("Access-log",
		zap.String("url", Url),
		zap.String("method", Method),
		zap.String("remote-ip", RemoteAddr))
}
func (lg Logger) HttpErrorLog(Url string, Method string, RemoteAddr string, HttpCode int, Details string) {
	lg.logger.Errorw("request error",
		zap.String("url", Url),
		zap.String("method", Method),
		zap.String("remote-ip", RemoteAddr),
		zap.Int("response-code", HttpCode),
		zap.String("details", Details))
}
