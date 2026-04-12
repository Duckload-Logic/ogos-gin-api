package email

import "context"

type Emailer interface {
	SendEmail(ctx context.Context, to, subject, body string) (bool, error)
	SendOTP(ctx context.Context, to, otp string) (bool, error)
}
