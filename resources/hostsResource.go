package resources

import (
	"bozosonparade/gsh"
	"encoding/json"
	"net/http"
	"sort"
	"strings"
)

// HostsResourceHandler handles requests to the /hosts/ path
func HostsResourceHandler(w http.ResponseWriter, r *http.Request) {
	mVals := r.URL.Query()
	strOp := mVals.Get("operation")
	var aRet []gsh.HostConfig

	if len(strOp) > 0 {
		for _, host := range gsh.CurrentConfig.Hosts {
			for _, strOpConf := range host.SupportedOperations {
				if strings.EqualFold(strOp, strOpConf) {
					aRet = append(aRet, host)
					break
				}
			}
		}
	} else {
		aRet = gsh.CurrentConfig.Hosts
	}

	// Now, sort by name to be nice
	sort.Sort(gsh.ByName(aRet))

	retVal, _ := json.MarshalIndent(aRet, "", "  ")
	w.Write(retVal)
}
