package kafka

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde/avro"
	"github.com/sirupsen/logrus"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"os"
	"os/signal"
	"syscall"
)

func Subscribe[T interface{}](topic string, onMessage func(message T) bool) error {
	var bootstrapServers = os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
	var group = os.Getenv("KAFKA_CONSUMER_GROUP")
	var schmaregistryUrl = os.Getenv("KAFKA_SCHEMAREGISTRY_CLIENT")

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  bootstrapServers,
		"group.id":           group,
		"session.timeout.ms": 6000,
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false,
	})

	defer c.Close()
	if err != nil {
		logrus.Errorln(os.Stderr, "Failed to create consumer: %s\n", err)
		return err
	}

	logrus.Infoln("Created Consumer %v\n", c)
	client, err := schemaregistry.NewClient(schemaregistry.NewConfig(schmaregistryUrl))

	if err != nil {
		logrus.Errorln("Failed to create schema registry client: %s\n", err)
		return err
	}

	deser, err := avro.NewGenericDeserializer(client, serde.ValueSerde, avro.NewDeserializerConfig())
	if err != nil {
		logrus.Errorln("Failed to create deserializer: %s\n", err)
		return err
	}

	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		logrus.Errorln("Failed to subscribe topic: %s\n", err)
		return err
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for true {
		select {
		case sig := <-sigchan:
			logrus.Infoln("Caught signal %v: terminating\n", sig)
			return nil
		default:
			ev := c.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:

				var value T
				err := deser.DeserializeInto(*e.TopicPartition.Topic, e.Value, &value)
				if err != nil {
					logrus.Errorln("Failed to deserialize payload: %s\n", err)
				} else {
					logrus.Infoln("%% Message on %s:\n%+v\n", e.TopicPartition, value)
					if onMessage(value) {
						go c.Commit()
						logrus.Infoln("%% Committed on %s:\n%+v\n", e.TopicPartition, value)
					}

					logrus.Warnln("%% Commit cancelled on %s:\n%+v\n", e.TopicPartition, value)
				}
				if e.Headers != nil {
					logrus.Infoln("%% Headers: %v\n", e.Headers)
				}
			case kafka.Error:
				logrus.Errorln("%% Error: %v: %v\n", e.Code(), e)
			default:
				logrus.Infoln("Ignored %v\n", e)
			}
		}
	}

	fmt.Printf("Closing consumer\n")
	return nil
}
