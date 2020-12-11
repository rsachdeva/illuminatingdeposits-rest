package reqlog

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

func Dump(r *http.Request) {
	output, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println("Error dumping request:", err)
		return
	}
	fmt.Println(string(output))
}
