package controllers

import (
	"documents"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"net/http"
	"time"
)

func Search(c *gin.Context) {
	word := c.Param("word")
	start := time.Now()
	pages := []documents.Page{}
	if indexes, ok := documents.IndexingFind(word); ok {
		if err := documents.PageCollection.Find(bson.M{"_id": bson.M{"$in": indexes}}).Select(bson.M{"url": 1, "domain": 1, "title": 1, "description": 1, "created_at": 1}).Sort("rank:-").Limit(c.MustGet("page_size").(int)).Skip(c.MustGet("page_skip").(int)).All(&pages); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "服务器错误",
				"error":   err.Error(),
			})
			return
		}
	}
	latency := time.Since(start)
	c.JSON(http.StatusOK, gin.H{
		"data": pages,
		"cost": latency.Nanoseconds(),
	})
}

func SearchDetail(c *gin.Context) {
	id := bson.ObjectIdHex(c.Param("id"))
	page := documents.Page{}
	if err := documents.PageCollection.FindId(id).One(&page); err == nil {
		c.HTML(http.StatusOK, "search-detail.tmpl", gin.H{
			"content": template.HTML(page.Content),
		})
	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}
}
