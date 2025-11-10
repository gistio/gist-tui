package git

type Git struct {
	WorkDir string
	Stash   []Stash
	Log     Log
	Remote  Remote
}
