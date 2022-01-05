package victoriametrics

import "testing"

func TestRun(t *testing.T) {
	p, err := NewPusher(SetPushURL("http://localhost:8428/api/v1/import/prometheus"))
	if err != nil {
		t.Error(err)
	}

	err = p.Push(`selfmerged{user="paul"} 1`)
	if err != nil {
		t.Error(err)
	}
}
