package reminds
import (
	"golem/config"
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"fmt"
	"os"
	"io/ioutil"
)
var reminds map[string]string
func importReminds(){
	jsonFile,err:=os.Open(config.Cfg.DataDir+"/reminds.json")
	if err !=nil {
		fmt.Println("error opening reminds file")
		return
	}
	fmt.Println("loaded reminds")
	defer jsonFile.Close()
	byteValue,_:=ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue),&reminds)
}
func init(){
	importReminds()
}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate){
	if m.Author.ID == s.State.User.ID {return}
	comm,param,isComm := config.GetCommand(m.Content)
	if !isComm {return}
	if comm=="remind" {
		if param != ""{
			if val,ok:=reminds["<@"+m.Author.ID+">"];ok {
				reminds["<@"+m.Author.ID+">"]=val+" "+param
			} else {
				reminds["<@"+m.Author.ID+">"]=param
			}
			s.ChannelMessageSend(m.ChannelID,"added remind")
		} else {
			if val,ok:=reminds["<@"+m.Author.ID+">"]; ok {
				if val != "" {
					s.ChannelMessageSend(m.ChannelID,val)
					reminds["<@"+m.Author.ID+">"] = ""
				}
			}
		}
	}
}
func exportReminds(){
	jsonString,err:=json.Marshal(reminds)
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile(config.Cfg.DataDir+"/reminds.json",jsonString,0644)
	if err!=nil {
		fmt.Println(err)
	}
}
