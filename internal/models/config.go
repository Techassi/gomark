package models

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/Techassi/gomark/internal/util"
)

// Config represents the structure of the config.json file.
type Config struct {
	DB struct {
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
	aP := util.GetAbsPath(p)

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
