package advertisement

import (
	"errors"
	"log/slog"
	"messenger/data/entities"
	"os"
	"time"
)

const categoriesBaseDir = "program-data/advertisement/"

var Instance = New(categoriesBaseDir)

func init() {
	var l = slog.With("module", "data/providers/advertisement", "method", "init")

	//Update the provider:
	l.Info("initialized", "error", tick(true))

	// Start parsing labels:
	go func() {
		for range time.Tick(time.Second) {
			if e := tick(false); e != nil {
				l.Error("(async) Update failed", "error", e)
			}
		}
	}()
}

// sampleCategories is the list of labels that would be created if the labels were not loaded.
var sampleCategories = []entities.AdCategoryLabel{
	"travel",
	"fashion",
	"clothes",
	"insurance",
	"politics",
	"movie",
	"fitness",
	"health",
	"computer",
	"education",
	"finance",
	"sports",
	"music",
	"food",
	"entertainment",
	"shopping",
	"gaming",
	"photography",
	"art",
	"real-estate",
	"automotive",
	"adventure",
	"cooking",
	"parenting",
	"lifestyle",
	"wellness",
	"marketing",
	"business",
	"investment",
	"environment",
	"skincare",
	"nutrition",
	"vacation",
	"home-decor",
	"pets",
	"gadgets",
	"startup",
}

func tick(allowCreate bool) (e error) {
	e = Instance.Update()
	if errors.Is(e, os.ErrNotExist) && allowCreate {
		e = createLabels(sampleCategories)
	}
	return
}

func createLabels(labels []entities.AdCategoryLabel) (e error) {
	for _, label := range labels {
		if e = Instance.CreateCategory(label); e != nil {
			return e
		}
	}
	// Update the provider:
	e = Instance.Update()
	return e
}
