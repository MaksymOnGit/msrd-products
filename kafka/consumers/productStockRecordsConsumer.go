package consumers

import (
	"context"
	"github.com/sirupsen/logrus"
	_ "gopkg.in/confluentinc/confluent-kafka-go.v1/kafka/librdkafka_vendor"
	"msrd-products/db"
	intrnalKafka "msrd-products/kafka"
	kafkaModels "msrd-products/kafka/models"
	"msrd-products/logic"
	"msrd-products/models"
	"time"
)

func LaunchProductStockRecordsConsumer(dbContext db.DbContext) {

	err := intrnalKafka.Subscribe[kafkaModels.PostgreSqlEvent]("MsrdStocks.public.stock_records", func(message kafkaModels.PostgreSqlEvent) bool {
		if message.Operation == "r" || message.Operation == "c" {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			prodRep := logic.NewProductsRepository(ctx, dbContext)
			product, err := prodRep.FindById(message.After.ProductId)

			if err != nil {
				return false
			}

			if product == nil {
				logrus.Warnln("Received non existing product id from MsrdStocks.public.stock_records: %s", message.After.ProductId)
				return false
			}

			if &message.After.ActualQuantity == product.Quantity {
				logrus.Warnln("Nothing to update MsrdStocks.public.stock_records: %s", message.After.ProductId)
				return false
			}

			product.Quantity = new(float32)
			*product.Quantity = *&(message.After.ActualQuantity)
			product, err = prodRep.UpdateByEvent(models.UpdateProductEvent{
				Id:       product.Id,
				Quantity: product.Quantity,
			})

			if err != nil {
				return false
			}

			logrus.Infoln("Updated the stock quantity of %s to %f", message.After.ProductId, message.After.ActualQuantity)
			logrus.Infoln(product)
		}
		return true
	})

	if err != nil {
		logrus.Errorln("Error in ProductStockRecordsConsumer: %s", err)
		return
	}
}
