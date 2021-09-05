package mysqlFormatter

func isNullToAvailable(tableName string, operation string) (string, bool){
	if len(allNullToTableOperations) == 0 {
		return "", false
	}
	if _, ok := allNullToTableOperations[tableName]; !ok {
		return "", false
	}
	if len(allNullToTableOperations[tableName]) == 0 {
		return "", false
	}
	if _, ok := allNullToTableOperations[tableName][operation]; !ok {
		return "", false
	}
	return allNullToTableOperations[tableName][operation], true
}