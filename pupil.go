package main

type Pupil struct {
	name           string
	isMale         bool
	group          int
	score          int
	prefs          []int
	groupBestScore int
	initialGroup   int
	locked         bool
	numOfMoves     int32
}

func (p *Pupil) IsMale() bool {
	return p.isMale
}
