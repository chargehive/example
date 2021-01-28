package main

import (
	"errors"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"os"
	"path/filepath"
)

const chiveFileName = ".chive.yaml"
const missingFieldValue = "!!CHANGE-ME!!"

var configFileName = "config.yaml"
var config conf

// configuration options
type conf struct {
	ProjectId         string `yaml:"projectId"`         // Name of the chargehive project
	PaymentAuthCdn    string `yaml:"paymentAuthCdn"`    // FQDN and/or port for the PaymentAuth service (usually "https://cdn.paymentauth.me:8823")
	PlacementToken    string `yaml:"placementToken"`    // Placement token to authorize the use of the chargehiveJS on a specific domain
	Currency          string `yaml:"currency"`          // Initial currency
	Country           string `yaml:"country"`           // Initial country
	HttpListen        string `yaml:"httpListen"`        // Port and/or address for the http webservice  (e.g. ":8080"/"localhost:8080"/"0.0.0.0:8080")
	HttpsListen       string `yaml:"httpsListen"`       // Port and/or address for the https webservice (e.g. ":8080"/"localhost:8080"/"0.0.0.0:8080")
	HttpsCertFilename string `yaml:"httpsCertFilename"` // File path to the pem certificate file
	HttpsKeyFilename  string `yaml:"httpsKeyFilename"`  // File path to the pem key file
	ApiHost           string `yaml:"apiHost"`           // Port and/or address for the chargehive api server for requests (e.g. ":8080"/"localhost:8080"/"0.0.0.0:8080")
	ApiAccessToken    string `yaml:"apiAccessToken"`    // Token used to authenticate requests to the api server
}

// config used if there is already chive config on the system
type chiveConf struct {
	ProjectId   string `yaml:"projectID"`
	AccessToken string `yaml:"accessToken"`
	ApiHost     string `yaml:"apiHost"`
}

func populateConfigFromChive(chc *chiveConf) error {
	if isMissingValue(config.ProjectId) {
		if isMissingValue(chc.ProjectId) {
			return errors.New("projectId")
		}
		config.ProjectId = chc.ProjectId
	}
	if isMissingValue(config.ApiAccessToken) {
		if isMissingValue(chc.AccessToken) {
			return errors.New("accessToken")
		}
		config.ApiAccessToken = chc.AccessToken
	}
	if isMissingValue(config.ApiHost) {
		if isMissingValue(chc.ApiHost) {
			return errors.New("apiHost")
		}
		config.ApiHost = chc.ApiHost
	}
	return nil
}

func loadChiveConfig() (*chiveConf, error) {
	var (
		dat           []byte
		chiveFilePath string
		chc           chiveConf
	)
	if dir, err := os.UserHomeDir(); err != nil {
		return nil, err
	} else if chiveFilePath, err = filepath.Abs(filepath.Join(dir, chiveFileName)); err != nil {
		return nil, err
	} else if dat, err = ioutil.ReadFile(chiveFilePath); err != nil {
		return nil, err
	} else if err = yaml.Unmarshal(dat, &chc); err != nil {
		return nil, err
	}
	return &chc, nil
}

func loadConfig(confFileName string) error {
	if _, err := os.Stat(confFileName); err == nil {
		// config file exists
		if dat, err := ioutil.ReadFile(confFileName); err != nil {
			return err
		} else if err = yaml.Unmarshal(dat, &config); err != nil {
			return err
		}
	} else {
		// missing config file, create new with defaults
		config = newConf()
		if j, err := yaml.Marshal(config); err != nil {
			return err
		} else if err = ioutil.WriteFile(confFileName, j, 0644); err != nil {
			return err
		}
		return errors.New("Created new config file to be edited then retry: " + confFileName)
	}
	return nil
}

func isMissingValue(val string) bool {
	return val == "" || val == missingFieldValue || val == "cert/example.com.pem" || val == "cert/example.com-key.pem"
}

func newConf() conf {
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
	}
}
