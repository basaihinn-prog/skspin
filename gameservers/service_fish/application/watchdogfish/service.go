package watchdogfish

import (
	PSF_ON_00001 "serve/service_fish/domain/fish/PSF-ON-00001"
	PSF_ON_00002 "serve/service_fish/domain/fish/PSF-ON-00002"
	PSF_ON_00003 "serve/service_fish/domain/fish/PSF-ON-00003"
	PSF_ON_00004 "serve/service_fish/domain/fish/PSF-ON-00004"
	PSF_ON_00005 "serve/service_fish/domain/fish/PSF-ON-00005"
	PSF_ON_00006 "serve/service_fish/domain/fish/PSF-ON-00006"
	PSF_ON_00007 "serve/service_fish/domain/fish/PSF-ON-00007"
	RKF_H5_00001 "serve/service_fish/domain/fish/RKF-H5-00001"
	"serve/service_fish/domain/maintain"
	PSFM_00001_98_1 "serve/service_fish/domain/probability/PSFM-00001-1/PSFM-00001-98-1"
	PSFM_00002_95_1 "serve/service_fish/domain/probability/PSFM-00002-1/PSFM-00002-95-1"
	PSFM_00002_96_1 "serve/service_fish/domain/probability/PSFM-00002-1/PSFM-00002-96-1"
	PSFM_00002_97_1 "serve/service_fish/domain/probability/PSFM-00002-1/PSFM-00002-97-1"
	PSFM_00002_98_1 "serve/service_fish/domain/probability/PSFM-00002-1/PSFM-00002-98-1"
	PSFM_00003_95_1 "serve/service_fish/domain/probability/PSFM-00003-1/PSFM-00003-95-1"
	PSFM_00003_96_1 "serve/service_fish/domain/probability/PSFM-00003-1/PSFM-00003-96-1"
	PSFM_00003_97_1 "serve/service_fish/domain/probability/PSFM-00003-1/PSFM-00003-97-1"
	PSFM_00003_98_1 "serve/service_fish/domain/probability/PSFM-00003-1/PSFM-00003-98-1"
	PSFM_00004_95_1 "serve/service_fish/domain/probability/PSFM-00004-1/PSFM-00004-95-1"
	PSFM_00004_96_1 "serve/service_fish/domain/probability/PSFM-00004-1/PSFM-00004-96-1"
	PSFM_00004_97_1 "serve/service_fish/domain/probability/PSFM-00004-1/PSFM-00004-97-1"
	PSFM_00004_98_1 "serve/service_fish/domain/probability/PSFM-00004-1/PSFM-00004-98-1"
	PSFM_00005_95_1 "serve/service_fish/domain/probability/PSFM-00005-1/PSFM-00005-95-1"
	PSFM_00005_96_1 "serve/service_fish/domain/probability/PSFM-00005-1/PSFM-00005-96-1"
	PSFM_00005_97_1 "serve/service_fish/domain/probability/PSFM-00005-1/PSFM-00005-97-1"
	PSFM_00005_98_1 "serve/service_fish/domain/probability/PSFM-00005-1/PSFM-00005-98-1"
	PSFM_00006_95_1 "serve/service_fish/domain/probability/PSFM-00006-1/PSFM-00006-95-1"
	PSFM_00006_96_1 "serve/service_fish/domain/probability/PSFM-00006-1/PSFM-00006-96-1"
	PSFM_00006_97_1 "serve/service_fish/domain/probability/PSFM-00006-1/PSFM-00006-97-1"
	PSFM_00006_98_1 "serve/service_fish/domain/probability/PSFM-00006-1/PSFM-00006-98-1"
	PSFM_00007_97_1 "serve/service_fish/domain/probability/PSFM-00007-1/PSFM-00007-97-1"
	PSFM_00008_97_1 "serve/service_fish/domain/probability/PSFM-00008-1/PSFM-00008-97-1"
	PSFM_00013_93_1 "serve/service_fish/domain/probability/PSFM-00013-1/PSFM-00013-93-1"
	PSFM_00013_94_1 "serve/service_fish/domain/probability/PSFM-00013-1/PSFM-00013-94-1"
	PSFM_00013_95_1 "serve/service_fish/domain/probability/PSFM-00013-1/PSFM-00013-95-1"
	PSFM_00013_96_1 "serve/service_fish/domain/probability/PSFM-00013-1/PSFM-00013-96-1"
	PSFM_00013_97_1 "serve/service_fish/domain/probability/PSFM-00013-1/PSFM-00013-97-1"
	PSFM_00013_98_1 "serve/service_fish/domain/probability/PSFM-00013-1/PSFM-00013-98-1"

	"serve/fish_comm/watchdog"

	"github.com/fsnotify/fsnotify"
)

