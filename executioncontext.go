package main

import (
	"fmt"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/tealeg/xlsx"
)

const NUM_OF_PREF = 3
const UNITE_VALUE = "איחוד"

const CELL_SUBGROUP = 7

const CELL_NAME = 2
const CELL_GENDER = 3
const CELL_PREF = 4
const CELL_INITIAL = 8
const CELL_INITIAL_WEIGHT = 9
const CELL_COMP_PROPOSAL = 10
const CELL_NUM_OF_MOVES = 11

type Candidate interface {
	Count() int
	GetGroup(pupil int) int
}

func (ec *ExecutionContext) Count() int {
	return ec.activePupilsCount
}

func (ec *ExecutionContext) GetGroup(pupil int) int {
	return ec.pupils[pupil].groupBestScore
}

//Constrain is a...
type ExecutionContext struct {
	id        string
	startTime time.Time
	endTime   time.Time
	done      bool
	result    string
	resultID  int
	Cancel    bool

	groupsCount       int
	iterationCount    int
	Constraints       []*SubGroupConstraint
	allGroup          *SubGroupConstraint
	maleGroup         *SubGroupConstraint
	pupils            []*Pupil
	activePupilsCount int
	dataExcel         *xlsx.File
	taskId            int
	file              string

	domainValues *ValuesDomainMain

	timeLimit int //in seconds

	MaxDepth            int
	currentLevelCount   int
	graceLevel          int
	sensitiveToOnlyLast int
	allowOnlyLastCount  int

	currentIteration             int
	prefFailCount                int
	prefThreashold               int64
	statusCandidate              []int
	bestSumOfSatisfiedPrefs      int
	bestSumOfSatisfiedFirstPrefs int
	bestCandidate                []int
	resultsCount                 int
	numOfUpgrades                int
	resultsScoreHistory          []int
	InitialErr                   string
}

var RunningExecutions map[string]*ExecutionContext

func NewExecutionContext() *ExecutionContext {
	ec := new(ExecutionContext)
	ec.startTime = time.Now()
	ec.done = false
	ec.Cancel = false
	ec.graceLevel = 0
	ec.sensitiveToOnlyLast = 0
	ec.allowOnlyLastCount = 3 //todo from conf

	t := time.Now().UnixNano()
	ec.id = fmt.Sprintf("%d", t)

	if RunningExecutions == nil {
		RunningExecutions = make(map[string]*ExecutionContext)
	}

	RunningExecutions[ec.id] = ec

	return ec
}

func FindExecutionContext(id string) *ExecutionContext {
	return RunningExecutions[id]
}

func (e *ExecutionContext) ID() string {
	return e.id
}

func (e *ExecutionContext) GetStatusHtml(clean bool) (string, string) {
	if e.done {
		return e.printHtml(e.resultID, clean), fmt.Sprintf("%f,resultID: %d, #of found results: %d</br>%s", e.endTime.Sub(e.startTime).Seconds(), e.resultID, e.resultsCount, slice2String(e.resultsScoreHistory))
	}
	//return fmt.Sprintf("Interaction Count:%d, progress: %d", e.currentIteration), fmt.Sprintf("%f", time.Now().Sub(e.startTime).Seconds(),
	//	e.currentIteration/e.iterationCount*100)

	sb := NewStringBuffer()
	sb.AppendFormat("Candidate Length:%d, Max: %d, elappsed: %f.0, iter: %s, graceLevel:%d, results so far: %d, upgrades: %d<br>\n<div dir=\"ltr\">%s</div></br>\n", e.currentIteration, e.MaxDepth, time.Now().Sub(e.startTime).Seconds(),
		FormatInt2String(e.currentLevelCount), e.graceLevel, e.resultsCount, e.numOfUpgrades, slice2String(e.statusCandidate))

	for _, c := range e.Constraints {

		sb.AppendFormat("%s - %d</br>\n", c.Description(), c.unsatisfiedCount)
	}
	sb.Append("\n</br><table border='1'><tr><th>#</th><th>name</th><th>Pref</th><th>Groups</th><th>Failed Preferences</th><th>Failed Groups<th></th></tr>\n")
	for i, p := range e.pupils {
		bgcolor := ""
		if p.unsatisfiedPrefsCount > 0 {
			bgcolor = `bgcolor="#FF0000"`
		}

		sb.AppendFormat("<tr><td>%d</td><td>%s</td><td %s>%s</td><td>%s</td><td>%s</td><td>%s</td><tr>\n",
			i, p.name, bgcolor, FormatInt642String(p.unsatisfiedPrefsCount), getPrefHtml(e, p), getGroupsHtml(e, p), getFailedGroupHTMLMessage(e, p.unsatisfiedGroupsCount))
	}
	sb.Append("</table><br/><br/>")

	sb.Append(strings.Replace(e.InitialErr, "\n", "<br/>\n", -1))

	return sb.ToString(), ""
}

