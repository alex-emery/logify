package parser

import (
	"testing"
)

func TestOperatorLog(t *testing.T) {
	_, err := NewOperatorLogs("../../testdata/cbopinfo/namespace/default/deployment/couchbase-operator/couchbase-operator.log")
	if err != nil {
		t.Fatal(err)
	}
}
