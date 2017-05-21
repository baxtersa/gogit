package gogit

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func getAccessToken() string {
	f, err := os.Open(".access-token")
	check(err)

	b := make([]byte, 40)
	n, err := f.Read(b)
	check(err)
	if n != 40 {
		panic("Invalid GitHub Access Token")
	}

	return string(b)
}

func Connect() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: getAccessToken()},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	user, _, err := client.Users.Get(ctx, "")
	check(err)

	d, err := json.MarshalIndent(user, "", "  ")
	check(err)
	fmt.Println(string(d))
}
