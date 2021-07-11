package main

import (
	"math"
	"math/rand"
)

/**
1. Get an initial solution S.
2. Get an initial temperature T > 0.
3. While not yet frozen do the following.
  3.1 Perform the following loop L times.
    3.1.1 Pick a random neighbor S' of S.
    3.1.2 Let Delta = cost (S')- cost (S).
    3.1.3 If Delta <= 0 (downhill move),
        Set S = S'.
    3.1.4 If Delta > 0 (uphill move),
        Set S = S' with probability e power -(Delta/T).
  3.2 Set T = rT (reduce temperature).
4. Return S
**/
func runSA(ec *ExecutionContext) {

	S := make([]int, ec.activePupilsCount)
	newS := make([]int, ec.activePupilsCount)
	getInitial(ec, &S)
	T := 1000
	L := 50000
	for T > 0 {
		costS := cost(ec, &S)
		for i := 0; i < L; i++ {
			neighbor(ec, &S, &newS)
			costNewS := cost(ec, &newS)
			delta := costNewS - costS
			if delta <= 0 ||
				probability(delta, T) {
				S = newS
				costS = costNewS

			}
			ec.currentIteration = i
		}
		T--
		ec.currentLevelCount = T
	}
	ec.bestCandidate = make([]int, len(S))
	copy(ec.bestCandidate, S)

	ec.Finish(false)
}

func getInitial(ec *ExecutionContext, S *[]int) {
	for i := 0; i < ec.activePupilsCount; i++ {
		startValue := randomDomainValue(ec)
		(*S)[i] = startValue
	}

}
func getInitialTemp(ec *ExecutionContext) int {
	return 100
}

func randomDomainValue(ec *ExecutionContext) int {
	return rand.Intn(ec.groupsCount)
}

func probability(delta int64, T int) bool {
	prob := 1 / math.Log(float64(delta)/float64(T))
	rnd := rand.Float64()

	return rnd < prob
}

func neighbor(ec *ExecutionContext, S *[]int, newS *[]int) {
	copy(*newS, *S)
	for {
		n1 := rand.Intn(len(*S))
		newVal := randomDomainValue(ec)
		//n2 := rand.Intn(len(*S))
		//(*newS)[n1], (*newS)[n2] = (*S)[n2], (*S)[n1]
		if (*S)[n1] != newVal {
			(*newS)[n1] = newVal
			return
		}
	}
}

func cost(ec *ExecutionContext, S *[]int) int64 {
	var cost int64
	groupCost := 500
	friendCostFactor := 50
	//go over all Constraints
	for _, g := range ec.Constraints {
		if !g.disabled && len(g.members) > 0 {
			g.calculateMembersCounts(ec, PartialAssignment(*S))
			pupilPrice := groupCost / len(g.members)

			if g.IsUnite {
				maxIndex := maxGroup(g.countForGroup)
				for i := 0; i < ec.groupsCount; i++ {
					if maxIndex != i {
						cost += int64(pupilPrice * g.countForGroup[i])
					}
				}
			} else {
				//seperation constraint
				for i := 0; i < ec.groupsCount; i++ {
					if float64(g.countForGroup[i]) > g.maxAllowed {
						cost += int64(float64(pupilPrice) * (float64(g.countForGroup[i]) - g.maxAllowed))
					}
					if float64(g.countForGroup[i]) < g.minAllowed {
						cost += int64(float64(pupilPrice) * (g.minAllowed - float64(g.countForGroup[i])))
					}
				}

			}

		}
	}

	for i := 0; i < ec.activePupilsCount; i++ {
		prefs := ec.pupils[i].origOrderPrefs
		amount := 0
		for j := 0; j < len(prefs); j++ {
			if prefs[j] < ec.activePupilsCount && (*S)[i] == (*S)[prefs[j]] {
				amount++
			}
		}
		cost += int64((len(prefs) - amount) * friendCostFactor)
	}

	return cost
}

func maxGroup(v []int) int {
	var maxIndex int
	var maxValue int
	for i, e := range v {
		if i == 0 || e > maxValue {
			maxValue = e
			maxIndex = i
		}
	}
	return maxIndex
}
