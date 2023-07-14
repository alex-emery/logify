package parser

import (
	"fmt"
	"testing"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
)

func TestParsePodYaml(t *testing.T) {
	filename := "../../testdata/cbopinfo/namespace/default/pod/cb-example-0000/cb-example-0000.yaml"
	pod, err := ParseResourceYaml[*corev1.Pod](filename)
	if err != nil {
		t.Fatal(err)
	}

	if pod == nil {
		t.Fatal(fmt.Errorf("parsed pod expected but not found"))
	}
	if pod.ObjectMeta.Name != "cb-example-0000" {
		t.Fatal(fmt.Errorf("failed to parse pod, didn't find expected name"))
	}

}

func TestParseCbOpInfo(t *testing.T) {
	logger, _ := zap.NewProduction()
	parser := New(logger)

	info, err := parser.Parse("../../testdata/cbopinfo")
	if err != nil {
		t.Fatal(err)
	}
	if info == nil {
		t.Fatal(fmt.Errorf("parsed cbopinfo expected but not found"))
	}

}
