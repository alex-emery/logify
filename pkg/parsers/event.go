package parser

import (
	"os"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
)

// searches recursively through a directory looking for files
// filepath.Glob doesn't handle **/events.yaml so
// we have a custom function to do it
func glob(dir string, filename string) ([]string, error) {

	files := []string{}
	err := filepath.Walk(dir, func(path string, _ os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Base(path) == filename {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

func ParseEvents(path string) ([]corev1.Event, error) {
	files, err := glob(path, "events.yaml")
	if err != nil {
		return nil, err
	}

	events := make([]corev1.Event, 0, len(files))
	for _, file := range files {
		event, err := ParseResourceYaml[[]corev1.Event](file)
		if err != nil {
			return nil, err
		}

		events = append(events, event...)
	}

	return events, nil
}
