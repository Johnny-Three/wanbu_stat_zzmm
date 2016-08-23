package main

import (
	//"flag"
	"fmt"
	"github.com/bitly/go-nsq"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	. "wanbu_stat_zzmm/src/calculate"
	. "wanbu_stat_zzmm/src/config"
	. "wanbu_stat_zzmm/src/socket"
)

var (
	Version  = "1.0.10PR1"
	consumer *nsq.Consumer
	err      error
)

func processmsg() error {

	LoadConfig()
	DbInit()
	LoadRules()

	go ZMRefresh() //刷数据接口

	//对接NSQ，消费上传消息
	consumer, err = NewConsummer("base_data_upload", "zzmm"+Trix)
	if err != nil {
		panic(err)
	}

	//Consumer运行，消费消息..
	go func(consumer *nsq.Consumer) {

		err := ConsumerRun(consumer, "base_data_upload", Consumerip+":"+Consumerport)
		if err != nil {
			panic(err)
		}
	}(consumer)

	for {

		uwd := <-Userwalkdata_chan
		fmt.Println("uid upload msg : ", uwd.Uid)
		//todo..过滤消息，LOAD文件中的UID，是否在这里。。
		StatZM(&uwd)
	}
	return nil
}

func init() {

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		sig := <-sigc
		switch sig {
		case syscall.SIGINT:
			fmt.Println("catch SIGINT ")
			os.Exit(1)
		case syscall.SIGQUIT:
			fmt.Println("handle SIGQUIT")
			os.Exit(1)
		}
	}()

}

func main() {

	args := os.Args

	if len(args) == 2 && (args[1] == "-v") {

		fmt.Println("看好了兄弟，现在的版本是【", Version, "】，可别弄错了")
		os.Exit(0)
	}

	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()

	go func() {
		processmsg()
	}()

	select {}
}
