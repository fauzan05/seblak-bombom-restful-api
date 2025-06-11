package entity

import (
	"database/sql"
	"seblak-bombom-restful-api/internal/helper/enum_state"
	"time"
)

type Payout struct {
	ID             uint64                  `gorm:"primary_key;column:id;autoIncrement"`
	UserId         uint64                  `gorm:"column:user_id"`
	XenditPayoutId sql.NullString          `gorm:"column:xendit_payout_id"`
	Amount         float32                 `gorm:"column:amount"`
	Currency       string                  `gorm:"column:currency"`
	Method         enum_state.PayoutMethod `gorm:"column:method"`
	Status         enum_state.PayoutStatus `gorm:"column:status"`
	Notes          string                  `gorm:"column:notes"`
	CreatedAt      time.Time               `gorm:"column:created_at;autoCreateTime;<-:create"`
	UpdatedAt      time.Time               `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	User           *User                   `gorm:"foreignKey:user_id;references:id"`
	XenditPayout   *XenditPayout           `gorm:"foreignKey:xendit_payout_id;references:id"`
}

func (u *Payout) TableName() string {
	return "payouts"
}
