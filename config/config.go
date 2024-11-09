package config

import (
	"log"
	"os"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

// Структура для хранения конфигурации
type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"database"`

	Server struct {
		Port    int `yaml:"port"`
		Timeout int `yaml:"timeout"`
	} `yaml:"server"`

	JWT struct {
		Secret     string        `yaml:"secret"`
		Expiration time.Duration `yaml:"expiration"`
	} `yaml:"jwt"`
}

// Объявляем переменные для Singleton
var (
	configInstance *Config
	once           sync.Once
)

// GetConfig предоставляет доступ к загруженной конфигурации
func GetConfig() *Config {
	return configInstance
}

func LoadConfig(filename string) *Config {
	once.Do(func() { // Используем sync.Once для загрузки только один раз
		file, err := os.Open(filename)
		if err != nil {
			log.Fatalf("Ошибка открытия конфигурационного файла: %v", err)
		}
		defer file.Close()

		decoder := yaml.NewDecoder(file)
		configInstance = &Config{}
		err = decoder.Decode(configInstance)
		if err != nil {
			log.Fatalf("Ошибка декодирования конфигурационного файла: %v", err)
		}

		// Преобразование JWT Expiration в Duration
		configInstance.JWT.Expiration, err = time.ParseDuration(configInstance.JWT.Expiration.String())
		if err != nil {
			log.Fatalf("Ошибка парсинга JWT Expiration: %v", err)
		}
	})
	return configInstance
}
