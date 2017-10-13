package spider

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"
	"bufio"
	"io"
)

const (
	JOB_CHANNAL_NUMBER    = 4096
	FETCH_THREAD_NUMBER   = 32
	RESULT_CHANNAL_NUMBER = 1024
	JOB_SAVE_PATH = "./data/jobs.txt"
)

type downloader struct {
	Output    chan *Page
	jobUrls   chan string
	isRunning int32
}

func newDownloader() *downloader {
	obj := downloader{
		Output:    make(chan *Page, RESULT_CHANNAL_NUMBER),
		jobUrls:   make(chan string, JOB_CHANNAL_NUMBER),
		isRunning: 0,
	}
	return &obj
}

func (downloader *downloader) AddUrl(url string) {
	downloader.jobUrls <- url
}

func (downloader *downloader) Start() {
	downloader.isRunning = 1
	for i := 0; i < FETCH_THREAD_NUMBER; i++ {
		go downloader.watchJobs()
	}
	//载入历史数据
	go func(){
		if f, err := os.Open(JOB_SAVE_PATH); err == nil {
			defer f.Close()
			buffer := bufio.NewReader(f)
			for {
				line, _, c := buffer.ReadLine()
				if c == io.EOF {
					break
				}
				if l := string(line); l != "" {
					downloader.AddUrl(l)
				}
			}
		} else {
			downloader.AddUrl("https://m.sohu.com/")
			downloader.AddUrl("https://sina.cn/index/feed?from=touch")
			downloader.AddUrl("http://3g.163.com/")
			downloader.AddUrl("http://3g.china.com/")
		}
	}()
}

func (downloader *downloader) Stop() bool {
	if atomic.LoadInt32(&downloader.isRunning) == 1 {
		atomic.StoreInt32(&downloader.isRunning, 0)
		log.Println("等待爬重器完成任务")
		downloader.saveJobs(JOB_SAVE_PATH)
		return true
	} else {
		return false
	}
}

func (downloader *downloader) watchJobs() {
	for url := range downloader.jobUrls {
		if content, err := downloader.fetchContent(url); err == nil {
			downloader.Output <- &Page{
				Url:     url,
				Content: content,
			}
		}
		if atomic.LoadInt32(&downloader.isRunning) == 0 {
			break
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
	if response.StatusCode == 301 {

	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "<nil>", err
	}
	return string(body), nil
}

func (downloader *downloader) Consume(callback func(*Page)) {
	go downloader.consumeContent(callback)
}

func (downloader *downloader) consumeContent(callback func(*Page)) {
	for content := range downloader.Output {
		go callback(content)
	}
}

func (downloader *downloader) saveJobs(filename string) {
	timer := time.NewTimer(time.Second * 60)
	defer timer.Stop()
	if f, err := os.Create(filename); err == nil {
		defer f.Close()
		LOOP:
		for {
			select {
			case url := <- downloader.jobUrls:
				f.WriteString(url)
				f.WriteString("\n")
			case <- timer.C:
				break LOOP
			}
		}
	}
}
