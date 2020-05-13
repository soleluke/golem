package help
import (
	"io/ioutil"
	"fmt"
	"os"
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"golem/config"
)

var commands map[string]string
func importCommands(){
	jsonFile,err:=os.Open(config.Cfg.ConfigDir+"/commands.json")
	if err!=nil {
		fmt.Println("error opening commands file",err)
		return
	}
	fmt.Println("Loaded command help")
	defer jsonFile.Close()
	byteValue,_:=ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue),&commands)
}
func Initialize() {
	importCommands()
}

func MessageCreate(s *discordgo.Session,m *discordgo.MessageCreate){
	comm,param,isComm:=config.GetCommand(m.Content)
	if !isComm {
		return
	}
	if comm == "help" {
		if v,found:=commands[param]; found {
			s.ChannelMessageSend(m.ChannelID,v)
		} else {
			s.ChannelMessageSend(m.ChannelID,"No command found")
		}
	}
}
