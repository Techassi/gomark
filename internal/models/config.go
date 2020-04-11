package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	cnst "github.com/Techassi/gomark/internal/constants"
	"github.com/Techassi/gomark/internal/util"
)

// Config represents the structure of the config.json file.
type Config struct {
	Domain  string `json:"domain"`
	WebRoot string `json:"web_root"`
	BaseURL string `json:"base_url"`
	UseSSL  bool   `json:"use_ssl"`
	DB      struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Database string `json:"database"`
	} `json:"db"`
	Server struct {
		Port int `json:"port"`
	} `json:"server"`
	Security struct {
		Jwt struct {
			Secret string `json:"secret"`
		} `json:"jwt"`
		TwoFA struct {
			Secret string `json:"secret"`
		} `json:"2fa"`
	} `json:"security"`
}

// Init reads the config.json file by the provided path and creates an Config
// instance.
func (c *Config) Init(p string) {
	aP := util.AbsolutePath(p)

	file, err := os.Open(aP)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	// Read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(file)

	// We unmarshal our byteArray which contains our
	// jsonFile's content into 'configuration' which we defined above
	err = json.Unmarshal(byteValue, c)
	if err != nil {
		panic(err)
	}
}

// SetURL sets all URLs based on the config's domain
func (c *Config) SetURL() {
	if c.UseSSL {
		c.BaseURL = fmt.Sprintf("%s%s/", cnst.WebHTTPSScheme, c.Domain)
		return
	}

	c.BaseURL = fmt.Sprintf("%s%s/", cnst.WebHTTPScheme, c.Domain)
}
