package server

import (
	"encoding/json"
	"github.com/chargehive/example/config"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

func GetWebhooks(c *gin.Context) {
	var result []webhookLog
	q, b := c.GetQuery("from")
	from, _ := strconv.ParseInt(q, 10, 64)
	if from == 0 || !b {
		result = webhookLogs
	} else {
		for _, v := range webhookLogs {
			if v.Received > from {
				result = append(result, v)
			}
		}
	}
	jsonResult, _ := json.Marshal(result)
	c.String(http.StatusOK, "%s", jsonResult)
}

func ClearWebhooks(c *gin.Context) {
	webhookLogs = []webhookLog{}
	c.String(http.StatusOK, "OK")
}

type webhookLog struct {
	Received int64
	Data     []byte
}

var webhookLogs []webhookLog

func startWebhookServer(cfg *config.Config) {
	// webhook server - always accepts all messages
	if !config.IsEmptyVal(cfg.WebhookListen) {
		webhookRouter := gin.Default()
		webhookRouter.NoRoute(func(c *gin.Context) {
			if rData, err := ioutil.ReadAll(c.Request.Body); err != nil {
				log.Printf("error reading webhook data: %s\n", err.Error())
			} else {
				webhookLogs = append(webhookLogs, webhookLog{Received: time.Now().UnixNano(), Data: rData})
				log.Printf("Received Webhook: %s\n", string(rData))
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
}
