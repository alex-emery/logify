package parser

import (
	"os"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
)

type LogArtifact struct {
	pods []*corev1.Pod
	pvc  []*corev1.PersistentVolumeClaim
}

type CbOpInfo struct {
	Namespaces map[string]*LogArtifact
}

func NewCbOpInfo() CbOpInfo {
	return CbOpInfo{
		Namespaces: make(map[string]*LogArtifact),
	}
}

// parses the namespace folder within a cbopinfo folder.
// i.e cbopinfo/namespaces/<namespace>
func ParseCbOpInfo(path string) (*CbOpInfo, error) {
	path += "/namespace"

	namespaces, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	cbopInfo := NewCbOpInfo()
	for _, namespace := range namespaces {
		if !namespace.IsDir() {
			continue
		}

		namespaceLogs, err := ParseNamespace(filepath.Join(path, namespace.Name()))
		if err != nil {
			// parsing is best effort so continue anywa and just let the user know it failed
			continue
		}

		cbopInfo.Namespaces[namespace.Name()] = namespaceLogs
	}

	return &cbopInfo, nil
}

func ParseNamespace(path string) (*LogArtifact, error) {
	resources, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	artifacts := &LogArtifact{}
	for _, resource := range resources {
		switch resource.Name() {
		case "pod":
			pods, err := ParsePods(filepath.Join(path, resource.Name()))
			if err != nil {
				continue
			}

			artifacts.pods = pods
		case "persistentvolumeclaim":
		}
	}

	return artifacts, nil
}
