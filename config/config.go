package config
import (
	"io/ioutil"
	"fmt"
	"os"
	"encoding/json"
	"regexp"
	"strings"
)
type Config struct {
	Token string `json:"token"`
	Prefix string `json:"prefix"`
	ConnectionString string `json:"connection_string"`
	DataDir string `json:"data_dir"`
	ConfigDir string `json:"config_dir"`
	AllowHistoryImport bool `json:"import_history"`
	MarkovOrder int `json:"markov_order"`
}
var Cfg Config
func Initialize(configFile string) {
	jsonFile, err := os.Open(configFile)
	if err != nil {
		fmt.Println("error opening config file ", err)
		return
	}
	fmt.Println("Config read successfully")
	defer jsonFile.Close()
	byteValue, _ :=ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue,&Cfg)
	fmt.Println("Token " + Cfg.Token)
}
func GetCommand(content string) (string,string,bool){
	content = strings.Replace(content,"<@!","<@",-1)
	r:= regexp.MustCompile(`^`+Cfg.Prefix+`(\w+)`)
	res:=r.FindStringSubmatch(content)
	if res != nil {
		return res[1],strings.TrimSpace(r.ReplaceAllString(content,"")),true
	}
	return "","",false 
}

