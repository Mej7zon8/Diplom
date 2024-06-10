package export_categories

import (
	"messenger/data/entities"
	userstore "messenger/data/store/user-store"
)

type exporter struct{}

type user struct {
	Username   entities.UserRef   `json:"username"`
	Name       string             `json:"name"`
	Email      string             `json:"email"`
	Categories map[string]float64 `json:"categories"`
}

func newExporter() exporter { return exporter{} }

func (exporter) Export() (res []user, e error) {
	data, e := userstore.Instance.GetAll()
	if e != nil {
		return nil, e
	}
	for _, u := range data {
		res = append(res, user{
			Username:   u.ID,
			Name:       u.Credentials.Name,
			Email:      u.Credentials.Email,
			Categories: u.Meta.Categories.Values,
		})
	}
	return res, nil
}
