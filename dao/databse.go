package dao

type Database struct {
	Type     string `yaml:"type"`
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
	Protocol string `yaml:"protocol"`
	IP       string `yaml:"ip"`
	Port     uint   `yaml:"port"`
	Library  string `yaml:"library"`
	CharSet  string `yaml:"char_set"`
}
