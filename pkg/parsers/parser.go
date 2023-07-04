package parser

import (
	"os"
	"path/filepath"

	couchbasev2 "github.com/couchbase/couchbase-operator/pkg/apis/couchbase/v2"
	corev1 "k8s.io/api/core/v1"

	"sigs.k8s.io/yaml"
)

type LogArtifact struct {
	pods []*corev1.Pod
	pvc  []*corev1.PersistentVolumeClaim
	cbc  []*couchbasev2.CouchbaseCluster
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
		resourcePath := filepath.Join(path, resource.Name())
		switch resource.Name() {
		case "pod":
			pods, err := ParseResourceFolder[corev1.Pod](resourcePath)
			if err != nil {
				continue
			}

			artifacts.pods = pods
		case "persistentvolumeclaim":
			pvcs, err := ParseResourceFolder[corev1.PersistentVolumeClaim](resourcePath)
			if err != nil {
				continue
			}

			artifacts.pvc = pvcs
		case "couchbasecluster":
			clusters, err := ParseResourceFolder[couchbasev2.CouchbaseCluster](resourcePath)
			if err != nil {
				continue
			}

			artifacts.cbc = clusters
		}
	}

	return artifacts, nil
}

func ParseResourceFolder[T any](path string) ([]*T, error) {
	resources, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	parsedResources := make([]*T, 0, len(resources))
	for _, resource := range resources {
		resourceYaml := filepath.Join(path, resource.Name(), resource.Name()+".yaml")
		resource, err := ParseResourceYaml[T](resourceYaml)
		if err != nil {
			continue
		}
		parsedResources = append(parsedResources, resource)

	}

	return parsedResources, nil
}

func ParseResourceYaml[T any](path string) (*T, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var resource = new(T)

	err = yaml.Unmarshal(file, resource)
	return resource, err
}
