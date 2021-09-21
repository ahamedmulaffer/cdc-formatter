package mysqlFormatter

import(
	// "fmt"
)

func isCreateOrUpdateByAvailable(tableName string, operation string, field string) (string, bool){
	if field != "CreatedBy" && field != "UpdatedBy" {
		return "", false
	}
	if len(allCreatedByOrUpdatedByTableOperations) == 0 {
		return "", false
	}
	if _, ok := allCreatedByOrUpdatedByTableOperations[tableName]; !ok {
		return "", false
	}
	if len(allCreatedByOrUpdatedByTableOperations[tableName]) == 0 {
		return "", false
	}
	if _, ok := allCreatedByOrUpdatedByTableOperations[tableName][operation]; !ok {
		return "", false
	}
	return allCreatedByOrUpdatedByTableOperations[tableName][operation], true
}