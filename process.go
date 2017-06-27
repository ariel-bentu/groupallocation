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
		ec.Cancel = false
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

	//prefs
	for i, v := range c {
		p := ec.pupils[i]

		if p.locked && p.initialGroup != v {
			return false
		}
		inRangeCount := 0
		refCount := 0

		for i = 0; i < len(p.prefs); i++ {
			if p.prefs[i] < k {
				inRangeCount++
				if c[p.prefs[i]] == v {
					refCount++
					break
				}
			}
		}

		if inRangeCount > 0 && inRangeCount == len(p.prefs) && refCount == 0 {
			return false
		}

	}

	//constraints
	for _, constraint := range ec.Constraints {
		csg, ok := constraint.(*SubGroupConstraint)

		if ok && !csg.ValidateNew(ec, c) {
			return false
		}
	}
	return true
}

func first(ec *ExecutionContext, c []int) []int {
	k := len(c)
	if k == len(ec.pupils) {
		return nil
	}
	return append(c, 0)
}

func next(ec *ExecutionContext, s []int) []int {
	k := len(s) - 1
	if s[k] == ec.groupsCount-1 {
		return nil
	}

	nextS := make([]int, len(s))
	copy(nextS, s)
	nextS[k]++

	return nextS
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
			members = slice2String(sg.Members)
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
