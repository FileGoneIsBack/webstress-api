package models
import (
	"encoding/json"
	"api/core/models/log"
	"path/filepath"
	"api/core/models/apis"
	"api/core/models/plans"
	"api/core/models/servers"
	"os"
)
var (
	Config *Conf = new(Conf)
)

type Conf struct {
	Name    		string `json:"name"`
	Secure			bool   `json:"secure"`
	Cert    		string `json:"cert"`
	Vers    		string `json:"version"`
	Key     		string `json:"key"`
	AutobuyKey     	string `json:"autobuy_key"`
	Database struct {
		Host 			string `json:"host"`
		Database 		string `json:"database"`
		Username 	    string `json:"username"`
		Password 		string `json:"password"`
	} `json:"database"`
	//Methods map[string]string `json:"methods"`
}

func ReloadConfigs() error {
	configFiles := map[string]interface{}{
		filepath.Join("assets/config", "config.json"):     &Config,
		filepath.Join("assets/config", "servers.json"):    &servers.Config,
		filepath.Join("assets/config", "plans.json"):      &plans.GeneralConfig,
		filepath.Join("assets/config", "apis.json"):       &apis.Apis,
	}
	for filePath, configPointer := range configFiles {
		if err := reloadConfigFile(filePath, configPointer); err != nil {
			return err
		}
	}
	log.Println("Configs reloaded successfully.")
	return nil
}

func reloadConfigFile(filePath string, target interface{}) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(content, target); err != nil {
		return err
	}
	log.Printf("Successfully reloaded %s", filePath)
	return nil 
}