package raspberrypi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

type RaspberryPi struct {
	URL string
}

func NewRaspberryPi() *RaspberryPi {
	return &RaspberryPi{URL: os.Getenv("RASPBERRY_PI_URL")}
}

func (p RaspberryPi) TurnOn() (int, error) {
	req, err := http.Post(fmt.Sprintf("%s/%s", p.URL, "on"), "application/json", nil)
	if err != nil {
		log.Error("%+v", err)
		return 0, err
	}

	return req.StatusCode, nil
}

func (p RaspberryPi) TurnOff() (int, error) {
	req, err := http.Post(fmt.Sprintf("%s/%s", p.URL, "off"), "application/json", nil)
	if err != nil {
		log.Error("%+v", err)
		return 0, err
	}

	return req.StatusCode, nil
}

func (p RaspberryPi) GetTemperature() (string, int, error) {
	req, err := http.Get(fmt.Sprintf("%s/%s", p.URL, "temperature"))
	if err != nil {
		log.Error("%+v", err)
		return "", 0, err
	}

	body, err := ioutil.ReadAll(req.Body)
	var tmpt float64
	if err := json.Unmarshal(body, &tmpt); err != nil {
		log.Error("%+v", err)
		return "", 0, err
	}
	if req.StatusCode != http.StatusOK {
		return "", req.StatusCode, err
	}

	return fmt.Sprintf("%.1f", tmpt), req.StatusCode, nil
}
