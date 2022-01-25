package victoriametrics

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

//curl -d 'foo{bar="baz"} 123' -X POST 'http://localhost:8428/api/v1/import/prometheus?timestamp=1626957551000'
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
	url := fmt.Sprintf("%s/api/v1/import/prometheus", p.host)
	_, err := p.httpPost(url, payload)
	return err
}

func (p *Pusher) PushWithTime(payload string, date time.Time) error {
	timestamp := fmt.Sprintf("%d000", date.Unix())
	postURL := fmt.Sprintf("%s/api/v1/import/prometheus?timestamp=%s", p.host, timestamp)
	_, err := p.httpPost(postURL, payload)
	return err
}

func (p *Pusher) Query(payload string) (result string, err error) {
	url := fmt.Sprintf("%s/api/v1/export", p.host)
	result, err = p.httpPost(url, payload)
	return result, err
}

func (p *Pusher) FirstContact() (bool, error) {
	result, err := p.Query(`match={__name__=~"first_contact"}`)
	if err != nil {
		return true, err
	}

	fmt.Println(result)
	if result != "" {
		return true, nil
	}

	err = p.Push("first_contact 1")
	if err != nil {
		return true, err
	}

	return false, nil
}

func (p *Pusher) setHost(pushURL string) error {
	p.host = pushURL
	return nil
}

func (p *Pusher) httpPost(postURL, payload string) (string, error) {
	reqbody := bytes.NewBuffer([]byte(payload))
	resp, err := http.Post(postURL, "application/x-www-form-urlencoded", reqbody)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode > 204 {
		return "", fmt.Errorf("status code: %d, msg: %s", resp.StatusCode, string(body))
	}

	return string(body), nil
}
