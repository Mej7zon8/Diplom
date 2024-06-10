package advertisement

import (
	"crypto/rand"
	"github.com/k773/utils"
	"github.com/k773/utils/maps"
	"messenger/data/entities"
	"os"
	"path/filepath"
	"slices"
	"sync"
)

type Provider struct {
	s sync.RWMutex

	// basePath is the path where all category directories are stored.
	basePath string

	categories map[entities.AdCategoryLabel]*Category
}

func New(basePath string) *Provider {
	return &Provider{
		basePath:   basePath,
		categories: make(map[entities.AdCategoryLabel]*Category),
	}
}

func (p *Provider) Update() error {
	p.s.Lock()
	defer p.s.Unlock()

	entries, e := os.ReadDir(p.basePath)
	if e != nil {
		return e
	}
	var passed = map[entities.AdCategoryLabel]struct{}{}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		var catLabel = entities.AdCategoryLabel(entry.Name())

		cat, ok := p.categories[catLabel]
		if !ok {
			cat = newCategory(filepath.Join(p.basePath, entry.Name()))
			p.categories[catLabel] = cat
		}
		// Update the category:
		if e = cat.update(); e != nil {
			return e
		}
		// Mark the category as passed:
		passed[catLabel] = struct{}{}
	}

	// Remove the categories that are not in the directory:
	for name := range p.categories {
		if _, ok := passed[name]; !ok {
			delete(p.categories, name)
		}
	}
	return nil
}

func (p *Provider) RandomAdForLabel(label entities.AdCategoryLabel) *Ad {
	p.s.RLock()
	defer p.s.RUnlock()

	if cat, ok := p.categories[label]; ok {
		return cat.randomAd()
	}
	return nil
}

func (p *Provider) RandomAdForLabels(labels []entities.AdCategoryLabel) *Ad {
	p.s.RLock()
	defer p.s.RUnlock()

	var okLabels []entities.AdCategoryLabel
	for _, label := range labels {
		if cat, ok := p.categories[label]; ok && len(cat.ads) != 0 {
			okLabels = append(okLabels, label)
		}
	}
	if len(okLabels) == 0 {
		return nil
	}
	// Choose a random label:
	var label = utils.RandomChoiceMust(rand.Reader, okLabels)
	return p.categories[label].randomAd()
}

func (p *Provider) ListLabels() []entities.AdCategoryLabel {
	p.s.RLock()
	defer p.s.RUnlock()
	var k = maps.Keys(p.categories)
	slices.Sort(k)
	return k
}

// CreateCategory creates a new category with the provided label.
// Note that the category will be added to the provider only after the next Update call.
func (p *Provider) CreateCategory(label entities.AdCategoryLabel) error {
	p.s.Lock()
	defer p.s.Unlock()

	if _, ok := p.categories[label]; ok {
		return nil
	}
	// Create the directory:
	var fp = filepath.Join(p.basePath, string(label))
	e := os.MkdirAll(fp, 0600)
	if e != nil {
		return e
	}
	return nil
}
