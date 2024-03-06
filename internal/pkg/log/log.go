package log

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/go-logr/logr"
	"k8s.io/klog/v2"
	"k8s.io/klog/v2/klogr"
)

const (
	logFilePerm   = 0666
	sigLogRotate  = syscall.Signal(0xa)
	sigDebugLevel = syscall.Signal(0xc)
)

func Ctx(ctx context.Context) logr.Logger {
	lazyInit()
	return defaultLogger
}

func Error(err error, msg string, keysAndValues ...interface{}) {
	lazyInit()
	if err == nil {
		return
	}

	defaultLogger.WithCallDepth(1).
		Error(fmt.Errorf("%+v", err), msg, keysAndValues...)
}

func Info(msg string, keysAndValues ...interface{}) {
	lazyInit()
	defaultLogger.WithCallDepth(1).Info(msg, keysAndValues...)
}

func Debug(msg string, keysAndValues ...interface{}) {
	lazyInit()
	defaultLogger.WithCallDepth(1).V(1).Info(msg, keysAndValues...)
}

func Fatal(msg string, keysAndValues ...interface{}) {
	lazyInit()
	defaultLogger.WithCallDepth(1).Info(msg, keysAndValues...)
	Flush()
	os.Exit(1)
}

var (
	logFile string
	logOnce sync.Once

	defaultLogger logr.Logger
)

func openLogFile(filename string) *os.File {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, logFilePerm)
	if err != nil {
		panic(err)
	}
	return f
}

// InitFlags 初始化 logger 参数.
func InitFlags() {
	klog.InitFlags(nil)
	flag.StringVar(&logFile, "log", "", "specify the path of log file")
}

func lazyInit() {
	logOnce.Do(func() {
		c := make(chan os.Signal, 2)
		signal.Notify(c, sigLogRotate, sigDebugLevel)

		defaultLogger = klogr.NewWithOptions(klogr.WithFormat(klogr.FormatKlog))
		var rotateHandler func()
		if logFile != "" {
			_ = flag.Set("one_output", "true")
			_ = flag.Set("logtostderr", "false")
			_ = flag.Set("stderrthreshold", "4")
			f := openLogFile(logFile)
			klog.SetOutput(f)

			rotateHandler = func() {
				nf := openLogFile(logFile)
				klog.SetOutput(nf)

				_ = f.Close()
				f = nf
			}
		} else {
			rotateHandler = func() {
				defaultLogger.Info("log to file not enable, rotate no effect")
			}
		}

		go func() {
			var debug = false
			for v := range c {
				switch v {
				case sigLogRotate:
					rotateHandler()
				case sigDebugLevel:
					level := "10"
					if debug {
						level = "0"
					}

					_ = flag.Set("v", level)
					debug = !debug
				}
			}
		}()
	})
}

func Flush() {
	klog.Flush()
}

func GetLogger() *logr.Logger {
	lazyInit()
	return &defaultLogger
}

func DebugEnabled() bool {
	lazyInit()
	return defaultLogger.V(1).Enabled()
}
