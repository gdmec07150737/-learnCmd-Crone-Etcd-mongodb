package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type TimePoint struct {
	StartTime int64 `bson:"startTime"`
	EndTime int64 `bson:"endTime"`
}

type LogRecord struct {
	JobName string `bson:"jobName"`	//任务名
	Command string `bson:"command"` //shell命令
	Err string `bson:"err"`	//脚本错误
	Content string `bson:"content"`	//脚本输出
	TimePoint TimePoint `bson:"timePoint"`	//执行时间
}

func main() {
	var (
		client *mongo.Client
		err error
		ctx context.Context
		cancelFunc context.CancelFunc
		database *mongo.Database
		collection *mongo.Collection
		insertResult *mongo.InsertOneResult
		record *LogRecord
		docId primitive.ObjectID
	)

	//1，建立链接
	ctx, cancelFunc = context.WithTimeout(context.TODO(), 5 * time.Second)
	defer cancelFunc()
	if client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017")); err != nil {
		fmt.Println(err)
		return
	}
	//2，选择数据库my_db
	database = client.Database("my_db")

	//3，选择表my_Collection
	collection = database.Collection("my_Collection")

	//测试数据库链接
	/*ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	fmt.Println(11111)
	fmt.Println(err)
	fmt.Println(22222)*/

	record = &LogRecord{
		JobName:   "pgc1",
		Command:   "echo hello",
		Err:       "",
		Content:   "hello",
		TimePoint: TimePoint{
			StartTime: time.Now().Unix(),
			EndTime:   time.Now().Unix() + 10,
		},
	}


	//使用结构体插入
	if insertResult, err = collection.InsertOne(context.TODO(), record); err != nil {
		fmt.Println(err)
		return
	}
	//拿着interface{}, 反射成objectID
	docId = insertResult.InsertedID.(primitive.ObjectID)
	fmt.Println("插入成功：",docId.Hex())


	/*//直接插入
	record = record
	if insertResult, err = collection.InsertOne(ctx, bson.D{{"name", "pi"}, {"value", 3.14159}}); err != nil {
		fmt.Println(err)
		return
	}
	docId = insertResult.InsertedID.(primitive.ObjectID)
	fmt.Println("插入成功：",docId.Hex())*/


	/*//批量插入
	record = record
	insertResult = insertResult
	docId = docId
	docs := []interface{}{
		bson.D{{"name", "Alice233"}},
		bson.D{{"name", "Bob233"}},
	}
	opts := options.InsertMany().SetOrdered(false)
	res, err := collection.InsertMany(context.TODO(), docs, opts)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("inserted documents with IDs %v\n", res.InsertedIDs)*/


}