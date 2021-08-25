package main

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	deliveryChan := make(chan kafka.Event)
	producer := NewKafkaProducer()
	Publish("transfer", "teste", producer, []byte("transfer"), deliveryChan)

	go DeliveryReport(deliveryChan) // async

	producer.Flush(1000)
}

func NewKafkaProducer() *kafka.Producer {

	configMap := &kafka.ConfigMap{
		"bootstrap.servers":   "golang-kafka-example_kafka_1:9092",
		"delivery.timeout.ms": "0",
		"acks":                "all",
		"enable.idempotence":  "true",
	}

	producer, err := kafka.NewProducer(configMap)
	if err != nil {
		log.Println(err.Error())
	}

	return producer
}

func Publish(msg string, topic string, producer *kafka.Producer, key []byte, deliveryChannel chan kafka.Event) error {
	messsage := &kafka.Message{
		Value: []byte(msg),
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Key: key,
	}
	err := producer.Produce(messsage, deliveryChannel)
	return err
}

func DeliveryReport(deliveryChan chan kafka.Event) {
	for e := range deliveryChan {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				fmt.Println("Erro ao enviar")
			} else {
				fmt.Println("Mensagem enviada", ev.TopicPartition)
			}
		}
	}
}
