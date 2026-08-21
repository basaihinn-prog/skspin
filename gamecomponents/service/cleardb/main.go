package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 数据库配置
type DatabaseConfig struct {
	Database    string
	Collections []string
}

// 日志颜色
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
)

// 日志函数
func logInfo(message string) {
	fmt.Printf("%s[INFO]%s %s\n", colorGreen, colorReset, message)
}

func logWarn(message string) {
	fmt.Printf("%s[WARN]%s %s\n", colorYellow, colorReset, message)
}

func logError(message string) {
	fmt.Printf("%s[ERROR]%s %s\n", colorRed, colorReset, message)
}

// 清空指定集合
func clearCollection(collection *mongo.Collection) error {
	// 删除所有文档
	_, err := collection.DeleteMany(context.Background(), bson.D{})
	return err
}

// 确认操作
func confirm() bool {
	var response string
	fmt.Print("确定要清空以上数据库和集合吗? (y/n): ")
	fmt.Scanln(&response)
	return strings.ToLower(response) == "y"
}

var (
	mongoUsername = "myAdmin"
	mongoPassword = "myAdminPassword1"
	mongoHost     = "172.31.2.6"
	mongoPort     = "27017"
)

// 创建备份
// 创建备份函数
// 创建备份函数
func createBackup(dbConfig DatabaseConfig, backupDir string) error {
	// 创建备份目录
	backupPath := fmt.Sprintf("%s/%s", backupDir, dbConfig.Database)
	err := os.MkdirAll(backupPath, os.ModePerm)
	if err != nil {
		return err
	}
	mongodumpPath := "D:\\mongodb\\mongodb\\bin\\mongodump.exe"

	// 备份指定集合
	for _, collection := range dbConfig.Collections {
		cmd := exec.Command(mongodumpPath,
			"--uri", fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=admin",
				mongoUsername,
				mongoPassword,
				mongoHost,
				mongoPort),
			"--db", dbConfig.Database,
			"--collection", collection,
			"--out", backupPath,
		)
		if len(dbConfig.Collections) == 1 && dbConfig.Collections[0] == "ALL" {
			cmd = exec.Command(mongodumpPath,
				"--uri", fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=admin",
					mongoUsername,
					mongoPassword,
					mongoHost,
					mongoPort),
				"--db", dbConfig.Database,
				"--out", backupPath,
			)
		}

		// 执行命令
		output, err := cmd.CombinedOutput()
		if err != nil {
			logError(fmt.Sprintf("创建 %s 集合备份失败: %v\nOutput: %s", collection, err, string(output)))
			return err
		}

		logInfo(fmt.Sprintf("成功创建 %s 集合备份：%s", collection, backupPath))
	}

	return nil
}

