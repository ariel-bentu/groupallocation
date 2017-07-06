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

type SortPupilList struct {
	list []int
	ec   *ExecutionContext
}

func (a *SortPupilList) Len() int      { return len(a.list) }
func (a *SortPupilList) Swap(i, j int) { a.list[i], a.list[j] = a.list[j], a.list[i] }
func (a *SortPupilList) Less(i, j int) bool {
	return a.ec.pupils[a.list[i]].name < a.ec.pupils[a.list[j]].name
}
