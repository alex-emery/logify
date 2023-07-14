package parser

import (
	"os"
	"path/filepath"

	couchbasev2 "github.com/couchbase/couchbase-operator/pkg/apis/couchbase/v2"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"

	"sigs.k8s.io/yaml"
)

type LogArtifact struct {
	Pods   map[string]*corev1.Pod
	Pvc    map[string]*corev1.PersistentVolumeClaim
	Cbc    map[string]*couchbasev2.CouchbaseCluster
	Events []corev1.Event
}

type CbOpInfo struct {
	Namespaces map[string]*LogArtifact
}

func newCbOpInfo() CbOpInfo {
	return CbOpInfo{
		Namespaces: make(map[string]*LogArtifact),
	}
}

type Parser struct {
	logger *zap.Logger
}

func New(logger *zap.Logger) *Parser {
	return &Parser{logger: logger}
}

func (p *Parser) Parse(path string) (*CbOpInfo, error) {
	path += "/namespace"

	namespaces, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	cbopInfo := newCbOpInfo()
	for _, namespace := range namespaces {
		if !namespace.IsDir() {
			continue
		}

		namespacePath := filepath.Join(path, namespace.Name())
		namespaceLogs, err := ParseNamespace(namespacePath)
		if err != nil {
			p.logger.Error("failed to parse resources", zap.String("namespace", namespace.Name()), zap.Error(err))
			continue
		}

		cbopInfo.Namespaces[namespace.Name()] = namespaceLogs

		events, err := ParseEvents(namespacePath)
		if err != nil {
			p.logger.Error("failed to parse events", zap.String("namespace", namespace.Name()), zap.Error(err))
			continue
		}

		cbopInfo.Namespaces[namespace.Name()].Events = events
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
			pods, err := ParseResourceFolder[*corev1.Pod](resourcePath)
			if err != nil {
				continue
			}

			artifacts.Pods = pods
		case "persistentvolumeclaim":
			pvcs, err := ParseResourceFolder[*corev1.PersistentVolumeClaim](resourcePath)
			if err != nil {
				continue
			}

			artifacts.Pvc = pvcs
		case "couchbasecluster":
			clusters, err := ParseResourceFolder[*couchbasev2.CouchbaseCluster](resourcePath)
			if err != nil {
				continue
			}

			artifacts.Cbc = clusters
		}
	}

	return artifacts, nil
}

func ParseResourceYaml[T any](path string) (T, error) {
	var resource = new(T)
	file, err := os.ReadFile(path)
	if err != nil {
		return *resource, err
	}

	err = yaml.Unmarshal(file, resource)
	return *resource, err
}

type Resource interface {
	GetName() string
}

func ParseResourceFolder[T Resource](path string) (map[string]T, error) {
	resources, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	parsedResources := make(map[string]T)
	for _, resource := range resources {
		resourceYaml := filepath.Join(path, resource.Name(), resource.Name()+".yaml")

		resource, err := ParseResourceYaml[T](resourceYaml)
		if err != nil {
			continue
		}

		parsedResources[resource.GetName()] = resource

	}

	return parsedResources, nil
}
