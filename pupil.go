package main

type Pupil struct {
	name                  string
	isMale                bool
	startGroup            int
	prefs                 []int
	groupBestScore        int
	lockedGroup           int
	locked                bool
	numOfMoves            int32
	groupsCount           int
	unsatisfiedPrefsCount int64
}

func (p *Pupil) IsMale() bool {
	return p.isMale
}

type ByGroupCount []*Pupil

func (a ByGroupCount) Len() int           { return len(a) }
func (a ByGroupCount) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByGroupCount) Less(i, j int) bool { return a[i].groupsCount > a[j].groupsCount }
