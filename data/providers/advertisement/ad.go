package advertisement

import (
	"encoding/json"
	"errors"
	"messenger/data/entities"
	adimage "messenger/data/providers/advertisement/ad-image"
	"os"
	"path/filepath"
	"time"
)

type Ad struct {
	dirPath string

	image   *adimage.Image
	config  entities.AdConfig
	modInfo modifiedInfo
}

func newAd(path string) *Ad {
	return &Ad{dirPath: path}
}

func (a *Ad) GetImage() *adimage.Image {
	return a.image
}

func (a *Ad) GetConfig() entities.AdConfig {
	return a.config
}

// the access to this method must be synchronized.
func (a *Ad) update() error {
	modInfo, isModified, err := a.isModified()
	if err != nil || !isModified {
		return err
	}

	// Load the image:
	if modInfo.image != a.modInfo.image {
		img, e := adimage.Load(a.getImagePath())
		if e != nil {
			return e
		}
		a.image = img
	}
	// Load the config:
	if modInfo.config != a.modInfo.config {
		cfgFileData, e := os.ReadFile(a.getConfigPath())
		if e != nil {
			return e
		}
		var cfg entities.AdConfig
		if e = json.Unmarshal(cfgFileData, &cfg); e != nil {
			return e
		}
		a.config = cfg
	}

	a.modInfo = modInfo
	return nil
}

type modifiedInfo struct {
	image  time.Time
	config time.Time
}

func (a *Ad) isModified() (modifiedInfo, bool, error) {
	imageStat, e := os.Stat(a.getImagePath())
	if e != nil {
		return modifiedInfo{}, false, e
	}
	configStat, e := os.Stat(a.getConfigPath())
	if e != nil {
		return modifiedInfo{}, false, e
	}

	if imageStat.ModTime().IsZero() || configStat.ModTime().IsZero() {
		return modifiedInfo{}, false, errors.New("advertisement: file modification time is zero")
	}

	var modInfo = modifiedInfo{
		image:  imageStat.ModTime(),
		config: configStat.ModTime(),
	}
	return modInfo, a.modInfo != modInfo, nil
}

func (a *Ad) getImagePath() string {
	return filepath.Join(a.dirPath, "image.png")
}

func (a *Ad) getConfigPath() string {
	return filepath.Join(a.dirPath, "config.json")
}
