package handler

import (
	"fmt"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

type PanicHandler struct{}

func NewPanicHandler() PanicHandler {
	return PanicHandler{}
}

func (i PanicHandler) Handle(c *gin.Context) {
	panic("Panic at path /panic")
}

func RecoveryHandler(c *gin.Context, err interface{}) {
	sentry.CaptureException(fmt.Errorf("catched a panic: %v", err))
	c.HTML(500, "error.tpl", gin.H{
		"title": "Internal server error",
		"err":   err,
	})
}

func (i PanicHandler) Log(c *gin.Context) {
	//Отправить лог в sentry
	sentry.CaptureMessage("log is handled and sent to sentry")
	c.JSON(http.StatusOK, gin.H{"msg": "sent to sentry"})
}
