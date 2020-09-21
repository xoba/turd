package taws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// NewSession creates a new default authenticated aws session.
func NewSession() (*session.Session, error) {
	return NewSessionFromProfile("")
}

// NewSessionFromProfile creates a new named authenticated aws session.
func NewSessionFromProfile(profile string) (*session.Session, error) {
	return session.NewSessionWithOptions(session.Options{
		Profile:           profile,
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{
			DisableRestProtocolURICleaning: aws.Bool(true),
		},
	})
}
