package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"encoding/json"

	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tealeg/xlsx"
)

func getSheet(file string, name string) (*xlsx.Sheet, *xlsx.File) {
	dataExcel, err := xlsx.OpenFile(file)
	if err != nil {
		return nil, nil
	}
	for _, sheet := range dataExcel.Sheets {
		if sheet.Name == name {
			return sheet, dataExcel
		}
	}

	return nil, nil
}

type IdNameJson struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	IsUnite    bool   `json:"isUnite"`
	IsInactive bool   `json:"isInactive"`
	IsMale     bool   `json:"isMale"`
}

type IdRefIDJson struct {
	Id    string `json:"id"`
	RefId string `json:"refId"`
}

func getPupilList(file string) ([]byte, error) {
	var pupils []IdNameJson

	pupilsSheet, _ := getSheet(file, "Pupils")

	if pupilsSheet != nil {
		for i := 1; i < len(pupilsSheet.Rows); i++ {
			name := pupilsSheet.Cell(i, CELL_NAME).String()
			if name == "" {
				break
			}
			newP := IdNameJson{Id: fmt.Sprintf("%d", i), Name: name}
			pupils = append(pupils, newP)
		}
	}
	return json.Marshal(pupils)
}

func getSubgroupList(file string) ([]byte, error) {
	var groups []IdNameJson

	groupsSheet, _ := getSheet(file, "Groups")
	if groupsSheet != nil {
		i := 1
		for ; !IsEmpty(groupsSheet.Cell(i, 2)); i++ {
			id := groupsSheet.Cell(i, 0).String()
			desc, _ := groupsSheet.Cell(i, 2).FormattedValue()
			newG := IdNameJson{Id: id, Name: desc}
			groups = append(groups, newG)
		}
	}
	return json.Marshal(groups)
}

func setSubGroupPupils(file string, groupId string, decoder *json.Decoder) error {
	pupilsSheet, dataFile := getSheet(file, "Pupils")

	//pupils := make([]IdRefIDJson, 0)
	var pupils []IdRefIDJson
	err := decoder.Decode(&pupils)
	if err != nil {
		return err
	}
	modified := false
	for i := 1; i < len(pupilsSheet.Rows); i++ {
		name := pupilsSheet.Cell(i, CELL_NAME).String()
		if name == "" {
			break
		}
		groups := pupilsSheet.Cell(i, CELL_SUBGROUP).String()
		newGroups, changed := mergeSubGroupToPupilsGroups(groups, groupId, i, pupils)
		if changed {
			//put it back to excel
			pupilsSheet.Cell(i, CELL_SUBGROUP).SetString(newGroups)
			modified = true
		}
	}

	if modified {
		return dataFile.Save("/Users/i022021/Dev/tmp/" + file)
	}
	return nil

}

func getSubGroupPupils(file string, groupID string) ([]byte, error) {
	var groupPupils []IdRefIDJson

	pupilsSheet, _ := getSheet(file, "Pupils")

	if pupilsSheet != nil {
		for i := 1; i < len(pupilsSheet.Rows); i++ {
			name := pupilsSheet.Cell(i, CELL_NAME).String()
			if name == "" {
				break
			}
			groups := pupilsSheet.Cell(i, CELL_SUBGROUP).String()
			if PupilInSubgroup(groups, groupID) {
				newP := IdRefIDJson{Id: groupID, RefId: fmt.Sprintf("%d", i)}
				groupPupils = append(groupPupils, newP)
			}

		}
	}
	return json.Marshal(groupPupils)
}

func PupilInSubgroup(groups string, groupID string) bool {
	subgroupsCellArray := strings.Split(groups, ",")
	for _, subGroupID := range subgroupsCellArray {
		if groupID == strings.TrimSpace(subGroupID) {
			return true
		}
	}
	return false
}

