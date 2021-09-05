package mysqlFormatter

import(
	// "fmt"
)

type sliceField []string
type mapField map[string]interface{}

func (fields sliceField) isMapIDToHPIDEnabled(tableName string, operation string, enabled bool) bool {
	if enabled {
		return false
	}
	if !canDoIDToHPID(tableName,operation, allIDToHPIDTableOperations) {
		return false
	}
	idAvailable := false
	hpidAvailable := false
	for _, v := range fields {
		if v == "ID" {
			idAvailable = true
		}
		if v == "HPID" {
			hpidAvailable = true
		}
	}
	if !idAvailable {
		return false
	}
	if !hpidAvailable {
		return false
	}
	return true
}


func (fields mapField)isMapIDToHPIDEnabled(tableName string, operation string, enabled bool) bool {
	// fmt.Println(enabled)
	if enabled {
		return false
	}
	if !canDoIDToHPID(tableName,operation, allIDToHPIDTableOperations) {
		return false
	}
	idAvailable := false
	hpidAvailable := false
	for k, _ := range fields {
		if k == "ID" {
			idAvailable = true
		}
		if k == "HPID" {
			hpidAvailable = true
		}
	}
	// fmt.Println(idAvailable, hpidAvailable)
	if !idAvailable {
		return false
	}
	if !hpidAvailable {
		return false
	}
	return true
}

func canDoIDToHPID(tableName string, operation string, dataMap map[string]map[string]bool) bool{
	if len(dataMap) == 0 {
		return false
	}
	if _, ok := dataMap[tableName]; !ok {
		return false
	}
	if len(dataMap[tableName]) == 0 {
		return false
	}
	if _, ok := dataMap[tableName][operation]; !ok {
		return false
	}
	return true
}