package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
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
	flagDryRun            = flag.Bool("dry-run", false, "Dry run")
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
	defer close(output)
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

type Semver struct {
	Major uint64
	Minor uint64
	Patch uint64
}

func ParseSemver(s string) (*Semver, error) {
	var major, minor, patch uint64
	_, err := fmt.Sscanf(s, "v%d.%d.%d", &major, &minor, &patch)
	if err != nil {
		return nil, err
	}

	return &Semver{
		Major: major,
		Minor: minor,
		Patch: patch,
	}, nil
}

func (s *Semver) Less(other *Semver) bool {
	if s.Major < other.Major {
		return true
	}
	if s.Major > other.Major {
		return false
	}

	if s.Minor < other.Minor {
		return true
	}
	if s.Minor > other.Minor {
		return false
	}

	if s.Patch < other.Patch {
		return true
	}
	return false
}

func (s *Semver) Equal(other *Semver) bool {
	return s.Major == other.Major && s.Minor == other.Minor && s.Patch == other.Patch
}

func (s *Semver) Greater(other *Semver) bool {
	return !s.Less(other) && !s.Equal(other)
}

func (s *Semver) IncrementPatch() {
	s.Patch++
}

func (s *Semver) String() string {
	return fmt.Sprintf("v%d.%d.%d", s.Major, s.Minor, s.Patch)
}

func findLatestTag(repo *git.Repository) (*Semver, error) {
	tagIter, err := repo.Tags()
	if err != nil {
		return nil, err
	}
	defer tagIter.Close()

	var tags []*Semver
	err = tagIter.ForEach(func(t *plumbing.Reference) error {
		if !strings.HasPrefix(t.Name().String(), "refs/tags/") {
			return nil
		}
		tagName := strings.TrimPrefix(t.Name().String(), "refs/tags/")
		tag, err := ParseSemver(tagName)
		if err != nil {
			log.Printf("Failed to parse tag: %v", err)
			return nil
		}
		tags = append(tags, tag)
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Less(tags[j])
	})

	return tags[len(tags)-1], nil
}

func parseAlreadySupportedVersion() ([]*Semver, error) {
	f, err := os.Open("ci/SUPPORTED_VERSIONS.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var versions []*Semver

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		if strings.TrimSpace(line) == "" {
			continue
		}

		ver, err := ParseSemver(line)
		if err != nil {
			return nil, err
		}

		versions = append(versions, ver)
	}

	return versions, nil
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

	alreadySupportedVersions, err := parseAlreadySupportedVersion()
	if err != nil {
		log.Fatalf("Failed to parse already supported versions: %v", err)
	}

	sort.Slice(alreadySupportedVersions, func(i, j int) bool {
		return alreadySupportedVersions[i].Less(alreadySupportedVersions[j])
	})

	var releaseVersions []*Semver
	for rel := range relChan {
		if *flagOldestReleaseDays > 0 && rel.GetPublishedAt().Before(time.Now().AddDate(0, 0, int(-*flagOldestReleaseDays))) {
			break
		}

		ver, err := ParseSemver(rel.GetTagName())
		if err != nil {
			log.Printf("Failed to parse version: %v", err)
			continue
		}
		releaseVersions = append(releaseVersions, ver)

		log.Printf("Release: %s, date: %s", rel.GetTagName(), rel.GetPublishedAt().Format(time.RFC3339))
	}

	sort.Slice(releaseVersions, func(i, j int) bool {
		return releaseVersions[i].Less(releaseVersions[j])
	})

	if releaseVersions[len(releaseVersions)-1].Greater(alreadySupportedVersions[len(alreadySupportedVersions)-1]) {
		log.Printf("New version detected: %s", releaseVersions[len(releaseVersions)-1])
		releaseVersionsStr := make([]string, len(releaseVersions))
		for i, v := range releaseVersions {
			releaseVersionsStr[i] = v.String()
		}
		slices.Reverse(releaseVersionsStr)
		os.WriteFile("ci/SUPPORTED_VERSIONS.txt", []byte(strings.Join(releaseVersionsStr, "\n")), 0644)
	} else {
		log.Printf("No new version detected, nothing to do")
		return
	}

	if *flagCommit {
		gitRepo, err := git.PlainOpen(".")
		if err != nil {
			log.Fatalf("Failed to open Git repository: %v", err)
		}

		w, err := gitRepo.Worktree()
		if err != nil {
			log.Fatalf("Failed to open Git worktree: %v", err)
		}

		_, err = w.Add("ci/SUPPORTED_VERSIONS.txt")
		if err != nil {
			log.Fatalf("Failed to add file: %v", err)
		}

		commitMessage := "[bot] Update version list and rebuild"
		if *flagGitCoauthor != "" && *flagGitCoauthorEmail != "" {
			commitMessage += "\n\nCo-authored-by: " + *flagGitCoauthor + " <" +
				*flagGitCoauthorEmail + ">"
		}

		if *flagDryRun {
			log.Printf("Dry run: committing with message: %s", commitMessage)
		} else {
			commit, err := w.Commit(commitMessage, &git.CommitOptions{
				Author: &object.Signature{Name: *flagGitName, Email: *flagGitEmail,
					When: time.Now()},
			})
			if err != nil {
				log.Fatalf("Failed to commit changes: %v", err)
			}

			_, err = gitRepo.CommitObject(commit)
			if err != nil {
				log.Fatalf("Failed to fetch commit object: %v", err)
			}
			log.Println("Successfully committed the changes")
		}

	}

	if *flagTag {
		if !*flagCommit {
			log.Fatalf("Tagging requires commit")
		}

		gitRepo, err := git.PlainOpen(".")
		if err != nil {
			log.Fatalf("Failed to open Git repository: %v", err)
		}

		latestTag, err := findLatestTag(gitRepo)
		if err != nil {
			log.Fatalf("Failed to find latest tag: %v", err)
		}

		fmt.Printf("Latest tag: %s\n", latestTag)

		latestTag.IncrementPatch()

		fmt.Printf("New tag: %s\n", latestTag)

		tagName := latestTag.String()

		headHash, err := gitRepo.Head()
		if err != nil {
			log.Fatalf("Failed to get HEAD: %v", err)
		}

		if *flagDryRun {
			log.Printf("Dry run: creating tag: %s", tagName)
		} else {
			tag, err := gitRepo.CreateTag(tagName, headHash.Hash(), &git.CreateTagOptions{
				Tagger: &object.Signature{Name: *flagGitName, Email: *flagGitEmail,
					When: time.Now()},
				Message: tagName,
			})

			if err != nil {
				log.Fatalf("Failed to create tag: %v", err)
			}

			log.Printf("Successfully created the tag: %v", tag)

		}
	}
	log.Printf("Done, version list updated and committed")
}
