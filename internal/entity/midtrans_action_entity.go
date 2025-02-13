package entity

import "seblak-bombom-restful-api/internal/helper"

type Action struct {
	ID                     uint64                `gorm:"primaryKey;column:id;autoIncrement"`
	MidtransCoreAPIOrderId uint64                `gorm:"column:midtrans_core_api_orders_id"`
	Name                   string                `gorm:"column:name"`
	Method                 helper.RequestMethod  `gorm:"column:method"`
	URL                    string                `gorm:"column:url"`
	MidtransCoreAPIOrder   *MidtransCoreAPIOrder `gorm:"foreignKey:midtrans_core_api_orders_id;references:id"`
}

func (u *Action) TableName() string {
	return "midtrans_actions"
}
