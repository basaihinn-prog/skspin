package PSFM_00001_1

import (
	"encoding/json"
	"os"
	"path/filepath"
	PSF_ON_00001_1 "serve/service_fish/domain/probability/PSF-ON-00001-1"
	"serve/service_fish/models"
)

const (
	DOCKER string = "/data/math"
	FOLDER string = "domain/probability/PSFM-00001-1/"
	RTP98  string = models.PSFM_00001_98_1
)

var mathInfo = &MathInfo{}

type MathInfo struct {
	FileType              string
	Docker                string
	Test                  string
	Folder                string
	File                  string
	PSF_ON_00001_1_BsMath *PSF_ON_00001_1.BsMath
	PSF_ON_00001_1_FsMath *PSF_ON_00001_1.FsMath
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
		if err := json.NewDecoder(m.Open()).Decode(&m.PSF_ON_00001_1_BsMath); err != nil {
			panic(err)
		}
	case models.FS:
		if err := json.NewDecoder(m.Open()).Decode(&m.PSF_ON_00001_1_FsMath); err != nil {
			panic(err)
		}
	}
}

func (m *MathInfo) Reload() {
	switch m.FileType {
	case models.BS:
		m.PSF_ON_00001_1_BsMath = nil
	case models.FS:
		m.PSF_ON_00001_1_FsMath = nil
	}

	m.Deserialization()
}
