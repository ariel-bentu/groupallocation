package main

type Pupil struct {
	id                     int
	name                   string
	isMale                 bool
	prefs                  []int
	incomingPrefs          []int
	groupBestScore         int
	numOfMoves             int32
	groupsCount            int
	unsatisfiedPrefsCount  int64
	unsatisfiedGroupsCount []int64

	prefInactive     bool
	uniteGroups      []int
	seperationGroups []int
}

func (p *Pupil) IsMale() bool {
	return p.isMale
}

type ByGroupCount []*Pupil

func (a ByGroupCount) Len() int      { return len(a) }
func (a ByGroupCount) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByGroupCount) Less(i, j int) bool {
	return a[i].groupsCount > a[j].groupsCount || (a[i].groupsCount == a[j].groupsCount && len(a[i].prefs) > len(a[j].prefs))
}

type SortPupilList struct {
	list []int
	ec   *ExecutionContext
}

func (a *SortPupilList) Len() int      { return len(a.list) }
func (a *SortPupilList) Swap(i, j int) { a.list[i], a.list[j] = a.list[j], a.list[i] }
func (a *SortPupilList) Less(i, j int) bool {
	return a.ec.pupils[a.list[i]].name < a.ec.pupils[a.list[j]].name
}
