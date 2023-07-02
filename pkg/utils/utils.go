package utils

import (
	"fmt"
	"io/ioutil"
	"os/exec"

	"gopkg.in/yaml.v3"
)

type Release struct {
	ReleaseName string                 `yaml:"releaseName"`
	RepoURL     string                 `yaml:"repoURL"`
	Chart       string                 `yaml:"chart"`
	Version     string                 `yaml:"version"`
	Namespace   string                 `yaml:"namespace"`
	Enabled     bool                   `yaml:"enabled,omitempty"`
	ValuesFile  string                 `yaml:"valuesFile,omitempty"`
	Values      map[string]interface{} `yaml:"values,omitempty"`
}

func ReadFile(inputFile string) []byte {
	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}
	return data
}

func ParseReleases(data []byte) []Release {
	var releases []Release
	err := yaml.Unmarshal(data, &releases)
	if err != nil {
		panic(err)
	}
	return releases
}

func MergeMaps(dst, src map[string]interface{}) {
	for k, v := range src {
		if _, ok := dst[k]; !ok {
			// If the key does not exist in dst, add it from src
			dst[k] = v
		} else {
			// If the key exists in dst, but the value is not a map, overwrite it with the value from src
			switch dstValue := dst[k].(type) {
			case map[string]interface{}:
				// If the value is a map, recursively merge the values
				srcValue, ok := v.(map[string]interface{})
				if ok {
					MergeMaps(dstValue, srcValue)
				}
			default:
				dst[k] = v
			}
		}
	}
}

func RunCommand(cmd string) {
	out, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		fmt.Printf("Error running command: %v\n", err)
		fmt.Printf("Command output: %s\n", out)
		panic(err)
	}
	//fmt.Printf("Command output: %s\n", string(out))
}
