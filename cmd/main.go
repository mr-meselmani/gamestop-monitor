package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"gamestop/internal"
)

func main() {
	monitor, err := internal.NewGamestopHandler("gamestop")

	if err != nil {
		log.Fatalf("can not boot monitor: %s", err.Error())
	}

	runCount := 0

	for {
		runCount++
		statusCode, err := monitor.Collect(context.Background(), "https://www.gamestop.de/SearchResult/QuickSearchAjax?platform=101&rootGenre=170&typeSorting=11&sDirection=Descending")

		if err != nil {
			log.Fatalf("failed to collect monitoring information: %s", err.Error())
		}

		if statusCode != http.StatusOK {
			log.Fatalf("unexpected status code received: %d", statusCode)
		}

		messages, err := monitor.Evaluate(context.Background())

		if err != nil {
			log.Fatalf("failed to evaluate monitoring information: %s", err.Error())
		}

		if len(messages) == 0 {
			log.Println("no messages received to ping from monitor for this run")
		}

		for index, message := range messages {
			log.Println(fmt.Sprintf("Index %d: Message: %v", index, message))
		}

		log.Println(fmt.Sprintf("monitor run %d done", runCount))

		if runCount >= 10 {
			log.Println(fmt.Sprintf("monitor done"))
			return
		}

		time.Sleep(5 * time.Second)
	}
}
