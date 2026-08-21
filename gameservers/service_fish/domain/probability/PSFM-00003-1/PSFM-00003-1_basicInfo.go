package PSFM_00003_1

import (
	"encoding/json"
	"os"
	"path/filepath"
	PSF_ON_00002_1 "serve/service_fish/domain/probability/PSF-ON-00002-1"
	"serve/service_fish/models"
)

const (
	DOCKER string = "/data/math"
	FOLDER string = "domain/probability/PSFM-00003-1/"
	RTP93  string = models.PSFM_00003_93_1
	RTP94  string = models.PSFM_00003_94_1
	RTP95  string = models.PSFM_00003_95_1
	RTP96  string = models.PSFM_00003_96_1
	RTP97  string = models.PSFM_00003_97_1
	RTP98  string = models.PSFM_00003_98_1
)

var mathInfo = &MathInfo{}

type MathInfo struct {
	FileType               string
	Docker                 string
	Test                   string
	Folder                 string
	File                   string
	PSF_ON_00002_1_BsMath  *PSF_ON_00002_1.BsMath
	PSF_ON_00002_1_DrbMath *PSF_ON_00002_1.DrbMath
	PSF_ON_00002_1_FsMath  *PSF_ON_00002_1.FsMath
}

func (m *MathInfo) Open() *os.File {
	var file *os.File
	var err error

	defer file.Close()

	if file, err = os.Open(filepath.Join(m.Docker, m.File)); err != nil {
		if file, err = os.Open(filepath.Join(m.Folder, m.File)); err != nil {
			if file, err = os.Open(m.File); err != nil {
				if file, err = os.Open(filepath.Join(m.Test, m.File)); err != nil {
					panic(err.Error())
				}
			}
		}
	}
	return file
}

func (m *MathInfo) Deserialization() {
	switch m.FileType {
	case models.BS:
		if err := json.NewDecoder(m.Open()).Decode(&m.PSF_ON_00002_1_BsMath); err != nil {
			panic(err)
		}
	case models.DRB:
		if err := json.NewDecoder(m.Open()).Decode(&m.PSF_ON_00002_1_DrbMath); err != nil {
			panic(err)
		}
	case models.FS:
		if err := json.NewDecoder(m.Open()).Decode(&m.PSF_ON_00002_1_FsMath); err != nil {
			panic(err)
		}
	}
}

func (m *MathInfo) Reload() {
	switch m.FileType {
	case models.BS:
		m.PSF_ON_00002_1_BsMath = nil
	case models.DRB:
		m.PSF_ON_00002_1_DrbMath = nil
	case models.FS:
		m.PSF_ON_00002_1_FsMath = nil
	}

	m.Deserialization()
}
