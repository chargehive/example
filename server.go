package main

import (
	"crypto/md5"
	"fmt"
	"github.com/chargehive/example/chargehive"
	"github.com/chargehive/proto/golang/chargehive/chtype"
	"github.com/gin-gonic/gin"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

func webserver(cfg *conf) {
	gin.SetMode(gin.TestMode)

	tmplH := gin.H{
		"placementToken": cfg.PlacementToken,
		"projectID":      cfg.ProjectId,
		"currency":       cfg.Currency,
		"country":        cfg.Country,
		"cdn":            cfg.PaymentAuthCdn,
	}

	// https server
	if !cfg.isEmptyVal(cfg.HttpsKeyFilename) && !cfg.isEmptyVal(cfg.HttpsCertFilename) && !cfg.isEmptyVal(cfg.HttpsListen) {
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

	// webhook server - always accepts all messages
	if !cfg.isEmptyVal(cfg.WebhookListen) {
		webhookRouter := gin.Default()
		webhookRouter.NoRoute(func(c *gin.Context) {
			if rData, err := ioutil.ReadAll(c.Request.Body); err != nil {
				log.Printf("error reading webhook data: %s\n", err.Error())
			} else {
				log.Printf("received webhook data: %v", string(rData))
			}
			c.String(http.StatusOK, "{\"message\":\"OK\"}")
		})
		go func() {
			if webhookErr := webhookRouter.Run(cfg.WebhookListen); webhookErr != nil {
				log.Printf("webhook server failed: %s\n", webhookErr)
			}
		}()
	} else {
		log.Println("Skipping webhook server, missing config data")
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
	router.SetFuncMap(template.FuncMap{"getRandString": getRandString})
	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(c *gin.Context) { c.HTML(http.StatusOK, "index.tmpl", h) })

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "%s", api.Ping("message"))
	})
	router.GET("/chargeCancel", func(c *gin.Context) {
		c.String(http.StatusOK, "%s", api.ChargeCancel("chargeid", chtype.Reason{
			Description:      "Cancel reason",
			ReasonType:       chtype.REASON_GENERIC,
			RequestorComment: "",
			RequestedBy:      chtype.ACTOR_TYPE_CHARGEHIVE,
		}))
	})

	router.GET("/chargeCapture", func(c *gin.Context) {
		c.String(http.StatusOK, "%s", api.ChargeCapture("chargeid", "USD", 5))
	})

	router.GET("/chargeRefund", func(c *gin.Context) {
		reason := chtype.Reason{
			Description:      "Cancel reason",
			ReasonType:       chtype.REASON_GENERIC,
			RequestorComment: "",
			RequestedBy:      chtype.ACTOR_TYPE_CHARGEHIVE,
		}
		var txns []*chargehive.ChargeRefundTransaction
		c.String(http.StatusOK, "%s", api.ChargeRefund("chargeid", "USD", 5, reason, txns))
	})

}

func getRandString() string {
	return fmt.Sprintf("%x", md5.Sum([]byte(strconv.FormatInt(time.Now().UnixNano(), 10))))
}
