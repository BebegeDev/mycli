package pipeline

type Deploy struct {
	Component Component `yaml:"component"`
}

type Component struct {
	Server_list []string `yaml:"serverList"`
	Login       string   `yaml:"login"`
	Password    string   `yaml:"password"`
	Steps       []Step   `yaml:"steps"`
}

type Step struct {
	Name   string                 `yaml:"name"`
	Params map[string]interface{} `yaml:"params"`
}
