package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

type CaseID struct {
	// Format is the string defining the format of the case ID.
	// It may contain any characters except for the placeholders:
	// - $i: the incremental number
	// - $y: the current year as a 4-digit number
	// - $m: the current month as a 2-digit number
	// - $d: the current day as a 2-digit number
	Format      string  `json:"format"`
	Name        string  `json:"name"`
	Incremental int     `json:"incremental"`
	storage     Storage `json:"-"`
}

func (c *CaseID) Next() string {
	c.Incremental++
	c.storage.Save()

	entry := c.Format
	entry = strings.ReplaceAll(entry, "$i", fmt.Sprintf("%02d", c.Incremental))
	entry = strings.ReplaceAll(entry, "$y", fmt.Sprintf("%d", time.Now().Year()))
	entry = strings.ReplaceAll(entry, "$m", fmt.Sprintf("%02d", time.Now().Month()))
	entry = strings.ReplaceAll(entry, "$d", fmt.Sprintf("%02d", time.Now().Day()))

	return entry
}

type Storage map[string]*CaseID

func LoadStorage() Storage {
	appdataPath, err := os.UserCacheDir()
	if err != nil {
		log.Println("Error getting appdata path: ", err)
		return nil
	}

	filePath := path.Join(appdataPath, "caseid.json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return Storage{}
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Println("Error opening caseid.json: ", err)
		return nil
	}

	defer file.Close()

	var storage Storage
	if err := json.NewDecoder(file).Decode(&storage); err != nil {
		log.Println("Error decoding caseid.json: ", err)
		return Storage{}
	}

	for _, c := range storage {
		c.storage = storage
	}

	return storage
}

func (s Storage) Create(name, format string) *CaseID {
	s[name] = &CaseID{Name: name, Format: format, Incremental: 0, storage: s}
	s.Save()

	return s[name]
}

func (s Storage) Remove(name string) {
	delete(s, name)
	s.Save()
}

func (s Storage) Save() {
	appdataPath, err := os.UserCacheDir()
	if err != nil {
		log.Println("Error getting appdata path: ", err)
		return
	}

	filePath := path.Join(appdataPath, "caseid.json")
	file, err := os.Create(filePath)
	if err != nil {
		log.Println("Error creating caseid.json: ", err)
		return
	}

	defer file.Close()

	if err := json.NewEncoder(file).Encode(s); err != nil {
		log.Println("Error saving caseid.json: ", err)
	}
}
