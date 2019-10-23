package places
import (
	"io/ioutil"
	"fmt"
	"os"
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"golem/config"
	"math/rand"
	"time"
	"strings"
)

var places []string
func importPlaces() {
	jsonFile,err := os.Open(config.Cfg.DataDir+"/places.json")
	if err != nil {
		fmt.Println("error opening doots file",err)
		return
	}
	fmt.Println("Loaded places")
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue),&places)
}
func init(){
	importPlaces()
	rand.Seed(time.Now().Unix())
}

func MessageCreate(s *discordgo.Session,m *discordgo.MessageCreate){
	comm,param,isComm:= config.GetCommand(m.Content)
	if !isComm {return}
	switch(comm){
		case "suggest":
			place := places[rand.Intn(len(places))]
			s.ChannelMessageSend(m.ChannelID,place)
			return
		case "addplace":
			places = append(places,param)
			exportPlaces()
			s.ChannelMessageSend(m.ChannelID,"Added "+ param+" to list")
			return
		case "list":
			msg := strings.Join(places,"\n")
			user,err:=s.UserChannelCreate(m.Author.ID)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID,"I can't talk to you: " + err.Error())
			}
			s.ChannelMessageSend(m.ChannelID,"I have PMed you the list of places")
			s.ChannelMessageSend(user.ID,msg)
			return
	}
}
func exportPlaces(){
	jsonString,err:=json.Marshal(places)
	if err != nil {
		fmt.Println(err)
	}
	err=ioutil.WriteFile(config.Cfg.DataDir+"/places.json",jsonString,0644)
	if err != nil {
		fmt.Println(err)
	}
}
