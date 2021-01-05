package reqlog

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

func Dump(r *http.Request, info string) {
	output, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println("Error dumping request:", err)
		return
	}
	fmt.Printf("Dumping entire request for %s..is %v\n", info, string(output))
}
