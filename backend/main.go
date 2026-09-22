package main

import (
	"groundTurn/src/config"
	"groundTurn/src/routes"
)

func main() {
	cfg := config.Load()
	routes.Start(":" + cfg.Port)
}
