package main

import (
	"fmt"

	"github.com/TestardR/cron_scheduler/internal/handler"
	"github.com/TestardR/cron_scheduler/internal/logger"
)

const appName = "cron_scheduler"

func main() {
	log := logger.New(appName)

	h, err := handler.New(log)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to initialize handler: %s", err))
	}

	if err := h.Start(); err != nil {
		log.Fatal(fmt.Sprintf("error occurred running the handler: %s", err))
	}
}
