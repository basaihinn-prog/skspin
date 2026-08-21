package watchdogfish

import (
	"path/filepath"
	PSFM_00002_95_1 "serve/service_fish/domain/probability/PSFM-00002-1/PSFM-00002-95-1"
	PSFM_00002_96_1 "serve/service_fish/domain/probability/PSFM-00002-1/PSFM-00002-96-1"
	PSFM_00002_97_1 "serve/service_fish/domain/probability/PSFM-00002-1/PSFM-00002-97-1"
	PSFM_00002_98_1 "serve/service_fish/domain/probability/PSFM-00002-1/PSFM-00002-98-1"

	"github.com/fsnotify/fsnotify"
)

func (s *service) psfm_on_00002_95_1_watchdog(event *fsnotify.Event) bool {
	switch event.Name {
	case filepath.Join(PSFM_00002_95_1.RTP95BS.Docker, PSFM_00002_95_1.RTP95BS.File):
		return s.mathReload(event, PSFM_00002_95_1.RTP95BS)

	case filepath.Join(PSFM_00002_95_1.RTP95FS.Docker, PSFM_00002_95_1.RTP95FS.File):
		return s.mathReload(event, PSFM_00002_95_1.RTP95FS)

	case filepath.Join(PSFM_00002_95_1.RTP95DRB.Docker, PSFM_00002_95_1.RTP95DRB.File):
		return s.mathReload(event, PSFM_00002_95_1.RTP95DRB)

	default:
		return false

	}
}

func (s *service) psfm_on_00002_96_1_watchdog(event *fsnotify.Event) bool {
	switch event.Name {
	case filepath.Join(PSFM_00002_96_1.RTP96BS.Docker, PSFM_00002_96_1.RTP96BS.File):
		return s.mathReload(event, PSFM_00002_96_1.RTP96BS)

	case filepath.Join(PSFM_00002_96_1.RTP96FS.Docker, PSFM_00002_96_1.RTP96FS.File):
		return s.mathReload(event, PSFM_00002_96_1.RTP96FS)

	case filepath.Join(PSFM_00002_96_1.RTP96DRB.Docker, PSFM_00002_96_1.RTP96DRB.File):
		return s.mathReload(event, PSFM_00002_96_1.RTP96DRB)

	default:
		return false

	}
}

func (s *service) psfm_on_00002_97_1_watchdog(event *fsnotify.Event) bool {
	switch event.Name {
	case filepath.Join(PSFM_00002_97_1.RTP97BS.Docker, PSFM_00002_97_1.RTP97BS.File):
		return s.mathReload(event, PSFM_00002_97_1.RTP97BS)

	case filepath.Join(PSFM_00002_97_1.RTP97FS.Docker, PSFM_00002_97_1.RTP97FS.File):
		return s.mathReload(event, PSFM_00002_97_1.RTP97FS)

	case filepath.Join(PSFM_00002_97_1.RTP97DRB.Docker, PSFM_00002_97_1.RTP97DRB.File):
		return s.mathReload(event, PSFM_00002_97_1.RTP97DRB)

	default:
		return false

	}
}

func (s *service) psfm_on_00002_98_1_watchdog(event *fsnotify.Event) bool {
	switch event.Name {
	case filepath.Join(PSFM_00002_98_1.RTP98BS.Docker, PSFM_00002_98_1.RTP98BS.File):
		return s.mathReload(event, PSFM_00002_98_1.RTP98BS)

	case filepath.Join(PSFM_00002_98_1.RTP98FS.Docker, PSFM_00002_98_1.RTP98FS.File):
		return s.mathReload(event, PSFM_00002_98_1.RTP98FS)

	case filepath.Join(PSFM_00002_98_1.RTP98DRB.Docker, PSFM_00002_98_1.RTP98DRB.File):
		return s.mathReload(event, PSFM_00002_98_1.RTP98DRB)

	default:
		return false

	}
}
