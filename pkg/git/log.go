package git

import (
	"fmt"
	"os/exec"
	"strings"
)

type Log struct {
	Commits []Commit
}

type Commit struct {
	ID,
	Author,
	Date,
	Name,
	CommitHash,
	Message string
}

func parseHash(str string) string {
	fields := strings.Fields(str)
	if len(fields) > 1 && fields[0] == "commit" {
		return fields[1]
	}
	return ""
}

func (g *Git) LogList() []Commit {
	if g.Log.Commits == nil {
		g.Log.Commits = []Commit{}
	}
	if len(g.Log.Commits) > 0 {
		return g.Log.Commits
	}
	list := []Commit{}
	cmd := exec.Command("bash", "-c", "git log | cat")
	cmd.Dir = g.WorkDir
	res, err := cmd.Output()
	if err != nil {
		return list
	}

	index := 0
	for _, item := range strings.Split(string(res), "\n") {
		if strings.HasPrefix(item, "commit") && len(list) > 0 {
			index += 1
			list = append(list, Commit{})
			list[index].ID = parseHash(item)
			list[index].CommitHash = parseHash(item)
			continue
		}
		if len(list)-1 < index || (index == 0 && len(list) == 0) {
			list = append(list, Commit{})
		}
		listItem := &list[index]
		if strings.HasPrefix(item, "commit") {
			listItem.ID = parseHash(item)
			list[index].CommitHash = parseHash(item)
		} else if strings.HasPrefix(item, "Author") {
			listItem.Author = item
		} else if strings.HasPrefix(item, "Date") {
			listItem.Date = item
		} else {
			listItem.Message += item
		}
	}

	g.Log.Commits = list
	return list
}

func (g *Git) ShowCommitDiff(hash string) string {
	if hash == "" {
		return ""
	}

	cmd := exec.Command("bash", "-c", fmt.Sprintf("git show %s --color=always | cat", hash))
	cmd.Dir = g.WorkDir
	res, err := cmd.Output()
	if err != nil {
		return ""
	}

	return string(res)
}

func (g *Git) FindLogByHash(hash string) (Commit, bool) {
	c := Commit{}
	ok := false
	if hash == "" {
		return c, ok
	}

	for _, Commit := range g.Log.Commits {
		if Commit.CommitHash == hash {
			c = Commit
			ok = true
			break
		}
	}
	return c, ok
}

func (g *Git) LogSearch(searchTerm string) []Commit {
	list := []Commit{}
	for _, item := range g.LogList() {
		if strings.Contains(item.Message, searchTerm) {
			list = append(list, item)
		}
	}
	return list
}
