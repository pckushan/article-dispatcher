package bootstrap

import (
	"article-dispatcher/internal/adaptors/cache"
	"article-dispatcher/internal/domain/adaptors/logger"
	"article-dispatcher/internal/http"
	"article-dispatcher/internal/pkg/configs"
	"article-dispatcher/internal/pkg/log"
	"article-dispatcher/internal/services"
	"fmt"
	sysLog "log"
	"os"
	"os/signal"
)

// Boot - initialize configurations, logger and create necessary dependency plugins
// init the router and serve the routes
func Boot() {
	initConfigs()
	l := initLogger()

	repo := cache.NewCache(l)
	articleService := services.NewArticleService(l, repo)

	r := &http.Router{
		Conf: &http.Config,
	}
	r.Init(l, articleService)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	go func() {
		<-signals
		if err := r.Stop(); err != nil {
			sysLog.Fatalln(fmt.Sprintf("failed to gracefully shutdown the server due to: %s", err))
		}
	}()

	err := r.Start()
	if err != nil {
		sysLog.Fatalln(fmt.Sprintf("failed to start the server due to: %s", err))
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
