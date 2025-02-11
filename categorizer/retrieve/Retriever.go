package retrieve

import (
	"context"
)

// Retriever : interface, defines the general method "retrieve" which is used to collect reconstructed tcp streams from a service
type Retriever interface {
	Retrieve(ctx context.Context, cancel context.CancelFunc, results chan<- Result)
}

type Result struct {
	Stream  string
	SrcPort uint16
}
