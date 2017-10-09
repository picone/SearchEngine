package spider

import (
	"ChienHo/SearchEngine/documents"
	"ChienHo/SearchEngine/utils/html"
	"github.com/huichen/sego"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

type analysis struct {
}

var (
	segmenter sego.Segmenter
)

func init() {
	segmenter.LoadDictionary("./data/dictionary.txt")
}

func newAnalysis() *analysis {
	obj := analysis{}
	return &obj
}

func (analysis *analysis) Watch(page *Page) {
	//分析超级链接
	for _, url := range html.GetHrefLinks(page.Content) {
		//去除锚点后面的东西避免重复
		url = html.RemoveUrlAnchor(url)
		//先判断有没有爬过,没有的话跟着爬下去
		if count, err := documents.PageCollection.Find(bson.M{"url": url}).Select(bson.M{"created_at": 1}).Count(); err == nil && count == 0 {
			producer.AddUrl(url)
		}
	}
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
	documents.PageCollection.Insert(p)
	//除去所有tags,方便做索引
	page.Content = html.RemoveHTMLTags(page.Content)
	//倒排索引,先分词
	segments := sego.SegmentsToSlice(segmenter.Segment([]byte(page.Content)), true)
	log.Println("分词内容:", segments)
}
