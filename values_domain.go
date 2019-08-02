package main

import (
	"fmt"
	"sync"
)

type ValuesDomain struct {
	values        []int
	allowedValues []int
	next          *ValuesDomain
	origin        int
	Restrictions  []int
}

type ValuesDomainMain struct {
	domainLength       int
	valuesDomainArray  []ValuesDomain
	onlyOneDomains     [][]int
	butOneDomains      [][]int
	valuesDomainPool   sync.Pool
	GroupsRestrictions []*GroupRestriction
}

type GroupRestriction struct {
	oneOnly bool
	value   int
	fromInx int
	gender  int
	next    *GroupRestriction
}

func NewValuesDomainMain(domainLength int, variableCount int, groupCount int) *ValuesDomainMain {

	vdm := ValuesDomainMain{domainLength: domainLength, valuesDomainArray: make([]ValuesDomain, variableCount)}
	for i, _ := range vdm.valuesDomainArray {
		vdm.valuesDomainArray[i].values = make([]int, domainLength)
		vdm.valuesDomainArray[i].InitValuesDomain()
	}

	vdm.GroupsRestrictions = make([]*GroupRestriction, groupCount)
	vdm.onlyOneDomains = make([][]int, domainLength)

	for i, _ := range vdm.onlyOneDomains {
		vdm.onlyOneDomains[i] = append(vdm.onlyOneDomains[i], i)
	}

	vdm.butOneDomains = make([][]int, domainLength)
	for i, _ := range vdm.butOneDomains {
		for j := 0; j < domainLength; j++ {
			if j != i {
				vdm.butOneDomains[i] = append(vdm.butOneDomains[i], j)
			}
		}
	}

	vdm.valuesDomainPool.New = vdm.NewValueDomain
	return &vdm
}

func (vd *ValuesDomainMain) NewValueDomain() interface{} {
	v := new(ValuesDomain)
	v.allowedValues = make([]int, vd.domainLength)
	v.values = make([]int, vd.domainLength)
	return v
}

func (vd *ValuesDomainMain) getPooledValuesDomain() *ValuesDomain {
	return vd.valuesDomainPool.Get().(*ValuesDomain)
}

func (vd *ValuesDomainMain) putPooledValuesDomain(v *ValuesDomain) {
	vd.valuesDomainPool.Put(v)
}

var debug int

func (vd *ValuesDomainMain) PushDomainRestriction(origin int, target int, values []int) (bool, bool) {
	if DebugVerbose {
		fmt.Printf("Push: %d, %d\n", origin, target)
	}
	
	v := vd.getPooledValuesDomain()
	

	v.origin = origin
	v.allowedValues = values
	effectiveValuesVd := vd.getEffectiveDomainValues(target)
	beforeCount := len(effectiveValuesVd.values)
	v.intersectWith(effectiveValuesVd.values)
	currNext := vd.valuesDomainArray[target].next
	v.next = currNext
	vd.valuesDomainArray[target].next = v

	vd.valuesDomainArray[origin].Restrictions = append(vd.valuesDomainArray[origin].Restrictions, target)
	return len(v.values) > 0, beforeCount > len(v.values)
}

func (gr *GroupRestriction) exists(pupilInx int, value int, onlyOne bool, gender int) bool {
	next := gr
	for next != nil {
		if next.fromInx <= pupilInx && next.gender == gender && next.value == value && next.oneOnly == onlyOne {
			return true
		}
		next = next.next
	}
	return false
}

