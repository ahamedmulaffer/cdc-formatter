package mongoFormatter

func isNullToAvailable(collectionName string, operation string) (string, bool){
	if len(allNullToCollectionOperations) == 0 {
		return "", false
	}
	if _, ok := allNullToCollectionOperations[collectionName]; !ok {
		return "", false
	}
	if len(allNullToCollectionOperations[collectionName]) == 0 {
		return "", false
	}
	if _, ok := allNullToCollectionOperations[collectionName][operation]; !ok {
		return "", false
	}
	return allNullToCollectionOperations[collectionName][operation], true
}