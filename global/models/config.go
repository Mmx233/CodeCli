package models

type Config struct {
	Default Default `yaml:"default"`
	Storage Storage `yaml:"storage"`
}

type Default struct {
	GitSite    string `yaml:"gitSite"`
	Username   string `yaml:"username"`
	CmdProgram string `yaml:"cmdProgram"`
}

type Storage struct {
	ProjectDir string `yaml:"projectDir"`
}
