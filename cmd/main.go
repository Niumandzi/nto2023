package main

import "github.com/niumandzi/nto2023/pkg/logging"

func main() {
	logging.Init()
	logger := logging.GetLogger()
	logger.Println("logger initialized")

}
