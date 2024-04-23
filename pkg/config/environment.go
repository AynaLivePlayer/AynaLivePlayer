package config

import (
	"os"
)

type EnvType string

const (
	Development EnvType = "development"
	Production  EnvType = "production"
)

func CurrentEnvironment() EnvType {
	t := EnvType(os.Getenv("ENV"))
	if t == Production {
		return Production
	}
	return Development
}
