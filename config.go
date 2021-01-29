package main

import (
	"errors"
	"fmt"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"os"
	"path/filepath"
)

const chiveFileName = ".chive.yaml"
const missingFieldValue = "!!CHANGE-ME!!"

var configFileName = "config.yaml"

// configuration options
type conf struct {
	ProjectId         string `yaml:"ProjectId"`         // Name of the chargehive project
	PaymentAuthCdn    string `yaml:"PaymentAuthCdn"`    // FQDN and/or port for the PaymentAuth service (usually "https://cdn.paymentauth.me:8823")
	PlacementToken    string `yaml:"PlacementToken"`    // Placement token to authorize the use of the chargehiveJS on a specific domain
	Currency          string `yaml:"Currency"`          // Initial currency
	Country           string `yaml:"Country"`           // Initial country
	HttpListen        string `yaml:"HttpListen"`        // Port and/or address for the http webservice  (e.g. ":8080"/"localhost:8080"/"0.0.0.0:8080")
	HttpsListen       string `yaml:"HttpsListen"`       // Port and/or address for the https webservice (e.g. ":8080"/"localhost:8080"/"0.0.0.0:8080")
	HttpsCertFilename string `yaml:"HttpsCertFilename"` // File path to the pem certificate file
	HttpsKeyFilename  string `yaml:"HttpsKeyFilename"`  // File path to the pem key file
	ApiHost           string `yaml:"ApiHost"`           // Port and/or address for the chargehive api server for requests (e.g. ":8080"/"localhost:8080"/"0.0.0.0:8080")
	ApiAccessToken    string `yaml:"ApiAccessToken"`    // Token used to authenticate requests to the api server
	WebhookListen     string `yaml:"WebhookReceiver"`   // Port and/or address to receive chargehive webhooks (e.g. ":8080"/"localhost:8080"/"0.0.0.0:8080")
}

func (c *conf) populateConfigFromChive(chc *chiveConf) error {
	if c.isEmptyVal(c.ProjectId) {
		if c.isEmptyVal(chc.ProjectId) {
			return errors.New("projectId")
		}
		c.ProjectId = chc.ProjectId
	}
	if c.isEmptyVal(c.ApiAccessToken) {
		if c.isEmptyVal(chc.AccessToken) {
			return errors.New("accessToken")
		}
		c.ApiAccessToken = chc.AccessToken
	}
	if c.isEmptyVal(c.ApiHost) {
		if c.isEmptyVal(chc.ApiHost) {
			return errors.New("apiHost")
		}
		c.ApiHost = chc.ApiHost
	}
	return nil
}

// config used if there is already chive config on the system
type chiveConf struct {
	ProjectId   string `yaml:"projectID"`
	AccessToken string `yaml:"accessToken"`
	ApiHost     string `yaml:"apiHost"`
}

func (chc *chiveConf) loadChiveConfig() error {
	var (
		dat           []byte
		chiveFilePath string
	)
	if dir, err := os.UserHomeDir(); err != nil {
		return err
	} else if chiveFilePath, err = filepath.Abs(filepath.Join(dir, chiveFileName)); err != nil {
		return err
	} else if dat, err = ioutil.ReadFile(chiveFilePath); err != nil {
		return err
	} else if err = yaml.Unmarshal(dat, chc); err != nil {
		return err
	}
	return nil
}

func (c *conf) saveConfig(confFileName string) error {
	j, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(confFileName, j, 0644)
}

func (c *conf) loadConfig(confFileName string) error {
	if _, err := os.Stat(confFileName); err == nil {
		// config file exists
		var dat []byte
		if dat, err = ioutil.ReadFile(confFileName); err != nil {
			return err
		} else if err = yaml.Unmarshal(dat, c); err != nil {
			return err
		} else if err = c.saveConfig(confFileName); err != nil {
			return fmt.Errorf("failed to save config `%s` back: %s", configFileName, err.Error())
		}
	} else {
		// missing config file, create new with defaults
		newConfig := conf{}.getDefault()
		if err = newConfig.saveConfig(confFileName); err != nil {
			return fmt.Errorf("missing `%s` config, and failed to write new config: %s", configFileName, err.Error())
		}
		return fmt.Errorf("missing config, created template config, edit and try again: %s", configFileName)
	}
	return nil
}

func (c conf) isEmptyVal(val string) bool {
	return val == "" || val == missingFieldValue || val == "cert/example.com.pem" || val == "cert/example.com-key.pem"
}

func (c conf) getDefault() conf {
	return conf{
		ProjectId:         "test-project",
		PlacementToken:    missingFieldValue,
		PaymentAuthCdn:    "https://cdn.paymentauth.me:8823",
		Currency:          "USD",
		Country:           "GB",
		HttpListen:        "0.0.0.0:9180",
		HttpsListen:       "0.0.0.0:4443",
		HttpsCertFilename: "cert/example.com.pem",
		HttpsKeyFilename:  "cert/example.com-key.pem",
		ApiHost:           "localhost:9050",
		ApiAccessToken:    missingFieldValue,
		WebhookListen:     "0.0.0.0:8092",
	}
}
