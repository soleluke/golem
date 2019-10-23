package doots
import (
	"io/ioutil"
	"fmt"
	"os"
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"golem/config"
	"regexp"
)

var doots map[string]int
func importDoots() {
	jsonFile,err := os.Open(config.Cfg.DataDir+"/doots.json")
	if err != nil {
		fmt.Println("error opening doots file",err)
		return
	}
	fmt.Println("Loaded doots")
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue),&doots)
}
func init(){
	importDoots()
}

func MessageCreate(s *discordgo.Session,m *discordgo.MessageCreate){
	comm,param,isComm:= config.GetCommand(m.Content)
	if !isComm {return}
	ra := regexp.MustCompile(`^<@[0-9]{18}>$`)
	switch(comm){	
		case "updoot":
			if !ra.MatchString(param) {return}
			if checkAuthor(param,m.Author.ID) {
				s.ChannelMessageSend(m.ChannelID,"No")
				return
			}
			modifyDoots(param,1)
			s.ChannelMessageSend(m.ChannelID,"doot doot")
			return
		case "downdoot":
			if !ra.MatchString(param) {return}
			if checkAuthor(param,m.Author.ID) {
				s.ChannelMessageSend(m.ChannelID,"No")
				return
			}
			modifyDoots(param,-1)
			s.ChannelMessageSend(m.ChannelID,"doot doot")
			return
		case "sidedoot":
			s.ChannelMessageSend(m.ChannelID,"~~doot doot~~")
			return
	}
}
func modifyDoots(p string,d int){
	if val,ok := doots[p]; ok {
		doots[p] = val + d
	} else {
		doots[p] = d
	}
	jsonString,err := json.Marshal(doots)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ioutil.WriteFile(config.Cfg.DataDir+"/doots.json",jsonString,0644)
	if err != nil {
		fmt.Println(err)
	}

}
func checkAuthor(p string,d string) bool {
	return p == "<@"+d+">"
}
