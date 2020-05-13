package tells
import (
	"io/ioutil"
	"fmt"
	"os"
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"golem/config"
	"strings"
	"time"
)

var tells map[string][]string
func importTells() {
	jsonFile,err := os.Open(config.Cfg.DataDir+"/tells.json")
	if err != nil {
		fmt.Println("error opening tells file",err)
		return
	}
	fmt.Println("Loaded tells")
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue),&tells)
}
func Initialize(){
	importTells()
}

func MessageCreate(s *discordgo.Session,m *discordgo.MessageCreate){
	if m.Author.ID == s.State.User.ID {
		return
	}
	if val,ok:= tells["<@"+m.Author.ID+">"]; ok {
		if len(val) > 0 {
			for _,t := range val {
				s.ChannelMessageSend(m.ChannelID, t)
			}
			tells["<@"+m.Author.ID+">"] = make([]string,0,0)
			exportTells()
		}
	}
	comm,param,isComm := config.GetCommand(m.Content)
	if !isComm {
		return
	}
	if comm == "tell" {
		params := strings.SplitN(param," ",2)
		target,tell:= params[0], params[1]
		t:=time.Now()
		tell = "<@" +m.Author.ID+"> said: '"+tell+"' on "+ t.Format("Mon Jan 2") +" at " + t.Format("15:04")
		if val,ok:=tells[target]; ok {
			tells[target] = append(val,tell)
		} else {
			tells[target] = []string{tell}
		}
		exportTells()
		s.ChannelMessageSend(m.ChannelID, "i'll tell " + target + " that next time i see them")
	}
}
func exportTells(){
	jsonString,err:=json.Marshal(tells)
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile(config.Cfg.DataDir+"/tells.json",jsonString,0644)
	if err != nil {
		fmt.Println(err)
	}
}
