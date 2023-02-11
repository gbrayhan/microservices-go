package middlewares

import "github.com/gin-gonic/gin"

func CommonHeaders(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, DELETE, GET, PUT")
	c.Header("Access-Control-Allow-Headers",
		"Content-Type, Depth, UserName-Agent, X-File-Size, X-Requested-With, If-Modified-Since, X-File-CompanyName, Cache-Control")
	c.Header("X-Frame-Options", "SAMEORIGIN")
	c.Header("Cache-Control", "no-cache, no-store")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")

}
