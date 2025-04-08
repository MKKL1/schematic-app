package metrics

import "time"

// Client is an interface for recording metrics.
type Client interface {
	// RecordCommandDuration records the execution time of a command handler.
	RecordCommandDuration(commandName string, success bool, duration time.Duration)
	RecordQueryDuration(queryName string, success bool, duration time.Duration)
}

func getStatusLabel(success bool) string {
	if success {
		return "success"
	}
	return "failure"
}
