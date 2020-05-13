package triggers
import (
	"io/ioutil"
	"fmt"
	"os"
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"golem/config"
	"regexp"
	"strings"
)

var triggers map[string]string
var triggerReg []*regexp.Regexp
var triggerRes []string
func importTriggers() {
	jsonFile,err := os.Open(config.Cfg.ConfigDir+"/triggers.json")
	if err != nil {
		fmt.Println("error opening triggers file",err)
		return
	}
	fmt.Println("Loaded triggers")
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue),&triggers)
	for r,val := range triggers {
		re,err:=regexp.Compile(r)
		if err != nil {
			fmt.Println("error compiling trigger " + r, err)
		} else {
			triggerReg = append(triggerReg,re)
			triggerRes = append(triggerRes,val)
		}
	}
}
func Initialize(){
	importTriggers()
}

func MessageCreate(s *discordgo.Session,m *discordgo.MessageCreate){
	if m.Author.ID == s.State.User.ID { return }
	for i,val:=range triggerReg {
		if val.MatchString(strings.ToLower(m.Content)) {
			s.ChannelMessageSend(m.ChannelID,triggerRes[i])
		}
	}	
}
