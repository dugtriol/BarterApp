package main

import "github.com/dugtriol/BarterApp/internal/app"

const configPath = "config/config.yaml"

func main() {
	app.Run(configPath)
}
