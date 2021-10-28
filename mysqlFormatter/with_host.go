package mysqlFormatter

import(
	"encoding/json"
	// "strings"
)


func (after AfterType) underscoreWithHostidFields(tableName string, operation string, field string, finalPayload map[string]interface{}, HostID string){
	if canDoAnUnderscore(tableName,operation,field,allUnderscoreWithHostidFields) {
		addUnderscoreToNormalFields(finalPayload, field, HostID)
	}
}

func (before BeforeType) underscoreWithHostidFields(tableName string, operation string, field string, finalPayload map[string]interface{}, HostID string){
	if canDoAnUnderscore(tableName,operation,field,allUnderscoreWithHostidFields) {
		addUnderscoreToNormalFields(finalPayload, field, HostID)
	}
}


func (after AfterType) underscoreWithHostidAssetFields(tableName string, operation string, field string, finalPayload map[string]interface{}, HostID string){
	if canDoAnUnderscore(tableName,operation,field,allUnderscoreWithHostidAssetFields) {
		addUnderscoreToAssetFields(finalPayload, field, HostID)
	}
}

func (before BeforeType) underscoreWithHostidAssetFields(tableName string, operation string, field string, finalPayload map[string]interface{}, HostID string){
	if canDoAnUnderscore(tableName,operation,field,allUnderscoreWithHostidAssetFields) {
		addUnderscoreToAssetFields(finalPayload, field, HostID)
	}
}

func canDoAnUnderscore(tableName string, operation string, field string, dataMap map[string]map[string]map[string]string) bool {
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
	if len(dataMap[tableName][operation]) == 0 {
		return false
	}
	if _, ok := dataMap[tableName][operation][field]; !ok {
		return false 
	}
	return true
}

func addUnderscoreToNormalFields(finalPayload map[string]interface{}, field string, HostID string){
	val, ok := canBeAssertToString(finalPayload[field])
	if ok {
		finalPayload[field] = val+"_"+HostID
	}
}

func addUnderscoreToAssetFields(finalPayload map[string]interface{}, field string, HostID string){
	val, ok := canBeAssertToString(finalPayload[field])
	if ok {
		var sliceAssetsMap []map[string]interface{}
		var sliceAssets []string
		err := json.Unmarshal([]byte(val), &sliceAssetsMap)
		if err != nil && val != ""{
			err := json.Unmarshal([]byte(val), &sliceAssets)
			if err != nil && val != ""{
				finalPayload[field] = HostID+"_"+val
				return
			}
			if len(sliceAssets) == 0 {
				return
			}
			for k,asset := range sliceAssets {
				if asset == "" {
					continue
				}
				sliceAssets[k] = HostID+"_"+asset
			}
			slice, _ := json.Marshal(sliceAssets)
			finalPayload[field] = string(slice)
			return
		}
		if len(sliceAssetsMap) == 0 {
			return
		}
		for _,asset := range sliceAssetsMap {
			if _, ok := asset["file_name"]; !ok {
				continue
			}
			fileName, ok := canBeAssertToString(asset["file_name"])
			if !ok {
				continue
			}
			asset["file_name"] = HostID+"_"+fileName
		}
		mapSlice, _ := json.Marshal(sliceAssetsMap)
		finalPayload[field] = string(mapSlice)
	}
}