func FormatInt2String(v int) string {
	str := strconv.Itoa(v)
	return FormatHuman(str)
}
func FormatInt642String(v int64) string {
	if v > 999 {
		stop()
	}
	str := strconv.FormatInt(v, 10)
	return FormatHuman(str)
}
func FormatHuman(str string) string {
	ret := ""
	counter := 0
	for i := len(str) - 1; i >= 0; i-- {
		counter++
		ret = str[i:i+1] + ret
		if counter == 3 && i > 0 {
			ret = "," + ret
			counter = 0
		}

	}
	return ret
}

func getGroupsHtml(e *ExecutionContext, p *Pupil) string {
	sb := NewStringBuffer()
	for i, v := range p.uniteGroups {
		sb.AppendFormat("%d: U - %s (B: %d, G:%d)<br/>", i+1, e.Constraints[v].desc, e.Constraints[v].boysCount, len(e.Constraints[v].members)-e.Constraints[v].boysCount)
	}

	for i, v := range p.seperationGroups {
		sb.AppendFormat("%d: S - %s (B: %d, G:%d)<br/>", i+1, e.Constraints[v].desc, e.Constraints[v].boysCount, len(e.Constraints[v].members)-e.Constraints[v].boysCount)
	}

	return sb.ToString()

}
func getPrefHtml(e *ExecutionContext, p *Pupil) string {
	sb := NewStringBuffer()
	for i, v := range p.origOrderPrefs {
		sb.AppendFormat("%d: %s (%d)<br/>", i+1, e.pupils[v].name, v)
	}
	return sb.ToString()

}
func getFailedGroupHTMLMessage(e *ExecutionContext, failedArray []int64) string {
	sb := NewStringBuffer()
	for i, v := range failedArray {
		if v > 0 {
			sb.AppendFormat("%s: %s<br/>", e.Constraints[i].desc, FormatInt642String(v))
		}
	}
	return sb.ToString()
}

func (e *ExecutionContext) Finish(persist bool) {
	e.done = true
	for i, p := range e.pupils {
		if i >= e.activePupilsCount {
			break
		}
		p.groupBestScore = e.bestCandidate[i]
	}
	e.endTime = time.Now()

	if persist {
		connect()
		r := db.QueryRow("select max(resultId) from taskResult")
		maxId := 1
		if r != nil {
			r.Scan(&maxId)
		}
		maxId++
		e.resultID = maxId

		stmt, _ := db.Prepare("INSERT INTO taskResult (resultId, tenant, task, runDate, duration, foundCount) values (?,?,?,?,?,?)")
		stmt.Exec(maxId, 0, e.taskId, e.endTime.Unix(), (e.endTime.Unix() - e.startTime.Unix()), e.resultsCount)

		stmt, _ = db.Prepare("INSERT INTO taskResultLines (resultId, groupId, pupilId) values (?,?,?)")

		for _, p := range e.pupils {
			stmt.Exec(maxId, p.groupBestScore, p.id)
		}
	}
}

func (e *ExecutionContext) ReadResults(id int) {
	connect()
	result, err := db.Query("select pupilId, groupId from taskResultLines B inner join taskResult A on (A.resultId = B.resultId) where task = ? and A.resultId = ?", e.taskId, id)
	if err != nil {
		panic(err)
	}
	var pupilId, groupId int
	for result.Next() {
		result.Scan(&pupilId, &groupId)
		pId := e.findPupilById(pupilId)
		e.pupils[pId].groupBestScore = groupId
	}
	e.resultID = id
}

func IsEmpty(c *xlsx.Cell) bool {
	var v, err = c.FormattedValue()
	return err != nil || v == ""
}
func MsgBox(text string) {
	log.Output(0, text)
}

func (ec *ExecutionContext) GetParam(name string) string {
	sheet := ec.getSheet("Configuration")
	for _, row := range sheet.Rows {
		if v, _ := row.Cells[0].FormattedValue(); v == name {
			p, _ := row.Cells[1].FormattedValue()
			return p
		}
	}
	return ""
}

func (ec *ExecutionContext) GetIntParam(name string) int {
	var p = ec.GetParam(name)
	var v, err = strconv.Atoi(p)
	if err == nil {
		return v
	}
	return 0

}

