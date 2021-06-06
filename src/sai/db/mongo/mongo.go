package mongo

import (
	"crypto/tls"
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"saiStorageMongo/src/github.com/fatih/color"
	"saiStorageMongo/src/gopkg.in/mgo.v2"
	"saiStorageMongo/src/gopkg.in/mgo.v2/bson"
	"saiStorageMongo/src/sai/common"
	"saiStorageMongo/src/sai/network/http"
	"saiStorageMongo/src/sai_storage/settings"
	"strconv"
	"time"
)

type Database struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type MongoDB interface {
	GetSession() *mgo.Session
	GetDatabaseConfig() *Database
	GetHosts() *[]string
	String() string
}

type LocalMongo struct {
	Config Database `json:"config"`
	Hosts  []string `json:"hosts"`
}

type AtlasMongo struct {
	Config           Database `json:"config"`
	Hosts            []string `json:"hosts"`
	ConnectionString string   `json:"connection_string"`
}

func (mi *LocalMongo) GetSession() *mgo.Session {
	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs: mi.Hosts,
	})
	if err != nil {
		panic(err)
	}

	return session
}

func (mi *LocalMongo) GetDatabaseConfig() *Database {
	return &mi.Config
}

func (mi *LocalMongo) GetHosts() *[]string {
	return &mi.Hosts
}

func (mi *LocalMongo) String() string {
	return fmt.Sprintln("Hosts", fmt.Sprintf("%+v", mi.Hosts))
}

func (mi *AtlasMongo) GetSession() *mgo.Session {
	if mi.ConnectionString == "" {
		dialInfo := &mgo.DialInfo{
			Addrs:    mi.Hosts,
			Database: mi.Config.Database,
			Username: mi.Config.Username,
			Password: mi.Config.Password,
		}

		tlsConfig := &tls.Config{}

		dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
			return conn, err
		}

		session, _ := mgo.DialWithInfo(dialInfo)

		return session
	} else {
		dialInfo, _ := mgo.ParseURL(mi.ConnectionString)

		tlsConfig := &tls.Config{}

		dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
			return conn, err
		}

		session, _ := mgo.DialWithInfo(dialInfo)

		return session
	}
}

func (mi *AtlasMongo) GetDatabaseConfig() *Database {
	return &mi.Config
}

func (mi *AtlasMongo) GetHosts() *[]string {
	return &mi.Hosts
}

func (mi *AtlasMongo) String() string {
	mglog := fmt.Sprintln("ConnectionString", fmt.Sprintf("%+v", mi.ConnectionString))
	mglog += fmt.Sprintln("Hosts", fmt.Sprintf("%+v", mi.Hosts))
	return mglog
}

var (
	MongoInstance MongoDB
)

