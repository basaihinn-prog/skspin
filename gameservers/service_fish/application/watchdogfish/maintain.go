package watchdogfish

import (
	"path/filepath"
	"serve/fish_comm/broadcaster"
	"serve/fish_comm/flux"
	jackpot_client "serve/fish_comm/jackpot-client"
	"serve/fish_comm/logger"
	"serve/fish_comm/mysql"
	Maintain "serve/service_fish/domain/maintain"

	"github.com/fsnotify/fsnotify"
)

func (s *service) maintain_watchdog(event *fsnotify.Event) bool {
	switch event.Name {
	case filepath.Join(Maintain.Service.Docker, Maintain.Service.File):
		logger.Service.Zap.Infow("WatchDog", "Event", event)

		if Maintain.Service.Reload() {
			if Maintain.Service.Shutdown() {
				flux.Send(broadcaster.ActionBroadcasterShutdown, event.Name, broadcaster.Service.Id, Maintain.Service.Message())
			}

			if Maintain.Service.DatabaseReload() {
				mysql.Repository.Reload()
			}
		}
		return true

	case filepath.Join(jackpot_client.Config.Docker, jackpot_client.Config.File):
		logger.Service.Zap.Infow("WatchDog", "Event", event)
		jackpot_client.Config.Reload()

		return true

	default:
		return false

	}
}
