package pg

import (
	"fmt"
	"github.com/tidwall/gjson"
	"encoding/json"
	"errors"
	"strings"
)

func CallPgSelectToJson(queryStr string, res interface{}) (err error) {
	var queryRes []byte

	err = Pg.QueryRow(queryStr).Scan(&queryRes)
	if err != nil {
		fmt.Printf("queryRes err %s\n", err)
		return
	}

	err = json.Unmarshal(queryRes, &res)
	if err != nil {
		return err
	}

	return nil
}

func CallPgFuncWithStruct(funcName string, jsonStruct, res interface{}) error  {
	jsonStr, err := json.Marshal(jsonStruct)
	if err != nil {
		return err
	}
	return CallPgFunc(funcName, jsonStr, res, nil)
}

func CallPgFunc(funcName string, jsonStr []byte, res interface{}, metaInfo interface{}) (err error) {

	var queryRes []byte
	var queryStr string

	if len(jsonStr) > 0 {
		jsonStrMod := strings.Replace(string(jsonStr), "'", "''", -1)
		queryStr = fmt.Sprintf("select * from %s('%s')", funcName, jsonStrMod)
	} else {
		queryStr = fmt.Sprintf("select * from %s()", funcName)
	}

	//fmt.Printf("funcName: %s, queryStr: %s\n", funcName, queryStr)

	err = Pg.QueryRow(queryStr).Scan(&queryRes)
	if err != nil {
		return
	}

	//fmt.Printf("funcName: %s, queryRes: %s\n", funcName, queryRes)

	return ParseResponseFromPostgresFunc(queryRes, res, metaInfo)
}

func ParseResponseFromPostgresFunc(queryRes []byte, tempRes interface{}, metaInfo interface{}) (err error) {
	ok := gjson.Get(fmt.Sprintf("%s", queryRes), "ok").Bool()
	if !ok {
		errMsg := gjson.Get(fmt.Sprintf("%s", queryRes), "message").Str
		err = errors.New(errMsg)
		return
	}

	err = json.Unmarshal([]byte(gjson.Get(fmt.Sprintf("%s", queryRes), "result").Raw), &tempRes)
	if err != nil {
		return err
	}
	if metaInfo != nil {
		err = json.Unmarshal([]byte(gjson.Get(fmt.Sprintf("%s", queryRes), "meta_info").Raw), &metaInfo)
		if err != nil {
			return err
		}
	}
	return nil
}
