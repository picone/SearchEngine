package main

import (
	"ChienHo/SearchEngine/spider"
	"github.com/gin-gonic/gin"
	"ChienHo/SearchEngine/controllers"
)

func main() {
	startSpider()
	startWeb()
}

func startWeb() {
	r := gin.Default()
	r.GET("/search/:word", controllers.Search)
	r.GET("/search/detail/:id", controllers.SearchDetail)
	r.Run("127.0.0.1:8080")
}

func startSpider() {
	spider.StartSpider()
}
