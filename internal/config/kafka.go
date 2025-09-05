package config

import (
	"fmt"
	"os"
	"strings"
)

const (
	kafkaAddressesEnv = "KAFKA_ADDRESSES"
	kafkaThemeEnv     = "THEME"
)

// KafkaConfig config для кафки
type KafkaConfig interface {
	Addresses() []string
	Theme() string
}

type kafkaConfig struct {
	addresses []string
	theme     string
}

func NewKafkaConfig() (KafkaConfig, error) {
	addresses := os.Getenv(kafkaAddressesEnv)
	if len(addresses) == 0 {
		return nil, fmt.Errorf("environment variable %s not set", kafkaAddressesEnv)
	}

	SliceAddresses := strings.Split(addresses, ",")

	theme := os.Getenv(kafkaThemeEnv)
	if len(theme) == 0 {
		return nil, fmt.Errorf("environment variable %s not set", kafkaThemeEnv)
	}

	return &kafkaConfig{
		addresses: SliceAddresses,
		theme:     theme,
	}, nil
}

func (k *kafkaConfig) Addresses() []string {
	return k.addresses
}

func (k *kafkaConfig) Theme() string {
	return k.theme
}
