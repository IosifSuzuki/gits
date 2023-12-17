package service

import (
	"encoding/json"
	"fmt"
	"gits/internal/container"
	"gits/internal/model/dto"
	"go.uber.org/zap"
	"net/http"
)

type Ip interface {
	GeoLocation(ip string) (*dto.GeoLocation, error)
}

type ip struct {
	container.Container
}

func NewIp(container container.Container) Ip {
	return &ip{
		container,
	}
}

func (i *ip) GeoLocation(ip string) (*dto.GeoLocation, error) {
	log := i.GetLogger()
	conf := i.GetConfig()

	url := fmt.Sprintf("http://ipinfo.io/%s?token=%s", ip, conf.IpInfo.Token)
	resp, err := http.Get(url)
	if err != nil {
		log.Error("make request to retrieve info about ip has failed", zap.Error(err))
		return nil, err
	}
	var geoLocation dto.GeoLocation
	if err := json.NewDecoder(resp.Body).Decode(&geoLocation); err != nil {
		log.Error("geo location model has failed map to struct", zap.Error(err))
		return nil, err
	}

	return &geoLocation, err
}
