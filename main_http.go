package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/willy182/boilerplate-go-cleanarch/config/postgres"
	"github.com/willy182/boilerplate-go-cleanarch/utils"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// HTTPDefaultPort , default port for HTTP Server
const HTTPDefaultPort = 8080

// Serve function for serving
func (hsi *HSIService) Serve() {
	defer func() {
		if r := recover(); r != nil {
			utils.Log(log.ErrorLevel, fmt.Sprint(r), "Serve()", "recover_server")
		}
	}()

	postgres.InitDB()

	g := gin.New()

	g.Use(gin.Recovery())

	article := g.Group("/v1")

	// version 4
	hsi.Article.Handler.V1.Mount(article)

	//start gin server
	var port uint16
	if portEnv, ok := os.LookupEnv("SITE_PORT"); ok {
		portInt, err := strconv.Atoi(portEnv)
		if err != nil {
			port = HTTPDefaultPort
		} else {
			port = uint16(portInt)
		}
	} else {
		port = HTTPDefaultPort
	}

	listenerPort := fmt.Sprintf(":%d", port)
	g.Run(listenerPort)
}
