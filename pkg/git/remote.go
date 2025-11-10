package git

import (
	"fmt"
	"os/exec"
	"strings"
)

type Remote struct {
	Url string
}

func (g *Git) GetRemoteUrl() string {
	if g.Remote.Url != "" {
		return g.Remote.Url
	}

	cmdText := fmt.Sprintf("config --get remote.origin.url")
	cmd := exec.Command("git", strings.Split(cmdText, " ")...)
	cmd.Dir = g.WorkDir
	res, err := cmd.Output()

	if err != nil {
		return ""
	}

	g.Remote.Url = string(res)
	g.Remote.Url = strings.Replace(g.Remote.Url, ".git", "", 1)
	g.Remote.Url = strings.Replace(g.Remote.Url, "\n", "", 1)
	return g.Remote.Url
}
