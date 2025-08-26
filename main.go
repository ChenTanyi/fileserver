package main

import (
	"flag"
	"fmt"
	"net"
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

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	group := router.Group(*baseUri)

	server.NewFileServer(group, "/", *directory)

	// Print local address
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Printf("[WARN] Unable to get local address. err: %s", err.Error())
	} else {
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if v4 := ip.To4(); v4 != nil {
				fmt.Printf("Listening at %s\n", v4.String())
			}
		}
	}

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
