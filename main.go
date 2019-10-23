package main
import (
	"fmt"
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
)
func init() {
}
func main() {
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
	fmt.Println(m.Content)
}
