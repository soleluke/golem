package markov
import (
	"fmt"
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"golem/config"
	"github.com/mb-14/gomarkov"
	"os"
	"regexp"
	"io/ioutil"
	"strings"
	"math/rand"
	"time"
)
var chain *gomarkov.Chain
var chainChan chan string
var chainLock bool
func importChain(){
	jsonFile,err:=os.Open(config.Cfg.DataDir+"/chain.json")
	if err!=nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue,_:= ioutil.ReadAll(jsonFile)
	err=json.Unmarshal([]byte(byteValue),&chain)
	if err !=nil {
		fmt.Println(err)
		chain = gomarkov.NewChain(config.Cfg.MarkovOrder)
		fmt.Println("created new markov chain")
	} else {
		fmt.Println("loaded markov chain from file")
	}
}
func exportChain(){
	for chainLock {
		time.Sleep(1)
	}
	chainLock = true
	jsonObj,err1:=json.MarshalIndent(chain,"  ","  ")
	if err1 !=nil {
		fmt.Println(err1)
	}
	err:=ioutil.WriteFile(config.Cfg.DataDir+"/chain.json",jsonObj,0644)
	if err != nil {
		fmt.Println(err)
	}
	chainLock = false
}
func init() {
	rand.Seed(time.Now().Unix())
	chainChan=make(chan string)
	chainLock = false
}
func Initialize() {
	importChain()
	go addChainLoop()
}

func MessageCreate(s *discordgo.Session,m *discordgo.MessageCreate){
	if m.Author.ID == s.State.User.ID { return }
	chainChan <- cleanString(m.ContentWithMentionsReplaced())
	r:=regexp.MustCompile(`^\@Golem`)
	res := r.MatchString(m.ContentWithMentionsReplaced())
	if res {
		arr:=strings.Split(strings.TrimSpace(strings.Replace(cleanString(m.ContentWithMentionsReplaced()),"@Golem","",-1))," ")
		var tokens []string
		for i:= 0; i < config.Cfg.MarkovOrder-1; i++ {
			tokens = append(tokens,gomarkov.StartToken)
		}
		if arr[0] != ""{
			tok:=arr[rand.Intn(len(arr))]
			tokens = append(tokens, tok)
		} else {
			tokens = append(tokens,gomarkov.StartToken)
		}
		for tokens[len(tokens)-1] != gomarkov.EndToken && len(tokens)<100 {
			for chainLock {
				time.Sleep(1)
			}
			chainLock = true
			next,_ := chain.Generate(tokens[len(tokens)-config.Cfg.MarkovOrder:])
			chainLock = false
			tokens = append(tokens,next)
		}
		resp:=""
		for _,token:=range(tokens) {
			if token!=gomarkov.StartToken && token!=gomarkov.EndToken {
				resp = resp+" "+token
			}
		}
		s.ChannelMessageSend(m.ChannelID,resp)
	}
}
func cleanString(val string) string{
	return strings.TrimSpace(strings.ReplaceAll(val,"\n",""))
	reg:=regexp.MustCompile(`[]`)
	return reg.ReplaceAllString(val,"")
}

func ImportValues(vals []string){
	for _,val:=range(vals){
		chainChan <- cleanString(val)
	}
	exportChain()
}
func addChainLoop() {
	for true {
		val:=<-chainChan
		if val != "" {
			strs:=strings.Split(val," ")
			if len(strs) > 0 {
				for chainLock {
					time.Sleep(1)
				}
				chainLock = true
				chain.Add(strs)
				chainLock = false
			}
		}
	}
}
