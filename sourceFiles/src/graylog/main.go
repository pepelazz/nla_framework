package graylog

import (
	"fmt"
	"github.com/pepelazz/nla_framework/types"
	"gopkg.in/aphistic/golf.v0"
)

var (
	Graylog *GraylogType
	appName string
)

type GraylogType struct {
	Client *golf.Client
}

func Init(config types.GraylogConfig) (err error) {
	Graylog = &GraylogType{}
	host := config.Host
	port := config.Port
	appName = config.AppName
	Graylog.Client, _ = golf.NewClient()
	err = Graylog.Client.Dial(fmt.Sprintf("udp://%s:%v", host, port))
	return
}

func (g *GraylogType) L() *golf.Logger {
	l, _ := g.Client.NewLogger()
	l.SetAttr("app", appName)
	return l
}

func (g *GraylogType) Close() error {
	return g.Client.Close()
}
