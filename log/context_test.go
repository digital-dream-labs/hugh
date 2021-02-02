package log

import (
	"context"
	"testing"
)

func Test(t *testing.T) {
	ctx := context.Background()

	// Empty tests first
	if FromContext(ctx) == nil {
		t.Log("Logger should never be nil")
		t.FailNow()
	}

	AddContextFields(ctx, Fields{"This does nothing": true})

	// Set logger
	lCtx := AddToContext(ctx)

	if lCtx == ctx {
		t.Log("Expected contexts to not be equal but they were")
		t.FailNow()
	}

	AddContextFields(lCtx, Fields{"Greeting": "Wie Gehts"})

	if FromContext(lCtx) == nil {
		t.Log("Logger best not be nil")
		t.FailNow()
	}

	// Set logger again
	if xCtx := AddToContext(lCtx); xCtx != lCtx {
		t.Log("Both contexts should be the same since we already set the logger")
		t.FailNow()
	}
}
