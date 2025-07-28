package pipeline

/*
Постепенная реализация скелета для YAML файла
Пока стремлюсь к этому структура следующая:
===========================================================|
deploy:
	componen:
		server_list: [] (обдумаю добавление список файлом)
		login: user
		password: ****
		steps:
			backup:
				src: "/srv/app"
				dst: "/archive"
				force: true
				owerwrite: true
				addDate: true
			build:
				src: "/new_build.zip"
				dst: "/srv/app"
			restart:
				service: []
				pool: []
===========================================================|

*/

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
