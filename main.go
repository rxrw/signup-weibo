// @Title        WeiboTask
// @Description  包括程序入口，日志初始化
// @Author       星辰
// @Update
package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"os"
	"sync"
	"weibo/client"
	"weibo/tasks"
)

// @title         main
// @description   程序入口
// @auth          星辰
// @param
// @return
func main() {
	var logPath string
	flag.StringVar(&logPath, "l", "./WeiboTask.log", "日志文件路径,默认为./WeiboTask.log")
	flag.Parse()

	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		println(err)
		initLog(nil)
	} else {
		initLog(logFile)
		defer logFile.Close()
	}
	err = LoadConfig()
	if err != nil {
		log.Println("配置文件加载失败:" + err.Error())
		os.Exit(6)
	}
	runOnce(false)

}

// @title         initLog
// @description   初始化日志
// @auth          星辰
// @param         logFile       *os.File   日志文件
// @return
func initLog(logFile *os.File) {
	log.SetFlags(log.LstdFlags)
	mBuffer = bytes.NewBufferString("")
	logIo := io.MultiWriter(os.Stdout, mBuffer)
	if logFile != nil {
		logIo = io.MultiWriter(logIo, logFile)
	}
	log.SetOutput(logIo)
}

// @title         runOnce
// @description   单次运行任务
// @auth          星辰
// @param         configPath       string   配置文件路径
// @param         reloadConfig     bool     执行前是否重载配置文件
// @return
func runOnce(reloadConfig bool) {
	if reloadConfig {
		err := LoadConfig()
		if err != nil {
			log.Println("配置文件加载失败:" + err.Error())
			os.Exit(6)
		}
	}
	defer sendToServerChan()
	wb := client.New(MyConfig.C, MyConfig.S, MyConfig.F)
	if wb.LoginByCookies(MyConfig.Cookies) {
		defer func() { MyConfig.Cookies = wb.GetCookies(); _ = SaveConfig() }()
		runTasks(wb)
	} else {
		log.Println("登录失败")
	}

}

// @title         runTasks
// @description   启动用户每日任务
// @auth          星辰
// @param         w          *WeiboClient.WeiboClient  微博客户端
// @return
func runTasks(w *client.WeiboClient) {
	var mywg sync.WaitGroup
	mywg.Add(5)
	go tasks.SuperCheckIn(w, &mywg)
	go tasks.ReceiveScore(w, &mywg)
	go tasks.RepostAndComment(w, &mywg)
	go tasks.AppSignIn(w, &mywg)
	go tasks.AppTaskEntry(w, &mywg)
	mywg.Wait()
}
