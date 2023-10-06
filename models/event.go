package models

type Event struct {
	Id     string      `json:"_id"`
	Key    string      `json:"_key"`
	Rev    string      `json:"_rev"`
	Author Author      `json:"author"`
	Group  string      `json:"group"`
	Msg    string      `json:"msg"`
	Time   string      `json:"time"`
	Type   string      `json:"type"`
	Params EventParams `json:"params"`
}

type EventParams struct {
	IndicatorToMoId int    `json:"indicator_to_mo_id"`
	Platform        string `json:"platform"`
	Period          Period `json:"period"`
}

type Period struct {
	Start   string `json:"start"`
	End     string `json:"end"`
	TypeKey string `json:"type_key"`
}
