package victoriametrics

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

//curl -d 'foo{bar="baz"} 123' -X POST 'http://localhost:8428/api/v1/import/prometheus'
//curl -G 'http://localhost:8428/api/v1/export' -d 'match={__name__=~"foo"}'

type VictoriaMetricOption func(*Pusher) error

func SetPushURL(pushURL string) VictoriaMetricOption {
	return func(p *Pusher) error {
		return p.setHost(pushURL)
	}
}

type Pusher struct {
	host string
}

func NewPusher(options ...VictoriaMetricOption) (pusher *Pusher, err error) {
	p := &Pusher{}
	for _, fn := range options {
		if fn == nil {
			continue
		}
		if err := fn(p); err != nil {
			return p, err
		}
	}

	return p, nil
}

func (p *Pusher) Push(payload string) error {
	return p.httpPost(p.host, payload)
}

func (p *Pusher) PushWithTime(payload string, date time.Time) error {
	timestamp := date.Unix()
	postURL := fmt.Sprintf("%s?timestamp=%d", p.host, timestamp)
	return p.httpPost(postURL, payload)
}

func (s *Pusher) setHost(pushURL string) error {
	s.host = pushURL
	return nil
}

func (s *Pusher) httpPost(postURL, payload string) error {
	responseBody := bytes.NewBuffer([]byte(payload))
	resp, err := http.Post(postURL, "application/x-www-form-urlencoded", responseBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode > 204 {
		return fmt.Errorf("status code: %d, msg: %s", resp.StatusCode, string(body))
	}

	return nil
}
