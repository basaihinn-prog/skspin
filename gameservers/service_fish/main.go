package main

import (
	_ "net/http/pprof"
	Version "serve/fish_comm/version"
	"serve/service_fish/application/watchdogfish"
	broadcasterFish "serve/service_fish/domain/broadcaster"
	"serve/service_fish/presentation/playerfish"

	"serve/fish_comm/broadcaster"
	"serve/fish_comm/flux"
	"serve/fish_comm/logger"
	"serve/fish_comm/mysql"
	"serve/fish_comm/player"
	"serve/fish_comm/watchdog"
)

func main() {

	flux.Start(logger.Service)
	watchdog.Service.Start(watchdogfish.Service)
	broadcaster.Service.Start(broadcasterFish.Service)
	Version.Service.Start("asd")
	//jackpot.Service.Start()

	// must place last
	player.Controller.Start(playerfish.ControllerFish)

	defer destroy()
}

func destroy() {
	logger.Service.Destroy()
	mysql.Repository.Destroy()
}