func SetMongoDBInv(mongoDB settings.DBConfig) {
	if mongoDB.Atlas.Enabled {
		MongoInstance = &AtlasMongo{
			Config: Database{
				Username: mongoDB.Atlas.Config.Username,
				Password: mongoDB.Atlas.Config.Password,
				Database: mongoDB.Atlas.Config.Database,
			},
			ConnectionString: mongoDB.Atlas.ConnectionString,
			Hosts:            mongoDB.Atlas.Hosts,
		}
	} else {
		MongoInstance = &LocalMongo{
			Config: Database{
				Username: mongoDB.Local.Config.Username,
				Password: mongoDB.Local.Config.Password,
				Database: mongoDB.Local.Config.Database,
			},
			Hosts: mongoDB.Local.Hosts,
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

		startMongoCmd := exec.Command("mongod")
		startMongoCmd.Start()
		d.Println("Mongod started. PID", startMongoCmd.Process.Pid)

		getStartedMongodCmd := exec.Command("/bin/sh", "-c", "ps ax | grep mongod")
		startedMongods, _ := getStartedMongodCmd.Output()
		fmt.Println(string(startedMongods))
	}
	// case "windows":
}

func TestMongodConnection() {
	d := color.New(color.FgMagenta, color.Bold)
	c1 := make(chan bool, 1)
	go func() {
		result := []interface{}{}
		err := Find(fmt.Sprint("test"), nil, nil, &result)
		if err != nil {
			fmt.Println(err)
		}

		c1 <- true
	}()

	select {
	case <-c1:
		d.Println("Mongod started correctly")
	case <-time.After(5 * time.Second):
		d.Println("Mongod not started correctly")
		fmt.Println(MongoInstance)
	}
}

func FindOne(collectionName string, selector map[string]interface{}, result *interface{}) (*common.Error) {
	session := MongoInstance.GetSession()
	defer session.Close()

	collection := session.DB(MongoInstance.GetDatabaseConfig().Database).C(collectionName)

	var document interface{}
	if id, exist := selector["_id"]; exist {
		if id.(string) == "" || len(id.(string)) != 24 {
			return http.BadRequestError()
		}
		if err := collection.FindId(bson.ObjectIdHex(fmt.Sprint(selector["_id"]))).One(&document); err != nil {
			return MongoDBError(err)
		}
	} else {
		if err := collection.Find(selector).One(&document); err != nil {
			return MongoDBError(err)
		}
	}

	if result != nil {
		*result = document
	}

	return nil
}

func Find(collectionName string, selector map[string]interface{}, options interface{}, result *[]interface{}) (*common.Error) {
	_ = options
	session := MongoInstance.GetSession()
	defer session.Close()

	collection := session.DB(MongoInstance.GetDatabaseConfig().Database).C(collectionName)
	theOptionsLimit :=  make(map[string]int)
	theOptionsSort :=  make(map[string]string)
	switch o := options.(type)  {
	case map[string]interface{}:
		for s, b := range o {
			theK := string(s)
			if theK == "limit" {
				val,_ := strconv.Atoi(b.(string));
				theOptionsLimit[theK] = val
			}
			if theK == "sort" {
				val, _ := b.(string);
				theOptionsSort[theK] = val
			}
		}
	default:
		fmt.Println("the type?????????",o)
	}
	fmt.Println("theOptionsSort",theOptionsSort)
	var documents []interface{}
	if id, exist := selector["_id"]; exist {
		if id.(string) == "" || len(id.(string)) != 24 {
			return http.BadRequestError()
		}
		if err := collection.FindId(bson.ObjectIdHex(fmt.Sprint(selector["_id"]))).All(&documents); err != nil {
			return MongoDBError(err)
		}
	} else {
		if  sort, sortExists := theOptionsSort["sort"]; sortExists {
			if limit, LimitExists := theOptionsLimit["limit"]; LimitExists {
				fmt.Println("Call with sort and limit")
				if err := collection.Find(selector).Limit(limit).Sort(sort).All(&documents); err != nil {
					return MongoDBError(err)
				}
			} else {
				fmt.Println("Call with sort ONLY")
				if err := collection.Find(selector).Sort(sort).All(&documents); err != nil {
					return MongoDBError(err)
				}
			}
		} else {
			if limit, exists := theOptionsLimit["limit"]; exists {
				fmt.Println("Call with limit ONLY LMIT",limit)
				if err := collection.Find(selector).Limit(limit).All(&documents); err != nil {
					return MongoDBError(err)
				}
			} else {
				if err := collection.Find(selector).All(&documents); err != nil {
					return MongoDBError(err)
				}
			}
		}
		/**
		if limit, exists := theOptionsLimit["limit"]; exists {
			fmt.Println("Call with limit ",limit)
			if err := collection.Find(selector).Limit(int(limit)).All(&documents); err != nil {
				return MongoDBError(err)
			}
		} else {
			if err := collection.Find(selector).All(&documents); err != nil {
				return MongoDBError(err)
			}
		}
		* **/
	}

	if result != nil {
		*result = append(*result, documents...)
	}

	return nil
}

func Insert(collectionName string, docs interface{}, result *[]interface{}) (*common.Error) {
	session := MongoInstance.GetSession()
	defer session.Close()

	collection := session.DB(MongoInstance.GetDatabaseConfig().Database).C(collectionName)
	fmt.Println(docs)
	if err := collection.Insert(docs); err != nil {
		return MongoDBError(err)
	}

	var document interface{};
	collection.Find(docs).One(&document)

	if result != nil {
		*result = append(*result, document)
	}

	return nil
}

//func Update(collectionName string, selector map[string]interface{}, update interface{}, options interface{}, result *[]interface{}) (*common.Error) {
//	session := MongoInstance.GetSession()
//	defer session.Close()
//
//	collection := session.DB(MongoInstance.GetDatabaseConfig().Database).C(collectionName)
//
//	switch options {
//	case nil:
//		options = "set"
//	case "":
//		options = "set"
//	case "all":
//		options = nil
//	}
//
//	updator := update
//	if options != nil {
//		update = bson.M{fmt.Sprintf("$" + fmt.Sprint(options)): update}
//	}
//
//	var documents []interface{};
//	if _, exist := selector["_id"]; exist {
//		if err := collection.UpdateId(bson.ObjectIdHex(fmt.Sprint(selector["_id"])), update); err != nil {
//			return MongoDBError(err)
//		}
//		collection.FindId(bson.ObjectIdHex(fmt.Sprint(selector["_id"]))).All(&documents)
//	} else {
//		if err := collection.Update(selector, update); err != nil {
//			return MongoDBError(err)
//		}
//		collection.Find(updator).All(&documents)
//	}
//
//	if result != nil {
//		*result = append(*result, documents...)
//	}
//
//	return nil
//}

func Update(collectionName string, selector map[string]interface{}, update interface{}, options interface{}, result *[]interface{}) (*common.Error) {
	session := MongoInstance.GetSession()
	defer session.Close()

	collection := session.DB(MongoInstance.GetDatabaseConfig().Database).C(collectionName)

	switch options {
	case nil:
		options = "set"
	case "":
		options = "set"
	case "all":
		options = nil
	}

	updator := update
	if options != nil {
		update = bson.M{fmt.Sprintf("$" + fmt.Sprint(options)): update}
	}

	var documents []interface{};
	if id, exist := selector["_id"]; exist {
		if id.(string) == "" || len(id.(string)) != 24 {
			return http.BadRequestError()
		}
		if err := collection.UpdateId(bson.ObjectIdHex(fmt.Sprint(selector["_id"])), update); err != nil {
			return MongoDBError(err)
		}
		collection.FindId(bson.ObjectIdHex(fmt.Sprint(selector["_id"]))).All(&documents)
	} else {
		if _, err := collection.UpdateAll(selector, update); err != nil {
			return MongoDBError(err)
		}
		collection.Find(updator).All(&documents)
	}

	if result != nil {
		*result = append(*result, documents...)
	}

	return nil
}

func Upsert(collectionName string, selector map[string]interface{}, update interface{}, options interface{}, result *[]interface{}) (*common.Error) {
	session := MongoInstance.GetSession()
	defer session.Close()

	collection := session.DB(MongoInstance.GetDatabaseConfig().Database).C(collectionName)

	switch options {
	case nil:
		options = "set"
	case "":
		options = "set"
	case "all":
		options = nil
	}

	updator := update
	if options != nil {
		update = bson.M{fmt.Sprintf("$" + fmt.Sprint(options)): update}
	}

	var documents []interface{};
	if id, exist := selector["_id"]; exist {
		if id.(string) == "" || len(id.(string)) != 24 {
			return http.BadRequestError()
		}
		if err := collection.UpdateId(bson.ObjectIdHex(fmt.Sprint(selector["_id"])), update); err != nil {
			return MongoDBError(err)
		}
		collection.FindId(bson.ObjectIdHex(fmt.Sprint(selector["_id"]))).All(&documents)
	} else {
		if _, err := collection.Upsert(selector, update); err != nil {
			return MongoDBError(err)
		}
		collection.Find(updator).All(&documents)
	}

	if result != nil {
		*result = append(*result, documents...)
	}

	return nil
}

func Remove(collectionName string, selector map[string]interface{}, result *[]interface{}) (*common.Error) {
	session := MongoInstance.GetSession()
	defer session.Close()

	collection := session.DB(MongoInstance.GetDatabaseConfig().Database).C(collectionName)

	var documents []interface{}
	if id, exist := selector["_id"]; exist {
		if id.(string) == "" || len(id.(string)) != 24 {
			return http.BadRequestError()
		}
		if err := collection.RemoveId(bson.ObjectIdHex(fmt.Sprint(selector["_id"]))); err != nil {
			return MongoDBError(err)
		}
	} else {
		if _, err := collection.RemoveAll(selector); err != nil {
			return MongoDBError(err)
		}
	}
	collection.Find(nil).All(&documents)

	if result != nil {
		*result = append(*result, documents...)
	}

	return nil
}

//func RemoveAll(collectionName string, selector interface{}, result *[]interface{}) (*common.Error) {
//	session := MongoInstance.GetSession()
//	defer session.Close()
//
//	collection := session.DB(MongoInstance.GetDatabaseConfig().Database).C(collectionName)
//
//	if _, err := collection.RemoveAll(selector); err != nil {
//		return MongoDBError(err)
//	}
//
//	var documents []interface{};
//	collection.Find(nil).All(&documents)
//
//	if result != nil {
//		*result = append(*result, documents...)
//	}
//
//	return nil
//}
