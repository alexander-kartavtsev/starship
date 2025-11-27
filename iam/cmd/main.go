package main

import (
	"fmt"
	"github.com/alexander-kartavtsev/starship/iam/internal/config"
)

const envPath = "deploy/compose/iam/.env"

func main() {
	err := config.Load(envPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

}
