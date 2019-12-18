package logger

import (
	"bytes"
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"runtime"
	"vcard/pkg/setting"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

var logger *logrus.Logger

func GetLogFilePath() string {
	return fmt.Sprintf("%s%s", setting.AppSetting.RuntimeRootPath, setting.AppSetting.LogSavePath)
}

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// stack returns a nicely formatted stack frame, skipping skip frames.
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}

func timeFormat(t time.Time) string {
	var timeString = t.Format("2006/01/02 - 15:04:05")
	return timeString
}

func Setup() {
	if logger != nil {
		return
	}

	logFilePath := GetLogFilePath()
	logFileName := setting.AppSetting.LogSaveName

	// 日志文件
	fileName := path.Join(logFilePath, logFileName)

	// 写入文件
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}

	// 实例化
	logger = logrus.New()

	// 设置输出
	logger.Out = src

	// 设置日志级别
	logger.SetLevel(logrus.DebugLevel)

	// 设置 rotatelogs
	logWriter, err := rotatelogs.New(
		// 分割后的文件名称
		fileName+".%Y%m%d.log",

		// 生成软链，指向最新日志文件
		rotatelogs.WithLinkName(fileName),

		// 设置最大保存时间(7天)
		rotatelogs.WithMaxAge(7*24*time.Hour),

		// 设置日志切割时间间隔(1天)
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}

	lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 新增 Hook
	logger.AddHook(lfHook)
}

func LoggerToFile() gin.HandlerFunc {

	return func(c *gin.Context) {

		defer func() {
			if err := recover(); err != nil {
				if logger != nil {
					stack := stack(3)
					httprequest, _ := httputil.DumpRequest(c.Request, false)
					logger.Errorf("[Recovery] %s panic recovered:\n%s\n%s\n%s%s", timeFormat(time.Now()), string(httprequest), err, stack, string([]byte{27, 91, 48, 109}))
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()

		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		// 日志格式
		logger.WithFields(logrus.Fields{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"client_ip":    clientIP,
			"req_method":   reqMethod,
			"req_uri":      reqUri,
		}).Info()
	}
}


func Recovery() gin.HandlerFunc {

	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				if logger != nil {
					stack := stack(3)
					httprequest, _ := httputil.DumpRequest(c.Request, false)
					logger.Errorf("[Recovery] %s panic recovered:\n%s\n%s\n%s%s", timeFormat(time.Now()), string(httprequest), err, stack, string([]byte{27, 91, 48, 109}))
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
	}
}
