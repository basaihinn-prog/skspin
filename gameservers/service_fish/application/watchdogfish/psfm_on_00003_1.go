package watchdogfish

import (
	"path/filepath"
	PSFM_00003_93_1 "serve/service_fish/domain/probability/PSFM-00003-1/PSFM-00003-93-1"
	PSFM_00003_94_1 "serve/service_fish/domain/probability/PSFM-00003-1/PSFM-00003-94-1"
	PSFM_00003_95_1 "serve/service_fish/domain/probability/PSFM-00003-1/PSFM-00003-95-1"
	PSFM_00003_96_1 "serve/service_fish/domain/probability/PSFM-00003-1/PSFM-00003-96-1"
	PSFM_00003_97_1 "serve/service_fish/domain/probability/PSFM-00003-1/PSFM-00003-97-1"
	PSFM_00003_98_1 "serve/service_fish/domain/probability/PSFM-00003-1/PSFM-00003-98-1"

	"github.com/fsnotify/fsnotify"
)

func (s *service) psfm_on_00003_93_1_watchdog(event *fsnotify.Event) bool {
	switch event.Name {
	case filepath.Join(PSFM_00003_93_1.RTP93BS.Docker, PSFM_00003_93_1.RTP93BS.File):
		return s.mathReload(event, PSFM_00003_93_1.RTP93BS)

	case filepath.Join(PSFM_00003_93_1.RTP93FS.Docker, PSFM_00003_93_1.RTP93FS.File):
		return s.mathReload(event, PSFM_00003_93_1.RTP93FS)

	case filepath.Join(PSFM_00003_93_1.RTP93DRB.Docker, PSFM_00003_93_1.RTP93DRB.File):
		return s.mathReload(event, PSFM_00003_93_1.RTP93DRB)

	default:
		return false
	}
}

func (s *service) psfm_on_00003_94_1_watchdog(event *fsnotify.Event) bool {
	switch event.Name {
	case filepath.Join(PSFM_00003_94_1.RTP94BS.Docker, PSFM_00003_94_1.RTP94BS.File):
		return s.mathReload(event, PSFM_00003_94_1.RTP94BS)

	case filepath.Join(PSFM_00003_94_1.RTP94FS.Docker, PSFM_00003_94_1.RTP94FS.File):
		return s.mathReload(event, PSFM_00003_94_1.RTP94FS)

	case filepath.Join(PSFM_00003_94_1.RTP94DRB.Docker, PSFM_00003_94_1.RTP94DRB.File):
		return s.mathReload(event, PSFM_00003_94_1.RTP94DRB)

	default:
		return false
	}
}

func (s *service) psfm_on_00003_95_1_watchdog(event *fsnotify.Event) bool {
	switch event.Name {
	case filepath.Join(PSFM_00003_95_1.RTP95BS.Docker, PSFM_00003_95_1.RTP95BS.File):
		return s.mathReload(event, PSFM_00003_95_1.RTP95BS)

	case filepath.Join(PSFM_00003_95_1.RTP95FS.Docker, PSFM_00003_95_1.RTP95FS.File):
		return s.mathReload(event, PSFM_00003_95_1.RTP95FS)

	case filepath.Join(PSFM_00003_95_1.RTP95DRB.Docker, PSFM_00003_95_1.RTP95DRB.File):
		return s.mathReload(event, PSFM_00003_95_1.RTP95DRB)

	default:
		return false
	}
}

func (s *service) psfm_on_00003_96_1_watchdog(event *fsnotify.Event) bool {
	switch event.Name {
	case filepath.Join(PSFM_00003_96_1.RTP96BS.Docker, PSFM_00003_96_1.RTP96BS.File):
		return s.mathReload(event, PSFM_00003_96_1.RTP96BS)

	case filepath.Join(PSFM_00003_96_1.RTP96FS.Docker, PSFM_00003_96_1.RTP96FS.File):
		return s.mathReload(event, PSFM_00003_96_1.RTP96FS)

	case filepath.Join(PSFM_00003_96_1.RTP96DRB.Docker, PSFM_00003_96_1.RTP96DRB.File):
		return s.mathReload(event, PSFM_00003_96_1.RTP96DRB)

	default:
		return false
	}
}

func (s *service) psfm_on_00003_97_1_watchdog(event *fsnotify.Event) bool {
	switch event.Name {
	case filepath.Join(PSFM_00003_97_1.RTP97BS.Docker, PSFM_00003_97_1.RTP97BS.File):
		return s.mathReload(event, PSFM_00003_97_1.RTP97BS)

	case filepath.Join(PSFM_00003_97_1.RTP97FS.Docker, PSFM_00003_97_1.RTP97FS.File):
		return s.mathReload(event, PSFM_00003_97_1.RTP97FS)

	case filepath.Join(PSFM_00003_97_1.RTP97DRB.Docker, PSFM_00003_97_1.RTP97DRB.File):
		return s.mathReload(event, PSFM_00003_97_1.RTP97DRB)

	default:
		return false
	}
}

func (s *service) psfm_on_00003_98_1_watchdog(event *fsnotify.Event) bool {
	switch event.Name {
	case filepath.Join(PSFM_00003_98_1.RTP98BS.Docker, PSFM_00003_98_1.RTP98BS.File):
		return s.mathReload(event, PSFM_00003_98_1.RTP98BS)

	case filepath.Join(PSFM_00003_98_1.RTP98FS.Docker, PSFM_00003_98_1.RTP98FS.File):
		return s.mathReload(event, PSFM_00003_98_1.RTP98FS)

	case filepath.Join(PSFM_00003_98_1.RTP98DRB.Docker, PSFM_00003_98_1.RTP98DRB.File):
		return s.mathReload(event, PSFM_00003_98_1.RTP98DRB)

	default:
		return false

	}
}
