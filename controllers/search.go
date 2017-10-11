package controllers

import (
	"ChienHo/SearchEngine/documents"
	"ChienHo/SearchEngine/indexing"
	mSegment "ChienHo/SearchEngine/utils/segment"
	"github.com/gin-gonic/gin"
	"github.com/huichen/sego"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

func Search(c *gin.Context) {
	word := c.Param("word")
	start := time.Now()
	segments := sego.SegmentsToSlice(mSegment.GetSegmenter().Segment([]byte(word)), true)
	ids := make(map[bson.ObjectId]bool)
	for _, segment := range segments {
		if indexes, ok := indexing.Find(segment); ok {
			for _, index := range indexes {
				ids[index] = true
			}
		}
	}
	pages := []documents.Page{}
	if len(ids) == 0 { //搜索引擎没找到结果

	} else {
		idSlice := make([]bson.ObjectId, len(ids))
		i := 0
		for id := range ids {
			idSlice[i] = id
			i++
		}
		if err := documents.PageCollection.Find(bson.M{"_id": bson.M{"$in": idSlice}}).Select(bson.M{"url": 1, "domain": 1, "title": 1, "description": 1, "created_at": 1}).All(&pages); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "服务器错误",
				"error":   err.Error(),
			})
		}
	}
	latency := time.Until(start)
	c.JSON(http.StatusOK, gin.H{
		"data": pages,
		"cost": latency.Nanoseconds(),
	})
}

func SearchDetail(c *gin.Context) {
	id := bson.ObjectIdHex(c.Param("id"))
	page := documents.Page{}
	if err := documents.PageCollection.FindId(id).One(&page); err == nil {
		c.JSON(http.StatusOK, gin.H{
			"data": page,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "页面不存在",
		})
	}
}
