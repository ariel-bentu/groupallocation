package main

import (
	"fmt"
)

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
	if p.locked && p.initialGroup != c[k-1] {
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
	for _, constraint := range ec.Constraints {
		csg, ok := constraint.(*SubGroupConstraint)

		if ok && !csg.ValidateNew(ec, c) {
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
	return append(c, ec.pupils[k].group)
}

func next(ec *ExecutionContext, s []int) []int {
	k := len(s) - 1
	s[k]++
	if s[k] == ec.groupsCount {
		s[k] = 0
	}
	if s[k] == ec.pupils[k].group {
		return nil
	}

	return s
}

func Process(ec *ExecutionContext) {

	for ec.Next() {
		satisfied := true
		for _, c := range ec.Constraints {
			sat1 := c.Validate(ec)
			satisfied = satisfied && sat1

			debugInfo(ec, c, sat1)

		}

		//debug
		if ec.currentIteration == 10000 {
			for i, p := range ec.pupils {
				fmt.Printf("Pupil: %d, name:%s, Score:%d, group:%d\n", i, p.name, p.score, p.group)
			}
		}

		if satisfied {
			break
		}
	}

	ec.Finish()
}

/*
* DEBUG ************************************
 */
/*
func slice2String(arr []int) string {
	return strings.Trim(strings.Replace(fmt.Sprint(arr), " ", ",", -1), "[]")
}
func array2String(arr [MAX_NUM_OF_GROUPS]int, n int) string {
	ret := ""
	for i := 0; i < n; i++ {
		ret += fmt.Sprintf(", %d", arr[i])
	}
	return ret
}
*/
func debugInfo(ec *ExecutionContext, c Constraint, satisfied bool) {
	if ec.currentIteration == 10000 {
		ba := "no"
		ga := "no"
		yesNo := "no"
		if satisfied {
			yesNo = "yes"
		}
		tpe := "pref"
		members := ""
		dist := ""
		sg, ok := c.(*SubGroupConstraint)
		if ok {
			if sg.IsUnite {
				tpe = "UniteGroup"
			} else {
				tpe = "SeprateGroup"
			}
			if sg.boyAlone {
				ba = "yes"
			}
			if sg.girlAlone {
				ga = "yes"
			}
			members = slice2String(sg.Members())
			dist = array2String(sg.countForGroup, ec.groupsCount)
			fmt.Printf("id: %d, Name:%s, Type:%s, members: %s, satisfied:%s, boy/girl-alone:%s/%s, dist:%s\n", c.ID(), c.Description(), tpe, members, yesNo, ba, ga, dist)
		} else {
			pref, ok := c.(*PrefConstraint)
			if ok {
				fmt.Printf("id: %d, Pupil:%d, Type:Pref %d, RefPupil: %d, satisfied:%s %d,%d\n", pref.ID(), pref.pupil, pref.prefPrio, pref.refPupil, yesNo, ec.pupils[pref.pupil].group, ec.pupils[pref.refPupil].group)
			}
		}
	}
}
