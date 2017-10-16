package main

import (
	"github.com/picone/SearchEngine/controllers"
	"github.com/gin-gonic/gin"
	"github.com/picone/SearchEngine/middlewares"
	"github.com/picone/SearchEngine/spider"
)

func main() {
	startSpider()
	startWeb()
}

func startWeb() {
	r := gin.Default()
	r.LoadHTMLGlob("./templates/*")
	r.GET("/search/:word", middlewares.GetPagination(), controllers.Search)
	r.GET("/search-detail/:id", controllers.SearchDetail)
	r.Static("/assets", "./assets")
	r.StaticFile("/", "./assets/pages/index.html")
	r.Run("127.0.0.1:8080")
}

func startSpider() {
	spider.StartSpider()
}
