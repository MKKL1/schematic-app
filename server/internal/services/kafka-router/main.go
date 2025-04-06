package main

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/bytedance/sonic"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lovoo/goka"
)

/*

Routes events.FileCreated event based on it's content type (image, minecraft-schematic,...) to different topics.

*/

// Define the broker addresses and streams.
var (
	brokers                             = []string{"localhost:9092"}
	fileUploadedStream      goka.Stream = "events.FileCreated"
	imageUploadedStream     goka.Stream = "events.ImageUploaded"
	schematicUploadedStream goka.Stream = "events.SchematicUploaded"
	group                   goka.Group  = "file-uploaded-group"
)

// FileUploaded represents the incoming event structure.
type FileUploaded struct {
	TempID   string            `json:"temp_id"`
	PermID   string            `json:"perm_id"`
	Existed  bool              `json:"existed"`
	Type     string            `json:"type"`
	Metadata map[string]string `json:"metadata"`
}

// FileUploadedCodec is a custom codec for the FileUploaded type.
type FileUploadedCodec struct{}

func (c *FileUploadedCodec) Encode(value interface{}) ([]byte, error) {
	return sonic.Marshal(value)
}

func (c *FileUploadedCodec) Decode(data []byte) (interface{}, error) {
	var fu FileUploaded
	if err := sonic.Unmarshal(data, &fu); err != nil {
		return nil, err
	}
	return &fu, nil
}

// process is called for each message arriving on the "file-uploaded" stream.
// It inspects the MIME type and emits an event to the appropriate output stream.
func process(ctx goka.Context, msg interface{}) {
	fu, ok := msg.(*FileUploaded)
	if !ok {
		log.Printf("unexpected message type %T", msg)
		return
	}

	switch fu.Type {
	case "image":
		ctx.Emit(imageUploadedStream, ctx.Key(), fu)
		log.Printf("Emitted %s event for key %s", imageUploadedStream, ctx.Key())
	case "minecraft-schematic":
		ctx.Emit(schematicUploadedStream, ctx.Key(), fu)
		log.Printf("Emitted %s event for key %s", schematicUploadedStream, ctx.Key())
	default:
		log.Printf("unsupported type: %s", fu.Type)
	}
}

func main() {
	// Define the processor group:
	// - Input from "file-uploaded" using our custom codec and process callback.
	// - Outputs to "image-uploaded" and "schematic-uploaded" using the same codec.
	g := goka.DefineGroup(group,
		goka.Input(fileUploadedStream, new(FileUploadedCodec), process),
		goka.Output(imageUploadedStream, new(FileUploadedCodec)),
		goka.Output(schematicUploadedStream, new(FileUploadedCodec)),
	)

	config := sarama.NewConfig()
	// If no offset is committed, start from the oldest message.
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	// Create a new processor.
	processor, err := goka.NewProcessor(brokers, g, goka.WithConsumerSaramaBuilder(goka.SaramaConsumerBuilderWithConfig(config)))
	if err != nil {
		log.Fatalf("Error creating processor: %v", err)
	}

	// Use a cancellable context and run the processor in a separate goroutine.
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error)
	go func() {
		done <- processor.Run(ctx)
	}()

	// Wait for an OS signal (SIGINT/SIGTERM) for graceful shutdown.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	log.Printf("Received signal: %v, shutting down...", sig)
	cancel()
	if err = <-done; err != nil {
		log.Fatalf("Processor shutdown with error: %v", err)
	}
	log.Println("Processor shutdown cleanly")
}
