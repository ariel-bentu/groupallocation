package main

import (
	"fmt"
	"strings"

	"encoding/json"

	"github.com/tealeg/xlsx"
)

func getSheet(file string, name string) (*xlsx.Sheet, *xlsx.File) {
	dataExcel, err := xlsx.OpenFile("/Users/i022021/Dev/tmp/" + file)
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
	Id   string `json:"id"`
	Name string `json:"name"`
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
