package html

import (
	"fmt"
	"gits/internal/helper"
	"gits/internal/model/app"
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

func NewAction(observable app.Observable) (*Action, error) {
	var (
		location *string
		user     *string
		browser  *string
	)
	if ip := observable.Ip; ip != nil {
		if ip.Country != nil && ip.Region != nil && ip.City != nil {
			location = helper.String(fmt.Sprintf("%s %s", *ip.Region, *ip.City))
		} else {
			location = helper.String("unknown")
		}
	} else {
		location = helper.String("unknown")
	}
	if account := observable.Account; account != nil {
		user = helper.String(account.Username)
	} else {
		user = helper.String("anonymous")
	}
	var (
		countryIso2Code *string
		flagUrl         *string
	)
	if country := observable.Ip.Country; country != nil {
		countryIso2Code = helper.String(strings.ToLower(*country))
	}
	if countryIso2Code != nil {
		flagUrl = helper.String(fmt.Sprintf("https://www.kidlink.org//icons/f0-%s.gif", *countryIso2Code))
	} else {
		flagUrl = helper.String("https://www.kidlink.org//icons/e.gif")
	}
	if observable.Browser != nil && observable.OS != nil && observable.OSVersion != nil {
		browser = helper.String(
			fmt.Sprintf("os: %s %s; browser: %s", *observable.OS, *observable.OSVersion, *observable.Browser),
		)
	}
	return &Action{
		IP:        observable.Ip.Ip,
		Location:  location,
		Path:      observable.Path,
		Browser:   browser,
		User:      user,
		Flag:      flagUrl,
		UpdatedAt: observable.UpdatedAt,
	}, nil
}
