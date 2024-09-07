package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type GitMetaEntry struct {
	URL  string `yaml:"url"`
	Type string `yaml:"type"`
	Dir  string `yaml:"dir"`
	User string `yaml:"user"`
	Path string `yaml:"path"`
	Host string `yaml:"host"`
}

func (g *Gitmeta) parseGMF(f *os.File) ([]GitMetaEntry, error) {
	records := []GitMetaEntry{}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return records, err
	}

	err = yaml.Unmarshal(data, &records)
	if err != nil {
		return records, err
	}

	return records, nil
}
