package mongoFormatter

import(
	// "fmt"
)

type AfterType string



func (after AfterType) loopRequiredFields(collectionName string, operation string, requiredFields []string, payloadAfter map[string]interface{}, finalPayload map[string]interface{}){
	for _, reqField := range requiredFields {
		if reqField == "id" || reqField == "_id"{
			finalPayload["id"] = extract_id(payloadAfter["_id"])
			continue
		}
		//check isNullToAvailable if value nil only
		if payloadAfter[reqField] == nil {
			val, ok := isNullToAvailable(collectionName, operation)
			if ok {
				finalPayload[reqField] = val
			}
			continue
		}
		finalPayload[reqField] = payloadAfter[reqField]
		if hostId, ok:= hostIdAvailable(collectionName, operation); ok {
			after.underscoreWithHostidFields(collectionName,operation,reqField,finalPayload,hostId)
			after.underscoreWithHostidArrayFields(collectionName,operation,reqField,finalPayload,hostId)
		}
		after.strConversion(collectionName,operation,reqField,finalPayload)
	}
}



func (after AfterType) isRequiredFieldsAvailable(collectionName string, operation string) ([]string, bool){
	if len(allRequiredAfterFields) == 0 {
		return nil, false
	}
	if _, ok := allRequiredAfterFields[collectionName]; !ok {
		return nil, false
	}
	if len(allRequiredAfterFields[collectionName]) == 0 {
		return nil, false
	}
	if _, ok := allRequiredAfterFields[collectionName][operation]; !ok {
		return nil, false
	}
	if len(allRequiredAfterFields[collectionName][operation]) == 0 {
		return nil, false
	}
	return allRequiredAfterFields[collectionName][operation], true
}

func (after AfterType) loopPayloadFields(collectionName string, operation string, payloadAfter map[string]interface{}, finalPayload map[string]interface{}) {
	for field, _ := range payloadAfter {
		if field == "id" || field == "_id"{
			finalPayload["id"] = extract_id(payloadAfter["_id"])
			continue
		}
		//check isNullToAvailable if value nil only
		if payloadAfter[field] == nil {
			val, ok := isNullToAvailable(collectionName, operation)
			if ok {
				finalPayload[field] = val
			}
			continue
		}
		finalPayload[field] = payloadAfter[field]
		if hostId, ok:= hostIdAvailable(collectionName, operation); ok {
			after.underscoreWithHostidFields(collectionName,operation,field,finalPayload,hostId)
			after.underscoreWithHostidArrayFields(collectionName,operation,field,finalPayload,hostId)
		}
		after.strConversion(collectionName,operation,field,finalPayload)
	}
}