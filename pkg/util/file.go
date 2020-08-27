package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	ml "jaegerperf/pkg/model"
	"os"
)

// IsFileExists checks the file availability
func IsFileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// IsDirExists checks the directory availability
func IsDirExists(dirname string) bool {
	info, err := os.Stat(dirname)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// CreateDir func
func CreateDir(dir string) {
	if !IsDirExists(dir) {
		os.MkdirAll(dir, os.ModePerm)
	}
}

// StoreJSON func
func StoreJSON(dir, filename string, data interface{}) error {
	CreateDir(dir)
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fmt.Sprintf("%s/%s.json", dir, filename), file, os.ModePerm)
	return err
}

// LoadJSON loads json from disk
func LoadJSON(dir, filename string, out interface{}) error {
	CreateDir(dir)
	b, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", dir, filename))
	if err != nil {
		return err
	}
	return json.Unmarshal(b, out)
}

// StoreFile func
func StoreFile(dir, filename string, data []byte) error {
	CreateDir(dir)
	return ioutil.WriteFile(fmt.Sprintf("%s/%s", dir, filename), data, os.ModePerm)
}

// WriteFile func
func WriteFile(dir, filename string, data []byte) error {
	CreateDir(dir)
	return ioutil.WriteFile(fmt.Sprintf("%s/%s", dir, filename), data, os.ModePerm)
}

// ReadFile func
func ReadFile(dir, filename string) ([]byte, error) {
	CreateDir(dir)
	return ioutil.ReadFile(fmt.Sprintf("%s/%s", dir, filename))
}

// ListFiles func
func ListFiles(dir string) ([]ml.File, error) {
	CreateDir(dir)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	items := make([]ml.File, 0)
	for _, file := range files {
		if !file.IsDir() {
			f := ml.File{
				Name:         file.Name(),
				Size:         file.Size(),
				ModifiedTime: file.ModTime(),
			}
			items = append(items, f)
		}
	}
	return items, nil
}
