package utils

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/glog"
	"github.com/urfave/cli/v2"
)

func AuthHeader(token string) http.Header {
	if len(token) != 0 {
		headers := http.Header{}
		headers.Add("Authorization", "Bearer "+string(token))
		return headers
	}
	glog.Warning("API Token not set and requested, capabilities might be limited.")
	return nil
}

func DaemonContext(cctx *cli.Context) context.Context {
	return context.Background()
}

func ReqContext(cctx *cli.Context) context.Context {
	tCtx := DaemonContext(cctx)

	ctx, done := context.WithCancel(tCtx)
	sigChan := make(chan os.Signal, 2)
	go func() {
		<-sigChan
		done()
	}()
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	return ctx
}
