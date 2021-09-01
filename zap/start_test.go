// @Author : Lik
// @Time   : 2021/1/26
package zap

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"testing"
	"time"
)

func TestStart(t *testing.T) {
	logger, _ := zap.NewProduction()
	// 如果有，刷新缓存
	defer logger.Sync()
	url := "www.baidu.com"
	logger.Info("failed to fetch URL",
		zap.String("url", url),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)
}

func TestNew(t *testing.T) {

	// 输出日志到文件
	file, _ := os.Create("./test.log")
	ws := zapcore.AddSync(file)

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = timeEncoder

	core := zapcore.NewCore(zapcore.NewConsoleEncoder(config),
		ws, zapcore.InfoLevel)

	logger := zap.New(core)
	logger.Debug("我不爱你")
	logger.Info("我爱您")
	logger.Warn("我爱您")
	logger.Error("我爱您")
	logger.DPanic("我爱您")
	logger.Panic("我爱您")
	//logger.Fatal("我爱您")

}

func TestLevel(t *testing.T) {
	lowPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level < zap.ErrorLevel && level >= zap.DebugLevel
	})
	highPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zap.ErrorLevel
	})

	prodEncoder := zap.NewProductionEncoderConfig()
	prodEncoder.EncodeTime = timeEncoder

	lowSyncer, lowClose, err := zap.Open("./level.txt")
	if err != nil {
		lowClose()
		return
	}

	highWriteSyncer, highClose, err := zap.Open("./high.txt")
	if err != nil {
		highClose()
		return
	}

	lowCore := zapcore.NewCore(zapcore.NewJSONEncoder(prodEncoder), lowSyncer, lowPriority)
	highCore := zapcore.NewCore(zapcore.NewJSONEncoder(prodEncoder), highWriteSyncer, highPriority)

	logger := zap.New(zapcore.NewTee(lowCore, highCore), zap.AddCaller())
	logger.Debug("i am debug", zap.String("key", "debug"))
	logger.Info("i am info", zap.String("key", "info"))
	logger.Error("i am error", zap.String("key", "error"))
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}
