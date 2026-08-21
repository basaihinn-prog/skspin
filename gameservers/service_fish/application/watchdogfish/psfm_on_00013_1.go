package watchdogfish

import (
	"path/filepath"
	PSFM_00013_93_1 "serve/service_fish/domain/probability/PSFM-00013-1/PSFM-00013-93-1"
	PSFM_00013_94_1 "serve/service_fish/domain/probability/PSFM-00013-1/PSFM-00013-94-1"
	PSFM_00013_95_1 "serve/service_fish/domain/probability/PSFM-00013-1/PSFM-00013-95-1"
	PSFM_00013_96_1 "serve/service_fish/domain/probability/PSFM-00013-1/PSFM-00013-96-1"
	PSFM_00013_97_1 "serve/service_fish/domain/probability/PSFM-00013-1/PSFM-00013-97-1"
	PSFM_00013_98_1 "serve/service_fish/domain/probability/PSFM-00013-1/PSFM-00013-98-1"

	"github.com/fsnotify/fsnotify"
)

func (s *service) psfm_on_00013_93_1_watchdog(event *fsnotify.Event) bool {
	switch event.Name {
	case filepath.Join(PSFM_00013_93_1.RTP93BS.Docker, PSFM_00013_93_1.RTP93BS.File):
		return s.mathReload(event, PSFM_00013_93_1.RTP93BS)
	case filepath.Join(PSFM_00013_93_1.RTP93FS.Docker, PSFM_00013_93_1.RTP93FS.File):
		return s.mathReload(event, PSFM_00013_93_1.RTP93FS)
	case filepath.Join(PSFM_00013_93_1.RTP93DRB.Docker, PSFM_00013_93_1.RTP93DRB.File):
		return s.mathReload(event, PSFM_00013_93_1.RTP93DRB)
	default:
		return false
	}
}

func (s *service) psfm_on_00013_94_1_watchdog(event *fsnotify.Event) bool {
	switch event.Name {
	case filepath.Join(PSFM_00013_94_1.RTP94BS.Docker, PSFM_00013_94_1.RTP94BS.File):
		return s.mathReload(event, PSFM_00013_94_1.RTP94BS)
	case filepath.Join(PSFM_00013_94_1.RTP94FS.Docker, PSFM_00013_94_1.RTP94FS.File):
		return s.mathReload(event, PSFM_00013_94_1.RTP94FS)
	case filepath.Join(PSFM_00013_94_1.RTP94DRB.Docker, PSFM_00013_94_1.RTP94DRB.File):
		return s.mathReload(event, PSFM_00013_94_1.RTP94DRB)
	default:
		return false
	}
}

func (s *service) psfm_on_00013_95_1_watchdog(event *fsnotify.Event) bool {
	switch event.Name {
	case filepath.Join(PSFM_00013_95_1.RTP95BS.Docker, PSFM_00013_95_1.RTP95BS.File):
		return s.mathReload(event, PSFM_00013_95_1.RTP95BS)
	case filepath.Join(PSFM_00013_95_1.RTP95FS.Docker, PSFM_00013_95_1.RTP95FS.File):
		return s.mathReload(event, PSFM_00013_95_1.RTP95FS)
	case filepath.Join(PSFM_00013_95_1.RTP95DRB.Docker, PSFM_00013_95_1.RTP95DRB.File):
		return s.mathReload(event, PSFM_00013_95_1.RTP95DRB)
	default:
		return false
	}
}

func (s *service) psfm_on_00013_96_1_watchdog(event *fsnotify.Event) bool {
	switch event.Name {
	case filepath.Join(PSFM_00013_96_1.RTP96BS.Docker, PSFM_00013_96_1.RTP96BS.File):
		return s.mathReload(event, PSFM_00013_96_1.RTP96BS)
	case filepath.Join(PSFM_00013_96_1.RTP96FS.Docker, PSFM_00013_96_1.RTP96FS.File):
		return s.mathReload(event, PSFM_00013_96_1.RTP96FS)
	case filepath.Join(PSFM_00013_96_1.RTP96DRB.Docker, PSFM_00013_96_1.RTP96DRB.File):
		return s.mathReload(event, PSFM_00013_96_1.RTP96DRB)
	default:
		return false
	}
}

func (s *service) psfm_on_00013_97_1_watchdog(event *fsnotify.Event) bool {
	switch event.Name {
	case filepath.Join(PSFM_00013_97_1.RTP97BS.Docker, PSFM_00013_97_1.RTP97BS.File):
		return s.mathReload(event, PSFM_00013_97_1.RTP97BS)
	case filepath.Join(PSFM_00013_97_1.RTP97FS.Docker, PSFM_00013_97_1.RTP97FS.File):
		return s.mathReload(event, PSFM_00013_97_1.RTP97FS)
	case filepath.Join(PSFM_00013_97_1.RTP97DRB.Docker, PSFM_00013_97_1.RTP97DRB.File):
		return s.mathReload(event, PSFM_00013_97_1.RTP97DRB)
	default:
		return false
	}
}

func (s *service) psfm_on_00013_98_1_watchdog(event *fsnotify.Event) bool {
	switch event.Name {
	case filepath.Join(PSFM_00013_98_1.RTP98BS.Docker, PSFM_00013_98_1.RTP98BS.File):
		return s.mathReload(event, PSFM_00013_98_1.RTP98BS)
	case filepath.Join(PSFM_00013_98_1.RTP98FS.Docker, PSFM_00013_98_1.RTP98FS.File):
		return s.mathReload(event, PSFM_00013_98_1.RTP98FS)
	case filepath.Join(PSFM_00013_98_1.RTP98DRB.Docker, PSFM_00013_98_1.RTP98DRB.File):
		return s.mathReload(event, PSFM_00013_98_1.RTP98DRB)
	default:
		return false
	}
}
