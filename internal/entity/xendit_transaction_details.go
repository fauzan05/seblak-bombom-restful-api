package entity

import (
	"time"
)

type XenditTransactionDetails struct {
	ID                  uint64              `gorm:"primary_key;column:id;autoIncrement"`
	XenditTransactionId string              `gorm:"column:xendit_transaction_id"`
	ReceiptId           string              `gorm:"column:receipt_id"`
	Source              string              `gorm:"column:source"`
	Name                string              `gorm:"column:name"`
	AccountDetails      string              `gorm:"column:account_details"`
	Created_At          time.Time           `gorm:"column:created_at;autoCreateTime;<-:create"`
	Updated_At          time.Time           `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	XenditTransactions  *XenditTransactions `gorm:"foreignKey:xendit_transaction_id;references:id"`
}

func (u *XenditTransactionDetails) TableName() string {
	return "xendit_transaction_details"
}
