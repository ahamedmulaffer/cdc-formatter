package mysqlFormatter

import(
	// "fmt"
)

type AfterType string



func (after AfterType) loopRequiredFields(tableName string, operation string, requiredFields []string, payloadAfter map[string]interface{}, finalPayload map[string]interface{}){
	idToHPIDFuncCalled := false
	for _, reqField := range requiredFields {
		if idToHPIDFuncCalled && reqField == "HPID" {
			continue
		}
		//check if map to id enbled
		if sliceField(requiredFields).isMapIDToHPIDEnabled(tableName, operation, idToHPIDFuncCalled) {
			idToHPIDFuncCalled = true
			if _, ok := payloadAfter["ID"]; ok {
				finalPayload["HPID"] = payloadAfter["ID"]
				delete(payloadAfter, "ID")
			}
			continue
		}
		//check isNullToAvailable if value nil only
		if payloadAfter[reqField] == nil {
			val, ok := isNullToAvailable(tableName, operation)
			if ok {
				finalPayload[reqField] = val
			}
			continue
		}
		finalPayload[reqField] = payloadAfter[reqField]
		if _, ok:= payloadAfter["HostID"]; ok {
			hostId := assertToString(payloadAfter["HostID"])
			after.underscoreWithHostidFields(tableName,operation,reqField,finalPayload,hostId)
			after.underscoreWithHostidAssetFields(tableName,operation,reqField,finalPayload,hostId)
		}
		after.timeFieldConversion(tableName,operation,reqField,finalPayload)
		after.strConversion(tableName,operation,reqField,finalPayload)
		//field type conversion
		//all value to string

	}
}



func (after AfterType) isRequiredFieldsAvailable(tableName string, operation string) ([]string, bool){
	if len(allRequiredAfterFields) == 0 {
		return nil, false
	}
	if _, ok := allRequiredAfterFields[tableName]; !ok {
		return nil, false
	}
	if len(allRequiredAfterFields[tableName]) == 0 {
		return nil, false
	}
	if _, ok := allRequiredAfterFields[tableName][operation]; !ok {
		return nil, false
	}
	if len(allRequiredAfterFields[tableName][operation]) == 0 {
		return nil, false
	}
	return allRequiredAfterFields[tableName][operation], true
}

func (after AfterType) loopPayloadFields(tableName string, operation string, payloadAfter map[string]interface{}, finalPayload map[string]interface{}) {
	idToHPIDFuncCalled := false
	for field, _ := range payloadAfter {
		if idToHPIDFuncCalled && field == "HPID" {
			continue
		}
		if mapField(payloadAfter).isMapIDToHPIDEnabled(tableName, operation, idToHPIDFuncCalled){
			idToHPIDFuncCalled = true
			if _, ok := payloadAfter["ID"]; ok {
				finalPayload["HPID"] = payloadAfter["ID"]
				delete(payloadAfter, "ID")
			}
			continue
		}
		//check isNullToAvailable if value nil only
		if payloadAfter[field] == nil {
			val, ok := isNullToAvailable(tableName, operation)
			if ok {
				finalPayload[field] = val
			}
			continue
		}
		finalPayload[field] = payloadAfter[field]
		if _, ok:= payloadAfter["HostID"]; ok {
			hostId := assertToString(payloadAfter["HostID"])
			after.underscoreWithHostidFields(tableName,operation,field,finalPayload,hostId)
			after.underscoreWithHostidAssetFields(tableName,operation,field,finalPayload,hostId)
		}
		after.timeFieldConversion(tableName,operation,field,finalPayload)
		after.strConversion(tableName,operation,field,finalPayload)
	}
}