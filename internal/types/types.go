package types

import "time"


type QuoteIngress struct {
	Name string `json:"name"`
}

type QuoteEgress struct {
	Message string `json:"message"`
}

type Register struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	UUID        string   `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	CreatedAt time.Time
}


type TestPayloadChild struct{
	MockCounter int `json:"mockCounter,omitempty"`
	MockId int `json:"mockId,omitempty"`
}

type TestPayload struct{
	Msg string `json:"msg"`
	Data TestPayloadChild `json:"data"`
}

type TestPayloadB struct{
	Foo string `json:"foo"`
	Blah int `json:"blah"`
}
