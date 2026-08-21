package watchdogfish

import (
	"path/filepath"
	PSFM_00001_98_1 "serve/service_fish/domain/probability/PSFM-00001-1/PSFM-00001-98-1"

	"github.com/fsnotify/fsnotify"
)

func (s *service) psfm_on_00001_98_1_watchdog(event *fsnotify.Event) bool {
	switch event.Name {
	case filepath.Join(PSFM_00001_98_1.RTP98BS.Docker, PSFM_00001_98_1.RTP98BS.File):
		return s.mathReload(event, PSFM_00001_98_1.RTP98BS)

	case filepath.Join(PSFM_00001_98_1.RTP98FS.Docker, PSFM_00001_98_1.RTP98FS.File):
		return s.mathReload(event, PSFM_00001_98_1.RTP98FS)

	default:
		return false

	}
}
