package s

import (
	"github.com/getsentry/sentry-go"
	"log"
)

type SentryLogger struct {
	//l sentry
}

func NewSentryLogger() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://e83ffb613ddd452d872870f10400c0d7@o1354434.ingest.sentry.io/6640945",
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
