package wtells
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

var wtells map[string][]string
func importwTells() {
	jsonFile,err := os.Open(config.Cfg.DataDir+"/wtells.json")
	if err != nil {
		fmt.Println("error opening wtells file",err)
		return
	}
	fmt.Println("Loaded wtells")
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue),&wtells)
}
func Initialize(){
	importwTells()
}

func MessageCreate(s *discordgo.Session,m *discordgo.MessageCreate){
	if m.Author.ID == s.State.User.ID {return}
	comm,param,isComm := config.GetCommand(m.Content)
	if !isComm {return}
	switch(comm) {
		case "wtell":
			params := strings.SplitN(param," ",2)
			target,tell:= params[0], params[1]
			t:=time.Now()
			tell = "<@" +m.Author.ID+"> said: '"+tell+"' on "+ t.Format("Mon Jan 2") +" at " + t.Format("15:04")
			if val,ok:=wtells[target]; ok {
				wtells[target] = append(val,tell)
			} else {
				wtells[target] = []string{tell}
			}
			exportwTells()
			s.ChannelMessageSend(m.ChannelID, "i'll tell " + target + " that next time i see them")
			return
		case "home":
			if val,ok:= wtells["<@"+m.Author.ID+">"]; ok {
				if len(val) > 0 {
				for _,t := range val {
					s.ChannelMessageSend(m.ChannelID, t)
				}
				wtells["<@"+m.Author.ID+">"] = make([]string,0,0)
				exportwTells()
				}
			}
			return
	}
}
func exportwTells(){
	jsonString,err:=json.Marshal(wtells)
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile(config.Cfg.DataDir+"/wtells.json",jsonString,0644)
	if err != nil {
		fmt.Println(err)
	}
}
