package views

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

func (tp *sampleTemplate) UnmarshalJSON(data []byte) error {
	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		fmt.Printf("Error while decoding %v\n", err)
		return err
	}
	tp.Input = v["input"].(string)
	tp.Output = v["output"].(string)
	return nil
}

func unZip(src, dst string) (map[string]string, error) {
	files := make(map[string]string)

	zr, err := zip.OpenReader(src)
	defer zr.Close()
	if err != nil {
		return nil, err
	}

	for _, file := range zr.File {

		if file.FileInfo().IsDir() {
			continue
		}
		if !strings.Contains(file.Name, ".in") && !strings.Contains(file.Name, ".out") {
			continue
		}
		if strings.Contains(file.Name, "/") {
			continue
		}

		path := filepath.Join(dst, file.Name)

		fr, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer fr.Close()

		fw, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0644))
		if err != nil {
			return nil, err
		}
		defer fw.Close()

		_, err = io.Copy(fw, fr)
		if err != nil {
			return nil, err
		}

		body, err := ioutil.ReadFile(path)
		files[file.Name] = string(body)
	}

	return files, err
}