func mergeSubGroupToPupilsGroups(groups string, groupId string, pupilId int, pupils []IdRefIDJson) (string, bool) {
	subgroupsCellArray := strings.Split(groups, ",")
	pupilIdStr := fmt.Sprintf("%d", pupilId)
	for i, subGroupID := range subgroupsCellArray {
		if groupId == strings.TrimSpace(subGroupID) {
			for _, p := range pupils {
				if p.RefId == pupilIdStr {
					//no change
					return "", false
				}
			}
			//need to remove the group from this pupil
			newGroups := ""
			for j, subGroupID := range subgroupsCellArray {
				if i != j {
					if j > 0 {
						newGroups += ","
					}
					newGroups += subGroupID
				}
			}
			return newGroups, true
		}
	}
	for _, p := range pupils {
		if p.RefId == pupilIdStr {
			if groups != "" {
				groups += ","
			}
			groups += groupId
			return groups, true
		}
	}
	//no change
	return "", false

}

func test() {
	connect()
	stmt, _ := db.Prepare("INSERT INTO pupils (tenant, task, id, name, gender) values (?,?,?,?,?)")
	stmt.Exec("ariel", 1, 1, "אריאל", 1)
}

var db *sql.DB = nil

func connect() {
	if db == nil {
		path := os.Getenv("GA_DBPATH")
		if path == "" {
			path = "/Users/i022021/Dev/tmp/groups.db"
		}
		var err error
		db, err = sql.Open("sqlite3", path)
		if err != nil {
			panic(err)
		}
	}
}

func checkErrPanic(err error, tx *sql.Tx) {
	if err != nil {
		tx.Rollback()
		panic(err)
	}
}

func uploadExcel(user *User, path string, taskName string, passcode string) {
	encrypt := passcode != ""
	var key, iv []byte
	if encrypt {
		key, iv = AESencryptInit(passcode)
	}

	tenant := user.getTenant()
	connect()
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	task := createTask(tx, tenant, taskName)

	pupilsSheet, _ := getSheet(path, "Pupils")

	stmt, err1 := tx.Prepare("INSERT INTO pupils (tenant, task, id, name, gender) values (?,?,?,?,?)")
	checkErrPanic(err1, tx)
	stmtSubGroupPupils, err2 := tx.Prepare("INSERT INTO subgroupPupils (tenant, task, groupId, pupilId) values (?,?,?,?)")
	checkErrPanic(err2, tx)

	var pupils []*Pupil

	if pupilsSheet != nil {
		for i := 1; i < len(pupilsSheet.Rows); i++ {
			name := pupilsSheet.Cell(i, CELL_NAME).String()
			if name == "" {
				break
			}
			if encrypt {
				name = AESencrypt(key, iv, name)
			}

			gender, _ := pupilsSheet.Cell(i, CELL_GENDER).Int()
			groups := pupilsSheet.Cell(i, CELL_SUBGROUP).String()
			stmt.Exec(tenant, task, i, name, gender)

			pupils = append(pupils, &Pupil{id: i, name: name})

			subgroupsCellArray := strings.Split(groups, ",")
			for _, subGroupID := range subgroupsCellArray {
				if subGroupID != "" {
					groupId, _ := strconv.Atoi(strings.TrimSpace(subGroupID))
					stmtSubGroupPupils.Exec(tenant, task, groupId, i)
				}
			}
		}
	}
	stmt.Close()
	stmtSubGroupPupils.Close()

	stmtPupilPrefs, err3 := tx.Prepare("INSERT INTO pupilPrefs (tenant, task, pupilId, refPupilId, priority) values (?,?,?,?,?)")
	checkErrPanic(err3, tx)

	for i := 1; i < len(pupilsSheet.Rows); i++ {
		for j := 0; j < NUM_OF_PREF; j++ {
			refPupil, _ := pupilsSheet.Cell(i, CELL_PREF+j).FormattedValue()
			if encrypt {
				refPupil = AESencrypt(key, iv, refPupil)
			}
			refId := findPupilId(pupils, refPupil)
			if refId != -1 {
				stmtPupilPrefs.Exec(user.getTenant(), task, i, refId, j)
			}
		}
	}
	stmtPupilPrefs.Close()

	stmt, _ = tx.Prepare(`INSERT INTO subgroups (tenant, task, id, name, sgtype, gendersensitive, speadevenly) values 
						  (?,?,?,?,?,?,?)`)

	groupsSheet, _ := getSheet(path, "Groups")
	if groupsSheet != nil {
		i := 1
		for ; !IsEmpty(groupsSheet.Cell(i, 2)); i++ {
			genderSensitve := 0
			speardToAll := 0
			sgtype := 0
			id := groupsSheet.Cell(i, 0).String()
			desc, _ := groupsSheet.Cell(i, 2).FormattedValue()
			sgtypeStr, _ := groupsSheet.Cell(i, 1).FormattedValue()
			isUnite := (sgtypeStr == UNITE_VALUE)
			if !isUnite {
				sgtype = 1
				genderSensitve, _ = groupsSheet.Cell(i, 3).Int()
				speardToAll, _ = groupsSheet.Cell(i, 4).Int()
			}
			stmt.Exec(tenant, task, id, desc, sgtype, genderSensitve, speardToAll)
		}
	}

	tx.Commit()

}

