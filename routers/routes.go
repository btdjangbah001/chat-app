package routers

import (
	"github.com/btdjangbah001/chat-app/auth"
	"github.com/btdjangbah001/chat-app/chat"
	"github.com/btdjangbah001/chat-app/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.Default())

	// Routes
	r.Any("/chat/:token", middlewares.AuthMiddleware(), chat.ChatHandler)
	r.POST("/signup", auth.RegisterUser)
	r.POST("/login", auth.LoginUser)

	return r
}