var Service = &service{}

type service struct{}

func (s *service) AddPathHandler(watcher *fsnotify.Watcher) {
	// Math
	// PSF-ON-00001
	watchdog.Service.Add(PSFM_00001_98_1.RTP98BS.Docker, watcher)
	watchdog.Service.Add(PSFM_00002_95_1.RTP95BS.Docker, watcher)
	watchdog.Service.Add(PSFM_00002_96_1.RTP96BS.Docker, watcher)
	watchdog.Service.Add(PSFM_00002_97_1.RTP97BS.Docker, watcher)
	watchdog.Service.Add(PSFM_00002_98_1.RTP98BS.Docker, watcher)

	// PSF-ON-00002
	watchdog.Service.Add(PSFM_00003_95_1.RTP95BS.Docker, watcher)
	watchdog.Service.Add(PSFM_00003_96_1.RTP96BS.Docker, watcher)
	watchdog.Service.Add(PSFM_00003_97_1.RTP97BS.Docker, watcher)
	watchdog.Service.Add(PSFM_00003_98_1.RTP98BS.Docker, watcher)

	// PSF-ON-00003
	watchdog.Service.Add(PSFM_00004_95_1.RTP95BS.Docker, watcher)
	watchdog.Service.Add(PSFM_00004_96_1.RTP96BS.Docker, watcher)
	watchdog.Service.Add(PSFM_00004_97_1.RTP97BS.Docker, watcher)
	watchdog.Service.Add(PSFM_00004_98_1.RTP98BS.Docker, watcher)

	// PSF-ON-00004
	watchdog.Service.Add(PSFM_00005_95_1.RTP95BS.Docker, watcher)
	watchdog.Service.Add(PSFM_00005_96_1.RTP96BS.Docker, watcher)
	watchdog.Service.Add(PSFM_00005_97_1.RTP97BS.Docker, watcher)
	watchdog.Service.Add(PSFM_00005_98_1.RTP98BS.Docker, watcher)

	// PSF-ON-00005
	watchdog.Service.Add(PSFM_00006_95_1.RTP95BS.Docker, watcher)
	watchdog.Service.Add(PSFM_00006_96_1.RTP96BS.Docker, watcher)
	watchdog.Service.Add(PSFM_00006_97_1.RTP97BS.Docker, watcher)
	watchdog.Service.Add(PSFM_00006_98_1.RTP98BS.Docker, watcher)

	// PSF-ON-00006
	watchdog.Service.Add(PSFM_00007_97_1.RTP97BS.Docker, watcher)

	// PSF-ON-00007
	watchdog.Service.Add(PSFM_00008_97_1.RTP97BS.Docker, watcher)

	// RKF-H5-00001
	watchdog.Service.Add(PSFM_00013_93_1.RTP93BS.Docker, watcher)
	watchdog.Service.Add(PSFM_00013_94_1.RTP94BS.Docker, watcher)
	watchdog.Service.Add(PSFM_00013_95_1.RTP95BS.Docker, watcher)
	watchdog.Service.Add(PSFM_00013_96_1.RTP96BS.Docker, watcher)
	watchdog.Service.Add(PSFM_00013_97_1.RTP97BS.Docker, watcher)
	watchdog.Service.Add(PSFM_00013_98_1.RTP98BS.Docker, watcher)

	// Script
	watchdog.Service.Add(PSF_ON_00001.FishPath.Docker, watcher)
	watchdog.Service.Add(PSF_ON_00002.FishPath.Docker, watcher)
	watchdog.Service.Add(PSF_ON_00003.FishPath.Docker, watcher)
	watchdog.Service.Add(PSF_ON_00004.FishPath.Docker, watcher)
	watchdog.Service.Add(PSF_ON_00005.FishPath.Docker, watcher)
	watchdog.Service.Add(PSF_ON_00006.FishPath.Docker, watcher)
	watchdog.Service.Add(PSF_ON_00007.FishPath.Docker, watcher)
	watchdog.Service.Add(RKF_H5_00001.FishPath.Docker, watcher)

	// Maintain
	watchdog.Service.Add(maintain.Service.Docker, watcher)
}