/*
func Initialize(file string) (*ExecutionContext, string) {
	ec := NewExecutionContext()
	ec.file = file
	ec.groupsCount = ec.GetIntParam("מספר כיתות")

	groupsSheet := ec.getSheet("Groups")
	err := NewStringBuffer()
	// initialize all sub groups from the "Groups" sheet
	i := 1
	for ; !IsEmpty(groupsSheet.Cell(i, 2)); i++ {
		id, _ := groupsSheet.Cell(i, 0).Int()
		v, _ := groupsSheet.Cell(i, 1).FormattedValue()
		isUnite := (v == UNITE_VALUE)
		desc, _ := groupsSheet.Cell(i, 2).FormattedValue()
		g := NewSubGroupConstraint(id, desc, isUnite, 0, ec.groupsCount)
		if !isUnite {
			genderSensitve, err := groupsSheet.Cell(i, 3).Int()
			if err == nil && genderSensitve == 1 {
				g.genderSensitive = true
			}
			speardToAll, err := groupsSheet.Cell(i, 4).Int()
			if err == nil && speardToAll == 1 {
				g.speadToAll = true
			}

		}
		ec.Constraints = append(ec.Constraints, g)
	}
	ec.maleGroup = NewSubGroupConstraint(i, "בנים", false, 70, ec.groupsCount)
	ec.maleGroup.genderSensitive = true
	ec.maleGroup.speadToAll = true
	ec.Constraints = append(ec.Constraints, ec.maleGroup)
	i++
	ec.allGroup = NewSubGroupConstraint(i, "כולם", false, 200, ec.groupsCount)
	ec.allGroup.genderSensitive = true
	ec.allGroup.speadToAll = true
	ec.Constraints = append(ec.Constraints, ec.allGroup)

	//Init Pupils
	pupilsSheet := ec.getSheet("Pupils")
	assign := 0

	for i := 1; i < len(pupilsSheet.Rows); i++ {

		if IsEmpty(pupilsSheet.Cell(i, CELL_NAME)) {
			break
		}
		p := new(Pupil)
		p.startGroup = assign
		assign++
		if assign == ec.groupsCount {
			assign = 0
		}
		ec.pupils = append(ec.pupils, p)
		name, _ := pupilsSheet.Cell(i, CELL_NAME).FormattedValue()
		p.name = name

		//gender:
		if v, _ := pupilsSheet.Cell(i, CELL_GENDER).Int(); v == 1 {
			p.isMale = true
		}

		v, _ := pupilsSheet.Cell(i, CELL_INITIAL).Int()
		if v > 0 {
			p.locked = true
			p.lockedGroup = v - 1
		}
	}

	for _, p := range ec.pupils {
		p.groupsCount = 0
	}

	initializePreferences(ec, pupilsSheet)
	InitializeGroupsMembers(ec, pupilsSheet)

	sort.Sort(ByGroupCount(ec.pupils))

	//since indexes moved, recreate
	initializePreferences(ec, pupilsSheet)
	InitializeGroupsMembers(ec, pupilsSheet)

	for _, c := range ec.Constraints {
		c.AfterInit(ec, err)
	}

	validateConflicts(ec, err)

	return ec, err.ToString()
}
*/
func Initialize2(user *User, taskId int) (*ExecutionContext, string) {
	ec := NewExecutionContext()
	ec.taskId = taskId
	ec.groupsCount = 3 //todo
	connect()

	err := NewStringBuffer()
	// initialize all sub groups from the "Groups" sheet
	groups, e := db.Query("select id, name, sgtype, gendersensitive, speadevenly , inactive, minAllowed, maxAllowed from subgroups where tenant=? and task=?", user.getTenant(), taskId)
	if e != nil {
		panic(e)
	}
	for groups.Next() {
		var id int
		var name string
		var sgtype int
		var gendersensitive int
		var speadevenly int
		var inactive int
		var minAllowed, maxAllowed int
		groups.Scan(&id, &name, &sgtype, &gendersensitive, &speadevenly, &inactive, &minAllowed, &maxAllowed)
		isUnite := sgtype == 0
		g := NewSubGroupConstraint(id, name, isUnite, 0, ec.groupsCount)
		if !isUnite {
			g.genderSensitive = (gendersensitive == 1)
			g.speadToAll = (speadevenly == 1)
		}
		if inactive == 1 {
			g.disabled = true
		}
		g.minAllowedGiven = minAllowed
		g.maxAllowedGiven = maxAllowed

		ec.Constraints = append(ec.Constraints, g)
	}
	groups.Close()

	ec.maleGroup = NewSubGroupConstraint(9999, "בנים", false, 70, ec.groupsCount)
	ec.maleGroup.genderSensitive = false
	ec.maleGroup.speadToAll = true
	ec.Constraints = append(ec.Constraints, ec.maleGroup)

	ec.allGroup = NewSubGroupConstraint(10000, "כולם", false, 200, ec.groupsCount)
	ec.allGroup.genderSensitive = true
	ec.allGroup.speadToAll = true
	ec.Constraints = append(ec.Constraints, ec.allGroup)

	//Init Pupils
	ec.activePupilsCount = 0
	pupils, e := db.Query("select id, name, gender, inactive, remarks from pupils where tenant=? and task=? order by inactive, id", user.getTenant(), taskId)
	if e != nil {
		panic(e)
	}
	//assign := 0

	for pupils.Next() {
		var id int
		var name string
		var gender int
		var inactive int
		var remarks string
		pupils.Scan(&id, &name, &gender, &inactive, &remarks)
		p := new(Pupil)
		//p.startGroup = assign

		//p.InitValuesDomain(ec.groupsCount)
		//p.optionsLeft = ec.groupsCount - 1

		//assign++
		//if assign == ec.groupsCount {
		//	assign = 0
		//}
		ec.pupils = append(ec.pupils, p)
		p.name = name
		p.id = id
		p.isMale = (gender == 1)
		p.inactive = (inactive == 1)
		if !p.inactive {
			ec.activePupilsCount++
		}
		p.remarks = remarks
	}

	for _, p := range ec.pupils {
		p.groupsCount = 0
		p.unsatisfiedGroupsCount = make([]int64, len(ec.Constraints))
	}

	ec.domainValues = NewValuesDomainMain(ec.groupsCount, ec.activePupilsCount, len(ec.Constraints))

	initializePreferences2(ec, user, taskId)
	InitializeGroupsMembers2(ec, user, taskId)

	sort.Sort(ByGroupCount(ec.pupils[0:ec.activePupilsCount]))

	//since indexes moved, recreate
	initializePreferences2(ec, user, taskId)
	InitializeGroupsMembers2(ec, user, taskId)

	for _, c := range ec.Constraints {
		c.AfterInit(ec, err)
	}

	validateConflicts(ec, err)

	//lockLargestUniteGroup(ec)

	return ec, err.ToString()
}

