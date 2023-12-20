package html

import (
	"fmt"
	stor "gits/internal/model/storage"
	"gits/internal/utils"
	"strings"
	"time"
)

type Action struct {
	IP        *string
	Location  *string
	Path      *string
	Browser   *string
	User      *string
	Flag      *string
	UpdatedAt *time.Time
}

func NewAction(observable *stor.Observable) (*Action, error) {
	return &Action{
		IP:        observable.Ip.Ip,
		Location:  getLocation(observable),
		Path:      observable.Path,
		Browser:   getHardwareInfo(observable),
		User:      getUsername(observable),
		Flag:      getFlagURL(observable),
		UpdatedAt: &observable.UpdatedAt,
	}, nil
}

func getLocation(observable *stor.Observable) *string {
	var location string
	if ip := observable.Ip; ip.Country != nil && ip.Region != nil && ip.City != nil {
		location = fmt.Sprintf("%s %s", *ip.Region, *ip.City)
	} else {
		location = "unknown"
	}

	return &location
}

func getUsername(observable *stor.Observable) *string {
	var user string
	if account := observable.Account; account != nil {
		user = account.Username
	} else {
		user = "anonymous"
	}

	return &user
}

func getFlagURL(observable *stor.Observable) *string {
	var flagURL *string

	country := observable.Ip.Country
	if country == nil {
		flagURL = utils.String("https://www.kidlink.org//icons/e.gif")
	} else {
		countryIso2Code := utils.String(strings.ToLower(*country))
		flagURL = utils.String(fmt.Sprintf("https://www.kidlink.org//icons/f0-%s.gif", *countryIso2Code))
	}

	return flagURL
}

func getHardwareInfo(observable *stor.Observable) *string {
	var browser *string

	if observable.Browser != nil && observable.OS != nil && observable.OSVersion != nil {
		hardwareInfo := fmt.Sprintf("os: %s %s; browser: %s", *observable.OS, *observable.OSVersion, *observable.Browser)
		browser = &hardwareInfo
	} else {
		browser = nil
	}

	return browser
}
