package mongoFormatter

import(
	"strings"
)


func (after AfterType) underscoreWithHostidFields(collectionName string, operation string, field string, finalPayload map[string]interface{}, HostID string){
	if canDoAnUnderscore(collectionName,operation,field,allUnderscoreWithHostidFields) {
		addUnderscoreToNormalFields(finalPayload, field, HostID)
	}
}

func (before BeforeType) underscoreWithHostidFields(collectionName string, operation string, field string, finalPayload map[string]interface{}, HostID string){
	if canDoAnUnderscore(collectionName,operation,field,allUnderscoreWithHostidFields) {
		addUnderscoreToNormalFields(finalPayload, field, HostID)
	}
}


func (after AfterType) underscoreWithHostidArrayFields(collectionName string, operation string, field string, finalPayload map[string]interface{}, HostID string){
	if canDoAnUnderscore(collectionName,operation,field,allUnderscoreWithHostidArrayFields) {
		addUnderscoreToArrayFields(finalPayload, field, HostID)
	}
}

func (before BeforeType) underscoreWithHostidArrayFields(collectionName string, operation string, field string, finalPayload map[string]interface{}, HostID string){
	if canDoAnUnderscore(collectionName,operation,field,allUnderscoreWithHostidArrayFields) {
		addUnderscoreToArrayFields(finalPayload, field, HostID)
	}
}

func canDoAnUnderscore(collectionName string, operation string, field string, dataMap map[string]map[string]map[string]string) bool {
	if len(dataMap) == 0 {
		return false
	}
	if _, ok := dataMap[collectionName]; !ok {
		return false
	}
	if len(dataMap[collectionName]) == 0 {
		return false
	}
	if _, ok := dataMap[collectionName][operation]; !ok {
		return false
	}
	if len(dataMap[collectionName][operation]) == 0 {
		return false
	}
	if _, ok := dataMap[collectionName][operation][field]; !ok {
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

func addUnderscoreToArrayFields(finalPayload map[string]interface{}, field string, HostID string){
	if sliceAssets, ok := finalPayload[field].([]interface{}); ok {
		if len(sliceAssets) == 0 {
			return
		}
		for k,asset := range sliceAssets {
			fileName, ok := canBeAssertToString(asset)
			if !ok {
				continue
			}
			sliceAssets[k] = fileName+"_"+HostID
		}
		finalPayload[field] = sliceAssets
		return
	}
	if val, ok := canBeAssertToString(finalPayload[field]); ok {
		var ids []string
		for _,v := range strings.Split(val, ",") {
			ids = append(ids, v+"_"+HostID)
		}
		finalPayload[field] = ids
		return
	}

}

func hostIdAvailable(collectionName string, operation string) (string, bool){
	if len(hostIdWithCollectioOperation) == 0 {
		return "", false
	}
	if _, ok := hostIdWithCollectioOperation[collectionName]; !ok {
		return "", false
	}
	if len(hostIdWithCollectioOperation[collectionName]) == 0 {
		return "", false
	}
	if _, ok := hostIdWithCollectioOperation[collectionName][operation]; !ok {
		return "", false
	}
	return hostIdWithCollectioOperation[collectionName][operation], true
}