package bootstrap

import (
	"article-dispatcher/internal/adaptors/cache"
	"article-dispatcher/internal/domain/adaptors/logger"
	"article-dispatcher/internal/http"
	"article-dispatcher/internal/pkg/configs"
	"article-dispatcher/internal/pkg/log"
	"article-dispatcher/internal/pkg/metrics"
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
	m := initMetrics(l)

	// plugin a cache to the repository
	repo := cache.NewCache(l)
	articleService := services.NewArticleService(l, repo)

	r := &http.Router{
		Conf: &http.Config,
	}
	r.Init(l, articleService, metrics.RequestLatency)

	// interrupt channel to stop the servers
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// channel to ensure graceful shutdown of the servers
	exitAll := make(chan bool, 1)

	go func() {
		<-signals
		if err := r.Stop(); err != nil {
			sysLog.Fatalf(fmt.Sprintf("failed to gracefully shutdown the server due to: %s", err))
		}
		if err := m.Stop(); err != nil {
			sysLog.Fatalf(fmt.Sprintf("failed to gracefully shutdown the metrics server due to: %s", err))
		}
		exitAll <- true
	}()

	go func() {
		err := r.Start()
		if err != nil {
			sysLog.Fatalf(fmt.Sprintf("failed to start the server due to: %s", err))
		}
	}()
	<-exitAll
}

// initConfigs - initialize configurations
func initConfigs() {
	err := configs.Load(
		new(http.RouterConfig),
		new(log.LoggerConfig),
		new(metrics.MetricConfig),
	)

	if err != nil {
		sysLog.Fatalln("error loading configs due to: ", err)
	}
}

// initLogger - init logger with log level defined in the environment
func initLogger() logger.Logger {
	l, err := log.NewLogger(log.Config.Level)
	if err != nil {
		sysLog.Fatalln("error loading new logger due to: ", err)
	}
	return l
}

// initMetrics - init metrics and start non-blocking metrics router
func initMetrics(l logger.Logger) *metrics.RouterMetrics {
	m := &metrics.RouterMetrics{
		Conf:   &metrics.Conf,
		Logger: l,
	}
	err := m.InitMetrics()
	if err != nil {
		sysLog.Fatal(err)
	}

	go func() {
		err = m.Start()
		if err != nil {
			sysLog.Fatal(err)
		}
	}()
	return m
}
