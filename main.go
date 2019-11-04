package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	port := flag.Int("p", 80, "port")
	directory := flag.String("d", ".", "directory")
	baseUri := flag.String("base", "", "base uri")
	flag.Parse()

	router := gin.Default()
	group := router.Group(*baseUri)

	group.StaticFS("/", http.Dir(*directory))

	router.Run(fmt.Sprintf(":%d", *port))
}
