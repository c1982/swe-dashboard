// +build integ

package victoriametrics

import (
	"fmt"
	"testing"
	"time"
)

func TestPushWithTimeInteg(t *testing.T) {
	p, err := NewPusher(SetPushURL("http://localhost:8428"))
	if err != nil {
		t.Error(err)
	}

	day := time.Hour * 24
	for i := 1; i < 10; i++ {
		pushtime := time.Now().Add(-1 * (day * time.Duration(i)))
		metric := fmt.Sprintf(`A1{user="paul"} %d`, i)
		err = p.PushWithTime(metric, pushtime)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestQueryInteg(t *testing.T) {
	p, err := NewPusher(SetPushURL("http://localhost:8428"))
	if err != nil {
		t.Error(err)
	}

	result, err := p.Query(`match={__name__=~"A1"}`)
	if err != nil {
		t.Error("query error", err)
	}

	fmt.Println(result)
}

func TestFirstContactInteg(t *testing.T) {
	p, err := NewPusher(SetPushURL("http://localhost:8428"))
	if err != nil {
		t.Error(err)
	}

	ok, err := p.FirstContact()
	if err != nil {
		t.Error(err)
	}

	t.Log("firstcontact: ", ok)
}
