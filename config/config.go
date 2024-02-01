package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
)

type Config struct {
	Db    *DatabaseConfig
	Cache *CacheConfig
	Image *ImageConfig
}

func NewConfig(envFile string) *Config {
	err := godotenv.Load(envFile)
	if err != nil {
		fmt.Println("Failed to load .env from file. trying to load .env from current directory")

		executable, errorLoadExecPath := os.Executable()
		if errorLoadExecPath != nil {
			panic(errorLoadExecPath)
		}

		exPath := filepath.Dir(executable)
		fmt.Println(exPath + "/.env")
		errorLoadExPath := godotenv.Load(exPath + "/.env")
		if errorLoadExPath != nil {
			panic(errorLoadExPath)
		}
	}

	dbConfig, errDb := LoadDatabaseConfig()
	if errDb != nil {
		panic(errDb)
	}

	cacheConfig, errCache := LoadCacheConfig()
	if errCache != nil {
		panic(errCache)
	}

	imageConfig, errImage := LoadImageConfig()
	if errImage != nil {
		panic(errImage)
	}

	return &Config{
		Db:    dbConfig,
		Cache: cacheConfig,
		Image: imageConfig,
	}
}
