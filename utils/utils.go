package utils
import (
	"github.com/bwmarrin/discordgo"
	"golem/config"
	"strings"
	"math/rand"
	"time"
	"regexp"
	"strconv"
	"fmt"
	"sort"
)

func Initialize(){
}
func init(){
	rand.Seed(time.Now().Unix())
}

func MessageCreate(s *discordgo.Session,m *discordgo.MessageCreate){
	if m.Author.ID == s.State.User.ID {	return}
	comm,param,isComm:=config.GetCommand(m.Content)
	if !isComm {return}
	switch(comm){
		case "ask":
			if strings.Contains(param," or ") {
				s.ChannelMessageSend(m.ChannelID,getChoice(param))
			} else {
				s.ChannelMessageSend(m.ChannelID,getChoice("yes or no"))
			}
			return
		case "coinflip":
			s.ChannelMessageSend(m.ChannelID,getChoice("heads or tails"))
			return
		case "roll":
			s.ChannelMessageSend(m.ChannelID,roll(param))
			return
	}
}
func getChoice(input string) string{
	choices := strings.Split(input," or ")
	return choices[rand.Intn(len(choices))]	
}
func rollDie(faces int,lim int,gt bool, mod int) int {
	die := 0
	for done:=true; done; done=withinLimit(die,lim,gt) {
		die = rand.Intn(faces) + 1
	}
	die = die + mod
	return die
}
func withinLimit(die int,lim int, gt bool) bool{
	if gt {
		return die <= lim
	} else {
		return die >= lim
	}
}
func roll(input string) string{
	r:= regexp.MustCompile(`\b(?P<number>[0-9]+)d(?P<faces>[0-9]+)(?P<drop>D[0-9]+[h,l])?(?P<lim>[<,>][0-9]+)?(?P<mod>[\+,\-][0-9]+)?(?P<times>x[0-9]+)?`)
	matches:=r.FindStringSubmatch(input)
	if matches == nil { return "invalid dice"}
	res:= make(map[string]string)
	for i,name:= range r.SubexpNames() {
		if i!=0 && name != ""{
			res[name] = matches[i]
		}
	}
	number,err:=strconv.Atoi(res["number"])
	if err != nil {fmt.Println(err)}
	faces,err:=strconv.Atoi(res["faces"])
	if err != nil {fmt.Println(err)}
	dropstr:=res["drop"]
	limstr:=res["lim"]
	modstr:=res["mod"]
	timestr:=res["times"]
	if number > 100 {return "I can't roll that many dice"}
	if faces > 100 {return "That's a big die"}
	rolls:=1
	if timestr != "" {
		r,err:= strconv.Atoi(timestr[1:])
		if err != nil {fmt.Println(err)}
		rolls=r
	}
	mod :=0
	if modstr != "" {
		md := modstr[:1]
		m,err:=strconv.Atoi(modstr[1:])
		if err!=nil {fmt.Println(err)}
		if md == "+" {
			mod = mod + m
		} else {
			mod = mod-m
		}
	}
	drop := 0
	high := true
	if dropstr != "" {
		re:=regexp.MustCompile(`D([0-9]+)([h,l])`)
		m:=re.FindStringSubmatch(dropstr)
		if m != nil {
			d,err := strconv.Atoi(m[1])
			if err!= nil {fmt.Println(err)}
			high = m[2] == "h"
			drop = d
		}
	}
	gt:=true
	lim:=0
	if limstr != "" {
		re:=regexp.MustCompile(`([<,>])([0-9]+)`)
		m:=re.FindStringSubmatch(limstr)
		if m!=nil {
			gt = m[1] == ">"
			l,err := strconv.Atoi(m[2])
			if err!=nil {fmt.Println(err)}
			lim = l
		}
	}
	results :=make([]string,0,0)
	for done:=true; done; done= len(results)<rolls {
		dice := make([]int,number)
		dstring:=""
		for i,_:=range dice {
			dice[i] = rollDie(faces,lim,gt,mod)
		}
		if drop > 0 {
			sort.Ints(dice)
			if high {
				dice = dice[:drop]
			} else {
				dice = dice[drop:]
			}
		}
		for i:=range dice {
			dstring+=" " + strconv.Itoa(dice[i])
		}
		dstring = strings.TrimSpace(dstring)
		s := strconv.Itoa(sum(dice)+mod)
		results = append(results,dstring + "("+s+")")
	}
	return strings.Join(results,"\n")
}
func sum(in []int) int {
	s:=0
	for i :=range in {
		s+=in[i]
	}
	return s
}
