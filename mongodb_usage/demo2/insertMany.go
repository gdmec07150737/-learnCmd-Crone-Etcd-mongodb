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
		insertResult *mongo.InsertManyResult
		record *LogRecord
		logArr []interface{}
		insertId interface{}
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

	record = &LogRecord{
		JobName:   "彭国朝1",
		Command:   "echo 你好！",
		Err:       "",
		Content:   "你好！",
		TimePoint: TimePoint{
			StartTime: time.Now().Unix(),
			EndTime:   time.Now().Unix() + 10,
		},
	}

	logArr = []interface{}{record, record, record}
	//使用结构体批量插入
	if insertResult, err = collection.InsertMany(context.TODO(), logArr); err != nil {
		fmt.Println(err)
		return
	}
	//snowflake：毫秒/微妙的当前时间 + 机器ID + 当前毫秒/微妙内的自增ID（每当毫秒变化，会重置为0，继续自增）
	for _, insertId = range insertResult.InsertedIDs {
		//拿着interface{}, 反射成objectID
		docId = insertId.(primitive.ObjectID)
		fmt.Println("插入成功,自增ID：",docId.Hex())
	}
}