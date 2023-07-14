package parser

import (
	"fmt"
	"testing"
)

func TestEventParsing(t *testing.T) {
	events, err := ParseEvents("../../testdata/cbopinfo/")
	if err != nil {
		t.Fatal(err)
	}

	if len(events) != 10 {
		t.Fatal(fmt.Errorf("expected 10 events but found %d", len(events)))
	}

}