func (s *service) FilterPathHandler(event *fsnotify.Event) {
	isDo := true

	if isDo {
		isDo = !s.psf_on_00001_watchDog(event)
	}

	if isDo {
		isDo = !s.psf_on_00002_watchDog(event)
	}

	if isDo {
		isDo = !s.psf_on_00003_watchDog(event)
	}

	if isDo {
		isDo = !s.psf_on_00004_watchDog(event)
	}

	if isDo {
		isDo = !s.psf_on_00005_watchDog(event)
	}

	if isDo {
		isDo = !s.psf_on_00006_watchDog(event)
	}

	if isDo {
		isDo = !s.psf_on_00007_watchDog(event)
	}

	if isDo {
		isDo = !s.rkf_h5_00001_watchDog(event)
	}

	if isDo {
		isDo = !s.psfm_on_00001_98_1_watchdog(event)
	}

	// PSFM-ON-00002
	if isDo {
		isDo = !s.psfm_on_00002_95_1_watchdog(event)
	}

	if isDo {
		isDo = !s.psfm_on_00002_96_1_watchdog(event)
	}

	if isDo {
		isDo = !s.psfm_on_00002_97_1_watchdog(event)
	}

	if isDo {
		isDo = !s.psfm_on_00002_98_1_watchdog(event)
	}

	// PSFM-ON-00003
	if isDo {
		isDo = !s.psfm_on_00003_93_1_watchdog(event)
	}

	if isDo {
		isDo = !s.psfm_on_00003_94_1_watchdog(event)
	}

	if isDo {
		isDo = !s.psfm_on_00003_95_1_watchdog(event)
	}

	if isDo {
		isDo = !s.psfm_on_00003_96_1_watchdog(event)
	}

	if isDo {
		isDo = !s.psfm_on_00003_97_1_watchdog(event)
	}

	if isDo {
		isDo = !s.psfm_on_00003_98_1_watchdog(event)
	}

	// PSFM-ON-00004
	if isDo {
		isDo = !s.psfm_on_00004_95_1_watchdog(event)
	}

	if isDo {
		isDo = !s.psfm_on_00004_96_1_watchdog(event)
	}

	if isDo {
		isDo = !s.psfm_on_00004_97_1_watchdog(event)
	}

	if isDo {
		isDo = !s.psfm_on_00004_98_1_watchdog(event)
	}

	// PSFM-ON-00005
	if isDo {
		isDo = !s.psfm_on_00005_95_1_watchdog(event)
	}

	if isDo {
		isDo = !s.psfm_on_00005_96_1_watchdog(event)
	}

	if isDo {
		isDo = !s.psfm_on_00005_97_1_watchdog(event)
	}

	if isDo {
		isDo = !s.psfm_on_00005_98_1_watchdog(event)
	}

	// PSFM-ON-00006
	if isDo {
		isDo = !s.psfm_on_00006_95_1_watchdog(event)
	}

	if isDo {
		isDo = !s.psfm_on_00006_96_1_watchdog(event)
	}

	if isDo {
		isDo = !s.psfm_on_00006_97_1_watchdog(event)
	}

	if isDo {
		isDo = !s.psfm_on_00006_98_1_watchdog(event)
	}

	// PSFM-ON-00007
	if isDo {
		isDo = !s.psfm_on_00007_97_1_watchdog(event)
	}

	// PSFM-ON-00008
	if isDo {
		isDo = !s.psfm_on_00008_97_1_watchdog(event)
	}

	// PSFM-ON-00013
	if isDo {
		isDo = !s.psfm_on_00013_93_1_watchdog(event)
	}

	if isDo {
		isDo = !s.psfm_on_00013_94_1_watchdog(event)
	}

	if isDo {
		isDo = !s.psfm_on_00013_95_1_watchdog(event)
	}

	if isDo {
		isDo = !s.psfm_on_00013_96_1_watchdog(event)
	}

	if isDo {
		isDo = !s.psfm_on_00013_97_1_watchdog(event)
	}

	if isDo {
		isDo = !s.psfm_on_00013_98_1_watchdog(event)
	}

	if isDo {
		isDo = !s.maintain_watchdog(event)
	}
}
