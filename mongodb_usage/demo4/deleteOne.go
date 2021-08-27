package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

//{"$lte":timestamp}
type TimeBeforeCond struct {
	Before int64 `bson:"$lte"`
}

//{"timePoint.startTime":{$lte:time.Now()}}
type DeleteCond struct {
	BeforeCond TimeBeforeCond `bson:"timePoint.startTime"`
}

func main() {
	var (
		client *mongo.Client
		err error
		ctx context.Context
		cancelFunc context.CancelFunc
		database *mongo.Database
		collection *mongo.Collection
		deleteCond *DeleteCond
		deleteResult *mongo.DeleteResult
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
	deleteCond = &DeleteCond{
		BeforeCond: TimeBeforeCond{
			//Before: 1626978800,
			Before: time.Now().Unix(),
		},
	}
	/*deleteResult, err = collection.DeleteOne(context.TODO(), bson.D{{
		"timePoint.startTime",
		bson.D{{
			"$lte",
			//time.Now().Unix(),
			1626949200,
		}},
	}})*/
	if deleteResult, err = collection.DeleteOne(context.TODO(), deleteCond); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("成功删除：", deleteResult.DeletedCount, "条文档")
}