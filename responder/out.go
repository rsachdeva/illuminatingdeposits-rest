package responder

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/pkg/errors"
)

// Output for outputting go values as json
func Output(w io.Writer, val interface{}) error {
	jInter, err := json.MarshalIndent(val, "", "\t")
	if err != nil {
		return errors.Wrap(err, "interestcal for all banks jsoninput marshalling")
	}
	fmt.Fprintf(w, "**** JSON Output with Indent ****\n\n%s\n\n", jInter)
	return nil
}
