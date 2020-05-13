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
var recents []string
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
	rand.Seed(time.Now().Unix())
}
func Initialize(){
	importGrabs()
}


func MessageCreate(s *discordgo.Session,m *discordgo.MessageCreate){
	comm,param,isComm:=config.GetCommand(m.Content)
	if !isComm {
		return
	}
	if comm=="grabr" || comm == "morn"  {
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
			s.ChannelMessageSend(m.ChannelID,getNonRecentGrab(arr))
		} else {
			arr := []string{}
			for k,v:=range grabs {
				for _,g := range v {
					if strings.Contains(g,param) {
						arr = append(arr,displayGrab(k,g))
					}
				}
			}
			s.ChannelMessageSend(m.ChannelID,getNonRecentGrab(arr))
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
func getNonRecentGrab(arr []string) string{
	var grab string
	isrec:=true
	for grab="";isrec;isrec=false {
		grab = arr[rand.Intn(len(arr))]
		isrec=isRecent(grab)
	}
	recents = append(recents,grab)
	if len(recents) > 5 {
		recents = recents[1:]
	}
	return grab
}
func isRecent(s string) bool {
	for _,item :=range(recents){
		if item == s { return true }
	}
	return false
}
func isMorn(comm string,s *discordgo.Session,chid string) bool {
	var tim,_=time.Parse("23:00","12:00")
	if time.Now().Before(tim) {
		return comm == "morn"
	}
	if comm == "morn" {
		s.ChannelMessageSend(chid,"too late, try .aft")
	}
	return false
}

func isAft(comm string,s *discordgo.Session,chid string) bool {
	var tim,_=time.Parse("23:00","12:00")
	if time.Now().After(tim) {
		return comm == "aft"
	}
	if comm == "aft" {
		s.ChannelMessageSend(chid,"too early, try .morn")
	}
	return false
}

func exportGrabs(){
	jsonString,err:=json.MarshalIndent(grabs,"","    ")
	if err != nil {
		fmt.Println(err)
	}
	err=ioutil.WriteFile(config.Cfg.DataDir+"/grabs.json",jsonString,0644)
	if err != nil {
		fmt.Println(err)
	}
}
