package types

type Notification struct {
	ToUserIds    []string
	ToCcUserIds  []string
	ToBccUserIds []string
	Subject      string
	Body         string
}
