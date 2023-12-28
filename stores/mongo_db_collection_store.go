package stores

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
)

type MongoDbCollectionStore struct {
	database        string
	collection      string
	client          *mongo.Client
	collectionStore *mongo.Collection
	context         context.Context
	docKKey         string
}

func (receiver MongoDbCollectionStore) All(results interface{}) error {
	return receiver.Get(0, OrderBy{Column: "_id", Order: 1}, results)
}

func (receiver MongoDbCollectionStore) Get(limit int64, orderBy OrderBy, results interface{}) error {
	findOptions := options.Find().SetLimit(limit)
	if len(orderBy.Column) > 0 {
		findOptions.SetSort(bson.D{{orderBy.Column, orderBy.Order}})
	}
	cursor, err := receiver.collectionStore.Find(receiver.context, bson.D{}, findOptions)
	if err != nil {
		return err
	}

	return cursor.All(receiver.context, results)
}

func (receiver MongoDbCollectionStore) Count() (int64, error) {
	return receiver.collectionStore.CountDocuments(receiver.context, bson.D{})
}

func (receiver MongoDbCollectionStore) Save(element interface{}) error {
	return receiver.SaveAll([]interface{}{element})
}

func (receiver MongoDbCollectionStore) SaveAll(elements []interface{}) error {
	slog.Debug("storing into mongo")
	opts := options.InsertMany().SetOrdered(false)
	_, err := receiver.collectionStore.InsertMany(receiver.context, elements, opts)
	return handleError(err)
}

func handleError(err error) error {
	var bulkWriteException mongo.BulkWriteException
	ok := errors.As(err, &bulkWriteException)
	if ok == true {
		var nonDuplicateErrors = make([]mongo.BulkWriteError, 0)
		for _, writeError := range bulkWriteException.WriteErrors {
			if writeError.Code != 11000 {
				nonDuplicateErrors = append(nonDuplicateErrors, writeError)
			}
		}
		if len(nonDuplicateErrors) > 0 {
			return mongo.BulkWriteException{WriteErrors: nonDuplicateErrors}
		} else {
			slog.Warn("duplicated items omitted")
			return nil
		}
	} else {
		return err
	}
}

func (receiver MongoDbCollectionStore) UpdateAllById(elements map[interface{}]interface{}) error {
	models := make([]mongo.WriteModel, 0)
	for id, element := range elements {
		models = append(models, mongo.NewUpdateOneModel().
			SetUpsert(true).
			SetFilter(bson.D{{"_id", id}}).
			SetUpdate(bson.D{{"$set", element}}))
	}

	_, err := receiver.collectionStore.BulkWrite(receiver.context, models, options.BulkWrite().SetOrdered(false))
	return err
}

func (receiver MongoDbCollectionStore) Close() error {
	return receiver.client.Disconnect(receiver.context)
}

func NewMongoDbStore(url string, port string, database string, collection string) (*MongoDbCollectionStore, error) {
	uri := fmt.Sprintf("mongodb://%s:%s", url, port)
	slog.Info(fmt.Sprintf("connecting to mongo: %s, database: %s, collection: %s", uri, database, collection))
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return &MongoDbCollectionStore{
		database:        database,
		collection:      collection,
		client:          client,
		collectionStore: client.Database(database).Collection(collection),
		context:         context.Background(),
		docKKey:         collection}, nil
}
