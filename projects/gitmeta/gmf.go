package main

import (
	"errors"
	"os"
)

type Cloneable struct {
	URL string
	Dir string
}

type Gitmeta struct {
	Cloneables []Cloneable
}

type PluginIf interface {
	Applicable(e GitMetaEntry) bool
	Expand(e GitMetaEntry) ([]Cloneable, error)
}

func (g *Gitmeta) FindPlugin(m GitMetaEntry) (PluginIf, error) {
	plugins := []PluginIf{
		NewGMFEntry(),
		NewGMFGithubUser(),
		NewGMFSSH(),
	}
	for _, plugin := range plugins {
		if plugin.Applicable(m) {
			return plugin, nil
		}
	}
	return nil, errors.New("No applicable plugin found")
}

func NewGitmeta() *Gitmeta {
	return &Gitmeta{
		Cloneables: []Cloneable{},
	}
}

func (g *Gitmeta) AddGMF(f *os.File) error {
	entries, err := g.parseGMF(f)
	if err != nil {
		return err
	}

	for _, e := range entries {
		p, err := g.FindPlugin(e)
		if err != nil {
			return err
		}
		cloneables, err := p.Expand(e)
		if err != nil {
			return err
		}
		for _, cloneable := range cloneables {
			g.Cloneables = append(g.Cloneables, cloneable)
		}
	}
	return nil
}

func (g *Gitmeta) AllCloneables() []Cloneable {
	return g.Cloneables
}
