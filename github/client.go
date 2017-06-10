package gogit

import (
	"context"
	"os"

	common "github.com/baxtersa/gogit/internal"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type Request interface {
	make() []string
}

type ReqRepo struct {
	client *github.Client
}
type ReqUser struct {
	client *github.Client
}
type ReqIssue struct {
	client *github.Client
}

func (req ReqRepo) make() []string {
	ctx := context.Background()
	repos, _, err := req.client.Repositories.List(ctx, "", nil)
	common.Check(err)

	strs := []string{}
	for _, repo := range repos {
		strs = append(strs, *repo.FullName)
	}
	return strs
}

func (req ReqUser) make() []string {
	ctx := context.Background()
	user, _, err := req.client.Users.Get(ctx, "")
	common.Check(err)

	return []string{*user.Name}
}

func (req ReqIssue) make() []string {
	ctx := context.Background()
	issues, _, err := req.client.Issues.ListByRepo(ctx, "plasma-umass", "Stopify", nil)
	common.Check(err)

	strs := []string{}
	for _, issue := range issues {
		strs = append(strs, *issue.Title)
	}
	return strs
}

type ReqChannels struct {
	Repo   chan bool
	User   chan bool
	Issue  chan bool
	Client *github.Client
}

type RespChannels struct {
	Repo  chan []string
	User  chan []string
	Issue chan []string
}

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

func Issues(client *github.Client) []*github.Issue {
	ctx := context.Background()
	issues, _, err := client.Issues.List(ctx, true, nil)
	common.Check(err)

	return issues
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

func HandleRequests(reqs *ReqChannels, resps *RespChannels) {
	for {
		select {
		case <-reqs.Repo:
			req := ReqRepo{client: reqs.Client}
			resps.Repo <- req.make()
		case <-reqs.User:
			req := ReqUser{client: reqs.Client}
			resps.User <- req.make()
		case <-reqs.Issue:
			req := ReqIssue{client: reqs.Client}
			resps.Issue <- req.make()
		}
	}
}

func Connect() *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: getAccessToken()},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}
