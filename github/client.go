package gogit

import (
	"context"
	"os"

	common "github.com/baxtersa/gogit/internal"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Generic type to emit a request to a GH client
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

// Make a request for repositories
func (req ReqRepo) make() []string {
	ctx := context.Background()
	var allRepos []*github.Repository
	opt := &github.RepositoryListOptions{}
	for {
		repos, resp, err := req.client.Repositories.List(ctx, "", opt)
		common.Check(err)

		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	strs := []string{}
	for _, repo := range allRepos {
		strs = append(strs, *repo.FullName)
	}
	return strs
}

// Make a request for the authenticated user
func (req ReqUser) make() []string {
	ctx := context.Background()
	user, _, err := req.client.Users.Get(ctx, "")
	common.Check(err)

	return []string{*user.Name}
}

// Make a request for issues
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

// Collection of channels to asynchronously prompt for GH client requests
type ReqChannels struct {
	Repo   chan bool
	User   chan bool
	Issue  chan bool
	Client *github.Client
}

// Collection of channels to asynchronously receive responses from
// GH client requests
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

// Asynchronously block and wait for prompts to make GH client requests
// params:
//   `reqs`: Channels to wait for prompt to make requests
//   `resps`: Channels to send results of completed requests
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

// Authenticate a GH client
// returns:
//   An authenticated GH client
func Connect() *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: getAccessToken()},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}
