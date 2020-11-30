package logger

import (
	"bluebell/settings"
	"net"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Init(cfg *settings.AppConfig) (err error) {
	writeSyncer := getWriteSyncer(cfg.Log.Filename, cfg.Log.MaxSize,
		cfg.Log.MaxSize, cfg.Log.MaxBackups, false)

	encoder := getEncoder()

	var level = new(zapcore.Level)
	err = level.UnmarshalText([]byte(viper.GetString("log.level")))
	if err != nil {
		return
	}
	var core zapcore.Core
	if cfg.App.Mode == "dev" {
		//进入开发模式，日志输出到终端
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

		core = zapcore.NewTee(
			zapcore.NewCore(encoder, writeSyncer, level),
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel),
		)
	} else {
		core = zapcore.NewCore(encoder, writeSyncer, level)
	}

	lg := zap.New(core, zap.AddCaller())
	//替换zab库中全局logger
	zap.ReplaceGlobals(lg)
	return
}

func getEncoder() zapcore.Encoder {
	encoderCofnig := zap.NewProductionEncoderConfig()
	encoderCofnig.EncodeTime = zapcore.ISO8601TimeEncoder
	return zapcore.NewJSONEncoder(encoderCofnig)
}

func getWriteSyncer(filename string, maxSize int, maxAge int,
	maxBackups int, compress bool) zapcore.WriteSyncer {

	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxAge:     maxAge,
		MaxBackups: maxBackups,
		Compress:   compress,
	}
	return zapcore.AddSync(lumberJackLogger)
}

//GinLogger,接收gin框架的默认日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		zap.L().Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost))
	}
}

//GinReconery recover掉项目可能出现的painc
func GinReconvery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				//check for a broken connection,as it is not realy a
				//condition that warrants a painc stack trace
				var brokenPipe bool
				if ne, ok := err.(net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					zap.L().Error(c.Request.URL.Path, zap.Any("error", err),
						zap.String("request", string(httpRequest)))
					//If the connection is dead,we can't write a status to it
					c.Error(err.(error))
					c.Abort()
					return
				}
				if stack {
					zap.L().Error("[Recovery from painc]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					zap.L().Error("[Recovery from painc]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)))
				}
			}
		}()
		c.Next()
	}
}
