package route

import (
	v1 "github.com/MyProject273/Join_Love/api/v1"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.RouterGroup, h v1.AuthHandler) {
	r.POST("/login", h.Login)
	r.POST("/signup", h.Signup)
	r.GET("/verify_email", h.VerifyEmail)
}
