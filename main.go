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

	var listenAddr string
	if strings.ContainsAny(*address, ":") {
		listenAddr = fmt.Sprintf("[%s]:%d", *address, *port)
	} else {
		listenAddr = fmt.Sprintf("%s:%d", *address, *port)
	}
	if err := router.Run(listenAddr); err != nil {
		panic(err)
	}
}
