package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// const dataFilePath = "tickdata.json"

func getDataFilePath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	exeDir := filepath.Dir(exePath)

	dataDir := filepath.Join(exeDir, ".tick")

	err = os.MkdirAll(dataDir, 0755)
	if err != nil {
		return "", err
	}

	return filepath.Join(dataDir, "tickdata.json"), nil
}

// load reads the tick data from the JSON storage file.
// If the file does not exist, it returns an initialized empty TickData struct without error.
func load() (TickData, error) {
	dataFilePath, err := getDataFilePath()
	if err != nil {
		return TickData{}, err
	}

	data, err := os.ReadFile(dataFilePath)
	if os.IsNotExist(err) {
		return TickData{
			Current: CurrentActivity{},
			History: []HistoryEntry{},
		}, nil
	}

	if err != nil {
		return TickData{}, err
	}

	if len(data) == 0 {
		return TickData{
			Current: CurrentActivity{},
			History: []HistoryEntry{},
		}, nil
	}

	var tickData TickData
	err = json.Unmarshal(data, &tickData)
	if err != nil {
		return TickData{}, err
	}

	return tickData, nil
}

// save writes the tick data formatted as indented JSON to the storage file.
func save(data TickData) error {
	dataFilePath, err := getDataFilePath()
	if err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(dataFilePath, bytes, 0644)
}