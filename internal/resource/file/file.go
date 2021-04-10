package file

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Resource interface {
	ReadJSONFile(fileName string, target interface{}) (err error)
}

func New() Resource {
	env := os.Getenv("XEYNSEENV")
	dir := "./files/xeynseJar_analytics"
	if env == "production" {
		dir = "/etc/xeynseJar_analytics"
	}
	return &resource{
		configDirPath: dir,
	}
}

type resource struct {
	configDirPath string
}

func (r *resource) ReadJSONFile(fileName string, target interface{}) (err error) {
	fileRead, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", r.configDirPath, fileName))
	if err != nil {
		return
	}
	if err = json.Unmarshal(fileRead, target); err != nil {
		return
	}
	return
}
