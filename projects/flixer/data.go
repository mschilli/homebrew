/////////////////////////////////////////
// data.go - Flixer Data Model
// Mike Schilli, 2026 (m@perlmeister.com)
/////////////////////////////////////////
package main

import (
  "gopkg.in/yaml.v2"
  "os"
  "sort"
)

type Pick struct {
  Date   string `yaml:"date"`
  Title  string `yaml:"title"`
  URL    string `yaml:"url"`
  Rating int    `yaml:"rating"`
}

type Picker struct {
  YAMLPath string
}

func NewPicker() *Picker {
  return &Picker{
    YAMLPath: "flixer.yaml",
  }
}

func (p *Picker) Load() []Pick {
  picks := []Pick{}

  b, err := os.ReadFile(p.YAMLPath)
  if err != nil {
    panic(err)
  }

  err = yaml.Unmarshal(b, &picks)
  if err != nil {
    panic(err)
  }

  p.Sort(picks)

  return picks
}

func (p *Picker) Sort(picks []Pick) {
  sort.Slice(picks, func(i, j int) bool {
    ri, rj := picks[i].Rating, picks[j].Rating

    // unrated first
    if ri == 0 && rj != 0 {
      return true
    }
    if rj == 0 && ri != 0 {
      return false
    }

    // high rank first
    return ri > rj
  })
}

func (p *Picker) Save(picks []Pick) {
  b, err := yaml.Marshal(picks)
  if err != nil {
    panic(err)
  }

  err = os.WriteFile(p.YAMLPath, b, 0644)
  if err != nil {
    panic(err)
  }
}
