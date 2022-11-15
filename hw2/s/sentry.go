package s

import (
	"log"

	"github.com/getsentry/sentry-go"
)

type SentryLogger struct {
	//l sentry
}

func NewSentryLogger() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "http://fe2abbd35ca84918990d6aa60a2315ee@localhost:9000/2",
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	log.Println("sentry is initialized sucessfully")
}
