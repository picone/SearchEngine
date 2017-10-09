package spider

import (
	"os"
	"os/signal"
	"syscall"
)

var (
	producer *downloader
	costumer *analysis
)

func init() {
	producer = newDownloader()
	costumer = newAnalysis()
}

func StartSpider() {
	//go signalHandler()//捕捉停止信号
	producer.Start()
	producer.Consume(costumer.Watch)
	//TODO 从这里爬起
	producer.AddUrl("https://m.sohu.com/")
}

func StopSpider() {
	producer.Stop()

}

func signalHandler() {
	sign := make(chan os.Signal)
	signal.Notify(sign, syscall.SIGINT)
	for {
		msg := <-sign
		if msg == syscall.SIGINT {
			//TODO 停止爬虫
			StopSpider()
		}
	}
}
