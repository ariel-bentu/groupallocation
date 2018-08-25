package main

import (
	"math"
	"time"
)

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

	//jump := accept(ec, c)
	//if jump == 0 {
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

	//} else {
	//	return jump
	//}

	jump := -1
	vd := ec.domainValues.getEffectiveDomainValues(80)
	if len(c) == 70 {
		stop2(vd.values)
	}

	domainNotEmpty := fc(ec, c, len(c)-1, len(c)-1)

	if ec.MaxDepth < len(c)-1 {
		ec.MaxDepth = len(c) - 1
	}

	if domainNotEmpty {

		ec.currentIteration = len(c) + 1
		nextCandidate := first(ec, c)

		for nextCandidate != nil && !ec.done {
			jump = bt(ec, nextCandidate)
			if jump == -1 || len(c) <= jump+1 {
				nextCandidate = next(ec, nextCandidate)
			} else {
				nextCandidate = nil
			}
		}
	}

	ec.domainValues.PopAllDomainRestriction(len(c) - 1)
	return jump
}

func accept(ec *ExecutionContext, c []int) int {
	//k := len(c)

	//constraints
	for _, csg := range ec.Constraints {
		if csg.IsUnite {
			//no need to check because of the domain retrictions
			continue
		}
		if !csg.ValidateNew(ec, PartialAssignment(c)) {
			csg.unsatisfiedCount++

			//try to jump back

			jump := -1

			return jump
		}
	}

	//prefs
	/*
		if ec.prefFailCount > 10000000 && ec.resultsCount == 0 {
			//can't pass the pref, start disabling
			for _, p := range ec.pupils {
				if p.unsatisfiedPrefsCount > ec.prefThreashold {
					ec.prefThreashold = p.unsatisfiedPrefsCount
				}
			}
			ec.prefFailCount = 0
		}
	*/

	for i := len(c) - 1; i >= 0; i-- {
		p := ec.pupils[i]

		if p.prefInactive {
			continue
		}
		satisfied, _, _ := prefStatus(p, c, i)
		if !satisfied {

			p.unsatisfiedPrefsCount++
			if p.unsatisfiedPrefsCount > 100000 {
				stop()
			}

			//optionsLeftForPupil := ec.domainValues.GetEffectiveDomainValues(i).AreOptionsLeft(c[i])
			ec.prefFailCount++
			jumpTo := -1

			vd := ec.domainValues.getEffectiveDomainValues(38)
			stop2(vd.values)
			return jumpTo
		}
		/*
			if ec.prefThreashold > 0 && p.unsatisfiedPrefsCount > ec.prefThreashold {
				p.prefInactive = true
				fmt.Printf("Disable prefs for %s", p.name)
			}
		*/

	}

	return 0
}

func prefStatus(p *Pupil, c []int, pupilIndex int) (bool, int, int) {
	k := len(c)
	if pupilIndex >= k {
		return true, 0, 0
	}

	inRangeCount := 0
	refCount := 0
	for j := len(p.prefs) - 1; j >= 0; j-- {
		if p.prefs[j] < k {
			inRangeCount++
			if c[p.prefs[j]] == c[pupilIndex] {
				refCount++
				break
			}
		} else {
			break
		}
	}
	satisfied := true

	if len(p.prefs) > 0 && inRangeCount == len(p.prefs) && refCount == 0 {
		satisfied = false
	}
	return satisfied, inRangeCount, refCount
}

func first(ec *ExecutionContext, c []int) []int {
	k := len(c)
	if k == len(ec.pupils) {
		return nil
	}
	preferedValue := 0
	if k > 0 {
		grp := ec.allGroup
		if ec.pupils[k-1].IsMale() {
			grp = ec.maleGroup
		}

		lowestCount := grp.countForGroup[0]
		for val, count := range grp.countForGroup[1:] {
			if count < lowestCount {
				lowestCount = count
				preferedValue = val + 1
			}
		}
	}

	startValue := ec.domainValues.FirstValue(k, preferedValue)
	if startValue == -1 {
		return nil
	}
	c = append(c, startValue)

	return c
}

