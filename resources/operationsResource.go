package resources

import (
	"bozosonparade/gsh"
	"encoding/json"
	"net/http"
	"sort"
)

// OperationsResourceHandler handles requests to the /hosts/ path
func OperationsResourceHandler(w http.ResponseWriter, r *http.Request) {
	aRet := gsh.CurrentConfig.Operations

	// Now, sort by name to be nice
	sort.Sort(gsh.OperationsByName(aRet))

	retVal, _ := json.MarshalIndent(aRet, "", "  ")
	w.Write(retVal)
}
