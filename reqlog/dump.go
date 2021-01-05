package reqlog

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

func Dump(r *http.Request, info string) {
	log.Printf("for %s simply req is %v ..\n", info, r)
	output, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println("Error dumping request:", err)
		return
	}
	fmt.Printf("Dumping entire request for %s..\n", info)
	fmt.Println(string(output))
}
