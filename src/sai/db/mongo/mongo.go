package mongo

import (
	"context"
	"fmt"
	"github.com/saiset-co/saiStorageMongo/src/github.com/fatih/color"
	"github.com/saiset-co/saiStorageMongo/src/sai/common"
	"github.com/saiset-co/saiStorageMongo/src/sai_storage/settings"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

type Database struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type DB interface {
	GetSession() *mongo.Client
	GetDatabaseConfig() *Database
	GetCollection(collectionName string) *mongo.Collection
}

type LocalMongo struct {
	Config Database `json:"config"`
	Host   string   `json:"host"`
}

func (mi *LocalMongo) GetCollection(collectionName string) *mongo.Collection {
	session := Instance.GetSession()
	return session.Database(Instance.GetDatabaseConfig().Database).Collection(collectionName)
}

func (mi *LocalMongo) GetDatabaseConfig() *Database {
	return &mi.Config
}

type AtlasMongo struct {
	Config Database `json:"config"`
	Host   string   `json:"host"`
}

func (mi *AtlasMongo) GetCollection(collectionName string) *mongo.Collection {
	session := Instance.GetSession()
	return session.Database(Instance.GetDatabaseConfig().Database).Collection(collectionName)
}

func (mi *AtlasMongo) GetDatabaseConfig() *Database {
	return &mi.Config
}

func (mi *LocalMongo) GetSession() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(
		"mongodb://" + mi.Host + "/" + mi.Config.Database,
	))

	if err != nil {
		panic("Connection to the mongo server can not be established.")
	}

	err = client.Connect(context.TODO())

	if err != nil {
		panic("Connection to the mongo server can not be established.")
	}

	return client
}

func (mi *AtlasMongo) GetSession() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		"mongodb+srv://"+mi.Config.Username+":"+mi.Config.Password+"@"+mi.Host+"/"+mi.Config.Database+"?ssl=true&authSource=admin&retryWrites=true&w=majority",
	))

	if err != nil {
		panic("Connection to the mongo server can not be established.")
	}

	return client
}

var Instance DB

func SetMongoDBInv(mongoDB settings.DBConfig) {
	if mongoDB.Atlas.Enabled {
		Instance = &AtlasMongo{
			Config: Database{
				Username: mongoDB.Atlas.Config.Username,
				Password: mongoDB.Atlas.Config.Password,
				Database: mongoDB.Atlas.Config.Database,
			},
			Host: mongoDB.Atlas.Host,
		}
	} else {
		Instance = &LocalMongo{
			Config: Database{
				Username: mongoDB.Local.Config.Username,
				Password: mongoDB.Local.Config.Password,
				Database: mongoDB.Local.Config.Database,
			},
			Host: mongoDB.Local.Host,
		}
	}
}

func StartMongod() {
	d := color.New(color.FgMagenta, color.Bold)
	d.Println("--------------------------------------------")
	d.Println("Starting mongod...")
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		getMongodCmd := exec.Command("/bin/sh", "-c", "ps ax | grep mongod")
		mongods, _ := getMongodCmd.Output()
		fmt.Println(string(mongods))

		getPidCmd := exec.Command("/bin/sh", "-c", "ps aux | pgrep mongod")
		pid, _ := getPidCmd.Output()

		if string(pid) != "" {
			d.Print("Killing old mongod process ", string(pid))
			killPsCmd := exec.Command("/bin/sh", "-c", "kill -9 "+string(pid))
			killPsCmd.Run()
		}

		startMongoCmd := exec.Command("mongod", "--bind_ip_all")
		startMongoCmd.Start()
		d.Println("Mongod started. PID", startMongoCmd.Process.Pid)

		getStartedMongodCmd := exec.Command("/bin/sh", "-c", "ps ax | grep mongod")
		startedMongods, _ := getStartedMongodCmd.Output()
		fmt.Println(string(startedMongods))
	}
	// case "windows":
}

