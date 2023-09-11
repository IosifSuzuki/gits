package service

import (
	"encoding/json"
	"fmt"
	"gits/internal/container"
	"gits/internal/model/app"
	"gits/internal/model/service/response"
	"go.uber.org/zap"
	"net/http"
)

type Ip interface {
	GeoLocation(ip string) (*app.GeoLocation, error)
}

type ip struct {
	container.Container
}

func NewIp(container container.Container) Ip {
	return &ip{
		container,
	}
}

func (i *ip) GeoLocation(ip string) (*app.GeoLocation, error) {
	log := i.GetLogger()
	conf := i.GetConfig()

	url := fmt.Sprintf("https://ipinfo.io/%s?token=%s", ip, conf.IpInfo.Token)
	resp, err := http.Get(url)
	if err != nil {
		log.Error("make request to retrieve info about ip has failed", zap.Error(err))
		return nil, err
	}
	var geoLocation response.GeoLocation
	if err := json.NewDecoder(resp.Body).Decode(&geoLocation); err != nil {
		log.Error("geo location model has failed map to struct", zap.Error(err))
		return nil, err
	}

	return &app.GeoLocation{
		Ip:       geoLocation.Ip,
		Hostname: geoLocation.Hostname,
		City:     geoLocation.City,
		Region:   geoLocation.Region,
		Country:  geoLocation.Country,
		Loc:      geoLocation.Loc,
		Org:      geoLocation.Org,
		Postal:   geoLocation.Postal,
		Timezone: geoLocation.Timezone,
	}, nil
}
