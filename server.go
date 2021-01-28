package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func webserver(cfg conf) {
	tmplH := gin.H{
		"placementToken": cfg.PlacementToken,
		"projectID":      cfg.ProjectId,
		"currency":       cfg.Currency,
		"country":        cfg.Country,
	}

	// https server
	if !isMissingValue(cfg.HttpsKeyFilename) && !isMissingValue(cfg.HttpsCertFilename) && !isMissingValue(cfg.HttpsListen) {
		httpsRouter := gin.Default()
		applyRoutes(httpsRouter, tmplH)
		go func() {
			if httpsErr := httpsRouter.RunTLS(cfg.HttpsListen, cfg.HttpsCertFilename, cfg.HttpsKeyFilename); httpsErr != nil {
				log.Printf("https server failed: %s\n", httpsErr)
			}
		}()
	} else {
		log.Println("Skipping https server, missing config data")
	}

	// http server
	httpRouter := gin.Default()
	applyRoutes(httpRouter, tmplH)
	if httpErr := httpRouter.Run(cfg.HttpListen); httpErr != nil {
		log.Printf("http server failed: %s\n", httpErr)
	}
}

func applyRoutes(router *gin.Engine, h map[string]interface{}) {
	router.StaticFS("/public", http.Dir("./public"))
	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", h)
	})
}
