package watchdogfish

import (
	"path/filepath"
	PSF_ON_00004 "serve/service_fish/domain/fish/PSF-ON-00004"

	"github.com/fsnotify/fsnotify"
)

func (s *service) psf_on_00004_watchDog(event *fsnotify.Event) bool {
	switch event.Name {
	case filepath.Join(PSF_ON_00004.FishPath.Docker, PSF_ON_00004.FishPath.File):
		return s.scriptReload(event, PSF_ON_00004.FishPath.Reload())

	case filepath.Join(PSF_ON_00004.Groups.Docker, PSF_ON_00004.Groups.File):
		return s.scriptReload(event, PSF_ON_00004.Groups.Reload())

	case filepath.Join(PSF_ON_00004.Objects.Docker, PSF_ON_00004.Objects.File):
		return s.scriptReload(event, PSF_ON_00004.Objects.Reload())

	case filepath.Join(PSF_ON_00004.ScriptA1.Docker, PSF_ON_00004.ScriptA1.File):
		return s.scriptReload(event, PSF_ON_00004.ScriptA1.Reload())

	case filepath.Join(PSF_ON_00004.ScriptA2.Docker, PSF_ON_00004.ScriptA2.File):
		return s.scriptReload(event, PSF_ON_00004.ScriptA2.Reload())

	case filepath.Join(PSF_ON_00004.ScriptA3.Docker, PSF_ON_00004.ScriptA3.File):
		return s.scriptReload(event, PSF_ON_00004.ScriptA3.Reload())

	case filepath.Join(PSF_ON_00004.ScriptB1.Docker, PSF_ON_00004.ScriptB1.File):
		return s.scriptReload(event, PSF_ON_00004.ScriptB1.Reload())

	case filepath.Join(PSF_ON_00004.ScriptB2.Docker, PSF_ON_00004.ScriptB2.File):
		return s.scriptReload(event, PSF_ON_00004.ScriptB2.Reload())

	case filepath.Join(PSF_ON_00004.ScriptB3.Docker, PSF_ON_00004.ScriptB3.File):
		return s.scriptReload(event, PSF_ON_00004.ScriptB3.Reload())

	case filepath.Join(PSF_ON_00004.ScriptC1.Docker, PSF_ON_00004.ScriptC1.File):
		return s.scriptReload(event, PSF_ON_00004.ScriptC1.Reload())

	case filepath.Join(PSF_ON_00004.ScriptC2.Docker, PSF_ON_00004.ScriptC2.File):
		return s.scriptReload(event, PSF_ON_00004.ScriptC2.Reload())

	case filepath.Join(PSF_ON_00004.ScriptC3.Docker, PSF_ON_00004.ScriptC3.File):
		return s.scriptReload(event, PSF_ON_00004.ScriptC3.Reload())

	case filepath.Join(PSF_ON_00004.ScriptD1.Docker, PSF_ON_00004.ScriptD1.File):
		return s.scriptReload(event, PSF_ON_00004.ScriptD1.Reload())

	case filepath.Join(PSF_ON_00004.ScriptD2.Docker, PSF_ON_00004.ScriptD2.File):
		return s.scriptReload(event, PSF_ON_00004.ScriptD2.Reload())

	case filepath.Join(PSF_ON_00004.ScriptD3.Docker, PSF_ON_00004.ScriptD3.File):
		return s.scriptReload(event, PSF_ON_00004.ScriptD3.Reload())

	default:
		return false
	}
}
