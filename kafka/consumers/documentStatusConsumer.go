package consumers

import (
	"context"
	"github.com/sirupsen/logrus"
	_ "gopkg.in/confluentinc/confluent-kafka-go.v1/kafka/librdkafka_vendor"
	"msrd-products/db"
	intrnalKafka "msrd-products/kafka"
	kafkaModels "msrd-products/kafka/documentStatusModels"
	"msrd-products/logic"
	"msrd-products/models"
	"time"
)

func LaunchDocumentStatusConsumer(dbContext db.DbContext) {

	err := intrnalKafka.Subscribe[kafkaModels.PostgreSqlEvent]("MsrdStocks.public.document_statuses", func(message kafkaModels.PostgreSqlEvent) bool {
		if message.Operation == "r" || message.Operation == "c" || message.Operation == "u" {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			docRep := logic.NewDocumentsRepository(ctx, dbContext)
			document, err := docRep.FindById(message.After.DocumentId)

			if err != nil {
				return false
			}

			if document == nil {
				logrus.Warnln("Received non existing document id from MsrdStocks.public.document_statuses: %s", message.After.DocumentId)
				return false
			}

			document, err = docRep.UpdateByEvent(models.UpdateDocumentEvent{
				Id:     document.Id,
				Status: message.After.Status,
			})

			if err != nil {
				return false
			}

			logrus.Infoln("Updated the status of %s to %f", message.After.DocumentId, message.After.Status)
			logrus.Infoln(document)
		}
		return true
	})

	if err != nil {
		logrus.Errorln("Error in ProductStockRecordsConsumer: %s", err)
		return
	}
}
