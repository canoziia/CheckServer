package main

import (
	"flag"
	"fmt"
	"strconv"
	"time"

	cfg "github.com/canoziia/checkserver/config"
	"github.com/canoziia/checkserver/logs"
	"github.com/canoziia/checkserver/mails"
	"github.com/canoziia/checkserver/services"
)

var (
	configPath string
	config     cfg.Config

	logPath    string
	serverList []*services.Server

	mailConf mails.MailConf

	cycleTime  int
	reportTime int
)

func init() {
	flag.StringVar(&configPath, "c", "config.ini", "配置文件路径")
	flag.Parse()
	config = cfg.GetConfig(configPath)

	var err error
	cycleTime, err = strconv.Atoi(config["common"]["cycle"])
	if err != nil {
		panic(err)
	}
	reportTime, err = strconv.Atoi(config["common"]["report"])
	if err != nil {
		panic(err)
	}

	logPath = config["common"]["log"]

	mailConf.Port = config["mail"]["port"]
	mailConf.Host = config["mail"]["host"]
	mailConf.Encrypt = config["mail"]["encrypt"]
	mailConf.Username = config["mail"]["username"]
	mailConf.Password = config["mail"]["password"]
	mailConf.Target = config["mail"]["target"]
	mailConf.Name = config["mail"]["name"]
	mails.LoadMailConf(mailConf)

	// logger = logs.NewLogger(logPath)
	logs.InitGlobalLogger(logPath)

	for key, serverMap := range config {
		if key != "mail" && key != "common" && key != "DEFAULT" {
			newServer := services.NewServer()
			newServer.ParseFromConfig(serverMap)
			serverList = append(serverList, newServer)
		}
	}
}

func succeedFunc(svr *services.Server, info logs.Info) string {
	return "Nothing"
}

func errorFunc(svr *services.Server, info logs.Info) string {
	mails.SendMail(fmt.Sprintf("[Error] [%s|%s %s:%s]", svr.Name, svr.Mode, svr.Host, svr.Port), info.Message, "html")
	return "Send Mail"
}

func failFunc(svr *services.Server, info logs.Info) string {
	mails.SendMail(fmt.Sprintf("[Failed] [%s|%s %s:%s]", svr.Name, svr.Mode, svr.Host, svr.Port), info.Message, "html")
	return "Send Mail"
}

func report(reportFunc func()) {
	for {
		logs.GlobalLogger.AllPrint("[report]")
		reportFunc()
		time.Sleep(time.Duration(reportTime) * time.Second)
	}
}

func reportFunc() {
	defer func() {
		if e := recover(); e != nil {
			logs.GlobalLogger.AllPrint("[failed] Report Failed")
		}
	}()
	err := mails.SendMail("[report]", "工作正常", "html")
	if err != nil {
		logs.GlobalLogger.AllPrint("[failed] Report Failed | " + err.Error())
	}
}

func main() {
	go report(reportFunc)
	for {
		fmt.Print("====================================================================\r\n")
		for _, server := range serverList {
			info := server.Check(failFunc, succeedFunc, errorFunc)
			logs.GlobalLogger.Record(info)
		}
		time.Sleep(time.Duration(cycleTime) * time.Second)
	}
}
