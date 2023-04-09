package main

import (
	"github.com/btdjangbah001/chat-app/models"
	"github.com/btdjangbah001/chat-app/routers"
)

func main() {
	models.ConnectDatabase()
	routers.SetupRouter().Run()
}