func findPupilId(pupils []*Pupil, name string) int {
	for _, p := range pupils {
		if p.name == name {
			return p.id
		}
	}
	return -1
}

func createTask(tx *sql.Tx, tenant string, taskName string) int {
	connect()
	r := db.QueryRow("select max(task) from task")
	maxId := 0
	if r != nil {
		r.Scan(&maxId)
	}
	_, err := db.Exec("insert into task (tenant, task, name, createDate) values (?,?,?,?)", tenant, maxId+1, taskName, time.Now().Unix())
	if err != nil {
		panic(err)
	}
	return maxId + 1
}

func deleteTask(user *User, taskId int) {
	connect()
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	_, err = tx.Exec("delete from task where tenant=? and task=?", user.getTenant(), taskId)
	checkErrPanic(err, tx)
	_, err = tx.Exec("delete from pupils where tenant=? and task=?", user.getTenant(), taskId)
	checkErrPanic(err, tx)
	_, err = tx.Exec("delete from subgroups where tenant=? and task=?", user.getTenant(), taskId)
	checkErrPanic(err, tx)
	_, err = tx.Exec("delete from subgroupPupils where tenant=? and task=?", user.getTenant(), taskId)
	checkErrPanic(err, tx)
	_, err = tx.Exec("delete from pupilPrefs where tenant=? and task=?", user.getTenant(), taskId)
	checkErrPanic(err, tx)

	tx.Commit()
}

func getTaskList(user *User) ([]byte, error) {
	var tasks []IdNameJson

	connect()
	res, err := db.Query("select task, name from task where tenant =? order by name", user.getTenant())
	if err != nil {
		panic(err)
	}
	for res.Next() {
		var id int
		var name string
		err = res.Scan(&id, &name)
		newP := IdNameJson{Id: fmt.Sprintf("%d", id), Name: name}
		tasks = append(tasks, newP)
	}
	return json.Marshal(tasks)
}

func getPupilList2(user *User, taskId int) ([]byte, error) {
	var pupils []IdNameJson

	connect()
	res, err := db.Query("select id, name, gender from pupils where tenant =? and task=? order by name", user.getTenant(), taskId)
	if err != nil {
		panic(err)
	}
	for res.Next() {
		var id int
		var name string
		var gender int
		isMale := false
		err = res.Scan(&id, &name, &gender)
		if gender == 1 {
			isMale = true
		}
		newP := IdNameJson{Id: fmt.Sprintf("%d", id), Name: name, IsMale: isMale}
		pupils = append(pupils, newP)
	}

	return json.Marshal(pupils)
}

func getPupilPrefs(user *User, taskId int, pupilId int) ([]byte, error) {
	var pupils []IdRefIDJson

	connect()
	res, err := db.Query("select pupilId, refPupilId from pupilPrefs where tenant =? and task=? and pupilId =? order by priority", user.getTenant(), taskId, pupilId)
	if err != nil {
		panic(err)
	}
	for res.Next() {
		var id int
		var refId string
		err = res.Scan(&id, &refId)

		newP := IdRefIDJson{Id: fmt.Sprintf("%d", id), RefId: refId}
		pupils = append(pupils, newP)
	}

	return json.Marshal(pupils)
}