//gender: 0 ignore, 1 boys, 2 girls
//returned array of changed
func (vd *ValuesDomainMain) pushGroupRestriction(ec *ExecutionContext, gInx int, origin int, target int, value int, onlyOne bool, gender int) (bool, []int) {
	gr := ec.domainValues.GroupsRestrictions[gInx]
	changedArray := []int{}
	if gr == nil || !gr.exists(target, value, onlyOne, gender) {

		g := ec.Constraints[gInx]
		for j := 0; j < len(g.members); j++ {
			memberInx := g.members[j]
			if memberInx <= target {
				break
			}
			//add restriction
			if gender == 1 && !ec.pupils[memberInx].IsMale() ||
				gender == 2 && ec.pupils[memberInx].IsMale() {
				continue
			}

			var domainNotEmpty, changed bool
			if onlyOne {
				domainNotEmpty, changed = ec.domainValues.PushDomainOnlyOneRestriction(origin, memberInx, value)
			} else {
				domainNotEmpty, changed = ec.domainValues.PushDomainButOneRestriction(origin, memberInx, value)
			}
			if changed {
				changedArray = append(changedArray, memberInx)
			}
			if !domainNotEmpty {
				ec.pupils[memberInx].unsatisfiedGroupsCount[gInx]++
				g.unsatisfiedCount++
				return false, changedArray
			}
		}
		vd.GroupsRestrictions[gInx] = &GroupRestriction{
			oneOnly: onlyOne,
			value:   value,
			fromInx: target,
			gender:  gender,
			next:    vd.GroupsRestrictions[gInx],
		}

		vd.valuesDomainArray[origin].Restrictions = append(vd.valuesDomainArray[origin].Restrictions, -1-gInx)
		//onlyOneStr := "F"
		//if onlyOne {
		//	onlyOneStr = "T"
		//}

		//fmt.Printf("Push o=%d, g=%d, 1:%s, val=%d", k, gInx, onlyOneStr, value)
	}

	return true, changedArray
}

func (vd *ValuesDomainMain) PopAllDomainRestriction(origin int) {

	if DebugVerbose {
		fmt.Printf("Pop: %d\n", origin)
	}

	valDomain := &vd.valuesDomainArray[origin]
	for _, v := range valDomain.Restrictions {
		if v >= 0 {
			valDomainRest := &vd.valuesDomainArray[v]
			if valDomainRest.next != nil {
				//pop the retriction
				old := valDomainRest.next
				valDomainRest.next = valDomainRest.next.next

				vd.putPooledValuesDomain(old)
			} else {
				print("problem")
			}
		} else {
			//group and is one based
			inx := -(v + 1)
			//fmt.Printf("Pop o=%d, g=%d, 1:%s", origin, inx, "?")
			if vd.GroupsRestrictions[inx] != nil {
				vd.GroupsRestrictions[inx] = vd.GroupsRestrictions[inx].next
			} else {
				print("problem2")
			}

		}
	}

	//reset
	valDomain.Restrictions = valDomain.Restrictions[:0]

}

func (vd *ValuesDomainMain) getEffectiveDomainValues(i int) *ValuesDomain {
	valDomain := &vd.valuesDomainArray[i]
	if valDomain.next != nil {
		return valDomain.next
	}
	return valDomain
}

func (vd *ValuesDomain) intersectWith(values []int) {
	newVal := vd.values[:0]
	for _, val := range values {
		for _, val2 := range vd.allowedValues {
			if val2 == val {
				newVal = append(newVal, val)
				break
			}
		}
	}
	vd.values = newVal
}

func (v *ValuesDomain) InitValuesDomain() {
	newVal := v.values[:0]
	for i := 0; i < cap(v.values); i++ {
		newVal = append(newVal, i)
	}
	v.values = newVal
}

func (v *ValuesDomainMain) FirstValue(inx, preferedVal int) int {
	vd := v.getEffectiveDomainValues(inx)
	for i, val := range vd.values {
		if val == preferedVal {
			vd.values[0], vd.values[i] = vd.values[i], vd.values[0]
			break
		}
	}
	if len(vd.values) > 0 {
		return vd.values[0]
	}
	return -1
}

func (v *ValuesDomainMain) NextValue(inx int, currVal int) int {
	vd := v.getEffectiveDomainValues(inx)
	found := false
	for _, val := range vd.values {
		if found {
			return val
		}
		if val == currVal {
			found = true
		}
	}
	//not found
	return -1
}

func (v *ValuesDomainMain) AreOptionsLeft(inx int, currInx int) bool {
	return v.NextValue(inx, currInx) != -1
}

func (vd *ValuesDomainMain) PushDomainButOneRestriction(origin int, target int, forbidVal int) (bool, bool) {
	return vd.PushDomainRestriction(origin, target, vd.butOneDomains[forbidVal])
}

func (vd *ValuesDomainMain) PushDomainOnlyOneRestriction(origin int, target int, allowedVal int) (bool, bool) {
	return vd.PushDomainRestriction(origin, target, vd.onlyOneDomains[allowedVal])
}