func fcSeperateGroups(ec *ExecutionContext, c []int, pupilIndex int, originalPupilIndex int) bool {
	k := len(c) - 1

	if pupilIndex < k {
		return true
	}
	if len(ec.pupils[k].seperationGroups) > 0 {
		for _, i := range ec.pupils[k].seperationGroups {
			g := ec.Constraints[i]
			if !g.disabled {
				boysLeft, girlsLeft := g.calculateMembersCounts(ec, PartialAssignment(c))
				maxAllowed := g.maxAllowed
				minAllowed := g.minAllowed

				if k == 68 && i == 15 {
					stop()
				}

				for j := 0; j < ec.groupsCount; j++ {
					if float64(g.countForGroup[j]) == math.Ceil(maxAllowed) {

						domainNotEmpty, _ := ec.domainValues.pushGroupRestriction(ec, i, originalPupilIndex, k, j, false, 0)
						if !domainNotEmpty {
							return false
						}
					}

					if math.Floor(minAllowed)-float64(boysLeft+girlsLeft)-float64(g.countForGroup[j]) == 0 {
						domainNotEmpty, _ := ec.domainValues.pushGroupRestriction(ec, i, originalPupilIndex, k, j, true, 0)
						if !domainNotEmpty {
							return false
						}
					}

					if !fcGenderMinCheck(ec, c, i, originalPupilIndex, k, j, 1, g.boysForGroup[j], boysLeft, g.minBoys, g.disallowZeroBoys) {
						return false
					}
					girlsForGroup := g.countForGroup[j] - g.boysForGroup[j]
					if !fcGenderMinCheck(ec, c, i, originalPupilIndex, k, j, 2, girlsForGroup, girlsLeft, g.minGirls, g.disallowZeroGirls) {
						return false
					}
				}

			}
		}
	}
	return true

}

func fcGenderMinCheck(ec *ExecutionContext, c []int, gInx int, originalPupilIndex int, pupilInx int, value int, gender int,
	genderForGroup int, genderLeft int, minAllowed int, disallowZero bool) bool {
	if genderForGroup+genderLeft <= minAllowed {
		if genderForGroup > 0 || disallowZero {
			return false //no chance to get to min boys required or be 0 boys
		} else if genderForGroup+genderLeft < minAllowed {
			//for boys to stay 0
			domainNotEmpty, _ := ec.domainValues.pushGroupRestriction(ec, gInx, originalPupilIndex, pupilInx, value, false, gender)
			if !domainNotEmpty {
				return false
			}

		}

	}
	return true
}

func fc(ec *ExecutionContext, c []int, pupilIndex int, originalPupilIndex int) bool {
	if pupilIndex < len(c)-1 {
		return true
	}
	domainNotEmpty := fcUniteGroups(ec, c, pupilIndex, originalPupilIndex)
	if domainNotEmpty {

		domainNotEmpty = fcPref(ec, c, pupilIndex, originalPupilIndex)
		if domainNotEmpty {
			return fcSeperateGroups(ec, c, pupilIndex, originalPupilIndex)
		}
	}
	return false
}

func fcValues(ec *ExecutionContext, c []int, pupilIndex int) (int, bool) {
	var value int
	if pupilIndex >= len(c) {
		vd := ec.domainValues.getEffectiveDomainValues(pupilIndex)
		if len(vd.values) > 1 {
			//todo
			return 0, true
		}
		value = vd.values[0]
	} else {
		value = c[pupilIndex]
	}
	return value, false
}

