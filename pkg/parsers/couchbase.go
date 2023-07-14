package parser

import (
	"bufio"
	"os"

	"github.com/perimeterx/marshmallow"
)

type OperatorLog struct {
	Level   string                 `json:"level"` //debug/info/error
	Ts      int64                  `json:"ts"`
	Logger  string                 `json:"logger"`
	Msg     string                 `json:"msg"`
	Cluster string                 `json:"cluster"`
	Extra   map[string]interface{} `json:"-"`
}

type OperatorLogs []OperatorLog

func NewOperatorLogs(path string) (OperatorLogs, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	logs := []OperatorLog{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		log := OperatorLog{}
		result, err := marshmallow.Unmarshal(line, &log, marshmallow.WithExcludeKnownFieldsFromMap(true))
		if err != nil {
			return nil, err
		}

		log.Extra = result
		logs = append(logs, log)
	}

	return logs, nil
}
