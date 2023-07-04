package parser

import (
	"os"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

// goes through the directory assuming everything is a folder containing
// a <name>.yaml for pod.
func ParsePods(path string) ([]*corev1.Pod, error) {
	pods, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	parsedPods := make([]*corev1.Pod, 0, len(pods))
	for _, pod := range pods {
		podYaml := filepath.Join(path, pod.Name(), pod.Name()+".yaml")
		pod, err := ParsePodYaml(podYaml)
		if err != nil {
			continue
		}
		parsedPods = append(parsedPods, pod)

	}

	return parsedPods, nil
}

func ParsePodYaml(path string) (*corev1.Pod, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	pod := &corev1.Pod{}

	err = yaml.Unmarshal(file, pod)
	return pod, err
}
