package entity

import "time"

type MidtransSnapOrder struct {
	ID          uint64    `gorm:"primary_key;column:id;autoIncrement"`
	OrderId     uint64    `gorm:"column:order_id"`
	Token       string    `gorm:"column:token"`
	RedirectUrl string    `gorm:"column:redirect_url"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime;<-:create"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	Order       *Order    `gorm:"foreignKey:order_id;references:id"`
}

func (m *MidtransSnapOrder) TableName() string {
	return "midtrans_snap_orders"
}
