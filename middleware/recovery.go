/**
 * Created by angelina on 2017/8/31.
 * Copyright © 2017年 yeeyuntech. All rights reserved.
 */

package middleware

import (
	"bytes"
	"fmt"
	"github.com/vannnnish/easyweb"
	"io"
	"io/ioutil"
	"log"
	"net/http/httputil"
	"runtime"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

func Recovery() easyweb.HandlerFunc {
	return RecoveryWithWriter(easyweb.DefaultErrorWriter)
}

func RecoveryWithWriter(out io.Writer) easyweb.HandlerFunc {
	var logger *log.Logger
	if out == nil {
		logger = log.New(out, "\n\n\x1b[31m", log.LstdFlags)
	}
	return func(context *easyweb.Context) {
		defer func() {
			if err := recover(); err != nil {
				if easyweb.RecoveryLogger != nil {
					stack := stack(3)
					httprequest, _ := httputil.DumpRequest(context.Request(), false)
					easyweb.RecoveryLogger.Error("\n\n\x1b[0;31;40m[Recovery] panic recovered:\n%s\n%s\n%s%s\n", string(httprequest), err, stack, reset)
				} else if logger != nil {
					stack := stack(3)
					httprequest, _ := httputil.DumpRequest(context.Request(), false)
					logger.Printf("[Recovery] panic recovered:\n%s\n%s\n%s%s\n", string(httprequest), err, stack, reset)
				}
				context.AbortWithStatus(500)
			}
		}()
		context.Next()
	}
}

func stack(skip int) []byte {
	buf := new(bytes.Buffer)
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
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

func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

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
