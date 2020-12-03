package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/chentanyi/fileserver/server"
	"github.com/gin-gonic/gin"
)

func main() {
	address := flag.String("a", "", "address")
	port := flag.Int("p", 80, "port")
	directory := flag.String("d", ".", "directory")
	baseUri := flag.String("base", "", "base uri")
	flag.Parse()

	router := gin.Default()
	group := router.Group(*baseUri)

	server.NewFileServer(group, "/", *directory)

	if strings.ContainsAny(*address, ":") {
		router.Run(fmt.Sprintf("[%s]:%d", *address, *port))
	} else {
		router.Run(fmt.Sprintf("%s:%d", *address, *port))
	}
}
