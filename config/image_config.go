package config

import (
	"fmt"
	"os"
)

type ImageConfig struct {
	ImagePath string
}

func validateImageConfig() error {
	if os.Getenv("IMAGE_PATH") == "" {
		return fmt.Errorf("IMAGE_PATH config is not set")
	}
	return nil
}

func LoadImageConfig() (*ImageConfig, error) {
	err := validateImageConfig()

	if err != nil {
		return nil, err
	}

	return &ImageConfig{
		ImagePath: os.Getenv("IMAGE_PATH"),
	}, err
}
