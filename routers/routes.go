package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/btdjangbah001/chat-app/auth"
	"github.com/btdjangbah001/chat-app/middlewares"
	"github.com/btdjangbah001/chat-app/chat"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Routes
	r.GET("/chat", middlewares.AuthMiddleware(), chat.ChatHandler)
	r.POST("/signup", auth.RegisterUser)
	r.POST("/login", auth.LoginUser)

	return r
}