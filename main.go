package main

import (
	"fmt"

	"github.com/btdjangbah001/chat-app/models"
	"github.com/gin-gonic/gin"
)

func main() {
	user := models.User{}
	fmt.Println(user)
	router := gin.Default()

	router.Run()
}