/*
func lockLargestUniteGroup(ec *ExecutionContext) {
	largest := -1
	highestCount := 0
	for i := 0; i < len(ec.Constraints); i++ {
		if ec.Constraints[i].IsUnite && len(ec.Constraints[i].members) > highestCount {
			largest = i
			highestCount = len(ec.Constraints[i].members)
		}
	}

	if largest > -1 {
		for i := 0; i < len(ec.Constraints[largest].members); i++ {
			pupilInx := ec.Constraints[largest].members[i]
			ec.pupils[pupilInx].SetDomainOnlyOne(0)
			ec.pupils[pupilInx].optionsLeft = 0
		}
	}
}
*/

func validateConflicts(ec *ExecutionContext, err *stringBuffer) {
	//check if two unite sub-groups overlap and if yes, tigh them together
	/*
		merged := true
		for merged {
			merged = false
			for i := 0; i < len(ec.Constraints); i++ {
				g1 := ec.Constraints[i]
				if g1.disabled {
					continue
				}
				if g1.IsUnite {
					for j := 0; j < len(ec.Constraints); j++ {
						g2 := ec.Constraints[j]
						if g2.disabled {
							continue
						}
						if i != j && g2.IsUnite {
							for _, m := range g2.members {
								if g1.IsMember(m) {
									//Found overlap
									err.AppendFormat("Unite group '%s' overlaps with unite group '%s' - at least %s is in both. merging groups</br>\n", g1.Description(), g2.Description(), ec.pupils[m].name)
									merged = true
									mergeSubGroups(ec, i, j, err)
									break
								}
							}
						}
					}
				}
				if merged {
					break
				}
			}
		}
	*/
	for i := 0; i < len(ec.Constraints); i++ {
		g1 := ec.Constraints[i]
		if g1.disabled {
			continue
		}
		if g1.IsUnite {
			for j := 0; j < len(ec.Constraints); j++ {
				g2 := ec.Constraints[j]
				if g2.disabled {
					continue
				}
				if i != j && !g2.IsUnite {
					boysIncluded := 0
					girlsIncluded := 0
					for _, m := range g2.members {
						if g1.IsMember(m) {
							if ec.pupils[m].IsMale() {
								boysIncluded++
							} else {
								girlsIncluded++
							}
						}
					}
					included := boysIncluded + girlsIncluded
					if included >= 2 {
						if included > int(math.Ceil(g2.maxAllowed)) {
							//found a conflict g2 is completed included in g1
							err.AppendFormat("Group '%s' is a unite group and is a too bigger subset of the seperatation group '%s' --> the group is being disabled</br>\n", g1.desc, g2.desc)
							g1.disabled = true
						}
						if boysIncluded > g2.boysCount-g2.minBoys {
							err.AppendFormat("Group '%s' is a unite group which includes %d boys, which prevents spreading the boys evenly in group '%s' --> boys even spearding is disabled</br>\n", g1.desc, boysIncluded, g2.desc)

							g2.minBoys = 0
						}
						if girlsIncluded > len(g2.members)-g2.boysCount-g2.minGirls {
							err.AppendFormat("Group '%s' is a unite group which includes %d girls, which prevents spreading the girls evenly in group '%s' --> girls even spearding is disabled</br>\n", g1.desc, girlsIncluded, g2.desc)

							g2.minGirls = 0
						}
					}
				}
			}

		} else if len(g1.members) == 2 {
			p1 := ec.pupils[g1.members[0]]
			p2 := ec.pupils[g1.members[1]]

			if len(p1.prefs) == 1 && p1.prefs[0] == g1.members[1] ||
				len(p2.prefs) == 1 && p2.prefs[0] == g1.members[0] {
				err.AppendFormat("Pupil '%s' and '%s' are members of '%s' - a seperation group, and have eachother as only preference --> disabling the group", p1.name, p2.name, g1.Description())
				g1.disabled = true
			}
		}

	}
}

