package config

import (
	"github.com/ghodss/yaml"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const chiveFileName = ".chive.yaml"
const missingFieldValue = "!!CHANGE-ME!!"
const configFileName = "config.yaml"

var configObj = Config{}
var configFilePath = ""

// configuration options
type Config struct {
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

func Init() {
	configObj.loadDefaults()

	// load config file if it exists
	if _, err := os.Stat(configObj.GetURI()); err == nil {
		if err = configObj.loadFromFile(configObj.GetURI()); err != nil {
			log.Printf("failed to load config file (%s): %s", configObj.GetURI(), err.Error())
		}
	}

	// populate missing fields with .chive.yaml if available
	if configObj.ProjectId == "" || configObj.ApiAccessToken == "" || configObj.ApiHost == "" {
		var chc chiveConfig
		if err := chc.loadChiveConfig(); err != nil {
			log.Printf("missing required fields in config and cannot load .chive.yaml: %s", err.Error())
		} else if err = configObj.loadFromChiveConfig(&chc); err != nil {
			log.Printf("missing required fields in config and cannot populate from .chive.yaml: %s", err.Error())
		}
	}

	// save config file
	if err := configObj.saveConfig(configObj.GetURI()); err != nil {
		log.Printf("failed to save config to: %s", err)
	}

}

func Get() *Config {
	return &configObj
}

func (c *Config) GetURI() string {
	if configFilePath == "" {
		if wd, err := os.Getwd(); err != nil {
			log.Fatal(err)
		} else if configFilePath, err = filepath.Abs(filepath.Join(wd, configFileName)); err != nil {
			log.Fatal(err)
		}
	}
	return configFilePath
}

func IsEmptyVal(val string) bool {
	return val == "" || val == missingFieldValue || val == "cert/example.com.pem" || val == "cert/example.com-key.pem"
}

func (c *Config) loadFromChiveConfig(chc *chiveConfig) error {
	if IsEmptyVal(c.ProjectId) {
		c.ProjectId = chc.ProjectId
	}
	if IsEmptyVal(c.ApiAccessToken) {
		c.ApiAccessToken = chc.AccessToken
	}
	if IsEmptyVal(c.ApiHost) {
		c.ApiHost = chc.ApiHost
	}
	return nil
}

func (c *Config) loadFromFile(confFileName string) error {
	if dat, err := ioutil.ReadFile(confFileName); err != nil {
		return err
	} else if err = yaml.Unmarshal(dat, c); err != nil {
		return err
	}
	return nil
}

func (c *Config) saveConfig(confFileName string) error {
	j, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(confFileName, j, 0644)
}

func (c *Config) loadDefaults() {
	c.ProjectId = "test-project"
	c.PlacementToken = missingFieldValue
	c.PaymentAuthCdn = "https://cdn.paymentauth.me:8823"
	c.Currency = "USD"
	c.Country = "GB"
	c.HttpListen = "0.0.0.0:9180"
	c.HttpsListen = "0.0.0.0:4443"
	c.HttpsCertFilename = "cert/example.com.pem"
	c.HttpsKeyFilename = "cert/example.com-key.pem"
	c.ApiHost = "localhost:9050"
	c.ApiAccessToken = missingFieldValue
	c.WebhookListen = "0.0.0.0:8092"
}
