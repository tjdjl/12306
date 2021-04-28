package main

import (
	"os"

	"12306.com/12306/common"
	"12306.com/12306/stations"
	"12306.com/12306/trains"
	"12306.com/12306/users"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

func main() {
	InitConfig()
	db := common.InitDB()
	db.AutoMigrate(&users.User{})
	db.AutoMigrate(&stations.Station{})
	db.AutoMigrate(&trains.Order{})
	db.AutoMigrate(&users.Passanger{})

	defer db.Close()

	r := gin.Default()
	r = CollectRoute(r) //注册路由

	port := viper.GetString("server.port")
	if port != "" {
		panic(r.Run(":" + port))
	}
	panic(r.Run())
}

func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
} //初始化配置文件
