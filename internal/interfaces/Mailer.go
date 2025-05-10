package interfaces

import "seblak-bombom-restful-api/internal/model"

type Mailer interface {
	Send(mail model.Mail) error
}