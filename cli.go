package main

import (
	gh "github.com/baxtersa/gogit/github"
	tui "github.com/baxtersa/gogit/ncurses"
)

func main() {
	client := gh.Connect()

	reqs := gh.ReqChannels{
		Repo:   make(chan bool),
		User:   make(chan bool),
		Issue:  make(chan bool),
		Client: client,
	}
	resps := gh.RespChannels{
		Repo:  make(chan []string),
		User:  make(chan []string),
		Issue: make(chan []string),
	}

	go gh.HandleRequests(&reqs, &resps)
	go tui.Interface(&reqs, &resps)

	defer tui.End()

	<-tui.Quit
}
