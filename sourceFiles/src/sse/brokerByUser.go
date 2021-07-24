package sse

import (
	"github.com/gin-gonic/gin"
	"github.com/pepelazz/nla_framework/utils"
	"net/http"
)

var brokerByUser map[string]broker

func AddConn(c *gin.Context)  {

	userId, ok := utils.ExtractUserIdString(c)
	if !ok {
		utils.HttpError(c, http.StatusBadRequest,"missed user_id")
		return
	}

	if b, ok := brokerByUser[userId]; ok {
		b.subscribe(c)
	} else {
		b := broker{
			make(map[chan string]bool),
			make(chan (chan string)),
			make(chan (chan string)),
			make(chan string, 10), // buffer 10 msgs and don't block sends,
		}
		b.handleEvents()
		brokerByUser[userId] = b
		b.subscribe(c)
	}
}

func SendJson(userId string, d interface{})  {
	if b, ok := brokerByUser[userId]; ok {
		go b.sendJSON(d)
	}
}