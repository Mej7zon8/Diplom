package advertisement

import (
	"messenger/data/entities"
	"messenger/data/providers/advertisement"
	userstore "messenger/data/store/user-store"
	"messenger/web/api"
	"net/http"
)

func Handle(w http.ResponseWriter, r *http.Request, user entities.UserRef) {
	api.WithTyped(w, r, func(r *http.Request) api.ResponseAny {
		switch r.FormValue("method") {
		case "get":
			return handleGet(w, r, user)
		default:
			return api.ResponseErrPathNotFound
		}
	})
}

type ad struct {
	// Image is a base64 encoded image.
	Image  string   `json:"image"`
	Config adConfig `json:"config"`
}

type adConfig struct {
	URL string `json:"url"`
}

// handleCheck checks whether the user is authenticated
func handleGet(w http.ResponseWriter, r *http.Request, userRef entities.UserRef) api.ResponseAny {
	// Get the advertisement for the user
	var u, e = userstore.Instance.Read(userRef)
	if e != nil {
		return api.ResponseBadRequest(e)
	}
	var allUserLabels []entities.AdCategoryLabel
	for k, v := range u.Meta.Categories.Values {
		if v < 10 {
			continue
		}
		allUserLabels = append(allUserLabels, entities.AdCategoryLabel(k))
	}

	adObj := advertisement.Instance.RandomAdForLabels(allUserLabels)
	if adObj == nil {
		return api.ResponseAny{Status: 204}
	}

	return api.ResponseAny{Status: 200, Data: ad{
		Image:  adObj.GetImage().Base64(),
		Config: adConfig{URL: adObj.GetConfig().URL},
	}}
}
