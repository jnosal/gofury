package fury

import (
	"encoding/json"
	"strings"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
)


func ProcessFile(config interface{}, file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	switch {
	case strings.HasSuffix(file, ".yaml") || strings.HasSuffix(file, ".yml"):
		return yaml.Unmarshal(data, config)
	case strings.HasSuffix(file, ".json"):
		return json.Unmarshal(data, config)
	default:
		return nil
	}
}