package main

import (
	"fmt"
	"math"
	"strings"
)

const MAX_NUM_OF_GROUPS = 20

type SubGroupConstraint struct {
	id                          int
	desc                        string
	IsUnite                     bool
	Level                       int
	weight                      int
	Members                     []int
	maxAllowed                  float64
	stillAllowOneToBeOneTooMany bool

	countForGroup [MAX_NUM_OF_GROUPS]int

	oversizedForGroup [MAX_NUM_OF_GROUPS]bool
	boysForGroup      [MAX_NUM_OF_GROUPS]int
	girlsForGroup     [MAX_NUM_OF_GROUPS]int
	boysCount         int
	minBoys           int
	minGirls          int
	girlAlone         bool
	boyAlone          bool

	satisfied bool
}

func NewSubGroupConstraint(id int, desc string, isUnite bool, weight int) Constraint {
	return &SubGroupConstraint{id: id, desc: desc, IsUnite: isUnite, weight: weight, Level: 1}
}

func (c *SubGroupConstraint) AfterInit(ec *ExecutionContext) {
	if !c.IsUnite {
		c.maxAllowed = float64(len(c.Members)) / float64(ec.groupsCount)
		//c.stillAllowOneToBeOneTooMany = c.maxAllowed > math.Floor(c.maxAllowed)

		ratio := c.boysCount / ec.groupsCount
		if ratio >= 3 {
			c.minBoys = 3
		} else if ratio >= 2 {
			c.minBoys = 2
		}

		/*		if c.boysCount > 7 {
					c.minBoys = 3
				} else if c.boysCount > 4 {
					c.minBoys = 2
				}
		*/
		girlsNum := len(c.Members) - c.boysCount
		ratio = girlsNum / ec.groupsCount
		if ratio >= 3 {
			c.minGirls = 3
		} else if ratio >= 2 {
			c.minGirls = 2
		}
		//c.minBoys = 0
		//c.minGirls = 0
	}
}

func (sgc *SubGroupConstraint) ID() int {
	return sgc.id
}
func (sgc *SubGroupConstraint) Description() string {
	return sgc.desc
}

func (sgc *SubGroupConstraint) Weight() int {
	return sgc.weight
}
func (c *SubGroupConstraint) Validate(ec *ExecutionContext) bool {
	if c.Members == nil || c.Level == 0 {
		return true
	}

	for i := 0; i < ec.groupsCount; i++ {
		c.countForGroup[i] = 0
		c.oversizedForGroup[i] = false
		c.boysForGroup[i] = 0
		c.girlsForGroup[i] = 0
		c.boysCount = 0

	}
	c.boyAlone = false
	c.girlAlone = false
	oversized := false

	for i := 0; i < len(c.Members); i++ {
		p := ec.pupils[c.Members[i]]
		if p.IsMale() {
			c.boysCount++
			c.boysForGroup[p.group]++
		} else {
			c.girlsForGroup[p.group]++
		}
		c.countForGroup[p.group]++
	}
	breakHere := false
	if c.IsUnite {
		largestCountIndex := c.findSmallOrLargeGroup(false, ec.groupsCount)
		for i := 0; i < ec.groupsCount; i++ {
			c.oversizedForGroup[i] = (i != largestCountIndex && c.countForGroup[i] > 0)
			oversized = oversized || c.oversizedForGroup[i]
		}
	} else {
		if c.countForGroup[0] == c.countForGroup[1] {
			if !breakHere {
				breakHere = true
			}
		}

		maxAllowed := float64(len(c.Members)) / float64(ec.groupsCount)
		stillAllowOneToBeOneTooMany := maxAllowed > math.Floor(maxAllowed)

		for i := 0; i < ec.groupsCount; i++ {
			diff := maxAllowed - float64(c.countForGroup[i])
			if diff <= 0 && diff > -0.5 {
				c.oversizedForGroup[i] = false
			} else if diff <= 0 && diff > -1 && stillAllowOneToBeOneTooMany {
				c.oversizedForGroup[i] = false
				stillAllowOneToBeOneTooMany = false
			} else {
				c.oversizedForGroup[i] = diff < 0
				oversized = oversized || c.oversizedForGroup[i]
			}

			if c.boysForGroup[i] < c.minBoys {
				c.boyAlone = true
			}
			if girls := len(c.Members) - c.boysCount; girls < c.minGirls {
				c.girlAlone = true
			}
		}

	}
	//updates the members' score
	smallestCountIndex := 0
	if oversized {
		if !c.IsUnite {
			smallestCountIndex = c.findSmallOrLargeGroup(true, ec.groupsCount)
		}

		for i := 0; i < ec.groupsCount; i++ {
			if c.oversizedForGroup[i] {
				if c.IsUnite {
					for j := 0; j < len(c.Members); j++ {
						p := ec.pupils[c.Members[j]]
						p.score -= c.weight
					}
				} else {
					if i != smallestCountIndex {
						//only deduct score to the amount of pupils equal to the diff between this group and the smallest
						for j := 0; j < c.countForGroup[i]-c.countForGroup[smallestCountIndex]; j++ {
							p := ec.pupils[c.Members[j]]
							p.score -= c.weight
						}
					}
				}
			}
		}
	}
	c.satisfied = !c.boyAlone && !c.girlAlone && !oversized
	return c.satisfied
}

