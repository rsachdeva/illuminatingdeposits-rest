package userauthn

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidAuthHeaderExpired(t *testing.T) {
	err := valid(authHeaderExpired())
	require.NotNil(t, err)
	require.Regexp(t, regexp.MustCompile("Token is expired"), err)
}

func TestValidAuthHeaderNotPresent(t *testing.T) {
	err := valid(authHeaderNotPresent())
	require.NotNil(t, err)
	require.Regexp(t, regexp.MustCompile("no authorization header"), err)
}

func authHeaderExpired() []string {
	return []string{
		"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Imdyb3d0aEBkcmlubm92YXRpb25zLnVzIiwicm9sZXMiOlsiVVNFUiJdLCJleHAiOjE2MTAwNTA2NjAsImlzcyI6ImdpdGh1Yi5jb20vcnNhY2hkZXZhL2lsbHVtaW5hdGluZ2RlcG9zaXRzLXJlc3QifQ.Q6bOd3qO-2sJoZLm0unR0XXefg8FnRaYZtyoUollE20",
	}
}

func authHeaderNotPresent() []string {
	return nil
}
