package main

import (
	"fmt"

	gh "github.com/baxtersa/gogit/github"
	tui "github.com/baxtersa/gogit/ncurses"
	"github.com/google/go-github/github"
)

type Channels struct {
	repo  chan *github.Repository
	user  chan *github.User
	issue chan *github.Issue
}

func main() {
	client := gh.Connect()
	fmt.Println(*gh.User(client).Name)
	for _, repo := range gh.Repositories(client) {
		fmt.Println(*repo.FullName)
	}
	go tui.Interface()
	defer tui.End()

loop:
	for {
		select {
		case <-tui.In:
			break loop
		}
	}
}
