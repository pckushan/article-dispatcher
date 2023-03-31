package bootstrap

import (
	"article-dispatcher/internal/pkg/configs"
	"article-dispatcher/internal/pkg/logs"
	sysLog "log"
)

func Boot() {
	initConfigs()
}

func initConfigs() {
	err := configs.Load(
		new(logs.LoggerConfig),
	)

	if err != nil {
		sysLog.Fatalln("error loading configs due to: ", err)
	}
}
