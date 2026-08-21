package watchdogfish

import (
	"path/filepath"
	RKF_H5_00001 "serve/service_fish/domain/fish/RKF-H5-00001"

	"github.com/fsnotify/fsnotify"
)

func (s *service) rkf_h5_00001_watchDog(event *fsnotify.Event) bool {
	switch event.Name {
	case filepath.Join(RKF_H5_00001.FishPath.Docker, RKF_H5_00001.FishPath.File):
		return s.scriptReload(event, RKF_H5_00001.FishPath.Reload())

	case filepath.Join(RKF_H5_00001.Groups.Docker, RKF_H5_00001.Groups.File):
		return s.scriptReload(event, RKF_H5_00001.Groups.Reload())

	case filepath.Join(RKF_H5_00001.Objects.Docker, RKF_H5_00001.Objects.File):
		return s.scriptReload(event, RKF_H5_00001.Objects.Reload())

	case filepath.Join(RKF_H5_00001.ScriptA1.Docker, RKF_H5_00001.ScriptA1.File):
		return s.scriptReload(event, RKF_H5_00001.ScriptA1.Reload())

	case filepath.Join(RKF_H5_00001.ScriptA2.Docker, RKF_H5_00001.ScriptA2.File):
		return s.scriptReload(event, RKF_H5_00001.ScriptA2.Reload())

	case filepath.Join(RKF_H5_00001.ScriptA3.Docker, RKF_H5_00001.ScriptA3.File):
		return s.scriptReload(event, RKF_H5_00001.ScriptA3.Reload())

	case filepath.Join(RKF_H5_00001.ScriptB1.Docker, RKF_H5_00001.ScriptB1.File):
		return s.scriptReload(event, RKF_H5_00001.ScriptB1.Reload())

	case filepath.Join(RKF_H5_00001.ScriptB2.Docker, RKF_H5_00001.ScriptB2.File):
		return s.scriptReload(event, RKF_H5_00001.ScriptB2.Reload())

	case filepath.Join(RKF_H5_00001.ScriptB3.Docker, RKF_H5_00001.ScriptB3.File):
		return s.scriptReload(event, RKF_H5_00001.ScriptB3.Reload())

	case filepath.Join(RKF_H5_00001.ScriptC1.Docker, RKF_H5_00001.ScriptC1.File):
		return s.scriptReload(event, RKF_H5_00001.ScriptC1.Reload())

	case filepath.Join(RKF_H5_00001.ScriptC2.Docker, RKF_H5_00001.ScriptC2.File):
		return s.scriptReload(event, RKF_H5_00001.ScriptC2.Reload())

	case filepath.Join(RKF_H5_00001.ScriptC3.Docker, RKF_H5_00001.ScriptC3.File):
		return s.scriptReload(event, RKF_H5_00001.ScriptC3.Reload())

	case filepath.Join(RKF_H5_00001.ScriptD1.Docker, RKF_H5_00001.ScriptD1.File):
		return s.scriptReload(event, RKF_H5_00001.ScriptD1.Reload())

	case filepath.Join(RKF_H5_00001.ScriptD2.Docker, RKF_H5_00001.ScriptD2.File):
		return s.scriptReload(event, RKF_H5_00001.ScriptD2.Reload())

	case filepath.Join(RKF_H5_00001.ScriptD3.Docker, RKF_H5_00001.ScriptD3.File):
		return s.scriptReload(event, RKF_H5_00001.ScriptD3.Reload())

	default:
		return false
	}
}
