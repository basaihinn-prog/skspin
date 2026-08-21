package PSFM_00013_1

import (
	"encoding/json"
	"os"
	"path/filepath"
	RKF_H5_00001_1 "serve/service_fish/domain/probability/RKF-H5-00001-1"
	"serve/service_fish/models"
)

const (
	DOCKER string = "/data/math"
	FOLDER string = "domain/probability/PSFM-00013-1/"
	RTP93  string = models.PSFM_00013_93_1
	RTP94  string = models.PSFM_00013_94_1
	RTP95  string = models.PSFM_00013_95_1
	RTP96  string = models.PSFM_00013_96_1
	RTP97  string = models.PSFM_00013_97_1
	RTP98  string = models.PSFM_00013_98_1
)

var mathInfo = &MathInfo{}

type MathInfo struct {
	FileType               string
	Docker                 string
	Test                   string
	Folder                 string
	File                   string
	RKF_H5_00001_1_BsMath  *RKF_H5_00001_1.BsMath
	RKF_H5_00001_1_DrbMath *RKF_H5_00001_1.DrbMath
	RKF_H5_00001_1_FsMath  *RKF_H5_00001_1.FsMath
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
		if err := json.NewDecoder(m.Open()).Decode(&m.RKF_H5_00001_1_BsMath); err != nil {
			panic(err)
		}
	case models.DRB:
		if err := json.NewDecoder(m.Open()).Decode(&m.RKF_H5_00001_1_DrbMath); err != nil {
			panic(err)
		}
	case models.FS:
		if err := json.NewDecoder(m.Open()).Decode(&m.RKF_H5_00001_1_FsMath); err != nil {
			panic(err)
		}
	}
}

func (m *MathInfo) Reload() {
	switch m.FileType {
	case models.BS:
		m.RKF_H5_00001_1_BsMath = nil
	case models.DRB:
		m.RKF_H5_00001_1_DrbMath = nil
	case models.FS:
		m.RKF_H5_00001_1_FsMath = nil
	}

	m.Deserialization()
}
