package database

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var dir string
var access = make(map[string]*sync.Mutex)

// StartService ...
func StartService(path string) {
	dir = path + "/"
	for _, file := range getFileList(path) {
		access[file.Name()] = &sync.Mutex{}
	}
}

// Create ...
func Create(name string, data map[string]interface{}) error {
	filename := name + ".json"
	if fileExits(filename) == true {
		return errors.New("Data name already exists")
	}
	file, err := os.OpenFile(dir+filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	file.Write([]byte(jsonData))
	access[filename] = &sync.Mutex{}
	return nil
}

// Read ...
func Read(name string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	filename := name + ".json"
	if fileExits(filename) == false {
		return nil, errors.New("No such name exists")
	}
	access[filename].Lock()
	defer access[filename].Unlock()
	file, err := os.Open(dir + filename)
	if err != nil {
		return nil, err
	}
	rawData, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(rawData, &data); err != nil {
		return nil, err
	}
	return data, nil
}

// Write ...
func Write(name string, data map[string]interface{}) error {
	filename := name+".json"
	file, err := os.OpenFile(dir+filename, os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if _, err := file.Write(jsonData); err != nil {
		return err
	}
	return nil
}

// Update ...
func Update(name string, data map[string]interface{}) error {
	filename := name + ".json"
	if fileExits(filename) == false {
		return errors.New("Data doesn't exits")
	}
	mainData, err := Read(name)
	if err != nil {
		return err
	}
	mergeMap(&data, &mainData)
	if err := Write(name, mainData); err != nil {
		return err
	}
	return nil
}

//------------------------------ PRIVATE

func getFileList(path string) []os.FileInfo {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	return files
}

func fileExits(filename string) bool {
	_, ok := access[filename]
	return ok
}

func mergeMap(mergeFrom, mergeTo *map[string]interface{}) {
	// Need to use stack for nested maps
	for k := range *mergeFrom {
		(*mergeTo)[k] = (*mergeFrom)[k]
	}
}
