package configure

import (
	"encoding/json"
	"os"

	log4plus "github.com/nGPU/include/log4go"
)

type ApplicationConfig struct {
	Name    string `json:"name"`
	Comment string `json:"comment"`
}

type WebConfig struct {
	Listen     string `json:"listen"`
	Domain     string `json:"domain"`
	ICPAddress string `json:"icpAddress"`
}

type Config struct {
	Application ApplicationConfig `json:"application"`
	Web         WebConfig         `json:"web"`
}

type Configure struct {
	config Config
}

var gConfigure *Configure

func (u *Configure) getConfig() error {
	funName := "getConfig"
	log4plus.Info("%s ---->>>>", funName)
	data, err := os.ReadFile("./config.json")
	if err != nil {
		log4plus.Error("%s ReadFile error=[%s]", funName, err.Error())
		return err
	}
	log4plus.Info("%s data=[%s]", funName, string(data))
	err = json.Unmarshal(data, &u.config)
	if err != nil {
		log4plus.Error("%s json.Unmarshal error=[%s]", funName, err.Error())
		return err
	}
	return nil
}

func SingletonConfigure() Config {
	if gConfigure == nil {
		gConfigure = &Configure{}
		_ = gConfigure.getConfig()
	}
	return gConfigure.config
}
