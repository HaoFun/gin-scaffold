package middlewares

import "github.com/gin-gonic/gin"

func LanguageMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.GetHeader("Accept-Language")
		if lang == "" {
			lang = "zh_tw"
		}
		c.Set("lang", lang)
		c.Next()
	}
}
