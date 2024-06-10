package message_categorization

import (
	"context"
	"testing"
)

func TestProcessor_Process(t *testing.T) {
	res, e := Instance.Process(context.Background(), "Пошли покушаем?")
	if e != nil {
		t.Fatal(e)
	}
	t.Log(res)
}
