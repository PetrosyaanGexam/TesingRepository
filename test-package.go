package TesingRepository

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"runtime"
	"strconv"
	"time"
)

func ShowMemoryUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	print("Alloc = ", m.Alloc, " b\t")
	print("TotalAlloc = ", m.TotalAlloc, " b\t")
	print("Sys = ", m.Sys, " b\t")
	print("HeaoInUS = ", m.HeapInuse, " b\t")
	print("NumGC = ", m.NumGC, "\n")
	print("Number of existing coroutines : ", strconv.Itoa(runtime.NumGoroutine()), "\n")
	//time.Sleep(time.Duration(30) * time.Second)
}

type TestModel1 struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

func (obj *TestModel1) GetCollectionName() string {
	return "test1"
}

func (obj *TestModel1) GetCollection(db *mongo.Database) mongo.Collection {
	return *db.Collection(obj.GetCollectionName())
}

func (obj *TestModel1) GetById(id string) (*TestModel1, error) {
	var db mongo.Database
	var client mongo.Client
	var clientOptions options.ClientOptions
	clientOptions = *options.Client().ApplyURI("mongodb://" + "localhost" + ":" + "27017")
	clientPtr, connection_err := mongo.Connect(context.TODO(), &clientOptions)
	if connection_err != nil {
		panic(connection_err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	if connection_err == nil {
		client = *clientPtr

		db = *client.Database("MongoDbConnectionMemLeek")
		err := client.Ping(context.TODO(), nil)
		if err != nil {
			println(err)
		}
	}
	//println("Connected to MongoDB!")

	idc, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": idc}
	collection := obj.GetCollection(&db)
	err := collection.FindOne(context.TODO(), filter).Decode(&obj)
	if err == nil {
		return obj, nil
	}
	return &TestModel1{}, errors.New("model - 1 no found")

}

func (obj *TestModel1) Insert(db *mongo.Database) error {
	//var conn *driver.Connection
	//var db *mongo.Database
	//var client *mongo.Client
	//var clientOptions *options.ClientOptions
	//clientOptions = options.Client().ApplyURI("mongodb://" + "localhost" + ":" + "27017")
	//ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	//defer cancel()
	//client, connection_err := mongo.Connect(ctx, clientOptions)
	//if connection_err != nil {
	//	panic(connection_err)
	//}
	//defer func() {
	//	if err := client.Disconnect(context.TODO()); err != nil {
	//		panic(err)
	//	}
	//}()
	//if connection_err == nil {
	//	db = client.Database("MongoDbConnectionMemLeek")
	//	err := client.Ping(context.TODO(), nil)
	//	if err != nil {
	//		println(err)
	//	}
	//println("Connected to MongoDB!")

	collection := obj.GetCollection(db)
	_, err := collection.InsertOne(context.TODO(), *obj)
	if err != nil {
		println(errors.New("model - 1 error on inserting"))
		return errors.New("model - 1 error on inserting")
	}
	//}
	return nil

	/*
		ip := "localhost"
		port := "27017"
		opts := options.Client()
		opts.SetDirect(true)
		opts.SetServerSelectionTimeout(1 * time.Second)
		opts.SetConnectTimeout(2 * time.Second)
		opts.SetSocketTimeout(2 * time.Second)
		opts.SetMaxConnIdleTime(1 * time.Second)
		opts.SetMaxPoolSize(1)
		url := fmt.Sprintf("mongodb://%s:%s/admin", ip, port)
		opts.ApplyURI(url)
		ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
		conn, err := mongo.Connect(ctx, opts)
		if err != nil {
			fmt.Printf("new %s:%s mongo connection error: %v\n", ip, port, err)
			return nil
		}
		defer conn.Disconnect(ctx)
		err = conn.Ping(ctx, nil)
		if err != nil {
			fmt.Printf("ping %s:%s ping error: %v\n", ip, port, err)
			return nil
		}
		db := conn.Database("MongoDbConnectionMemLeek")

		collection := obj.GetCollection(*db)
		_, err = collection.InsertOne(context.TODO(), *obj)
		if err != nil {
			println(errors.New("model - 1 error on inserting"))
			return errors.New("model - 1 error on inserting")
		}

		return nil
	*/
}
