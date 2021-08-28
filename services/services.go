package services

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/canoziia/checkserver/logs"
)

// Server 服务器结构体
type Server struct {
	Name       string
	Host       string
	Port       string
	Mode       string
	Timeout    int
	Conn       net.Conn
	ReExecTime int
	Status     struct {
		Code       int // 0是过程出错, 1是正常, 2是异常
		Message    string
		LastChange time.Time
		NeedExec   bool
	}
}

// NewServer 创建新的server结构
func NewServer() *Server {
	return new(Server)
}

// ParseFromConfig 写入配置并赋初值
func (svr *Server) ParseFromConfig(conf map[string]string) {
	defer func() {
		if e := recover(); e != nil {
			err := errors.New("Parse Config Failed: \r\n    " + fmt.Sprint(e))
			panic(err)
		}
	}()
	svr.Name = conf["name"]
	svr.Host = conf["host"]
	svr.Port = conf["port"]
	svr.Mode = conf["mode"]
	var err error
	svr.Timeout, err = strconv.Atoi(conf["timeout"])
	if err != nil {
		panic(err)
	}
	svr.ReExecTime, err = strconv.Atoi(conf["reexectime"])
	if err != nil {
		panic(err)
	}
	svr.Status.Code = 0
	svr.Status.Message = "init"
	svr.Status.LastChange = time.Now()
	svr.Status.NeedExec = false
}

// TryConnect 尝试连接此server
func (svr *Server) TryConnect() (err error) {
	svr.Conn, err = net.DialTimeout(svr.Mode, fmt.Sprintf("%s:%s", svr.Host, svr.Port), time.Duration(svr.Timeout)*time.Second)
	if err != nil {
		svr.ChangeStatus(2, fmt.Sprintf("Connect Error [%s|%s %s:%s] | %s", svr.Name, svr.Mode, svr.Host, svr.Port, err.Error()))
		err = errors.New(svr.Status.Message)
		return
	}
	defer svr.Conn.Close()
	svr.ChangeStatus(1, "Online")
	return
}

func (svr *Server) ChangeStatus(code int, message string) {
	if !(svr.Status.Code == code && svr.Status.LastChange.Add(time.Duration(svr.ReExecTime)*time.Second).After(time.Now())) {
		svr.Status.LastChange = time.Now()
		svr.Status.NeedExec = true
	}
	svr.Status.Code = code
	svr.Status.Message = message
}

func (svr *Server) ReExec(afunc func(*Server, logs.Info) string, info logs.Info) {
	if svr.Status.NeedExec {
		res := afunc(svr, info)
		logs.GlobalLogger.AllPrint(fmt.Sprintf("[exec] %s | [%s|%s %s:%s]", res, svr.Name, svr.Mode, svr.Host, svr.Port))
		svr.Status.NeedExec = false
	}
}

func (svr *Server) Check(failFunc func(*Server, logs.Info) string, succeedFunc func(*Server, logs.Info) string, errorFunc func(*Server, logs.Info) string) (info logs.Info) {
	defer func() {
		if e := recover(); e != nil {
			info.Code = 0
			info.Message = fmt.Sprintf("Connect Failed | [%s|%s %s:%s]", svr.Name, svr.Mode, svr.Host, svr.Port)
			svr.ReExec(failFunc, info)
		}
	}()
	err := svr.TryConnect()
	if err != nil {
		info.Code = 2
		info.Message = err.Error()
		svr.ReExec(errorFunc, info)
	} else {
		info.Code = 1
		info.Message = fmt.Sprintf("Connect Success [%s|%s %s:%s]", svr.Name, svr.Mode, svr.Host, svr.Port)
		svr.ReExec(succeedFunc, info)
	}
	return
}
