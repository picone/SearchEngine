package spider

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"github.com/picone/SearchEngine/utils/redigo"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	FETCH_ROUTINE_NUMBER     = 16
	ANALYSIS_ROUTINE_NUMBER = 64
	RESULT_CHANNAL_NUMBER   = 1024
	JOB_SAVE_KEY            = "fetch_urls"
)

type downloader struct {
	Output             chan *Page
	spiderStopSignal   chan bool
	analysisStopSignal chan bool
}

func newDownloader() *downloader {
	obj := downloader{
		Output:             make(chan *Page, RESULT_CHANNAL_NUMBER),
		spiderStopSignal:   make(chan bool),
		analysisStopSignal: make(chan bool),
	}
	return &obj
}

func (downloader *downloader) AddUrl(url string) {
	conn := redigo.GetConnection()
	defer conn.Close()
	conn.Send("SADD", JOB_SAVE_KEY, url)
}

func (downloader *downloader) Start() {
	if downloader.GetUrlCount() == 0 {
		downloader.AddUrl("https://m.sohu.com/")
		downloader.AddUrl("https://sina.cn/index/feed?from=touch")
		downloader.AddUrl("http://3g.163.com/")
		downloader.AddUrl("http://3g.china.com/")
	}
	for i := 0; i < FETCH_ROUTINE_NUMBER; i++ {
		go downloader.watchJobs()
	}
}

func (downloader *downloader) Stop() {
	//放入线程数量次的停止信号,等到无阻塞时说明爬虫线程都已经停止了
	for i := 0; i < FETCH_ROUTINE_NUMBER; i++ {
		downloader.spiderStopSignal <- true
	}
	log.Println("爬虫例程已停止,请等待分析器完成")
	for i := 0; i < ANALYSIS_ROUTINE_NUMBER; i++ {
		downloader.analysisStopSignal <- true
	}
}

func (downloader *downloader) watchJobs() {
	conn := redigo.GetConnection()
	defer conn.Close()
LOOP:
	for {
		select {
		case <-downloader.spiderStopSignal:
			break LOOP //收到停止信号退出循环
		default:
			if url, err := redis.String(conn.Do("SPOP", JOB_SAVE_KEY)); err == nil {
				if content, err := downloader.fetchContent(url); err == nil {
					downloader.Output <- &Page{
						Url:     url,
						Content: content,
					}
				}
			}
		}
	}
}

func (downloader *downloader) fetchContent(url string) (string, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "<nil>", err
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (iPhone, U, CPU iPhone OS 4_3_3 like Mac OS X, en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5")
	client := http.DefaultClient
	client.Timeout = time.Second * 30
	response, err := client.Do(request)
	if err != nil {
		return "<nil>", err
	}
	defer response.Body.Close()
	if response.StatusCode == 200 {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return "<nil>", err
		}
		return string(body), nil
	} else {
		return "<nil>", errors.New("status code check error")
	}
}

func (downloader *downloader) Consume(callback func(*Page)) {
	for i := 0; i < ANALYSIS_ROUTINE_NUMBER; i++ {
		go downloader.consumeContent(callback)
	}
}

func (downloader *downloader) consumeContent(callback func(*Page)) {
LOOP:
	for content := range downloader.Output {
		select {
		case <-downloader.analysisStopSignal:
			break LOOP
		default:
			callback(content)
		}
	}
}

func (downloader *downloader) GetUrlCount() (count uint64) {
	conn := redigo.GetConnection()
	defer conn.Close()
	if reply, err := redis.Values(conn.Do("SCARD", JOB_SAVE_KEY)); err == nil {
		redis.Scan(reply, &count)
	}
	return
}
