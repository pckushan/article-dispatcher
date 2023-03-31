package bootstrap

import (
	"article-dispatcher/internal/domain/adaptors/logger"
	"article-dispatcher/internal/http"
	"article-dispatcher/internal/pkg/configs"
	"article-dispatcher/internal/pkg/log"
	"fmt"
	sysLog "log"
	"os"
	"os/signal"
)

func Boot() {
	initConfigs()
	l := initLogger()

	r := &http.Router{
		Conf: &http.Config,
	}
	r.Init(l)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	go func() {
		<-signals
		if err := r.Stop(); err != nil {
			sysLog.Fatalln(fmt.Sprintf("failed to gracefully shutdown the server: %s", err))
		}
	}()

	err := r.Start()
	if err != nil {
		sysLog.Fatal(err)
	}

}

func initConfigs() {
	err := configs.Load(
		new(http.RouterConfig),
		new(log.LoggerConfig),
	)

	if err != nil {
		sysLog.Fatalln("error loading configs due to: ", err)
	}
}

func initLogger() logger.Logger {
	l, err := log.NewLogger(log.Config.Level)
	if err != nil {
		sysLog.Fatalln("error loading new logger due to: ", err)
	}
	return l
}
