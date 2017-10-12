package middlewares

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetPagination() gin.HandlerFunc {
	return func(context *gin.Context) {
		var (
			pageSize int   = 20
			page     int64 = 1
		)
		if pageSizeStr, ok := context.GetQuery("page_size"); ok {
			if pageSizeTmp, err := strconv.ParseInt(pageSizeStr, 10, 32); err == nil {
				pageSize = int(pageSizeTmp)
			}
		}
		context.Set("page_size", pageSize)
		if pageStr, ok := context.GetQuery("page"); ok {
			if pageTmp, err := strconv.ParseInt(pageStr, 10, 32); err == nil {
				page = pageTmp
			}
		}
		context.Set("page_skip", int(page-1)*pageSize)
	}
}
