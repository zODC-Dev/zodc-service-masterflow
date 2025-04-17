package types

type Notification struct {
	ToUserIds    []string
	ToCcUserIds  []string
	ToBccUserIds []string
	ToBccEmails  []string
	ToCcEmails   []string
	ToEmails     []string
	Subject      string
	Body         string
}
