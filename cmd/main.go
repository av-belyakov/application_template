package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/av-belyakov/application_template/internal/appname"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGINT)

	go func() {
		<-ctx.Done()

		fmt.Printf("Module '%s' is stop", appname.GetName())

		stop()
	}()

	app(ctx)
}
