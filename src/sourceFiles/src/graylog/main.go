package graylog

import (
	"github.com/pepelazz/projectGenerator/src/types"
	"gopkg.in/aphistic/golf.v0"
	"fmt"
)

var (
	Graylog *GraylogType
)

type GraylogType struct {
	Client *golf.Client
}

func Init(config types.GraylogConfig) (err error) {
	Graylog = &GraylogType{}
	host := config.Host
	port := config.Port
	Graylog.Client, _ = golf.NewClient()
	err = Graylog.Client.Dial(fmt.Sprintf("udp://%s:%v", host, port))
	return
}

func (g *GraylogType) L() (*golf.Logger) {
	l, _ := g.Client.NewLogger()
	l.SetAttr("app", "fourPl")
	return l
}

func (g *GraylogType) Close() error {
	return g.Client.Close()
}
