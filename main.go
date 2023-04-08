package main

import (
	"fmt"
	"github.com/btdjangbah001/chat-app/model"
	"github.com/gin-gonic/gin"
)

func main() {
	user := model.User{}
	fmt.Println(user)
	router := gin.Default()

	router.Run()
}
