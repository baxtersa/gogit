package gogit

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/baxtersa/gogit/internal"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func getAccessToken() string {
	f, err := os.Open(".access-token")
	gogit.Check(err)

	b := make([]byte, 40)
	n, err := f.Read(b)
	gogit.Check(err)
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
	gogit.Check(err)

	d, err := json.MarshalIndent(user, "", "  ")
	gogit.Check(err)
	fmt.Println(string(d))
}
