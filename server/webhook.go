package server

import (
	"encoding/json"
	"github.com/chargehive/example/config"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
)

func startWebhookServer(cfg *config.Config) {
	// webhook server - always accepts all messages
	if !config.IsEmptyVal(cfg.WebhookListen) {
		webhookRouter := gin.Default()
		webhookRouter.NoRoute(func(c *gin.Context) {
			if rData, err := ioutil.ReadAll(c.Request.Body); err != nil {
				log.Printf("error reading webhook data: %s\n", err.Error())
			} else {
				hook := parseWebhook(rData)
				Webhooks[hook.Created] = hook
				log.Printf("Type:%s  -   %s", hook.Type, hook.Data)
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

var Webhooks = map[int64]Webhook{}

type Webhook struct {
	Uuid           string
	Type           string
	Data           string
	JsonData       interface{}
	Checksum       string
	Verification   string
	Created        int64
	WebhookVersion string
	ProjectId      string
}

func parseWebhook(rData []byte) Webhook {
	var hook Webhook
	err := json.Unmarshal(rData, &hook)
	if err != nil {
		log.Printf("failed to unmarshal webhook: %s", err.Error())
	}

	err = json.Unmarshal([]byte(hook.Data), &hook.JsonData)
	if err != nil {
		log.Printf("failed to unmarshal webhook data: %s", err.Error())
	}
	return hook
}
