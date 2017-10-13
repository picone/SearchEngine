package main

import (
	"ChienHo/SearchEngine/controllers"
	"ChienHo/SearchEngine/spider"
	"github.com/gin-gonic/gin"
	"ChienHo/SearchEngine/middlewares"
)

func main() {
	//startSpider()
	startWeb()
}

func startWeb() {
	r := gin.Default()
	r.GET("/search/:word", middlewares.GetPagination(), controllers.Search)
	r.GET("/search-detail/:id", controllers.SearchDetail)
	r.Static("/assets", "./assets")
	r.StaticFile("/", "./assets/pages/index.html")
	r.Run("127.0.0.1:8080")
}

func startSpider() {
	spider.StartSpider()
}
