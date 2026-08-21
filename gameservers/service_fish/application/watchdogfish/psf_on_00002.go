package watchdogfish

import (
	"path/filepath"
	PSF_ON_00002 "serve/service_fish/domain/fish/PSF-ON-00002"

	"github.com/fsnotify/fsnotify"
)

func (s *service) psf_on_00002_watchDog(event *fsnotify.Event) bool {

	switch event.Name {

	case filepath.Join(PSF_ON_00002.FishPath.Docker, PSF_ON_00002.FishPath.File):
		return s.scriptReload(event, PSF_ON_00002.FishPath.Reload())

	case filepath.Join(PSF_ON_00002.Groups.Docker, PSF_ON_00002.Groups.File):
		return s.scriptReload(event, PSF_ON_00002.Groups.Reload())

	case filepath.Join(PSF_ON_00002.Objects.Docker, PSF_ON_00002.Objects.File):
		return s.scriptReload(event, PSF_ON_00002.Objects.Reload())

	case filepath.Join(PSF_ON_00002.ScriptA1.Docker, PSF_ON_00002.ScriptA1.File):
		return s.scriptReload(event, PSF_ON_00002.ScriptA1.Reload())

	case filepath.Join(PSF_ON_00002.ScriptA2.Docker, PSF_ON_00002.ScriptA2.File):
		return s.scriptReload(event, PSF_ON_00002.ScriptA2.Reload())

	case filepath.Join(PSF_ON_00002.ScriptA3.Docker, PSF_ON_00002.ScriptA3.File):
		return s.scriptReload(event, PSF_ON_00002.ScriptA3.Reload())

	case filepath.Join(PSF_ON_00002.ScriptB1.Docker, PSF_ON_00002.ScriptB1.File):
		return s.scriptReload(event, PSF_ON_00002.ScriptB1.Reload())

	case filepath.Join(PSF_ON_00002.ScriptB2.Docker, PSF_ON_00002.ScriptB2.File):
		return s.scriptReload(event, PSF_ON_00002.ScriptB2.Reload())

	case filepath.Join(PSF_ON_00002.ScriptB3.Docker, PSF_ON_00002.ScriptB3.File):
		return s.scriptReload(event, PSF_ON_00002.ScriptB3.Reload())

	case filepath.Join(PSF_ON_00002.ScriptC1.Docker, PSF_ON_00002.ScriptC1.File):
		return s.scriptReload(event, PSF_ON_00002.ScriptC1.Reload())

	case filepath.Join(PSF_ON_00002.ScriptC2.Docker, PSF_ON_00002.ScriptC2.File):
		return s.scriptReload(event, PSF_ON_00002.ScriptC2.Reload())

	case filepath.Join(PSF_ON_00002.ScriptC3.Docker, PSF_ON_00002.ScriptC3.File):
		return s.scriptReload(event, PSF_ON_00002.ScriptC3.Reload())

	case filepath.Join(PSF_ON_00002.ScriptD1.Docker, PSF_ON_00002.ScriptD1.File):
		return s.scriptReload(event, PSF_ON_00002.ScriptD1.Reload())

	case filepath.Join(PSF_ON_00002.ScriptD2.Docker, PSF_ON_00002.ScriptD2.File):
		return s.scriptReload(event, PSF_ON_00002.ScriptD2.Reload())

	case filepath.Join(PSF_ON_00002.ScriptD3.Docker, PSF_ON_00002.ScriptD3.File):
		return s.scriptReload(event, PSF_ON_00002.ScriptD3.Reload())

	default:
		return false

	}
}
