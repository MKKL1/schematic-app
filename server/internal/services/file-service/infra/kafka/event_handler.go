package kafka

import (
	"log"

	"github.com/ThreeDotsLabs/watermill/message"
)

func FileUploadedHandler(msg *message.Message) {
	log.Printf("Received FileUploaded event: %s", string(msg.Payload))
}

func FileDeletedHandler(msg *message.Message) {
	log.Printf("Received FileDeleted event: %s", string(msg.Payload))
}
