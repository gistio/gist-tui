package git

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/google/uuid"
)

type Stash struct {
	Message,
	Diff,
	ID,
	Ref string
	Index int
}

func (g *Git) StashList() []Stash {
	if len(g.Stash) > 0 {
		return g.Stash
	}

	cmd := exec.Command("bash", "-c", "git stash list | cat")
	cmd.Dir = g.WorkDir
	res, err := cmd.Output()
	list := []Stash{}
	if err != nil {
		return list
	}
	for index, item := range strings.Split(string(res), "\n") {
		if item == "" {
			continue
		}
		list = append(list, Stash{Message: item, ID: uuid.New().String(), Index: index})
	}
	if len(g.Stash) == 0 {
		g.Stash = list
	}
	return list
}

func (g *Git) StashDiff(index int) string {
	cmdText := fmt.Sprintf("stash show -p stash@{%d} --color=always", index)
	cmd := exec.Command("git", strings.Split(cmdText, " ")...)
	cmd.Dir = g.WorkDir
	res, err := cmd.Output()
	if err != nil {
		return ""
	}

	return string(res)
}

func (g *Git) StashApplyIndex(index int) bool {
	cmdText := fmt.Sprintf("stash apply stash@{%d}", index)
	cmd := exec.Command("git", strings.Split(cmdText, " ")...)
	cmd.Dir = g.WorkDir
	_, err := cmd.Output()
	if err != nil {
		return false
	}

	return true
}

func (g *Git) StashSearch(term string) []Stash {
	list := []Stash{}

	for _, item := range g.StashList() {
		if strings.Contains(strings.ToLower(item.Message), strings.ToLower(term)) {
			list = append(list, item)
		}
	}
	return list
}

func (g *Git) GetStash(index int) (Stash, bool) {
	if index > len(g.Stash)-1 || index < 0 {
		return Stash{}, false
	}
	return g.Stash[index], true
}

func (g *Git) GetStashById(id string) (Stash, int) {
	index := -1
	for i, item := range g.Stash {
		if item.ID == id {
			return item, i
		}
	}
	return Stash{}, index
}
