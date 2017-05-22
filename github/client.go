package gogit

import (
	"context"
	"os"

	common "github.com/baxtersa/gogit/internal"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func getAccessToken() string {
	f, err := os.Open(".access-token")
	common.Check(err)

	b := make([]byte, 40)
	n, err := f.Read(b)
	common.Check(err)
	if n != 40 {
		panic("Invalid GitHub Access Token")
	}

	return string(b)
}

func Repositories(client *github.Client) []*github.Repository {
	ctx := context.Background()
	repos, _, err := client.Repositories.List(ctx, "", nil)
	common.Check(err)

	return repos
}

func User(client *github.Client) *github.User {
	ctx := context.Background()
	user, _, err := client.Users.Get(ctx, "")
	common.Check(err)

	return user
}

func Connect() *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: getAccessToken()},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}
