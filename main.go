package main

import (
	"ChienHo/SearchEngine/spider"
	"github.com/gin-gonic/gin"
)

func main() {
	startSpider()
	startWeb()
}

func startWeb() {
	r := gin.Default()
	r.Run("127.0.0.1:8888")
}

func startSpider() {
	spider.StartSpider()
}
