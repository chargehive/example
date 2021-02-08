package server

import (
	"crypto/md5"
	"fmt"
	"github.com/chargehive/example/chargehive"
	"github.com/chargehive/example/client"
	"github.com/chargehive/example/config"
	"github.com/chargehive/proto/golang/chargehive/chtype"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func Start(cfg *config.Config) {
	gin.SetMode(gin.TestMode)
	startHttpsServer(cfg)
	startWebhookServer(cfg)
	// start http server
	httpRouter := gin.Default()
	applyRoutes(httpRouter, cfg)
	if httpErr := httpRouter.Run(cfg.HttpListen); httpErr != nil {
		log.Printf("http server failed: %s\n", httpErr)
	}
}

func startHttpsServer(cfg *config.Config) {
	if !config.IsEmptyVal(cfg.HttpsKeyFilename) && !config.IsEmptyVal(cfg.HttpsCertFilename) && !config.IsEmptyVal(cfg.HttpsListen) {
		httpsRouter := gin.Default()
		applyRoutes(httpsRouter, cfg)
		go func() {
			if httpsErr := httpsRouter.RunTLS(cfg.HttpsListen, cfg.HttpsCertFilename, cfg.HttpsKeyFilename); httpsErr != nil {
				log.Printf("https server failed: %s\n", httpsErr)
			}
		}()
	} else {
		log.Println("Skipping https server, missing config data")
	}
}

func applyRoutes(router *gin.Engine, cfg *config.Config) {
	router.StaticFS("/static", http.Dir("./static"))
	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(c *gin.Context) {
		host := strings.Split(c.Request.Host, ":")[0]
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"placementToken": cfg.PlacementToken,
			"projectID":      cfg.ProjectId,
			"currency":       cfg.Currency,
			"country":        cfg.Country,
			"cdn":            cfg.PaymentAuthCdn,
			"randString":     fmt.Sprintf("%x", md5.Sum([]byte(strconv.FormatInt(time.Now().UnixNano(), 10)))),
			"httpLink":       fmt.Sprintf("http://%s:%s", host, strings.SplitAfter(cfg.HttpListen, ":")[1]),
			"httpsLink":      fmt.Sprintf("https://%s:%s", host, strings.SplitAfter(cfg.HttpsListen, ":")[1]),
		})
	})

	controls := router.Group("/controls")
	controls.POST("/ping", ping)
	controls.POST("/create", chargeCreate)
	controls.POST("/capture", chargeCapture)
	controls.POST("/refund", chargeRefund)
	controls.POST("/cancel", chargeCancel)
	controls.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "controls.tmpl", gin.H{"result": ""})
	})

	router.GET("/webhooks", GetWebhooks)
	router.DELETE("/webhooks", ClearWebhooks)
}

func ping(c *gin.Context) {
	c.String(http.StatusOK, client.Ping(c, c.PostForm("value")))
}

func chargeCreate(c *gin.Context) {
	units, _ := strconv.ParseInt(c.PostForm("units"), 10, 64)
	c.String(http.StatusOK, client.ChargeCreate(c, c.PostForm("currency"), units, c.PostForm("paymentMethodId")))
}

func chargeCapture(c *gin.Context) {
	units, _ := strconv.ParseInt(c.PostForm("units"), 10, 64)
	c.String(http.StatusOK, client.ChargeCapture(c, c.PostForm("chargeId"), c.PostForm("currency"), units))
}

func chargeCancel(c *gin.Context) {
	c.String(http.StatusOK, client.ChargeCancel(c, c.PostForm("chargeId"), chtype.Reason{
		Description:      c.PostForm("description"),
		ReasonType:       chtype.REASON_GENERIC,
		RequestorComment: "",
		RequestedBy:      chtype.ACTOR_TYPE_CHARGEHIVE,
	}))
}

func chargeRefund(c *gin.Context) {
	units, _ := strconv.ParseInt(c.PostForm("units"), 10, 64)
	var txns []*chargehive.ChargeRefundTransaction
	c.String(http.StatusOK, client.ChargeRefund(c, c.PostForm("chargeId"), c.PostForm("currency"), units, chtype.Reason{
		Description:      "Refund reason",
		ReasonType:       chtype.REASON_GENERIC,
		RequestorComment: "",
		RequestedBy:      chtype.ACTOR_TYPE_CHARGEHIVE,
	}, txns))
}
