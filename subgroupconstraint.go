package main

import (
	"fmt"
	"math"
	"strings"
)

const MAX_NUM_OF_GROUPS = 20

type SubGroupConstraint struct {
	id              int
	desc            string
	IsUnite         bool
	Level           int
	weight          int
	members         []int
	originalCount   int
	maxAllowed      float64
	minAllowed      float64
	genderSensitive bool
	speadToAll      bool

	countForGroup []int

	oversizedForGroup []bool
	boysForGroup      []int
	girlsForGroup     []int
	boysCount         int
	minBoys           int
	minGirls          int
	girlAlone         bool
	boyAlone          bool
	disallowZeroBoys  bool
	disallowZeroGirls bool

	satisfied        bool
	unsatisfiedCount int64
	disabled         bool
}

func NewSubGroupConstraint(id int, desc string, isUnite bool, weight int, groupsCount int) *SubGroupConstraint {
	g := SubGroupConstraint{id: id, desc: desc, IsUnite: isUnite, weight: weight, Level: 1}
	g.countForGroup = make([]int, groupsCount)
	g.boysForGroup = make([]int, groupsCount)
	g.girlsForGroup = make([]int, groupsCount)
	return &g
}

func (c *SubGroupConstraint) AfterInit(ec *ExecutionContext, err *stringBuffer) {
	if !c.IsUnite {
		howManyGroupsToSpread := float64(ec.groupsCount)
		c.minAllowed = 0
		if !c.speadToAll {
			howManyGroupsToSpread = 2
		}
		c.maxAllowed = float64(len(c.members)) / howManyGroupsToSpread
		if c.speadToAll {
			c.minAllowed = c.maxAllowed - 1
		}

		//c.stillAllowOneToBeOneTooMany = c.maxAllowed > math.Floor(c.maxAllowed)
		if c.boysCount/ec.groupsCount >= 2 && c.genderSensitive {
			c.minBoys = 2
			//} else if c.boysCount/ec.groupsCount > 2 && c.genderSensitive {
			//	c.minBoys = 2
		} else {
			c.minBoys = 0
		}
		/*		if c.boysCount > 7 {
					c.minBoys = 3
				} else if c.boysCount > 4 {
					c.minBoys = 2
				}
		*/
		girlsNum := len(c.members) - c.boysCount
		if girlsNum/ec.groupsCount >= 2 && c.genderSensitive {
			c.minGirls = 2
			//} else if girlsNum/ec.groupsCount > 2 && c.genderSensitive {
			//	c.minGirls = 2
		} else {
			c.minGirls = 0
		}

		if ec.groupsCount == 2 {
			if c.maxAllowed-float64(c.boysCount+c.minGirls) < 0 {
				c.disallowZeroBoys = true
			}
			if c.maxAllowed-float64(len(c.members)-c.boysCount+c.minBoys) < 0 {
				c.disallowZeroGirls = true
			}

		}

		//c.minBoys = 0
		//c.minGirls = 0
	} else {
		c.originalCount = len(c.members)

		//add to members pupils that all their perfs are in this unite group
		for pupilInx, p := range ec.pupils {
			count := 0
			if len(p.prefs) > 0 && !c.IsMember(pupilInx) {
				for i := 0; i < len(p.prefs); i++ {
					if c.IsMember(p.prefs[i]) {
						count++
					} else {
						break
					}
				}
				if count == len(p.prefs) {
					//add this pupil to c group
					c.AddMember(pupilInx, ec)
					//todo add message
					err.AppendFormat("Adding '%s' to group '%s' as all its preferences in that unite group",
						ec.pupils[pupilInx].name, c.desc)
				}
			}
		}
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

/*
func (c *SubGroupConstraint) Validate(ec *ExecutionContext) bool {
	if c.members == nil || c.Level == 0 {
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

	for i := 0; i < len(c.members); i++ {
		p := ec.pupils[c.members[i]]
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

		maxAllowed := float64(len(c.members)) / float64(ec.groupsCount)
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
			if girls := len(c.members) - c.boysCount; girls < c.minGirls {
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
					for j := 0; j < len(c.members); j++ {
						p := ec.pupils[c.members[j]]
						p.score -= c.weight
					}
				} else {
					if i != smallestCountIndex {
						//only deduct score to the amount of pupils equal to the diff between this group and the smallest
						for j := 0; j < c.countForGroup[i]-c.countForGroup[smallestCountIndex]; j++ {
							p := ec.pupils[c.members[j]]
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
*/
func (c *SubGroupConstraint) checkOnlyGroup(ec *ExecutionContext, grp int) bool {
	for i := 0; i < ec.groupsCount; i++ {
		if i != grp && c.countForGroup[i] > 0 {
			return false
		}
	}
	return true
}

func (c *SubGroupConstraint) checkNoGroup(ec *ExecutionContext, grp int) bool {
	return c.countForGroup[grp] == 0
}

func (c *SubGroupConstraint) ValidateNew(ec *ExecutionContext, candidate Candidate) bool {
	if len(c.members) == 0 || c.disabled {
		c.satisfied = true
		return true
	}
	k := candidate.Count()
	maxAllowed := 34.0 //todo from config
	if c == ec.allGroup {

	} else {
		maxAllowed = (1 + float64(ec.graceLevel)/10) * c.maxAllowed
	}
	minAllowed := (1 - float64(ec.graceLevel)/10) * c.minAllowed

	if k < 2 || !c.IsUnite && k <= int(math.Ceil(maxAllowed)) {
		c.satisfied = true
		return true
	}

	boysLeft, girlsLeft := c.calculateMembersCounts(ec, candidate)
	/*
		if c.desc == "1" {
			c.satisfied = c.checkOnlyGroup(ec, 0)
			if !c.satisfied {
				c.unsatisfiedCount++
			}
			return c.satisfied
		}
		if c.desc == "2" {
			c.satisfied = c.checkOnlyGroup(ec, 1)
			if !c.satisfied {
				c.unsatisfiedCount++
			}
			return c.satisfied
		}
		if c.desc == "not3" && ec.groupsCount > 2 {
			c.satisfied = c.checkNoGroup(ec, 2)
			if !c.satisfied {
				c.unsatisfiedCount++
			}
			return c.satisfied
		}
	*/
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
		c.satisfied = true
		return true
	}

	for i := 0; i < ec.groupsCount; i++ {
		if float64(c.countForGroup[i]) > math.Ceil(maxAllowed) {
			return false
		}

		//diff := c.maxAllowed - float64(c.countForGroup[i])
		//if diff <= -1 { //diff < -0.5 {
		//more pupils in one group than allowed
		//	return false
		//}

		//if diff > 0 && float64(left) < math.Floor(diff) { // diff - 0.5 {
		if math.Floor(minAllowed)-float64(boysLeft+girlsLeft)-float64(c.countForGroup[i]) > 0 {
			return false
		}

		if (c.boysForGroup[i] > 0 || c.disallowZeroBoys) && c.boysForGroup[i]+boysLeft < c.minBoys {
			return false
		}

		if girls := c.countForGroup[i] - c.boysForGroup[i]; (girls > 0 || c.disallowZeroGirls) && girls+girlsLeft < c.minGirls {
			return false
		}

	}
	c.satisfied = true
	return true
}

func (c *SubGroupConstraint) calculateMembersCounts(ec *ExecutionContext, candidate Candidate) (int, int) {
	k := candidate.Count()
	for i := 0; i < ec.groupsCount; i++ {
		c.countForGroup[i] = 0
		c.boysForGroup[i] = 0
		c.girlsForGroup[i] = 0
	}

	boysLeft := 0
	girlsLeft := 0
	for i := 0; i < len(c.members); i++ {
		p := ec.pupils[c.members[i]]
		if c.members[i] < k {
			group := candidate.GetGroup(c.members[i])

			if p.IsMale() {
				c.boysForGroup[group]++
			} else {
				c.girlsForGroup[group]++
			}
			c.countForGroup[group]++
		} else {
			if p.IsMale() {
				boysLeft++
			} else {
				girlsLeft++
			}
		}
	}
	return boysLeft, girlsLeft
}

func (c *SubGroupConstraint) IsMember(pupil int) bool {
	for _, m := range c.members {
		if m == pupil {
			return true
		}
	}
	return false
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

	c.members = append(c.members, pupilInx)

	//sort highest first
	for i := len(c.members) - 2; i >= 0; i-- {
		if c.members[i+1] > c.members[i] {
			temp := c.members[i+1]
			c.members[i+1] = c.members[i]
			c.members[i] = temp
		} else {
			break
		}
	}

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
func array2String(arr []int, n int) string {
	ret := ""
	for i := 0; i < n; i++ {
		ret += fmt.Sprintf(", %d", arr[i])
	}
	return ret
}
