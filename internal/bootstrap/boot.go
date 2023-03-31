package bootstrap

import (
	"article-dispatcher/internal/domain/adaptors/logger"
	"article-dispatcher/internal/pkg/configs"
	"article-dispatcher/internal/pkg/log"
	sysLog "log"
)

func Boot() {
	initConfigs()
	_ = initLogger()

}

func initConfigs() {
	err := configs.Load(
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
