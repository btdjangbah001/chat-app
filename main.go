package main

import (
	"github.com/btdjangbah001/chat-app/routers"
)

func main() {
	routers.SetupRouter().Run()
}
