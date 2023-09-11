package storage

import "gorm.io/gorm"

type Ip struct {
	gorm.Model
	Ip         *string
	Hostname   *string
	City       *string
	Region     *string
	Country    *string
	Loc        *string
	Org        *string
	Postal     *string
	Timezone   *string
	Observable []Observable `gorm:"foreignKey:IpId"`
}

func (c *Ip) TableName() string {
	return "ip"
}
