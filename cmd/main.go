package main

import (
	"flag"
	"fmt"
	"log"
	"sort"

	parser "github.com/alex-emery/logify/pkg/parsers"
	"go.uber.org/zap"
)

func main() {
	path := flag.String("path", "", "filepath of the cbopinfo")

	flag.Parse()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	logParser := parser.New(logger)
	cbopinfo, err := logParser.Parse(*path)
	if err != nil {
		log.Panic(err)
	}

	for _, logs := range cbopinfo.Namespaces {

		// var filteredPods []DisplayPod

		// for _, pod := range logs.Pods {
		// 	if pod.Labels["couchbase_server"] == "true" {
		// 		//TODO: filter for just problem conditions
		// 		displayPod := DisplayPod{pod: pod}
		// 		for _, volume := range pod.Spec.Volumes {
		// 			if volume.VolumeSource.PersistentVolumeClaim == nil {
		// 				continue
		// 			}
		// 			foundPvc, ok := logs.Pvc[volume.VolumeSource.PersistentVolumeClaim.ClaimName]
		// 			if ok {
		// 				displayPod.pvc = append(displayPod.pvc, foundPvc)
		// 			}
		// 		}
		// 		filteredPods = append(filteredPods, displayPod)
		// 	}
		// }

		sort.Slice(logs.Events, func(i, j int) bool {
			return logs.Events[i].LastTimestamp.After(logs.Events[j].LastTimestamp.Time)
		})

		for _, event := range logs.Events {
			fmt.Println(event.InvolvedObject.Kind + " " + event.InvolvedObject.Name)
			fmt.Println(event.LastTimestamp)
			fmt.Println(event.Message)
			fmt.Println("---------")
		}
	}
}

// type DisplayPod struct {
// 	pod *v1.Pod
// 	pvc []*v1.PersistentVolumeClaim
// }
