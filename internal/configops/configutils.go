package configops

import (
	"fmt"

	"github.com/spf13/viper"
)

// Поверяем флаг конфига
func ConfigExists(configPath string) error {
	viper.SetConfigFile(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("ошибка при чтении конфига: %w", err)
	}
	return nil
}

// Считываем конфиг в структуру
func ConfigRead(section string, out interface{}) error {
	err := viper.UnmarshalKey(section, out)
	if err != nil {
		return fmt.Errorf("ошибка в секции %s: %w", section, err)
	}
	return nil
}
