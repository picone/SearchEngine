package spider

import (
	"ChienHo/SearchEngine/utils/redigo"
	"errors"
	"github.com/garyburd/redigo/redis"
	"io/ioutil"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

const (
	FETCH_THREAD_NUMBER   = 16
	RESULT_CHANNAL_NUMBER = 1024
	JOB_SAVE_KEY          = "fetch_urls"
)

type downloader struct {
	Output               chan *Page
	stopSignal           chan bool
	analysisThreadNumber int64
}

func newDownloader() *downloader {
	obj := downloader{
		Output:     make(chan *Page, RESULT_CHANNAL_NUMBER),
		stopSignal: make(chan bool),
	}
	return &obj
}

func (downloader *downloader) AddUrl(url string) {
	conn := redigo.GetConnection()
	defer conn.Close()
	conn.Send("LPUSH", JOB_SAVE_KEY, url)
}

func (downloader *downloader) Start() {
	for i := 0; i < FETCH_THREAD_NUMBER; i++ {
		go downloader.watchJobs()
	}
	if downloader.GetUrlCount() == 0 {
		downloader.AddUrl("https://m.sohu.com/")
		downloader.AddUrl("https://sina.cn/index/feed?from=touch")
		downloader.AddUrl("http://3g.163.com/")
		downloader.AddUrl("http://3g.china.com/")
	}
}

func (downloader *downloader) Stop() {
	//放入线程数量次的停止信号,等到无阻塞时说明爬虫线程都已经停止了
	for i := 0; i < FETCH_THREAD_NUMBER; i++ {
		downloader.stopSignal <- true
	}
	log.Println("爬虫线程已停止,请等待分析器完成")
	for downloader.analysisThreadNumber > 0 {

	}
}

func (downloader *downloader) watchJobs() {
	conn := redigo.GetConnection()
	defer conn.Close()
LOOP:
	for {
		select {
		case <-downloader.stopSignal:
			break LOOP //收到停止信号退出循环
		default:
			if reply, err := redis.Values(conn.Do("BRPOP", JOB_SAVE_KEY, 10)); err == nil {
				var queueName, url string
				if _, err := redis.Scan(reply, &queueName, &url); err == nil {
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
	go downloader.consumeContent(callback)
}

func (downloader *downloader) consumeContent(callback func(*Page)) {
	for content := range downloader.Output {
		atomic.AddInt64(&downloader.analysisThreadNumber, 1)
		go func() {
			callback(content)
			atomic.AddInt64(&downloader.analysisThreadNumber, -1)
		}()
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
