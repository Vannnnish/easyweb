/**
 * Created by angelina on 2017/8/25.
 */

package middleware

import (
	"fmt"
	"github.com/mattn/go-isatty"
	"easyweb"
	"io"
	"os"
	"strings"
	"time"
)

// ascii码表
// 27 : ESC 换码(溢出) 0x1B
// 48 : 0
// 49 : 1
// 50 : 2
// 51 : 3
// 52 : 4
// 53 : 5
// 54 : 6
// 55 : 7
// 57 : 9
// 59 : ;
// 91 : [
// 109 : m
var (
	green        = string([]byte{27, 91, 57, 55, 59, 52, 50, 109}) // 0x1B[97;42m
	white        = string([]byte{27, 91, 57, 48, 59, 52, 55, 109}) // 0x1B[90;47m
	yellow       = string([]byte{27, 91, 57, 55, 59, 52, 51, 109}) // 0x1B[97;43m
	red          = string([]byte{27, 91, 57, 55, 59, 52, 49, 109}) // 0x1B[97;41m
	blue         = string([]byte{27, 91, 57, 55, 59, 52, 52, 109}) // 0x1B[97;44m
	magenta      = string([]byte{27, 91, 57, 55, 59, 52, 53, 109}) // 0x1B[97;45m 紫红色
	cyan         = string([]byte{27, 91, 57, 55, 59, 52, 54, 109}) // 0x1B[97;46m 青蓝色
	reset        = string([]byte{27, 91, 48, 109})                 // 0x1B[0m
	disableColor = false
)

// 禁止颜色输出
func DisableConsoleColor() {
	disableColor = true
}

func Logger() easyweb.HandlerFunc {
	return LoggerWithWriter(easyweb.DefaultWriter)
}

// logger中间件
func LoggerWithWriter(out io.Writer, notLogged ...string) easyweb.HandlerFunc {
	isTerm := true
	if w, ok := out.(*os.File); !ok ||
		(os.Getenv("TERM") == "dumb" || (!isatty.IsTerminal(w.Fd()) && !isatty.IsCygwinTerminal(w.Fd()))) ||
		disableColor {
		isTerm = false
	}
	// 过滤掉不需要记录的path
	var skip map[string]struct{}
	if length := len(notLogged); length > 0 {
		skip = make(map[string]struct{}, length)
		for _, path := range notLogged {
			skip[path] = struct{}{}
		}
	}
	return func(c *easyweb.Context) {
		start := time.Now()
		path := c.Request().URL.Path
		raw := c.Request().URL.RawQuery
		c.Next()
		if _, ok := skip[path]; !ok {
			end := time.Now()
			latency := end.Sub(start)
			// 只打印时间大于2s的请求
			if latency.Seconds() <= 2 {
				return
			}
			if strings.Contains(path, "/web/coin/ws") || strings.Contains(path, "/app/coin/ws") {
				return
			}
			method := c.Request().Method
			statusCode := c.Response().Status()
			clientIP := c.ClientIP()
			var statusColor, methodColor, resetColor string
			if isTerm {
				statusColor = colorForStatus(statusCode)
				methodColor = colorForMethod(method)
				resetColor = reset
			}
			if raw != "" {
				path = path + "?" + raw
			}
			responseSize := c.Response().Size()
			fmt.Fprintf(out, "[EasyWeb] %v |%s %3d %s| %13v | %10d | %10s |%s %-7s %s %s\n",
				end.Format("2006/01/02 - 15:04:05"),
				statusColor, statusCode, resetColor,
				latency,
				responseSize,
				clientIP,
				methodColor, method, resetColor,
				path,
			)
		}
	}
}

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return green
	case code >= 300 && code < 400:
		return white
	case code >= 400 && code < 500:
		return yellow
	default:
		return red
	}
}

func colorForMethod(method string) string {
	switch method {
	case "GET":
		return blue
	case "POST":
		return cyan
	case "PUT":
		return yellow
	case "DELETE":
		return red
	case "PATCH":
		return green
	case "HEAD":
		return magenta
	case "OPTIONS":
		return white
	default:
		return reset
	}
}
