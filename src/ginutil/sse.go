package ginutil

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

// ChannelToServerSentEvents consumes events from a channel and sends them as Server Sent Events
func ChannelToServerSentEvents(ctx *gin.Context, sourceChan chan interface{}) {
	resp := ctx.Writer
	h := resp.Header()
	h.Set("Cache-Control", "no-cache")
	h.Set("Connection", "keep-alive")
	h.Set("Content-Type", "text/event-stream")
	h.Set("X-Accel-Buffering", "no")

	resp.WriteHeader(http.StatusOK)
	resp.WriteHeaderNow()
	resp.Flush()

	ctx.Stream(func(w io.Writer) bool {
		if state, ok := <-sourceChan; ok {
			ctx.SSEvent("message", state)
			return true
		}
		return false
	})
}
