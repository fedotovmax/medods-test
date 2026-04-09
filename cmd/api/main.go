package main

import (
	"github.com/fedotovmax/medods-test/internal/core/logger/zap"
	httpServer "github.com/fedotovmax/medods-test/internal/core/transport/http/server"
)

func main() {

	log, err := zap.New(zap.NewConfigMust())

	if err != nil {
		panic(err)
	}

	defer log.Stop()

	serverHTPP, err := httpServer.New(httpServer.NewConfigMust(), log)

	if err != nil {
		panic(err)
	}

	serverHTPP.RouteGroup("/api", func(r httpServer.Router) {

	})

}