func mergeSubGroups(ec *ExecutionContext, g1 int, g2 int, err *stringBuffer) {
	ec.Constraints[g1].desc = fmt.Sprintf("מיזוג: '%s' עם '%s'", ec.Constraints[g1].desc, ec.Constraints[g2].desc)
	for _, m := range ec.Constraints[g2].members {
		if !ec.Constraints[g1].IsMember(m) {
			ec.Constraints[g1].AddMember(m, ec)
		}
	}

	//removeConstraint(g2, ec)
	ec.Constraints[g2].disabled = true

	ec.Constraints[g1].AfterInit(ec, err)
}

func removeConstraint(cIndex int, ec *ExecutionContext) {
	a := ec.Constraints
	i, j := cIndex, cIndex+1

	copy(a[i:], a[j:])
	for k, n := len(a)-j+i, len(a); k < n; k++ {
		a[k] = nil
	}
	ec.Constraints = a[:len(a)-j+i]
}

func initializePreferences(ec *ExecutionContext, pupilsSheet *xlsx.Sheet) {

	for i := 1; i <= len(ec.pupils); i++ {

		name, _ := pupilsSheet.Cell(i, CELL_NAME).FormattedValue()
		pIndex := ec.findPupil(name)
		p := ec.pupils[pIndex]
		p.prefs = nil
		//preferences
		//		pupilIndex := i_0based
		var pupilPrefConstraint [3]int

		for j := 0; j < NUM_OF_PREF; j++ {
			refPupil, _ := pupilsSheet.Cell(i, CELL_PREF+j).FormattedValue()

			refIndex := ec.findPupil(refPupil)
			if refIndex != -1 {
				//name := fmt.Sprintf("Pref %d", j+1)
				//pref := NewPrefConstraint(len(ec.Constraints), name, prefPrioWeight[j], pupilIndex,
				//	refIndex, j)
				//ec.Constraints = append(ec.Constraints, pref)
				p.prefs = append(p.prefs, refIndex)
				pupilPrefConstraint[j] = len(ec.Constraints) - 1
				ec.pupils[refIndex].groupsCount++
			} else {
				//adds the weight of the missing constrains to #1,2
				if j == 1 {
					//	ec.Constraints[pupilPrefConstraint[0]].(*PrefConstraint).weight += prefPrioWeight[1] + prefPrioWeight[2]
				} else if j == 2 {
					//	ec.Constraints[pupilPrefConstraint[0]].(*PrefConstraint).weight += prefPrioWeight[2] / 2
					//	ec.Constraints[pupilPrefConstraint[1]].(*PrefConstraint).weight += prefPrioWeight[2] / 2
				}

				break
			}

		}
	}
}

func InitializeGroupsMembers(ec *ExecutionContext, pupilsSheet *xlsx.Sheet) {

	for _, c := range ec.Constraints {
		c.members = nil
		c.boysCount = 0
	}

	for i := 1; i <= len(ec.pupils); i++ {

		name, _ := pupilsSheet.Cell(i, CELL_NAME).FormattedValue()
		pupilIndex := ec.findPupil(name)
		p := ec.pupils[pupilIndex]

		ec.allGroup.AddMember(pupilIndex, ec)
		if p.IsMale() {
			ec.maleGroup.AddMember(pupilIndex, ec)
			p.groupsCount++
		}

		grps, _ := pupilsSheet.Cell(i, CELL_SUBGROUP).FormattedValue()
		if grps != "" {
			subgroupsCellArray := strings.Split(grps, ",")
			for _, subGroupID := range subgroupsCellArray {
				subGroupIdInt, _ := strconv.Atoi(strings.TrimSpace(subGroupID))
				if subGroupIdInt < 0 {
					//todo
					MsgBox("תת קבוצה מיוצגת על ידי מספר מ - 1 עד ")
					return
				}
				if subGroupIdInt > 0 {
					grpIndex := ec.findGroup(subGroupIdInt)
					if grpIndex == -1 {
						MsgBox(fmt.Sprintf("could not find group %d", subGroupIdInt))
					}
					sg := ec.Constraints[grpIndex]
					sg.AddMember(pupilIndex, ec)
					p.groupsCount++
				}
			}
		}
	}
}

