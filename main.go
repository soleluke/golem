package main
import (
	"fmt"
	"sync"
	"flag"
	"time"
	"os"
	"os/signal"
	"syscall"
	"github.com/bwmarrin/discordgo"
	"golem/grabs"
	"golem/config"
	"golem/doots"
	"golem/tells"
	"golem/utils"
	"golem/places"
	"golem/wtells"
	"golem/triggers"
	"golem/reminds"
	"golem/markov"
	"golem/help"
)
var configFile string
var clientRead bool
func init() {
	flag.StringVar(&configFile,"config","","config file")
	clientRead = false
}
func main() {
	flag.Parse()
	fmt.Println("Using config file " + configFile)
	if configFile !="" {
		config.Initialize(configFile)
	} else {
		config.Initialize("/usr/local/etc/golem/config.json")
	}
	grabs.Initialize()
	doots.Initialize()
	tells.Initialize()
	utils.Initialize()
	places.Initialize()
	wtells.Initialize()
	triggers.Initialize()
	reminds.Initialize()
	markov.Initialize()
	help.Initialize()
	dg,err := discordgo.New("Bot " + config.Cfg.Token)
	if err != nil {
		fmt.Println("error creating discord session ", err)
		return
	}
	dg.AddHandler(grabs.MessageCreate)
	dg.AddHandler(doots.MessageCreate)
	dg.AddHandler(tells.MessageCreate)
	dg.AddHandler(utils.MessageCreate)
	dg.AddHandler(places.MessageCreate)
	dg.AddHandler(wtells.MessageCreate)
	dg.AddHandler(triggers.MessageCreate)
	dg.AddHandler(reminds.MessageCreate)
	dg.AddHandler(markov.MessageCreate)
	dg.AddHandler(help.MessageCreate)
	if config.Cfg.AllowHistoryImport {
		dg.AddHandler(messageCreate)
	}	
	err = dg.Open()
	if err != nil {
		fmt.Println("error creating connection " ,err)
		return
	}
	fmt.Println("Running as "+dg.State.User.Username+"("+dg.State.User.ID+"), press CTRL-C to exit")
	sc:=make(chan os.Signal,1)
	signal.Notify(sc,syscall.SIGINT,syscall.SIGTERM,os.Interrupt,os.Kill)
	<-sc
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate){
	if (m.Content == ".importmessages") {
		guild,_:=s.Guild(m.Message.GuildID)
		var wg sync.WaitGroup
		wg.Add(len(guild.Channels))
		for _,channel := range(guild.Channels) {
			go func(id string,name string){
				defer wg.Done()
				var content []string
				for clientRead {
					time.Sleep(1)
				}
				clientRead = true
				msgs,_ := s.ChannelMessages(id,100,"","","")
				clientRead = false
				if len(msgs)>0 {
					loop := true
					for loop {
						loop = len(msgs)==100
						for _,msg:=range(msgs){
							s:=msg.ContentWithMentionsReplaced()
							for _,e:=range(msg.Attachments) {
								s = s+ " "+e.URL
							}
							content=append(content,s)
						}
						if(len(msgs)>0){
							lastid:=msgs[len(msgs)-1].ID
							msgs,_ = s.ChannelMessages(id,100,lastid,"","")
						}
					}
					markov.ImportValues(content)
					fmt.Println("finished " + name)
				}
			} (channel.ID,channel.Name)
		}
		wg.Wait()
		fmt.Println("finished all channels")
	}
}
