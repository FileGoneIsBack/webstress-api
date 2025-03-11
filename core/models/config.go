package models
import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"api/core/models/apis"
	"api/core/models/plans"
	"api/core/models/servers"
)
var (
	Config *Conf = new(Conf)
	logger  = log.New(os.Stderr, "[init] ", log.Ltime|log.Lshortfile)
)

type Conf struct {
	Name    		string `json:"name"`
	Domain  		string `json:"domain"`
	Secure  		bool   `json:"secure"`
	Cert    		string `json:"cert"`
	Vers    		string `json:"version"`
	Key     		string `json:"key"`
	Database struct {
		Host 			string `json:"host"`
		Database 		string `json:"database"`
		Username 	    string `json:"username"`
		Password 		string `json:"password"`
	} `json:"database"`
	Autobuy struct {
		Key 		string `json:"key"`
		Email 		string `json:"email"`
	} `json:"autobuy"`
    FreeUser struct {
        Enabled1    bool   `json:"enabled"`
        Concurrents string `json:"concurrents"`
        Duration   string `json:"duration"`
    } `json:"freeuser"`
	Fake struct {
		Users 	   	int    `json:"users"`
		Attacks    	int    `json:"attacks"`
	} `json:"fake"`
	Server struct {
		Enabled 	bool   `json:"enabled"`
		SSH    		string `json:"ssh"`
		Telnet 		string `json:"telnet"`
	} `json:"cnc"`
	Bot struct {
		Auth 	string   `json:"key"`
		URL		string	 `json:"url"`
	} `json:"bot"`
	Methods map[string]string `json:"methods"`
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
	logger.Printf("Successfully reloaded %s", filePath)
	return nil
}