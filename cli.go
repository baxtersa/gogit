package main

import (
	"fmt"

	gh "github.com/baxtersa/gogit/github"
	tui "github.com/baxtersa/gogit/ncurses"
)

func main() {
	client := gh.Connect()
	fmt.Println(*gh.User(client).Name)
	for _, repo := range gh.Repositories(client) {
		fmt.Println(*repo.FullName)
	}
	tui.Interface()
}