func initializePreferences2(ec *ExecutionContext, user *User, taskId int) {

	for _, p := range ec.pupils {
		p.prefs = nil
		p.origOrderPrefs = nil
		p.incomingPrefs = nil
	}

	prefs, err := db.Query("select pupilId, refPupilId, priority from pupilPrefs where tenant=? and task=? order by pupilId, priority",
		user.getTenant(), taskId)
	if err != nil {
		panic(err)
	}

	for prefs.Next() {
		var pupilId int
		var refPupilId int
		var priority int
		prefs.Scan(&pupilId, &refPupilId, &priority)

		pupilIndex := ec.findPupilById(pupilId)
		pupilRefIndex := ec.findPupilById(refPupilId)
		if pupilIndex >= 0 && pupilRefIndex >= 0 {
			p := ec.pupils[pupilIndex]
			if !(ec.pupils[pupilRefIndex].inactive) {
				p.prefs = append(p.prefs, pupilRefIndex)
			}
			p.origOrderPrefs = append(p.origOrderPrefs, pupilRefIndex)
			if len(p.origOrderPrefs) > 3 {
				doNoth()
			}

			// //sort highest first - needed for the algorithm
			for i := len(p.prefs) - 2; i >= 0; i-- {
				if p.prefs[i+1] > p.prefs[i] {
					temp := p.prefs[i+1]
					p.prefs[i+1] = p.prefs[i]
					p.prefs[i] = temp
				} else {
					break
				}
			}
		}

	}

	for i, p := range ec.pupils[0:ec.activePupilsCount] {
		for _, pref := range p.prefs {
			ec.pupils[pref].incomingPrefs = append(ec.pupils[pref].incomingPrefs, i)
		}
	}

}
func doNoth() {}

func InitializeGroupsMembers2(ec *ExecutionContext, user *User, taskId int) {

	for _, c := range ec.Constraints {
		c.members = nil
		c.boysCount = 0
	}

	stmt, err := db.Prepare("select groupId from subgroupPupils where tenant=? and task=? and pupilId=?")
	if err != nil {
		panic(err)
	}
	mailIndex := ec.findGroup(ec.maleGroup.ID())
	allIndex := ec.findGroup(ec.allGroup.ID())
	for _, p := range ec.pupils[0:ec.activePupilsCount] {
		p.uniteGroups = []int{}
		p.seperationGroups = []int{}
		pupilIndex := ec.findPupil(p.name)

		ec.allGroup.AddMember(pupilIndex, ec)
		p.seperationGroups = append(p.seperationGroups, allIndex)
		if p.IsMale() {
			ec.maleGroup.AddMember(pupilIndex, ec)
			p.seperationGroups = append(p.seperationGroups, mailIndex)
			p.groupsCount++
		}
		res, err := stmt.Query(user.getTenant(), taskId, p.id)
		if err != nil {
			panic(err)
		}
		for res.Next() {
			var groupId int

			res.Scan(&groupId)
			grpIndex := ec.findGroup(groupId)
			sg := ec.Constraints[grpIndex]
			sg.AddMember(pupilIndex, ec)
			if sg.IsUnite {
				p.uniteGroups = append(p.uniteGroups, grpIndex)
			} else {
				p.seperationGroups = append(p.seperationGroups, grpIndex)
			}
			p.groupsCount++
		}
	}
}

func (ec *ExecutionContext) findPupil(name string) int {
	for i, p := range ec.pupils {
		if p.name == name {
			return i
		}
	}
	return -1
}

func (ec *ExecutionContext) findPupilById(id int) int {
	for i, p := range ec.pupils {
		if p.id == id {
			return i
		}
	}
	return -1
}

func (ec *ExecutionContext) getSheet(name string) *xlsx.Sheet {
	if ec.dataExcel == nil {
		//MsgBox("Code must be called from a Data excel")
		var err error
		ec.dataExcel, err = xlsx.OpenFile("/Users/i022021/Dev/tmp/" + ec.file)
		if err != nil {
			return nil
		}
	}
	for _, sheet := range ec.dataExcel.Sheets {
		if sheet.Name == name {
			return sheet
		}
	}

	return nil
}

func (ec *ExecutionContext) findGroup(id int) int {
	for i, c := range ec.Constraints {
		if c.ID() == id {
			return i
		}

	}
	return -1
}

