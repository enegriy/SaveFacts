package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"kpi-test/models"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const events_body_filter string = `{
    "filter": {
        "field": {
            "key": "type",
            "sign": "LIKE",
            "values": [
                "MATRIX_REQUEST"
            ]
        }
    },
    "sort": {
        "fields": [
            "time"
        ],
        "direction": "desc"
    },
    "limit": 10
}`

func main() {

	cookie, err := getCookie()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	events, err := getEvents(cookie)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	facts, err := convertEventsToFacts(events)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	saveFacts(cookie, facts)
}

func pretty_string(str string) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}

func saveFacts(cookie string, facts []*models.FactMysql) (bool, error) {

	fmt.Println("7. Сохраняем факты development.kpi-drive.ru/_api/facts/save_fact")

	for i := 0; i < len(facts); i++ {
		fact := facts[i]
		supertag, _ := json.Marshal(fact.Supertags)

		form := url.Values{
			"period_start":            {fact.PeriodStart},
			"period_end":              {fact.PeriodEnd},
			"key_period":              {fact.FactTime},
			"indicator_to_mo_id":      {strconv.Itoa(fact.IndicatorToMoId)},
			"indicator_to_mo_fact_id": {strconv.Itoa(fact.IndicatorToFactId)},
			"value":                   {strconv.Itoa(fact.Value)},
			"fact_time":               {fact.FactTime},
			"is_plan":                 {strconv.Itoa(fact.IsPlan)},
			"auth_user_id":            {strconv.Itoa(fact.AuthUserId)},
			"comment":                 {fact.Comment},
			"supertags":               {string(supertag[:])}}

		req, err := http.NewRequest("POST", "https://development.kpi-drive.ru/_api/facts/save_fact", strings.NewReader(form.Encode()))
		if err != nil {
			return false, err
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Cookie", cookie)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return false, err
		}

		if resp.StatusCode == 200 {
			defer resp.Body.Close()
			resp_body, _ := io.ReadAll(resp.Body)

			fmt.Println("Данные факта сохранены: Ответ от сервера " + string(resp_body))
		}
	}

	return true, nil

}

func convertEventsToFacts(events []*models.Event) ([]*models.FactMysql, error) {

	fmt.Println("5. Конвертируем события в факты")

	facts := make([]*models.FactMysql, len(events))
	for i := 0; i < len(events); i++ {
		event := events[i]
		fact := models.FactMysql{}

		factTime, err := time.Parse(time.RFC3339, event.Time)

		if err != nil {
			return nil, err
		}

		fact.PeriodEnd = event.Params.Period.End
		fact.PeriodStart = event.Params.Period.Start
		fact.KeyPeriod = event.Params.Period.TypeKey
		fact.IndicatorToMoId = 315914 //event.Params.IndicatorToMoId
		fact.IndicatorToFactId = 0
		fact.Value = 80
		fact.FactTime = factTime.Format("2006-01-02")
		fact.IsPlan = 0
		fact.Supertags = []models.SuperTag{}
		fact.AuthUserId = 2
		fact.Comment = fmt.Sprintf("indicator_to_mo_id:%d; platform:%s", event.Params.IndicatorToMoId, event.Params.Platform)

		fact.Supertags = append(fact.Supertags,
			models.SuperTag{
				Value: event.Author.UserName,
				Tag: models.Tag{
					Id:           2,
					Name:         "Клиент",
					Key:          "client",
					ValuesSource: 0,
				},
			})

		facts[i] = &fact
	}

	fmt.Println("6. Конвертация успешно завершена")

	return facts, nil
}

func getEvents(cookie string) ([]*models.Event, error) {

	fmt.Println("3. Запрашиваем события development.kpi-drive.ru/_api/events")

	req, err := http.NewRequest("GET",
		"https://development.kpi-drive.ru/_api/events",
		bytes.NewReader([]byte(events_body_filter)))

	if err != nil {
		return nil, err
	}

	req.Header.Add("Cookie", cookie)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	resp_body, _ := io.ReadAll(resp.Body)

	mes := models.Message{}
	err = json.Unmarshal([]byte(resp_body), &mes)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	count_events := len(mes.Data.Rows)
	events := make([]*models.Event, count_events)

	for i := 0; i < count_events; i++ {

		jsonData, _ := json.Marshal(mes.Data.Rows[i])
		event := &models.Event{}
		json.Unmarshal(jsonData, event)
		events[i] = event

	}

	fmt.Printf("4. Получены события. Количество %d \n", count_events)

	return events, nil
}

func getCookie() (string, error) {
	fmt.Println("1. Авторизуемся на development.kpi-drive.ru  ")

	resp, err := http.Get(
		"https://development.kpi-drive.ru/_api/auth/login?login=admin&password=admin")

	if err != nil {
		return "", err
	}

	cookie := resp.Header.Get("Set-Cookie")
	if cookie == "" {
		return "", errors.New("Не удалось получить Cookie")
	}

	i := strings.Index(cookie, ";")

	if i > -1 {
		cookie = string(cookie[0:i])
	} else {
		return "", errors.New("Не удалось получить Cookie")
	}

	fmt.Println("2. Сессия открыта " + cookie)
	return cookie, nil
}
