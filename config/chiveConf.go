package config

import (
	"github.com/ghodss/yaml"
	"io/ioutil"
	"os"
	"path/filepath"
)

// config used if there is already chive config on the system
type chiveConfig struct {
	ProjectId   string `yaml:"projectID"`
	AccessToken string `yaml:"accessToken"`
	ApiHost     string `yaml:"apiHost"`
}

func (chc *chiveConfig) loadChiveConfig() error {
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
