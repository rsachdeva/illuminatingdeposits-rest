package testserver

import (
	"context"
	"testing"

	"google.golang.org/grpc/test/bufconn"
)

type clientResult struct {
	Listener *bufconn.Listener
	// PosgresClient *
}

func InitRestHttpServer(ctx context.Context, t *testing.T, allowPurge bool) *clientResult {

	return nil
}
