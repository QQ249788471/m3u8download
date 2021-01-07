package util

import (
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type NodeConfig struct {
	// 节点名称，建议使用部署机器的 ip
	NodeName string `json:"node_name"`

	// 下载路径
	DownloadPath string `json:"download_path"`

	// 日志路径
	LogPath string `json:"log_path"`

	// 轮询任务时间间隔，单位秒
	Interval int `json:"interval"`

	// 最大任务数量
	MaxTaskCount int `json:"max_task_count"`

	// ts 下载线程数量
	TsDownloadCount int `json:"ts_download_count,omitempty"`

	// ts 重试下载线程数量
	TsRetryDownloadCount int `json:"ts_retry_download_count,omitempty"`

	// ts 下载最大重试次数
	MaxRetry int `json:"max_retry,omitempty"`

	// ts 下载时向前最大偏移时间，单位 h
	SearchTime int `json:"search_time,omitempty"`

	// ts 下载的超时时间，单位秒
	Timeout int `json:"timeout,omitempty"`

	// 数据库配置
	Db Database `json:"database"`
}

type Database struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
	DbName   string `json:"db_name"`
}

// 配置文件实例
var Conf NodeConfig

// 数据库连接实例
var Db *gorm.DB

// http client
var Client *http.Client

func init() {

	f, err := ioutil.ReadFile("./config/node.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(f, &Conf)
	if err != nil {
		log.Fatal(err)
	}

	// 初始化 http client
	Client = &http.Client{
		Timeout: time.Second * time.Duration(Conf.Timeout),
	}

	// 节点名称最大只能是 128 字符，多的被截断
	if len(Conf.NodeName) > 128 {
		Conf.NodeName = Conf.NodeName[:128]
	}

	// ts 下载和重试下载线程数量
	if Conf.TsDownloadCount <= 0 {
		Conf.TsDownloadCount = Conf.MaxTaskCount
	}
	if Conf.TsRetryDownloadCount <= 0 {
		Conf.TsRetryDownloadCount = Conf.TsDownloadCount / 2
	}

	// ts 下载最大重试次数
	if Conf.MaxRetry <= 0 {
		Conf.MaxRetry = 3
	}

	// 默认向前检索 24h 的 ts 下载任务
	if Conf.SearchTime <= 0 {
		Conf.SearchTime = 24
	}

	// ts 下载超时时间限制
	if Conf.Timeout <= 0 {
		Conf.Timeout = 30
	}

	// 初始化所有文件夹
	initPath()

	// 连接 mysql
	connectDb()
}

func connectDb() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		Conf.Db.Username,
		Conf.Db.Password,
		Conf.Db.Host,
		Conf.Db.DbName)

	// 打开日志文件
	//gormLogPath := GetLogPath("gorm")
	//fd, err := OpenFile(fmt.Sprintf("%s/%s", gormLogPath, "gorm.log"))
	//if err != nil {
	//	log.Fatal(err)
	//}

	newLogger := logger.New(
		// log.New(fd, "\r\n", log.LstdFlags), // io writer
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Warn, // Log level
			Colorful:      true,        // 禁用彩色打印
		},
	)

	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})

	if err != nil {
		log.Fatal(err)
	}
}

func initPath() {

	// 为 app 创建日志目录
	appLogPath := GetLogPath("app")
	if !PathIsExist(appLogPath) {
		if err := Mkdir(appLogPath); err != nil {
			log.Fatal(err)
		}
	}

	// 为 gorm 创建日志目录
	gormLogPath := GetLogPath("gorm")
	if !PathIsExist(gormLogPath) {
		if err := Mkdir(gormLogPath); err != nil {
			log.Fatal(err)
		}
	}

	// 创建 Download 目录
	if !PathIsExist(Conf.DownloadPath) {
		if err := Mkdir(Conf.DownloadPath); err != nil {
			log.Fatal(err)
		}
	}

}

func (c NodeConfig) String() string {
	str, _ := json.Marshal(c)
	return string(str)
}
