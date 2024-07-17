package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/google/go-github/v63/github"
)

var (
	flagOwner             = flag.String("owner", "gotify", "Owner")
	flagRepo              = flag.String("repo", "server", "Repository")
	flagOldestReleaseDays = flag.Uint64("days", 0, "Oldest release days")
	flagCommit            = flag.Bool("commit", false, "Update version list and commit")
	flagTag               = flag.Bool("tag", false, "Tag")
	flagGitName           = flag.String("name", "upstream-follower", "Git user name")
	flagGitEmail          = flag.String("email", "upstream-follower@localhost", "Git user email")
	flagGitCoauthor       = flag.String("coauthor", "", "Git coauthor")
	flagGitCoauthorEmail  = flag.String("coauthor-email", "", "Git coauthor email")
)

func dumpFlags() {
	log.Printf("owner: %s", *flagOwner)
	log.Printf("repo: %s", *flagRepo)
	log.Printf("days: %d", *flagOldestReleaseDays)
	log.Printf("commit: %t", *flagCommit)
	log.Printf("tag: %t", *flagTag)
	log.Printf("name: %s", *flagGitName)
	log.Printf("email: %s", *flagGitEmail)
	log.Printf("coauthor: %s", *flagGitCoauthor)
	log.Printf("coauthor-email: %s", *flagGitCoauthorEmail)
}

func autoEnumPages[T any](
	ctxFactory func() (context.Context, func()),
	queryFunc func(ctx context.Context, pagination *github.ListOptions) (ret []*T, resp *github.Response, err error),
	output chan<- *T,
) {
	var pagination *github.ListOptions
	ctx, cancel := ctxFactory()
	defer cancel()

	for {
		ret, resp, err := queryFunc(ctx, pagination)
		if err != nil {
			log.Fatalf("autoEnumPages: %v", err)
		}

		for _, v := range ret {
			output <- v
		}

		if resp.NextPage == 0 {
			break
		}

		pagination = &github.ListOptions{
			Page: resp.NextPage,
		}

		ctx, cancel = ctxFactory()
		defer cancel()
	}
}

func main() {
	flag.Parse()

	dumpFlags()

	client := github.NewClient(nil)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	relChan := make(chan *github.RepositoryRelease)
	go autoEnumPages(func() (context.Context, func()) {
		return ctx, cancel
	}, func(ctx context.Context, pagination *github.ListOptions) ([]*github.RepositoryRelease, *github.Response, error) {
		ret, resp, err := client.Repositories.ListReleases(ctx, *flagOwner, *flagRepo, pagination)
		if err != nil {
			return nil, nil, err
		}
		return ret, resp, nil
	}, relChan)

	for rel := range relChan {
		if *flagOldestReleaseDays > 0 && rel.GetPublishedAt().Before(time.Now().AddDate(0, 0, int(-*flagOldestReleaseDays))) {
			break
		}

		log.Printf("Release: %s, date: %s", rel.GetTagName(), rel.GetPublishedAt().Format(time.RFC3339))
	}

}
