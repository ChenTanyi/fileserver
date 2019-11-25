package main

import (
	"flag"
	"fmt"

	"github.com/chentanyi/fileserver/server"
	"github.com/gin-gonic/gin"
)

func main() {
	port := flag.Int("p", 80, "port")
	directory := flag.String("d", ".", "directory")
	baseUri := flag.String("base", "", "base uri")
	flag.Parse()

	router := gin.Default()
	group := router.Group(*baseUri)

	server.NewFileServer(group, "/", *directory)

	router.Run(fmt.Sprintf(":%d", *port))
}
