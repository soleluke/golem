package grabs
import (
	"io/ioutil"
	"fmt"
	"os"
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"golem/config"
	"regexp"
	"strings"
	"math/rand"
	"time"
)

var grabs map[string][]string
func importGrabs() {
	jsonFile,err := os.Open(config.Cfg.DataDir+"/grabs.json")
	if err != nil {
		fmt.Println("error opening grabs file",err)
		return
	}
	fmt.Println("Loaded grabs")
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue),&grabs)
}
func init(){
	importGrabs()
	rand.Seed(time.Now().Unix())
}


func MessageCreate(s *discordgo.Session,m *discordgo.MessageCreate){
	comm,param,isComm:=config.GetCommand(m.Content)
	if !isComm {
		return
	}
	if comm=="grabr" {
		ra := regexp.MustCompile(`^<@[0-9]{18}>$`)
		if ra.MatchString(param) {
			s.ChannelMessageSend(m.ChannelID,displayGrab(param,grabs[param][rand.Intn(len(grabs[param]))]))
		} else if len(param)==0 {
			arr := []string{}
			for k,v := range grabs {
				for _,g := range v {
					arr = append(arr,displayGrab(k,g))
				}
			}
			s.ChannelMessageSend(m.ChannelID,arr[rand.Intn(len(arr))])
		} else {
			arr := []string{}
			for k,v:=range grabs {
				for _,g := range v {
					if strings.Contains(g,param) {
						arr = append(arr,displayGrab(k,g))
					}
				}
			}
			s.ChannelMessageSend(m.ChannelID,arr[rand.Intn(len(arr))])
		}
		return
	}
	if comm=="igrab" {
		importGrabs()
		s.ChannelMessageSend(m.ChannelID,"Refreshed Grabs from file")
		return
	}
	if comm=="grab" {
		messages,err:=s.ChannelMessages(m.ChannelID,2,"","","")
		if err != nil {
			fmt.Println(err)
		}
		msg :=messages[1]
		if val,ok:= grabs["<@"+msg.Author.ID+">"]; ok {
			grabs["<@" + msg.Author.ID+">"] = append(val,msg.Content)
		} else {
			grabs["<@" + msg.Author.ID+">"] = []string{msg.Content}
		}
		s.ChannelMessageSend(m.ChannelID,"Grab Successful!")
		exportGrabs()
		return
	}
}
func displayGrab(a string, g string) string{
	return a+" : " + g
}
func exportGrabs(){
	jsonString,err:=json.Marshal(grabs)
	if err != nil {
		fmt.Println(err)
	}
	err=ioutil.WriteFile(config.Cfg.DataDir+"/grabs.json",jsonString,0644)
	if err != nil {
		fmt.Println(err)
	}
}
