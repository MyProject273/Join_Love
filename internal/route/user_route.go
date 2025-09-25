package route

import (
	v1 "github.com/MyProject273/Join_Love/api/v1"
	db "github.com/MyProject273/Join_Love/internal/db/sqlc"
	"github.com/MyProject273/Join_Love/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RegisterUserRoutes(r *gin.RouterGroup, h v1.UserHandler, redis *redis.Client, store db.Store) {
	r.GET("/:id", middleware.RequirePermission("read_users", redis, store), h.GetUser)
	r.POST("/create", middleware.RequirePermission("create_users", redis, store), h.CreateUserByAdmin)
	r.PATCH("/update/:id", middleware.RequirePermission("update_users", redis, store), h.UpdateUser)
	r.PATCH("/active/:id", middleware.RequirePermission("update_active_users", redis, store), h.UpdateUserActiveStatus)
	r.DELETE("/delete/:id", middleware.RequirePermission("delete_users", redis, store), h.DeleteUser)
}
