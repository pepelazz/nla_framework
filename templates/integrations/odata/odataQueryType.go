package odata

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"[[.Config.LocalProjectPath]]/pg"
)

type (
	odataQueryType struct {
		Id      string
		DocType string
		Format  string
		Select  []string
		Expand  []string
		Filter  []string
		Limit int
	}
	resultMsgType struct {
		Title string `json:"title"`
		Result []string `json:"result"`
		Errors []string `json:"errors"`
		Duration string `json:"duration"`
	}
)

func (q *odataQueryType) buildQuery() string {
	idStr := fmt.Sprintf("(%s)", q.Id)
	baseUrl := fmt.Sprintf("%s/%s%s?", odataConfig.Url, q.DocType, idStr)
	if len(q.Format) > 0 {
		baseUrl += "$format=" + q.Format
	} else {
		baseUrl += "$format=atom"
	}
	if len(q.Select) > 0 {
		baseUrl += "&$select=" + strings.Join(q.Select, ",")
	}
	if len(q.Expand) > 0 {
		baseUrl += "&$expand=" + strings.Join(q.Expand, ",")
	}
	if len(q.Filter) > 0 {
		baseUrl += "&$filter=" + strings.Join(q.Filter, ",")
	}
	if q.Limit>0 {
		baseUrl += fmt.Sprintf("&$top=%v", q.Limit)
	}
	return baseUrl
}

// получение данных из odata
func odataCallByUrl(url, method, formatType string, res interface{}, reqBody []byte) error {

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	// This one line implements the authentication required for the task.
	req.SetBasicAuth(odataConfig.Login, odataConfig.Password)

	// Make request and show output.
	httpRes, err := client.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return err
	}

	//fmt.Printf("body %s\n", body)

	if formatType == "json" {
		err = json.Unmarshal(body, &res)
	} else {
		err = xml.Unmarshal(body, &res)
	}
	if err != nil {
		return err
	}
	return nil
}

// внесение изменений в odata
func postAuthOdataByUrl(url string, jsonData io.Reader, res interface{}) error {

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, jsonData)
	if err != nil {
		return err
	}

	// This one line implements the authentication required for the task.
	req.SetBasicAuth(odataConfig.Login, odataConfig.Password)

	// Make request and show output.
	httpRes, err := client.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return err
	}

	//fmt.Printf("body %s\n", body)
	return json.Unmarshal(body, &res)
}

func newResultMsgType(title string) resultMsgType  {
	return resultMsgType{
		Title: 	title,
		Result:   []string{},
		Errors:   []string{},
		Duration: "",
	}
}

func (r *resultMsgType) addErr(msg string)  {
	r.Errors = append(r.Errors, msg)
}

func (r *resultMsgType) addResult(msg string)  {
	r.Result = append(r.Result, msg)
}

func (r *resultMsgType) setDuration(msg string)  {
	r.Duration = msg
}

func saveResultMsgToPg(userId string, title string, res []resultMsgType) error {
	msg := ""
	for _, v := range res {
		msg = fmt.Sprintf("%s<strong>%s</strong><br>", msg, v.Title)
		for _, r := range v.Result {
			msg = fmt.Sprintf("%s - %s<br>", msg, r)
		}
		for _, r := range v.Errors {
			msg = fmt.Sprintf(`%s - <strong>ошибка:</strong> %s<br>`, msg, r)
		}
		msg = fmt.Sprintf(`%s - <small>время синхронизации: %s</small><br><br>`, msg, v.Duration)
	}
	jsonStr, _ := json.Marshal(map[string]interface{}{"id": -1, "user_id": userId, "title": title, "data": map[string]string{"message": msg}})
	return pg.CallPgFunc("message_update", jsonStr, nil, nil)
}