package main

type MaintainerInfo struct {
	Name  string `yaml:"name"`
	Email string `yaml:"email" binding:"required,email"`
}

type Metadata struct {
	Title       string           `yaml:"title"`
	Version     string           `yaml:"version"`
	Company     string           `yaml:"company"`
	Website     string           `yaml:"website"`
	Source      string           `yaml:"source"`
	License     string           `yaml:"license"`
	Maintainers []MaintainerInfo `yaml:"maintainers" binding:"required"`
	Description string           `yaml:"description"`
}