func setPupilPrefs(user *User, taskId int, pupilId int, decoder *json.Decoder) error {
	var pupils []IdRefIDJson
	err := decoder.Decode(&pupils)
	if err != nil {
		return err
	}

	connect()
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	_, err = tx.Exec("delete from pupilPrefs where tenant=? and task=? and pupilId=?", user.getTenant(), taskId, pupilId)
	if err != nil {
		panic(err)
	}

	stmt, err := tx.Prepare("insert into pupilPrefs (tenant, task, pupilId, refPupilId, priority) values (?,?,?,?,?)")

	for i, p := range pupils {
		_, err = stmt.Exec(user.getTenant(), taskId, pupilId, p.RefId, i)
		if err != nil {
			tx.Rollback()
			stmt.Close()
			panic(err)
		}
	}
	stmt.Close()
	tx.Commit()
	return nil
}

func getSubgroupList2(user *User, taskId int) ([]byte, error) {
	var groups []IdNameJson

	connect()
	res, err := db.Query("select id, name, sgtype, inactive from subgroups where tenant=? and task=?", user.getTenant(), taskId)
	if err != nil {
		panic(err)
	}
	for res.Next() {
		var id int
		var name string
		var sgtype int
		var inactive int
		res.Scan(&id, &name, &sgtype, &inactive)
		newG := IdNameJson{Id: fmt.Sprintf("%d", id), Name: name, IsUnite: (sgtype == 0), IsInactive: (inactive == 1)}
		groups = append(groups, newG)

	}
	return json.Marshal(groups)
}

func getIntFromBool(val bool) int {
	if val {
		return 1
	}
	return 0
}

func createNewSubgroup(user *User, taskId int, name string) {
	connect()
	r := db.QueryRow("select max(id) from subgroups")
	groupId := 1
	if r != nil {
		r.Scan(&groupId)
		groupId++

	}

	_, err := db.Exec(`insert into subgroups (tenant, task, id, name, sgtype, gendersensitive, speadevenly, inactive) values 
				      (?, ?, ?, ?, 0, 0, 0, 0)`,
		user.getTenant(), taskId, groupId, name)
	if err != nil {
		panic(err)
	}

}
func updateSubgroup(user *User, taskId int, groupId int, isUnite bool, isInactive bool) {
	connect()
	sgType := 1
	if isUnite {
		sgType = 0
	}
	_, err := db.Exec("update subgroups set sgtype=?, inactive=? where tenant =? and task=? and id=?",
		sgType, getIntFromBool(isInactive), user.getTenant(), taskId, groupId)
	if err != nil {
		panic(err)
	}
}

func getSubGroupPupils2(user *User, taskId int, groupId int) ([]byte, error) {
	var groupPupils []IdRefIDJson
	groupStr := fmt.Sprintf("%d", groupId)
	connect()
	res, err := db.Query("select pupilId from subgroupPupils A inner join pupils B on (A.tenant = B.tenant and A.task = B.task and pupilId = B.id) where A.tenant=? and A.task=? and A.groupId=? order by B.name", user.getTenant(), taskId, groupId)
	if err != nil {
		panic(err)
	}
	for res.Next() {
		var refId int

		res.Scan(&refId)
		newG := IdRefIDJson{Id: groupStr, RefId: fmt.Sprintf("%d", refId)}
		groupPupils = append(groupPupils, newG)

	}
	return json.Marshal(groupPupils)

}

func setSubGroupPupils2(user *User, taskId int, groupId int, decoder *json.Decoder) error {
	var pupils []IdRefIDJson
	err := decoder.Decode(&pupils)
	if err != nil {
		return err
	}

	connect()
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	_, err = tx.Exec("delete from subgroupPupils where tenant=? and task=? and groupId=?", user.getTenant(), taskId, groupId)
	if err != nil {
		panic(err)
	}

	stmt, err := tx.Prepare("insert into subgroupPupils (tenant, task, groupId, pupilId) values (?,?,?,?)")

	for _, p := range pupils {
		_, err = stmt.Exec(user.getTenant(), taskId, groupId, p.RefId)
		if err != nil {
			tx.Rollback()
			stmt.Close()
			panic(err)
		}
	}
	stmt.Close()
	tx.Commit()
	return nil
}
