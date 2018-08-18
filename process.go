package main

import "time"
import "fmt"

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

func bt(ec *ExecutionContext, c []int) int {

	ec.currentLevelCount++

	if ec.currentLevelCount > 10000000 && ec.graceLevel < 6 {
		ec.currentLevelCount = 0
		ec.graceLevel++
	}

	if time.Now().Sub(ec.startTime).Seconds() > float64(ec.timeLimit) && ec.resultsCount > 0 {

		ec.Finish()
		return -1
	}

	if ec.Cancel {
		return -1
	}
	jump := accept(ec, c)
	if jump == 0 {
		ec.statusCandidate = c
		if len(c) == len(ec.pupils) {
			//found a solution

			//stores this candidate and continue
			_, sumOfSat, someOfSatFirsts := getPreferencesScore(ec, PartialAssignment(c))
			if ec.resultsCount == 0 ||
				sumOfSat > ec.bestSumOfSatisfiedPrefs ||
				sumOfSat == ec.bestSumOfSatisfiedPrefs && someOfSatFirsts > ec.bestSumOfSatisfiedFirstPrefs {
				ec.resultsScoreHistory = append(ec.resultsScoreHistory, ec.bestSumOfSatisfiedPrefs)
				ec.bestSumOfSatisfiedPrefs = sumOfSat
				ec.bestSumOfSatisfiedFirstPrefs = someOfSatFirsts
				ec.bestCandidate = make([]int, len(c))
				copy(ec.bestCandidate, c)
			}
			ec.resultsCount++
			return -1
		}

	} else {
		return jump
	}

	ec.currentIteration = len(c) + 1
	nextCandidate := first(ec, c)
	jump = -1
	for nextCandidate != nil && !ec.done {
		jump = bt(ec, nextCandidate)
		if jump == -1 || len(c) <= jump+1 {
			nextCandidate = next(ec, nextCandidate)
		} else {
			nextCandidate = nil
		}
	}
	return jump
}

func accept(ec *ExecutionContext, c []int) int {
	k := len(c)

	//p := ec.pupils[k-1]

	//constraints
	for _, csg := range ec.Constraints {

		if !csg.ValidateNew(ec, PartialAssignment(c)) {
			csg.unsatisfiedCount++

			//try to jump back
			jump := -1
			for i := 0; i < len(csg.members); i++ {
				if csg.members[i] == k-1 { //this is the member who failed it
					for ec.pupils[csg.members[i]].optionsLeft == 0 {
						if i == len(csg.members)-1 {
							jump = csg.members[i] - 1
							break
						} else {
							i++
							jump = csg.members[i]
						}
					}
					break
				}
			}

			return jump
		}
	}

	//prefs

	if ec.prefFailCount > 10000000 && ec.resultsCount == 0 {
		//can't pass the pref, start disabling
		for _, p := range ec.pupils {
			if p.unsatisfiedPrefsCount > ec.prefThreashold {
				ec.prefThreashold = p.unsatisfiedPrefsCount
			}
		}
		ec.prefFailCount = 0
	}

	for i := len(c) - 1; i >= 0; i-- {
		p := ec.pupils[i]

		if p.prefInactive {
			continue
		}

		if ec.prefThreashold > 0 && p.unsatisfiedPrefsCount > ec.prefThreashold {
			//	p.prefInactive = true
			fmt.Printf("Disable prefs for %s", p.name)
		}

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

				ec.prefFailCount++
				jumpTo := -1
				for j := 0; j < len(p.prefs); j++ {
					if ec.pupils[p.prefs[j]].optionsLeft == 0 {
						if j < len(p.prefs)-1 {
							jumpTo = p.prefs[j+1]
						} else if p.optionsLeft > 0 {
							jumpTo = i
						} else {
							jumpTo = p.prefs[0] - 1
							if i < jumpTo {
								//jumpTo = i - 1
							}
						}
					} else {
						break
					}
				}

				if jumpTo >= 0 && i > jumpTo && p.optionsLeft > 0 {
					return i
				}
				return jumpTo
			}

		}
	}

	//ec.prefFailCount = 0

	return 0
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
	if !ec.pupils[k].locked {
		ec.pupils[k].optionsLeft = ec.groupsCount - 1
	}
	return append(c, ec.pupils[k].startGroup)
}

func next(ec *ExecutionContext, s []int) []int {
	k := len(s) - 1

	if ec.pupils[k].optionsLeft == 0 {
		return nil
	}

	s[k]++
	ec.pupils[k].optionsLeft--

	if s[k] == ec.groupsCount {
		s[k] = 0
	}
	//if s[k] == ec.pupils[k].startGroup {
	//	return nil
	//}

	return s
}

/*
	returns (total preferences, ammount satisfied, amount of first pref satisfied)
*/
func getPreferencesScore(ec *ExecutionContext, candidate Candidate) (int, int, int) {
	firstSatisfiedSum := 0
	totalSatisfiedSum := 0
	total := 0
	k := candidate.Count()
	for i := 0; i < k; i++ {
		p := ec.pupils[i]
		g := candidate.GetGroup(i)
		for j := 0; j < len(p.prefs); j++ {
			if p.prefs[j] < k {
				total++
				if g == candidate.GetGroup(p.prefs[j]) {
					totalSatisfiedSum++
					if j == 0 {
						firstSatisfiedSum++
					}
				}
			}
		}
	}

	return total, totalSatisfiedSum, firstSatisfiedSum

}
