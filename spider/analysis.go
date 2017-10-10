package spider

import (
	"ChienHo/SearchEngine/documents"
	"ChienHo/SearchEngine/indexing"
	"ChienHo/SearchEngine/utils/html"
	mSegment "ChienHo/SearchEngine/utils/segment"
	"github.com/huichen/sego"
	"gopkg.in/mgo.v2/bson"
	"time"
	"log"
	"gopkg.in/mgo.v2"
)

type analysis struct {
}

func newAnalysis() *analysis {
	obj := analysis{}
	return &obj
}

func (analysis *analysis) Watch(page *Page) {
	log.Println("分析连接:", page.Url)
	meta := html.ParseMeta(page.Content)
	p := documents.Page{
		Id:          bson.NewObjectId(),
		Url:         page.Url,
		Content:     page.Content,
		Domain:      html.GetDomain(page.Url),
		Title:       html.ParseTitle(page.Content),
		Description: meta["description"],
		Keyword:     meta["keywords"],
		CreatedAt:   time.Now(),
	}
	//保存到数据库中中
	if err := documents.PageCollection.Insert(p); err != nil {
		switch err.(type) {
		case *mgo.LastError:
			if err.(*mgo.LastError).Code == 11000 {
				return//已插入过,不用再索引
			}
		}
	}
	//分析超级链接
	for _, url := range html.GetHrefLinks(page.Content) {
		//去除锚点后面的东西避免重复
		url = html.RemoveUrlAnchor(url)
		//先判断有没有爬过,没有的话跟着爬下去
		if err := documents.PageCollection.Update(bson.M{"url": url}, bson.M{"$inc": bson.M{"rank": 1}}); err == mgo.ErrNotFound {
			producer.AddUrl(url)
		}
	}
	//除去所有tags,方便做索引
	page.Content = html.RemoveHTMLTags(page.Content)
	//倒排索引,先分词
	segments := sego.SegmentsToSlice(mSegment.GetSegmenter().Segment([]byte(page.Content)), true)
	for _, segment := range segments {
		indexing.Add(segment, p.Id)
	}
}