var colors = []string{"green", "yellow", "red", "blue"}

const INACTIVE_COLOR = "gray"

func getColor(group int) string {
	if group >= 0 && group < len(colors) {
		return colors[group]
	}
	return ""
}

func (ec *ExecutionContext) printHtml(resultID int, clean bool) string {
	res := NewStringBuffer()
	res.Clear()

	for _, c := range ec.Constraints {
		c.calculateMembersCounts(ec, ec)
		c.ValidateNew(ec, ec)
	}

	res.Append("<head>")
	res.Append("<link href=\"https://fonts.googleapis.com/icon?family=Material+Icons\" rel=\"stylesheet\" >")

	res.Append("<link rel=\"stylesheet\" type=\"text/css\" href=\"/index.css\">")
	res.Append("<script src=\"jquery.min.js\"></script>")
	res.Append("<script>")
	res.AppendFormat("\nMAX_GROUPS=%d;", ec.groupsCount)
	res.AppendFormat("\nresultID=%d;", resultID)
	res.Append("</script>")

	res.Append("<script src=\"/results.js\"></script>")
	res.Append("</head>")
	res.Append("<body onload='setEvents()'>")

	res.Append("<H2> כיתות </H2></br>\n")
	res.Append("<table id='groups' border=1><tr><th>#</th>")
	for i := 0; i < ec.groupsCount; i++ {
		colorName := getColor(i)
		res.AppendFormat("<th bgcolor=\"%s\">כיתה %d</th>", colorName, i+1)
	}
	res.Append("</tr>\n")

	var list [10][]int //todo
	for i, p := range ec.pupils[0:ec.activePupilsCount] {
		list[p.groupBestScore] = append(list[p.groupBestScore], i)
	}

	max := 0
	for i := 0; i < ec.groupsCount; i++ {
		sort.Sort(&SortPupilList{list: list[i], ec: ec})
		if len(list[i]) > max {
			max = len(list[i])
		}
	}
	for i := 0; i < max; i++ {
		res.AppendFormat("<tr><td>%d</td>", i+1)
		for j := 0; j < ec.groupsCount; j++ {
			name := ""
			pid := -1
			if len(list[j]) > i {
				name = ec.pupils[list[j][i]].name
				pid = ec.pupils[list[j][i]].id //list[j][i]
			}
			res.AppendFormat("<td class=\"moveable\" pid= \"%d\" gid=\"%d\" name=\"encryptedCell\">%s</td>", pid, j, name)
		}
		res.Append("</tr>")
	}
	res.Append("</table></br></br>\n\n")
	if !clean {

		//Print list of Pupils and their preferences
		res.Append("<H2> מפת העדפות </H2></br>\n")
		//res.AppendFormat("<P>Sum of Matching Pref: %d, Sum of Matching Firsts: %d<P></br>\n", ec.bestSumOfSatisfiedPrefs, ec.bestSumOfSatisfiedFirstPrefs)
		res.Append("<table border='1'><tr><th>#</th><th>שם</th><th>בחירה 1</th><th>בחירה 2</th><th>בחירה 3</th><th>הערות</th></tr>\n")

		numOfAll := 0
		numOfNone := 0
		numOfOneLess := 0
		numOfTwoLess := 0
		numOfFirst := 0
		numOfSecond := 0
		numOfThird := 0
		numOfFirstAndSecond := 0
		numOfFirstAndThird := 0
		numOfSecondAndThird := 0
		numOfThree := 0

		for inx, p := range ec.pupils[0:ec.activePupilsCount] {
			colorName := getColor(p.groupBestScore)

			//name
			res.Append(fmt.Sprintf("<tr><td>%d</td><td bgcolor=%s name=\"encryptedCell\">%s</td>", inx+1, colorName, p.name))
			//preferences
			numOfMatch := 0
			got1 := false //len(p.origOrderPrefs) > 0 && p.groupBestScore == ec.pupils[p.origOrderPrefs[0]].groupBestScore
			got2 := false //len(p.origOrderPrefs) > 1 && p.groupBestScore == ec.pupils[p.origOrderPrefs[1]].groupBestScore
			got3 := false //len(p.origOrderPrefs) > 2 && p.groupBestScore == ec.pupils[p.origOrderPrefs[2]].groupBestScore

			for i := 0; i < len(p.origOrderPrefs); i++ {
				refP := ec.pupils[p.origOrderPrefs[i]]
				colorPref := getColor(refP.groupBestScore)
				if p.origOrderPrefs[i] >= ec.activePupilsCount {
					//this is an inactive pupil
					colorPref = INACTIVE_COLOR
				}
				if colorName == colorPref {
					numOfMatch++
					if i == 0 {
						got1 = true
					} else if i == 1 {
						got2 = true
					} else if i == 2 {
						got3 = true
					}
				}

				res.Append(fmt.Sprintf("<td bgcolor=\"%s\" name=\"encryptedCell\">%s</td>", colorPref, refP.name))

			}

			if numOfMatch == len(p.prefs) {
				numOfAll++
			} else if numOfMatch+1 == len(p.prefs) {
				numOfOneLess++
			} else if numOfMatch+2 == len(p.prefs) {
				numOfTwoLess++
			}

			if got1 && got2 && got3 {
				numOfThree++
			} else if got1 && !got2 && !got3 {
				numOfFirst++
			} else if !got1 && got2 && !got3 {
				numOfSecond++
			} else if !got1 && !got2 && got3 {
				numOfThird++
			} else if got1 && got3 {
				numOfFirstAndThird++
			} else if got1 && got2 {
				numOfFirstAndSecond++
			} else if got2 && got3 {
				numOfSecondAndThird++
			} else if !got1 && !got2 && !got3 {
				numOfNone++
			}
			for i := len(p.origOrderPrefs); i < NUM_OF_PREF; i++ {
				res.Append("<td>ללא העדפות</td>")
			}

			//הערות
			res.Append("<td>")
			res.Append(p.remarks)
			if !got1 && !got2 && got3 {
				res.Append(fmt.Sprintf("<p style='background-color:red;'>רק שלישית</p>,"))
			}
			if len(p.origOrderPrefs) > 0 && numOfMatch == 0 {
				res.Append(fmt.Sprintf("<p style='background-color:red;'>לא קיבל אף עדיפות</p>"))
			}
			res.Append("</td>")
			res.Append("</tr>\n")

		}
		res.Append("</table></br>\n")
		res.AppendFormat("Num of Pupil got all pref: %d, one less: %d, two less: %d</br>\n", numOfAll, numOfOneLess, numOfTwoLess)
		res.AppendFormat("1:%d, 2:%d, 3:%d\n</br>1+2:%d, 1+3:%d, 2+3:%d\n</br>all3:%d, none:%d\n",
			numOfFirst, numOfSecond, numOfThird, numOfFirstAndSecond, numOfFirstAndThird, numOfSecondAndThird, numOfThree, numOfNone)

		//groups
		res.Append("<h1>קבוצות</h1></br>\n")

		res.Append("<table border=\"1\"><tr><th>#</th><th>שם</th><th>סוג</th><th  width='35%'>חברי הקבוצה</th>")
		for i := 0; i < ec.groupsCount; i++ {
			colorName := getColor(i)
			res.AppendFormat("<th bgcolor=\"%s\">כיתה %d (מס' בנים) (מספר בנות)</th>", colorName, i+1)
		}
		res.Append("</tr>\n")

		for _, c := range ec.Constraints {
			groupType := "איחוד"
			if !c.IsUnite {
				groupType = "פירוד"
				if c.speadToAll {
					groupType = "פירוד מלא"
				}
			}
			members := ""
			for _, m := range c.members {
				color := getColor(ec.pupils[m].groupBestScore)
				gender := "&#9792;"
				if ec.pupils[m].IsMale() {
					gender = "&#9794;"
				}
				members += fmt.Sprintf("<span style=\"background-color: %s\">%s %s</span>, ", color, ec.pupils[m].name, gender)
			}

			res.AppendFormat("<tr><td>%d</td><td>%s</td><td>%s</td><td>%s</td>", c.ID(), c.Description(), groupType, members)
			for i := 0; i < ec.groupsCount; i++ {

				res.AppendFormat("<td>%d - (%d)(%d)</td>", c.countForGroup[i], c.boysForGroup[i], c.countForGroup[i]-c.boysForGroup[i])
			}
			res.Append("</tr>\n")
		}

		res.Append("</table></br>\n")

		//list of Conflicts

		res.Append("<h1>Conflicts</h1></br>\n")

		for _, c := range ec.Constraints {
			if msg := c.Message(ec); msg != "" {
				res.AppendFormat("* %s </br>\n", msg)
			}
		}

		for _, p := range ec.pupils {
			count := 0
			for i := 0; i < len(p.prefs); i++ {
				if p.groupBestScore == ec.pupils[p.prefs[i]].groupBestScore {
					count++
				}
			}
			if count == 0 && len(p.prefs) > 0 {
				res.AppendFormat("%s קיבל/ה 0 מתוך %d העדפות</br>\n", p.name, len(p.prefs))
			}
		}
		/*
			for _, cc := range ec.ConstraintsAll {
				c, ok := cc.(*PrefConstraint)
				if ok {
					if msg := c.Message(ec); msg != "" {
						res.AppendFormat("* %s </br>\n", msg)
					}
				}
			}
		*/
	}
	res.Append("</body>")

	return res.ToString()

}
