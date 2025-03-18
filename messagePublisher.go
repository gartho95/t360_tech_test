package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
)

func publishMessage(topicTitle string, pubSubPacketBody SearchResponse) {

	if googleCloudProjectID == "" {
		jsonData, err := json.MarshalIndent(pubSubPacketBody, "", "  ")
		if err != nil {
			log.Printf("Error marshalling JSON: %v", err)
			return
		}
		log.Printf("%s %s", "Publishing body to pubsub with data: ", string(jsonData))
		return
	}

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT_ID"))
	if err != nil {
		log.Fatalf("Failed to create Pub/Sub client: %v", err)
	}
	defer client.Close()

	topic := client.Topic(topicTitle)

	jsonData, err := json.Marshal(pubSubPacketBody)
	if err != nil {
		log.Fatalf("Error marshalling message to JSON: %v", err)
	}
	result := topic.Publish(ctx, &pubsub.Message{
		Data: []byte(jsonData),
	})

	messageID, err := result.Get(ctx)
	if err != nil {
		log.Fatalf("Failed to publish message: %v", err)
	}

	log.Printf("Published message with ID: %s", messageID)
}
