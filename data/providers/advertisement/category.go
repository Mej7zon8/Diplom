package advertisement

import (
	"os"
	"path/filepath"
)

type Category struct {
	path string
	ads  map[string]*Ad
}

func newCategory(path string) *Category {
	return &Category{path: path, ads: make(map[string]*Ad)}
}

// update performs a reload of the Category entries.
// the access to this method must be synchronized.
func (c *Category) update() error {
	entries, e := os.ReadDir(c.path)
	if e != nil {
		return e
	}
	var passed = map[string]struct{}{}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		ad, ok := c.ads[entry.Name()]
		if !ok {
			ad = newAd(filepath.Join(c.path, entry.Name()))
			c.ads[entry.Name()] = ad
		}
		// Update the ad:
		if e = ad.update(); e != nil {
			return e
		}
		// Mark the ad as passed:
		passed[entry.Name()] = struct{}{}
	}

	// Remove the ads that are not in the directory:
	for name := range c.ads {
		if _, ok := passed[name]; !ok {
			delete(c.ads, name)
		}
	}
	return nil
}

func (c *Category) randomAd() *Ad {
	for _, ad := range c.ads {
		return ad
	}
	return nil
}