func fcUniteGroups(ec *ExecutionContext, c []int, pupilIndex, originalPupilIndex int) bool {
	if pupilIndex < len(c)-1 {
		return true
	}

	value, multiple := fcValues(ec, c, pupilIndex)
	if multiple {
		return true
	}

	if len(ec.pupils[pupilIndex].uniteGroups) > 0 {
		for _, i := range ec.pupils[pupilIndex].uniteGroups {
			if !ec.Constraints[i].disabled {

				domainNotEmpty, changed := ec.domainValues.pushGroupRestriction(ec, i, originalPupilIndex, pupilIndex, value, true, 0)
				if !domainNotEmpty {
					return false
				}

				for _, changedInx := range changed {
					domainNotEmpty := fc(ec, c, changedInx, originalPupilIndex)
					if !domainNotEmpty {
						return false
					}
				}
			}
		}
	}
	return true
}

func fcPref(ec *ExecutionContext, c []int, pupilIndex int, originalPupilIndex int) bool {
	if pupilIndex < len(c)-1 {
		return true
	}

	value, multiple := fcValues(ec, c, pupilIndex)
	if multiple {
		return true
	}
	prefCount := len(ec.pupils[pupilIndex].prefs)

	_, inRangePrefCount, foundPrefCount := prefStatus(ec.pupils[pupilIndex], c, pupilIndex)

	if prefCount > 0 && foundPrefCount == 0 &&
		prefCount-inRangePrefCount == 1 {
		domainNotEmpty, changed := ec.domainValues.PushDomainOnlyOneRestriction(originalPupilIndex, ec.pupils[pupilIndex].prefs[0], value)
		if !domainNotEmpty {
			ec.pupils[pupilIndex].unsatisfiedPrefsCount++
			return false
		}

		if changed {
			domainNotEmpty := fc(ec, c, ec.pupils[pupilIndex].prefs[0], originalPupilIndex)
			if !domainNotEmpty {
				return false
			}
		}

	}

	for _, inPerf := range ec.pupils[pupilIndex].incomingPrefs {

		if inPerf > pupilIndex && ec.pupils[inPerf].prefs[0] == pupilIndex &&
			pupilIndex < len(c) {
			//some pupil has this pupil has the highest index perf
			//collect all this pupils prefs assignments:
			allowedValues := []int{}
			for _, pref := range ec.pupils[inPerf].prefs {
				allowedValues = append(allowedValues, c[pref])
			}

			//if len(ec.pupils[inPerf].prefs) == 1 && inPerf > k {
			domainNotEmpty, changed := ec.domainValues.PushDomainRestriction(originalPupilIndex, inPerf, allowedValues)
			if !domainNotEmpty {
				ec.pupils[pupilIndex].unsatisfiedPrefsCount++
				return false
			}
			if changed {
				domainNotEmpty := fc(ec, c, inPerf, originalPupilIndex)
				if !domainNotEmpty {
					return false
				}
			}

		}

		if inPerf < pupilIndex {
			//this pupil has been assigned and now another one of its preferences is assigned,
			//if it is one before last, restrict the last!
			prefCount := len(ec.pupils[inPerf].prefs)

			_, inRangePrefCount, foundPrefCount = prefStatus(ec.pupils[inPerf], c, inPerf)
			if prefCount > 0 && foundPrefCount == 0 &&
				prefCount-inRangePrefCount == 1 {
				value, multiple := fcValues(ec, c, inPerf)
				if multiple {
					return true
				}
				domainNotEmpty, changed := ec.domainValues.PushDomainOnlyOneRestriction(originalPupilIndex, ec.pupils[inPerf].prefs[0], value)
				if !domainNotEmpty {

					ec.pupils[inPerf].unsatisfiedPrefsCount++
					return false
				}
				if changed {
					domainNotEmpty := fc(ec, c, ec.pupils[inPerf].prefs[0], originalPupilIndex)
					if !domainNotEmpty {
						return false
					}
				}

			}
		}

	}

	return true
}

func next(ec *ExecutionContext, c []int) []int {
	k := len(c) - 1

	nextVal := ec.domainValues.NextValue(k, c[k])
	if nextVal == -1 {
		return nil
	}
	c[k] = nextVal

	return c
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
