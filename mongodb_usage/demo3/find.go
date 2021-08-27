package main

import (
	"context"
	"fmt"
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

type FindByJobName struct {
	JobName string `bson:"jobName"`	//任务名查询条件
}

func main() {
	var (
		client *mongo.Client
		err error
		ctx context.Context
		cancelFunc context.CancelFunc
		database *mongo.Database
		collection *mongo.Collection
		record *LogRecord
		cond *FindByJobName
		logArr []interface{}
		skip int64
		limit int64
		ops *options.FindOptions
		cursor *mongo.Cursor
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
	cond = &FindByJobName{JobName: "彭国朝1"}
	skip = 0
	limit = 10
	ops = &options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
	}
	if cursor, err = collection.Find(context.TODO(), cond, ops); err != nil {
		fmt.Println(err)
		return
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		record = &LogRecord{}
		if err = cursor.Decode(record); err != nil {
			fmt.Println(err)
			return
		}
		//打印查询到的数据
		fmt.Println(*record)
		logArr = append(logArr, *record)
	}
	fmt.Println("--------------------")
	fmt.Println(logArr)
}
