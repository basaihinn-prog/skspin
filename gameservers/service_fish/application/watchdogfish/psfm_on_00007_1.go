package watchdogfish

import (
	"path/filepath"
	PSFM_00007_97_1 "serve/service_fish/domain/probability/PSFM-00007-1/PSFM-00007-97-1"

	"github.com/fsnotify/fsnotify"
)

func (s *service) psfm_on_00007_97_1_watchdog(event *fsnotify.Event) bool {
	switch event.Name {
	case filepath.Join(PSFM_00007_97_1.RTP97BS.Docker, PSFM_00007_97_1.RTP97BS.File):
		return s.mathReload(event, PSFM_00007_97_1.RTP97BS)

	case filepath.Join(PSFM_00007_97_1.RTP97FS.Docker, PSFM_00007_97_1.RTP97FS.File):
		return s.mathReload(event, PSFM_00007_97_1.RTP97FS)

	case filepath.Join(PSFM_00007_97_1.RTP97DRB.Docker, PSFM_00007_97_1.RTP97DRB.File):
		return s.mathReload(event, PSFM_00007_97_1.RTP97DRB)

	default:
		return false
	}
}
