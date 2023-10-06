package models

type FactMysql struct {
	PeriodStart       string     `json:"period_start"`
	PeriodEnd         string     `json:"period_end"`
	KeyPeriod         string     `json:"key_period"`
	IndicatorToMoId   int        `json:"indicator_to_mo_id"`
	IndicatorToFactId int        `json:"indicator_to_mo_fact_id"`
	Value             int        `json:"value"`
	FactTime          string     `json:"fact_time"`
	IsPlan            int        `json:"is_plan"`
	Supertags         []SuperTag `json:"supertags"`
	AuthUserId        int        `json:"auth_user_id"`
	Comment           string     `json:"comment"`
}
type SuperTag struct {
	Tag   Tag    `json:"tag"`
	Value string `json:"value"`
}

type Tag struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Key          string `json:"key"`
	ValuesSource int    `json:"values_source"`
}
