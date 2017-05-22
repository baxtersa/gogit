package main

import (
	client "github.com/baxtersa/gogit/github"
	tui "github.com/baxtersa/gogit/ncurses"
)

func main() {
	client.Connect()
	tui.Interface()
}