func main() {
	// 数据库配置
	databases := []DatabaseConfig{
		{
			Database: "betlog",
			Collections: []string{
				"ALL",
			},
		},
		{
			Database: "reports",
			Collections: []string{
				"ALL",
			},
		},
		{
			Database: "GameAdmin",
			Collections: []string{
				"SlotsPoolHistory", "WhitUser", "GameConfig", "AdminOperator", "MaintenanceLog", "PlayerRTPControl", "RickRule",
			},
		},
		{
			Database: "game",
			Collections: []string{
				"orders", "ModifyGoldLog", "Players",
			},
		},
	}

	// MongoDB 连接信息

	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%s", mongoUsername, mongoPassword, mongoHost, mongoPort)

	// 获取所有数据库
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		logError(fmt.Sprintf("连接到 MongoDB 失败: %v", err))
		return
	}
	defer client.Disconnect(context.Background())

	databaseList, err := client.ListDatabaseNames(context.Background(), bson.D{})
	if err != nil {
		logError(fmt.Sprintf("获取数据库列表失败: %v", err))
		return
	}

	// 过滤以 "pg_" 开头的数据库
	var pgDatabases []string
	for _, db := range databaseList {
		if strings.HasPrefix(db, "pg_") {
			pgDatabases = append(pgDatabases, db)
		}
	}
	for _, pgdb := range pgDatabases {
		databases = append(databases, DatabaseConfig{
			Database: pgdb,
			Collections: []string{
				"BetHistory", "psidMap", "players",
			}})
	}

	// 创建备份目录
	backupDir := fmt.Sprintf("D:\\mongodb%s", time.Now().Format("2006-01-02"))
	for _, dbConfig := range databases {
		logInfo(fmt.Sprintf("处理数据库: %s", dbConfig.Database))

		// 创建备份
		err = createBackup(dbConfig, backupDir)
		if err != nil {
			continue
		}
	}

	// 显示将要清空的数据库和集合
	fmt.Println("将要清空以下数据库和集合：")
	for _, dbConfig := range databases {
		fmt.Printf("数据库: %s\n", dbConfig.Database)

		// 如果是 "ALL"，获取所有集合
		var collections []string
		if len(dbConfig.Collections) == 1 && dbConfig.Collections[0] == "ALL" {
			collections, err = getAllCollections(client, dbConfig.Database)
			if err != nil {
				logError(fmt.Sprintf("获取 %s 数据库的集合失败: %v", dbConfig.Database, err))
				continue
			}
		} else {
			collections = dbConfig.Collections
		}

		for _, collection := range collections {
			fmt.Printf("  - 集合: %s\n", collection)
		}
	}

	// 确认操作
	if !confirm() {
		logWarn("操作已取消")
		os.Exit(0)
	}

	// 处理每个数据库
	for _, dbConfig := range databases {
		logInfo(fmt.Sprintf("处理数据库: %s", dbConfig.Database))

		// 连接到 MongoDB
		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
		if err != nil {
			logError(fmt.Sprintf("连接到 MongoDB 失败: %v", err))
			continue
		}
		defer client.Disconnect(context.Background())

		// 获取集合列表
		var collections []string
		if len(dbConfig.Collections) == 1 && dbConfig.Collections[0] == "ALL" {
			collections, err = getAllCollections(client, dbConfig.Database)
			if err != nil {
				logError(fmt.Sprintf("获取 %s 数据库的集合失败: %v", dbConfig.Database, err))
				continue
			}
		} else {
			collections = dbConfig.Collections
		}

		// 清空每个集合
		for _, collectionName := range collections {
			collection := client.Database(dbConfig.Database).Collection(collectionName)

			logInfo(fmt.Sprintf("正在清空集合: %s", collectionName))

			if err := clearCollection(collection); err != nil {
				logError(fmt.Sprintf("清空 %s.%s 失败: %v", dbConfig.Database, collectionName, err))
			} else {
				logInfo(fmt.Sprintf("成功清空 %s.%s", dbConfig.Database, collectionName))
			}
		}
	}
	ClearAdminPermission(mongoURI)
	ClearAdminUser(mongoURI)
	logInfo("所有操作完成")
}

// 获取数据库中的所有集合
func getAllCollections(client *mongo.Client, dbName string) ([]string, error) {
	// 获取数据库
	database := client.Database(dbName)

	// 列出所有集合
	collections, err := database.ListCollectionNames(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}

	return collections, nil
}

// 处理GameAdmin表中的    AdminUser表与AdminPermission表
func ClearAdminPermission(mongoURI string) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		logError(fmt.Sprintf("连接到 MongoDB 失败: %v", err))
		return
	}
	defer client.Disconnect(context.Background())

	collection := client.Database("GameAdmin").Collection("AdminPermission")
	_, err = collection.DeleteMany(context.TODO(), bson.M{"_id": bson.M{"$gt": 3}})

	if err != nil {
		logError(fmt.Sprintf("清空 %s.%s 失败: %v", "GameAdmin", "AdminPermission", err))
	} else {
		logInfo(fmt.Sprintf("成功清空 %s.%s", "GameAdmin", "AdminPermission"))
	}
}

func ClearAdminUser(mongoURI string) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		logError(fmt.Sprintf("连接到 MongoDB 失败: %v", err))
		return
	}
	defer client.Disconnect(context.Background())

	collection := client.Database("GameAdmin").Collection("AdminUser")
	_, err = collection.DeleteMany(context.TODO(), bson.M{"_id": bson.M{"$gt": 1}})

	if err != nil {
		logError(fmt.Sprintf("清空 %s.%s 失败: %v", "GameAdmin", "AdminUser", err))
	} else {
		logInfo(fmt.Sprintf("成功清空 %s.%s", "GameAdmin", "AdminUser"))
	}
}