func TestMongoConnection() {
	err := Instance.GetSession().Ping(context.TODO(), nil)

	if err != nil {
		fmt.Println("Err: ", err)
	}

	//var result interface{}
	//var results []interface{}
	//var results2 []interface{}
	//
	//_ = Insert("test", map[string]interface{}{"field1": "value1"}, nil)
	//_ = Find("test", map[string]interface{}{"field": "value"}, nil, &results)
	//_ = Update("test", map[string]interface{}{"field1": "value1"}, map[string]interface{}{"$set": map[string]interface{}{"field1": 1}}, nil, nil)
	//_ = Update("test", map[string]interface{}{"field1": "value2"}, map[string]interface{}{"$set": map[string]interface{}{"field1": 2}}, nil, nil)
	//_ = FindOne("test", map[string]interface{}{"field1": "value2"}, &result)
	//fmt.Println("FindOne: ", result)
	//_ = Remove("test", map[string]interface{}{"field1": "value2"}, nil)
	//_ = Find("test", map[string]interface{}{"field1": map[string]interface{}{"$exists": true}}, nil, &results)
	//fmt.Println("Find: ", results)
	//
	//
	//_ = Find("test", map[string]interface{}{"field1": map[string]interface{}{"$exists": true}}, map[string]interface{}{"limit": "1", "skip": "2"}, &results2)
	//fmt.Println("Find: ", results2)

	fmt.Println("Mongo server reached!")
}

func FindOne(collectionName string, selector map[string]interface{}, result *map[string]interface{}) *common.Error {
	collection := Instance.GetCollection(collectionName)
	cur, err := collection.Find(context.TODO(), selector)

	if err != nil {
		return MongoDBError(err)
	}

	for cur.Next(context.TODO()) {
		var elem *map[string]interface{}
		err1 := cur.Decode(&elem)

		if err1 != nil {
			return MongoDBError(err1)
		}

		result = elem
		break
	}

	if err2 := cur.Err(); err2 != nil {
		return MongoDBError(err2)
	}

	_ = cur.Close(context.TODO())

	return nil
}

func Find(collectionName string, selector map[string]interface{}, inputOptions interface{}, result *[]interface{}) *common.Error {
	theOptionsLimit := make(map[string]int)
	theOptionsSkip := make(map[string]int)
	theOptionsSort := make(map[string]string)

	switch o := inputOptions.(type) {
	case map[string]interface{}:
		for s, b := range o {
			theK := string(s)
			if theK == "limit" {
				val, _ := strconv.Atoi(b.(string))
				theOptionsLimit[theK] = val
			}
			if theK == "sort" {
				val, _ := b.(string)
				theOptionsSort[theK] = val
			}
			if theK == "skip" {
				val, _ := strconv.Atoi(b.(string))
				theOptionsSkip[theK] = val
			}
		}
		break
	}

	requestOptions := options.Find()

	if sort, sortExists := theOptionsSort["sort"]; sortExists {
		requestOptions.SetSort(sort)
	}

	if skip, skipExists := theOptionsSkip["skip"]; skipExists {
		requestOptions.SetSkip(int64(skip))
	}

	if limit, LimitExists := theOptionsLimit["limit"]; LimitExists {
		requestOptions.SetLimit(int64(limit))
	}

	collection := Instance.GetCollection(collectionName)
	cur, err := collection.Find(context.TODO(), selector, requestOptions)

	if err != nil {
		return MongoDBError(err)
	}

	for cur.Next(context.TODO()) {
		var elem map[string]interface{}
		err1 := cur.Decode(&elem)

		fmt.Println("Elem: ", elem)

		if err1 != nil {
			return MongoDBError(err1)
		}

		*result = append(*result, elem)
	}

	if err2 := cur.Err(); err2 != nil {
		return MongoDBError(err2)
	}

	_ = cur.Close(context.TODO())

	return nil
}

func Insert(collectionName string, doc interface{}, result *[]interface{}) *common.Error {
	collection := Instance.GetCollection(collectionName)
	_, err := collection.InsertOne(context.TODO(), doc)
	if err != nil {
		return MongoDBError(err)
	}

	return nil
}

func Update(collectionName string, selector map[string]interface{}, update interface{}, inputOptions interface{}, result *[]interface{}) *common.Error {
	collection := Instance.GetCollection(collectionName)
	_, err := collection.UpdateMany(context.TODO(), selector, update)
	if err != nil {
		return MongoDBError(err)
	}

	return nil
}

func Upsert(collectionName string, selector map[string]interface{}, update interface{}, inputOptions interface{}, result *[]interface{}) *common.Error {
	collection := Instance.GetCollection(collectionName)
	requestOptions := options.Update().SetUpsert(true)
	_, err := collection.UpdateMany(context.TODO(), selector, update, requestOptions)
	if err != nil {
		return MongoDBError(err)
	}

	return nil
}

func Remove(collectionName string, selector map[string]interface{}, result *[]interface{}) *common.Error {
	collection := Instance.GetCollection(collectionName)
	_, err := collection.DeleteOne(context.TODO(), selector)
	if err != nil {
		return MongoDBError(err)
	}

	return nil
}
