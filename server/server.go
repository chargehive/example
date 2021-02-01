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

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "%s", client.Get().Ping("message"))
	})
	router.GET("/chargeCancel", func(c *gin.Context) {
		c.String(http.StatusOK, "%s", client.Get().ChargeCancel("chargeid", chtype.Reason{
			Description:      "Cancel reason",
			ReasonType:       chtype.REASON_GENERIC,
			RequestorComment: "",
			RequestedBy:      chtype.ACTOR_TYPE_CHARGEHIVE,
		}))
	})

	router.GET("/chargeCapture", func(c *gin.Context) {
		c.String(http.StatusOK, "%s", client.Get().ChargeCapture("chargeid", "USD", 5))
	})

	router.GET("/chargeRefund", func(c *gin.Context) {
		reason := chtype.Reason{
			Description:      "Cancel reason",
			ReasonType:       chtype.REASON_GENERIC,
			RequestorComment: "",
			RequestedBy:      chtype.ACTOR_TYPE_CHARGEHIVE,
		}
		var txns []*chargehive.ChargeRefundTransaction
		c.String(http.StatusOK, "%s", client.Get().ChargeRefund("chargeid", "USD", 5, reason, txns))
	})

}
