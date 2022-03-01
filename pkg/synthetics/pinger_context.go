package synthetics

import (
	"context"
	"time"
)

type PingerContext struct {
	context.Context
	PingInterval time.Duration
	URLs MultipleUrls
	BindAddr string
}