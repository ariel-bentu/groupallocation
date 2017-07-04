package main

type PartialAssignment []int

func (pa PartialAssignment) Count() int {
	return len(pa)
}

func (pa PartialAssignment) GetGroup(pupil int) int {
	return pa[pupil]
}

func RunBackTrack(ec *ExecutionContext) {
	bt(ec, first(ec, nil))
	if ec.done {
		ec.Cancel = false
	} else {
		ec.Cancel = true
	}
}

//c is a candidate: a slice of 0..k (k<num of pupils), each represents the group for a pupil

func bt(ec *ExecutionContext, c []int) {

	if ec.Cancel {
		return
	}
	if accept(ec, c) {
		ec.statusCandidate = c
		if len(c) == len(ec.pupils) {
			//found a solution
			for i, p := range ec.pupils {
				p.groupBestScore = c[i]
			}
			ec.Finish()
			return
		}

	} else {
		return
	}
	ec.currentIteration = len(c) + 1
	nextCandidate := first(ec, c)
	for nextCandidate != nil && !ec.done {
		bt(ec, nextCandidate)
		nextCandidate = next(ec, nextCandidate)
	}
}

func accept(ec *ExecutionContext, c []int) bool {
	k := len(c)

	p := ec.pupils[k-1]
	if p.locked && p.lockedGroup != c[k-1] {
		//return false
	}

	//prefs
	for i := len(c) - 1; i >= 0; i-- {
		p = ec.pupils[i]

		if len(p.prefs) > 0 {

			inRangeCount := 0
			refCount := 0

			for j := 0; j < len(p.prefs); j++ {
				if p.prefs[j] < k {
					inRangeCount++
					if c[p.prefs[j]] == c[i] {
						refCount++
						break
					}
				} else {
					break
				}
			}

			if inRangeCount > 0 && inRangeCount == len(p.prefs) && refCount == 0 {
				p.unsatisfiedPrefsCount++
				return false
			}

		}
	}

	//constraints
	for _, csg := range ec.Constraints {

		if !csg.ValidateNew(ec, PartialAssignment(c)) {
			csg.unsatisfiedCount++
			return false
		}
	}
	return true
}

/*
func getBackJump(arr []int, currentPupil int) int {
	currentMax := -1
	for j := 0; j < len(arr); j++ {
		if arr[j] < currentPupil {
			currentMax = max(currentMax, arr[j])
		}
	}

	return currentMax
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
*/
func first(ec *ExecutionContext, c []int) []int {
	k := len(c)
	if k == len(ec.pupils) {
		return nil
	}
	return append(c, ec.pupils[k].startGroup)
}

func next(ec *ExecutionContext, s []int) []int {
	k := len(s) - 1
	s[k]++
	if s[k] == ec.groupsCount {
		s[k] = 0
	}
	if s[k] == ec.pupils[k].startGroup {
		return nil
	}

	return s
}
