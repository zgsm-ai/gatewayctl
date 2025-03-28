package main

import (
	"flag"

	"github.com/zgsm-ai/gatewayctl/internal"
	"github.com/zgsm-ai/gatewayctl/internal/pkg/config"
	"github.com/zgsm-ai/gatewayctl/internal/pkg/logger"
	"github.com/zgsm-ai/gatewayctl/internal/store"
	"github.com/zgsm-ai/gatewayctl/internal/store/postgres"
)

func main() {
	var envConf = flag.String("conf", "config/config.yaml", "config path, eg: -conf ./config/config.yaml")
	flag.Parse()

	config.InitConfig(*envConf)
	logger.InitLogger(logger.NewOptsFromConfig())

	db, err := postgres.GetDBInstance()
	if err != nil {
		panic(err)
	}
	store.NewPluginModel(db)

	internal.InitRouter(db)
}
