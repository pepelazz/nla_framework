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
	Attrs  map[string]string // список дополнительных аттрибутов для проекта
}

func Init(config types.GraylogConfig, attrs map[string]string) (err error) {
	Graylog = &GraylogType{}
	host := config.Host
	port := config.Port
	appName = config.AppName
	Graylog.Client, _ = golf.NewClient()
	Graylog.Attrs = attrs
	err = Graylog.Client.Dial(fmt.Sprintf("udp://%s:%v", host, port))
	return
}

func (g *GraylogType) L() *golf.Logger {
	l, _ := g.Client.NewLogger()
	l.SetAttr("application_name", appName)
	for k, v := range g.Attrs {
		l.SetAttr(k, v)
	}
	return l
}

func (g *GraylogType) Close() error {
	return g.Client.Close()
}