func (c *SubGroupConstraint) ValidateNew(ec *ExecutionContext, candidate []int) bool {
	if c.Members == nil || c.Level == 0 {
		return true
	}

	for i := 0; i < ec.groupsCount; i++ {
		c.countForGroup[i] = 0
		c.oversizedForGroup[i] = false
		c.boysForGroup[i] = 0
		c.girlsForGroup[i] = 0

	}
	c.boyAlone = false
	c.girlAlone = false
	k := len(candidate)
	left := 0
	for i := 0; i < len(c.Members); i++ {
		if c.Members[i] < k {
			p := ec.pupils[c.Members[i]]
			group := candidate[c.Members[i]]

			if p.IsMale() {
				c.boysForGroup[group]++
			} else {
				c.girlsForGroup[group]++
			}
			c.countForGroup[group]++
		} else {
			left++
		}
	}

	count := 0
	if c.IsUnite {
		for i := 0; i < ec.groupsCount; i++ {
			if c.countForGroup[i] > 0 {
				if count > 0 {
					return false
				}
				count++
			}
		}
		return true
	}

	//stillAllowOneToBeOneTooMany := c.stillAllowOneToBeOneTooMany
	for i := 0; i < ec.groupsCount; i++ {
		if float64(c.countForGroup[i]) > math.Ceil(c.maxAllowed) {
			return false
		}

		diff := c.maxAllowed - float64(c.countForGroup[i])
		if diff <= -1 { //diff < -0.5 {
			return false
		}

		if diff > 0 && float64(left) < math.Floor(diff) { // diff - 0.5 {
			return false
		}

		//todo many groups
		if c.boysForGroup[i]+left < c.minBoys {
			return false
		}

		if girls := len(c.Members) - c.boysCount; girls+left < c.minGirls {
			return false
		}

	}
	return true
}

func (c *SubGroupConstraint) Message(ec *ExecutionContext) string {
	if c.satisfied {
		return ""
	}

	sb := NewStringBuffer()
	sb.Clear()
	sb.Append(c.Description())
	if c.IsUnite {
		sb.Append(" - לא מאוחדת. ")
	} else {
		sb.Append(" - לא מפוזרת אחיד. ")
	}
	for i := 0; i < ec.groupsCount; i++ {
		sb.AppendFormat("%d : %d ; ", i+1, c.countForGroup[i])
	}
	if c.boyAlone {
		sb.Append(" בן לבד.")
	}
	if c.girlAlone {
		sb.Append(" בת לבד.")
	}

	return sb.ToString()
}

func (c *SubGroupConstraint) AddMember(pupilInx int, ec *ExecutionContext) {
	p := ec.pupils[pupilInx]

	c.Members = append(c.Members, pupilInx)
	if p.IsMale() {
		c.boysCount++
	}
}

func (c *SubGroupConstraint) findSmallOrLargeGroup(isSmallest bool, groupsCount int) int {
	count := 0
	if isSmallest {
		count = 1000000
	}
	countIndex := -1
	for i := 0; i < groupsCount; i++ {
		if (isSmallest && c.countForGroup[i] < count) || (!isSmallest && c.countForGroup[i] > count) {
			count = c.countForGroup[i]
			countIndex = i
		}
	}
	return countIndex
}

func (sg *SubGroupConstraint) printOneInfo(ec *ExecutionContext) string {
	tpe := ""
	if sg.IsUnite {
		tpe = "UniteGroup"
	} else {
		tpe = "SeprateGroup"
	}
	ba, ga := "", ""
	if sg.boyAlone {
		ba = "yes"
	}
	if sg.girlAlone {
		ga = "yes"
	}
	yesNo := "no"
	if sg.satisfied {
		yesNo = "yes"
	}

	//members = slice2String(sg.Members)
	dist := array2String(sg.countForGroup, ec.groupsCount)
	return fmt.Sprintf("id: %d, Name:%s, Type:%s, satisfied:%s, boy/girl-alone:%s/%s, dist:%s<br>\n", sg.ID(), sg.Description(), tpe, yesNo, ba, ga, dist)

}

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
