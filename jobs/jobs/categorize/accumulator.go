package categorize

import "messenger/data/entities"

type accumulator struct {
	user2label2value map[entities.UserRef]map[string]float64
}

func newAccumulator() *accumulator {
	return &accumulator{
		user2label2value: make(map[entities.UserRef]map[string]float64),
	}
}

func (a *accumulator) add(user entities.UserRef, label string, value float64) {
	if _, ok := a.user2label2value[user]; !ok {
		a.user2label2value[user] = make(map[string]float64)
	}
	a.user2label2value[user][label] += value
}

func (a *accumulator) export() map[entities.UserRef]map[string]float64 {
	return a.user2label2value
}
