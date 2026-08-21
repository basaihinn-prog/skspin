package watchdogfish

import (
	"reflect"
	"serve/fish_comm/flux"
	"serve/fish_comm/logger"
	"serve/service_fish/domain/fish"

	"github.com/fsnotify/fsnotify"
)

func (s *service) mathReload(event *fsnotify.Event, file interface{}) bool {
	logger.Service.Zap.Infow("WatchDog", "Event", event)
	reflect.ValueOf(file).MethodByName("Reload").Call([]reflect.Value{})

	return true
}

func (s *service) scriptReload(event *fsnotify.Event, check bool) bool {
	logger.Service.Zap.Infow("WatchDog", "Event", event)
	if check {
		flux.Send(fish.ActionFishRestart, event.Name, fish.Service.Id, "")
	}

	return true
}
