package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"

	"github.com/kingwel-xie/k2/common/config"
)

func TlsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     config.SslConfig.Domain,
		})
		err := secureMiddleware.Process(c.Writer, c.Request)
		if err != nil {
			return
		}
		c.Next()
	}
}
