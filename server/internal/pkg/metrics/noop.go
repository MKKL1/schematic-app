package metrics

import (
	"time"
)

type NoOp struct{}

func (n NoOp) RecordCommandDuration(commandName string, success bool, duration time.Duration) {
}

func (n NoOp) RecordQueryDuration(queryName string, success bool, duration time.Duration) {
}

var _ Client = (*NoOp)(nil)
