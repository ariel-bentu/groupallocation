package main

import (
	"fmt"
	"log"
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
	return len(ec.pupils)
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
	Cancel    bool

	groupsCount    int
	iterationCount int
	Constraints    []*SubGroupConstraint
	allGroup       *SubGroupConstraint
	maleGroup      *SubGroupConstraint
	pupils         []*Pupil
	dataExcel      *xlsx.File

	currentIteration     int
	levelChangeIteration int
	statusCandidate      []int
	bestTotal            int
}

var RunningExecutions map[string]*ExecutionContext

func NewExecutionContext() *ExecutionContext {
	ec := new(ExecutionContext)
	ec.startTime = time.Now()
	ec.done = false
	ec.Cancel = false
	ec.bestTotal = -99999

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

func (e *ExecutionContext) GetStatusHtml() (string, string) {
	if e.done {
		return e.printHtml(), fmt.Sprintf("%f", e.endTime.Sub(e.startTime).Seconds())
	}
	//return fmt.Sprintf("Interaction Count:%d, progress: %d", e.currentIteration), fmt.Sprintf("%f", time.Now().Sub(e.startTime).Seconds(),
	//	e.currentIteration/e.iterationCount*100)

	sb := NewStringBuffer()
	sb.AppendFormat("Candidate Length:%d, elappsed: %f.0<br>\n%s</br>\n", e.currentIteration, time.Now().Sub(e.startTime).Seconds(),
		slice2String(e.statusCandidate))

	for _, c := range e.Constraints {

		sb.AppendFormat("%s - %d</br>\n", c.Description(), c.unsatisfiedCount)
	}
	sb.Append("\n</br>-----------------------------------</br>\n")
	for _, p := range e.pupils {
		sb.AppendFormat("%s - %d</br>\n", p.name, p.unsatisfiedPrefsCount)
	}

	return sb.ToString(), ""
}
func (e *ExecutionContext) Finish() {
	e.done = true
	e.endTime = time.Now()
}

func IsEmpty(c *xlsx.Cell) bool {
	var v, err = c.String()
	return err != nil || v == ""
}
func MsgBox(text string) {
	log.Output(0, text)
}

func (ec *ExecutionContext) GetParam(name string) string {
	sheet := ec.getSheet("Configuration")
	for _, row := range sheet.Rows {
		if v, _ := row.Cells[0].String(); v == name {
			p, _ := row.Cells[1].String()
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

func Initialize() *ExecutionContext {
	ec := NewExecutionContext()

	ec.groupsCount = ec.GetIntParam("מספר כיתות")
	ec.iterationCount = ec.GetIntParam("מספר איטרציות")

	var prefPrioWeight [3]int
	prefPrioWeight[0] = ec.GetIntParam("משקל העדפה ראשונה")
	prefPrioWeight[1] = ec.GetIntParam("משקל העדפה שניה")
	prefPrioWeight[2] = ec.GetIntParam("משקל העדפה שלישית")

	groupsSheet := ec.getSheet("Groups")

	// initialize all sub groups from the "Groups" sheet
	i := 1
	for ; !IsEmpty(groupsSheet.Cell(i, 2)); i++ {
		id, _ := groupsSheet.Cell(i, 0).Int()
		v, _ := groupsSheet.Cell(i, 1).String()
		desc, _ := groupsSheet.Cell(i, 2).String()
		if IsEmpty(groupsSheet.Cell(i, 3)) {
			//todo
			MsgBox(fmt.Sprint("חובה להזין משקל לקבוצה - ראה שורה %d", i))
			return nil
		}
		w, _ := groupsSheet.Cell(i, 3).Int()
		g := NewSubGroupConstraint(id, desc, v == UNITE_VALUE, w)
		ec.Constraints = append(ec.Constraints, g)
	}
	ec.maleGroup = NewSubGroupConstraint(i, "בנים", false, 70)
	ec.Constraints = append(ec.Constraints, ec.maleGroup)
	i++
	ec.allGroup = NewSubGroupConstraint(i, "כולם", false, 200)
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
		name, _ := pupilsSheet.Cell(i, CELL_NAME).String()
		p.name = name
		//		pupilIndex := len(ec.pupils) - 1

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

	//for _, p := range ec.pupils {
	//		p.groupsCount = 0
	//}

	//since indexes moved, recreate
	initializePreferences(ec, pupilsSheet)
	InitializeGroupsMembers(ec, pupilsSheet)

	for _, c := range ec.Constraints {
		c.AfterInit(ec)
	}

	validateConflicts(ec)

	//sort.Sort(ByMembers(ec.Constraints))
	//ec.Reshuffel()
	//ec.balancePupils()

	return ec
}

func validateConflicts(ec *ExecutionContext) {
	for i := 0; i < len(ec.Constraints); i++ {
		g1 := ec.Constraints[i]
		if g1.IsUnite {
			for j := 0; j < len(ec.Constraints); j++ {
				g2 := ec.Constraints[j]
				if i != j && !g2.IsUnite {
					boysIncluded := 0
					girlsIncluded := 0
					if g1.ID() == 23 && g2.ID() == 3 {
						girlsIncluded = 0
					}
					for _, m := range g2.Members() {
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
						if included > len(g2.Members())/2 {
							//found a conflict g2 is completed included in g1
							fmt.Printf("Group %d is unite group and is a too bigger subset of the seperatation group '%d' - it is being disabled</br>\n", g1.ID(), g2.ID())
							g1.disabled = true
						}
						if boysIncluded > g2.boysCount-g2.minBoys {
							fmt.Printf("Group %d is unite group which include %d boys, which prevents spreading the boys evenly in group %d  - boys even spearding is disabled</br>\n", g1.ID(), boysIncluded, g2.ID())

							g2.minBoys = 0
						}
						if girlsIncluded > len(g2.Members())-g2.boysCount-g2.minGirls {
							fmt.Printf("Group %d is unite group which include %d girls, which prevents spreading the girls evenly in group %d  - girls even spearding is disabled</br>\n", g1.ID(), girlsIncluded, g2.ID())

							g2.minGirls = 0
						}
					}
				}
			}

		} else if len(g1.Members()) == 2 {
			p1 := ec.pupils[g1.Members()[0]]
			p2 := ec.pupils[g1.Members()[1]]

			if len(p1.prefs) == 1 && p1.prefs[0] == g1.Members()[1] ||
				len(p2.prefs) == 1 && p2.prefs[0] == g1.Members()[0] {
				fmt.Printf("Pupil '%s' and '%s' are members of '%s' - a seperation group and have eachother as preference", p1.name, p2.name, g1.Description())

			}
		}

	}
}

func initializePreferences(ec *ExecutionContext, pupilsSheet *xlsx.Sheet) {

	for i := 1; i <= len(ec.pupils); i++ {

		name, _ := pupilsSheet.Cell(i, CELL_NAME).String()
		pIndex := ec.findPupil(name)
		p := ec.pupils[pIndex]
		p.prefs = nil
		//preferences
		//		pupilIndex := i_0based
		var pupilPrefConstraint [3]int

		for j := 0; j < NUM_OF_PREF; j++ {
			refPupil, _ := pupilsSheet.Cell(i, CELL_PREF+j).String()

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

		name, _ := pupilsSheet.Cell(i, CELL_NAME).String()
		pupilIndex := ec.findPupil(name)
		p := ec.pupils[pupilIndex]

		ec.allGroup.AddMember(pupilIndex, ec)
		if p.IsMale() {
			ec.maleGroup.AddMember(pupilIndex, ec)
			p.groupsCount++
		}

		grps, _ := pupilsSheet.Cell(i, CELL_SUBGROUP).String()
		if grps != "" {
			subgroupsCellArray := strings.Split(grps, ",")
			for _, subGroupID := range subgroupsCellArray {
				subGroupIdInt, _ := strconv.Atoi(strings.TrimSpace(subGroupID))
				if subGroupIdInt < 1 {
					//todo
					MsgBox("תת קבוצה מיוצגת על ידי מספר מ - 1 עד ")
					return
				}
				grpIndex := ec.findGroup(subGroupIdInt)
				sg := ec.Constraints[grpIndex]
				sg.AddMember(pupilIndex, ec)
				p.groupsCount++
			}
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

func (ec *ExecutionContext) getSheet(name string) *xlsx.Sheet {
	if ec.dataExcel == nil {
		//MsgBox("Code must be called from a Data excel")
		var err error
		ec.dataExcel, err = xlsx.OpenFile("c:/temp/file4.xlsx")
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

var colors = []string{"green", "blue", "red", "yellow"}

func getColor(group int) string {
	if group >= 0 && group < len(colors) {
		return colors[group]
	}
	return ""
}

func (ec *ExecutionContext) printHtml() string {
	res := NewStringBuffer()
	res.Clear()

	for _, c := range ec.Constraints {
		c.ValidateNew(ec, ec)
	}

	//Print list of Pupils and their preferences
	res.Append("<H2> List </H2></br>\n")
	res.Append("<table><tr><th>#</th><th>שם</th><th>בחירה 1</th><th>בחירה 2</th><th>בחירה 3</th></tr>\n")

	for inx, p := range ec.pupils {
		colorName := getColor(p.groupBestScore)

		//name
		res.Append(fmt.Sprintf("<tr><td>%d</td><td bgcolor=%s>%s</td>", inx+1, colorName, p.name))
		//preferences

		for i := 0; i < len(p.prefs); i++ {
			refP := ec.pupils[p.prefs[i]]
			colorPref := getColor(refP.groupBestScore)
			res.Append(fmt.Sprintf("<td bgcolor=\"%s\">%s</td>", colorPref, refP.name))
		}

		for i := len(p.prefs); i < NUM_OF_PREF; i++ {
			res.Append("<td>אין</td>")
		}
		res.Append("</tr>\n")

	}
	res.Append("</table></br>\n")

	//groups
	res.Append("<h1>Groups</h1></br>\n")

	res.Append("<table border=\"1\"><tr><th>#</th><th>שם</th><th>סוג</th>")
	for i := 0; i < ec.groupsCount; i++ {
		res.AppendFormat("<th>כיתה %d (מס' בנים)</th>", i+1)
	}
	res.Append("</tr>\n")

	for _, c := range ec.Constraints {
		groupType := "איחוד"
		if !c.IsUnite {
			groupType = "פירוד"
		}
		res.AppendFormat("<tr><td>%d</td><td>%s</td><td>%s</td>", c.ID(), c.Description(), groupType)
		for i := 0; i < ec.groupsCount; i++ {
			res.AppendFormat("<td>%d - (%d)</td>", c.countForGroup[i], c.boysForGroup[i])
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
	return res.ToString()

}
