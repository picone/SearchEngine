package spider

import (
	"io/ioutil"
	"net/http"
)

const (
	JOB_CHANNAL_NUMBER    = 5
	FETCH_THREAD_NUMBER   = 32
	RESULT_CHANNAL_NUMBER = 1024
)

var (
	jobUrls = make(chan string, JOB_CHANNAL_NUMBER)
)

type downloader struct {
	Output chan *Page
}

func newDownloader() *downloader {
	obj := downloader{}
	obj.Output = make(chan *Page, RESULT_CHANNAL_NUMBER)
	return &obj
}

func (downloader *downloader) AddUrl(url string) {
	jobUrls <- url
}

func (downloader *downloader) Start() {
	for i := 0; i < FETCH_THREAD_NUMBER; i++ {
		go downloader.watchJobs()
	}
}

func (downloader *downloader) Stop() {
	//TODO 停止spider

}

func (downloader *downloader) watchJobs() {
	for url := range jobUrls {
		if content, err := downloader.fetchContent(url); err == nil {
			downloader.Output <- &Page{
				Url:     url,
				Content: content,
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
	response, err := http.DefaultClient.Do(request)
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
