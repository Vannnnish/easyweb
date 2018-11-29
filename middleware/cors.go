/**
 * Created by angelina on 2018/9/25.
 * Copyright © 2018年 yeeyuntech. All rights reserved.
 */

package middleware

import (
	"gitlab.yeeyuntech.com/yee/easyweb"
	"strings"
)

// 跨域设置
func CORSMiddleware() easyweb.HandlerFunc {
	f := func(c *easyweb.Context) {
		c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
		origin := c.Request().Header.Get("origin")
		if origin == "" {
			origin = c.Request().Header.Get("referer")
			if origin == "" {
				origin = "*"
			} else {
				if strings.HasSuffix(origin, "/") {
					origin = origin[:len(origin)-1]
				}
			}
		}
		c.Response().Header().Set("Access-Control-Allow-Origin", origin)
		c.Response().Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		c.Next()
	}
	return f
}
