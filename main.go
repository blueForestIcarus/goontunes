package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/robfig/cron/v3"
	"gopkg.in/yaml.v2"
)

var configfile = ".config"
var helpurl = "https://github.com/blueforesticarus/goontunes"

type State struct {
	Discord   DiscordApp
	Spotify   SpotifyApp // might not make sense to be list
	Playlists []*Playlist

	CachePath string
	em        *EntryManager
	lib       *Library

	plumber *Plumber
}

var global State

func main() {

	bytes, err := ioutil.ReadFile(configfile)
	if err != nil {
		fmt.Printf("Error, cannot open config file %v\n", configfile)
		return
	}

	err = yaml.UnmarshalStrict(bytes, &global)
	if err != nil {
		fmt.Printf("Could not parse config %v\n", err)
		fmt.Printf("Follow instructions at %s\n", helpurl)
		return
	}

	err = os.MkdirAll(global.CachePath, 0755)
	if err != nil {
		fmt.Printf("Could not create cachepath %s, %v\n", global.CachePath, err)
		return
	}

	global.em = new_EntryManager()
	global.em.cachepath = global.CachePath + "/entries"
	global.em.load()

	global.lib = new_Library()
	global.lib.cachepath = global.CachePath + "/library"
	global.lib.load()

	global.plumber = new_Plumber()
	go global.Spotify.connect()
	go global.Discord.connect()
	//global.Spotify.init_playlists() now in spotify.connect

	c := cron.New()
	c.AddFunc("@every 1h", func() {
		fmt.Printf("1hr timer\n")
		global.plumber.rescan()
		global.plumber.j_playlist_task.Trigger()
	})
	c.Start()

	exitSignal := make(chan os.Signal)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
}

// std lib devs are clowning
func contains_a_fucking_string(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
func fucking_index(slice []string, item string) int {
	for i := range slice {
		if slice[i] == item {
			return i
		}
	}
	return -1
}
