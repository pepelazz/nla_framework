package odata

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"[[LocalProjectPath]]/pg"
	"[[LocalProjectPath]]/utils"
	"time"
	[[range .Integrations.Odata.Import]]
	"[[.]]"
	[[- end]]
)

type (
	[[DocNameCamel]]Type struct {
	[[- range .Flds]]
	[[- if IsOdataFld .]]
		[[ToCamel .Name]] 			[[GetOdataFldType .]]  		`json:"[[GetOdataFldName .]]" xml:"[[GetOdataFldName .]]"`
	[[- end]]
	[[- end]]
	[[.Integrations.Odata.Hooks.TypeAddFlds]]
	}

	[[DocNameCamel]]ForPgType struct {
	Id                  int           `json:"id"`
	[[- range .Flds]]
	[[- if IsOdataFld .]]
		[[ToCamel .Name]] 		[[.GoType]]  		`json:"[[.Name]]"`
	[[- end]]
	[[- end]]
	[[.Integrations.Odata.Hooks.PgTypeAddFlds]]
	}
)

func Start[[DocNameCamel]]Sync(c *gin.Context) {
	go func() {
		resMsg := sync[[DocNameCamel]]With1C()
		// message с результатами
		userId, _ := utils.ExtractUserIdString(c)
		err := saveResultMsgToPg(userId, "Синхронизация с 1С: [[.NameRu]]", resMsg)
		if err != nil {
			fmt.Printf("Start[[DocNameCamel]]Sync saveResultMsgToPg error: %s\n", err)
		}
	}()
	utils.HttpSuccess(c, "ok")
}

func sync[[DocNameCamel]]With1C() ([]resultMsgType) {
	start := time.Now()
	resMsg := newResultMsgType("Синхронизация: [[.NameRu]]")
	resList, err := get[[DocNameCamel]]()
	if err != nil {
		fmt.Printf("sync[[DocNameCamel]]With1C get[[DocNameCamel]] error: %s\n", err)
		resMsg.addErr(err.Error())
		return []resultMsgType{resMsg}
	}
	cnt := 0 //счетчик записей
	for _, v := range resList {
		jsonStr, _ := json.Marshal(v)
		err := pg.CallPgFunc("[[.Name]]_update", jsonStr, nil, nil)
		if err != nil {
			fmt.Printf("sync[[DocNameCamel]]With1C [[.Name]]_update error: %s jsonStr: %s\n", err, jsonStr)
			resMsg.Errors = append(resMsg.Errors, fmt.Sprintf("%s uuid: %s", err, v.Uuid))
			continue
		}
		cnt++
	}
	resMsg.addResult(fmt.Sprintf("синхронизировано записей: <strong>%v</strong>", cnt))
	elapsed := time.Since(start)
	resMsg.setDuration(fmt.Sprintf("%s", elapsed))
	return []resultMsgType{resMsg}
}

func get[[DocNameCamel]]() ([][[DocNameCamel]]ForPgType, error) {
	res := [][[DocNameCamel]]ForPgType{}
	start := time.Now()
	odataQuery := odataQueryType{
		DocType: "[[GetOdataName]]",
		Format:  "json",
		Select:  []string{[[range GetOdataFldNames]]"[[.]]",[[end]][[range .Integrations.Odata.Hooks.UrlAddFlds]]"[[.]]",[[end]]},
		Expand:  []string{},
		Filter:  []string{[[range .Integrations.Odata.Filter]][[.]],[[end]]},
		//Limit:   50,
	}
	targetUrl := odataQuery.buildQuery()
	//fmt.Printf("targetUrl %s\n", targetUrl)
	tempRes := struct {
		Value [][[DocNameCamel]]Type `json:"value"`
	}{}
	err := odataCallByUrl(targetUrl, "GET", "json", &tempRes, nil)
	if err != nil {
		return res, err
	}
	elapsed := time.Since(start)
	fmt.Printf("get[[DocNameCamel]] len: %v took time: %s\n", len(tempRes.Value), elapsed)
	for _, v := range tempRes.Value {
		c := [[DocNameCamel]]ForPgType{}
		[[- range .Flds]]
		[[- if IsOdataFld .]]
		c.[[ToCamel .Name]] = v.[[ToCamel .Name]]
		[[- end]]
		[[- end]]
		[[.Integrations.Odata.Hooks.ConvertAddFlds]]
		res = append(res, c)
	}
	return res, nil
}

func Sync[[DocNameCamel]]With1CDebug(c *gin.Context) {

	resList := [][[DocNameCamel]]ForPgType{}
	start := time.Now()
	odataQuery := odataQueryType{
		DocType: "[[GetOdataName]]",
		Format:  "json",
		Select:  []string{[[range GetOdataFldNames]]"[[.]]",[[end]][[range .Integrations.Odata.Hooks.UrlAddFlds]]"[[.]]",[[end]]},
		Expand:  []string{},
		Filter:  []string{[[range .Integrations.Odata.Filter]][[.]],[[end]]},
		Limit:   [[if .Integrations.Odata.Filter]] 0 [[else]] 100 [[end]],
	}

	targetUrl := odataQuery.buildQuery()
	//fmt.Printf("targetUrl %s\n", targetUrl)
	tempRes := struct {
		Value [][[DocNameCamel]]Type `json:"value"`
	}{}
	err := odataCallByUrl(targetUrl, "GET", "json", &tempRes, nil)
	if err != nil {
		utils.HttpError(c, 400, err.Error())
		return
	}
	elapsed := time.Since(start)
	fmt.Printf("get[[DocNameCamel]] len: %v took time: %s\n", len(tempRes.Value), elapsed)
	for _, v := range tempRes.Value {
		c := [[DocNameCamel]]ForPgType{}
		[[- range .Flds]]
		[[- if IsOdataFld .]]
		c.[[ToCamel .Name]] = v.[[ToCamel .Name]]
		[[- end]]
		[[- end]]
		[[.Integrations.Odata.Hooks.ConvertAddFlds]]
		resList = append(resList, c)
	}

	cnt := 0 //счетчик записей
	for _, v := range resList {
		jsonStr, _ := json.Marshal(v)
		err := pg.CallPgFunc("[[.Name]]_update", jsonStr, nil, nil)
		if err != nil {
			fmt.Printf("sync[[DocNameCamel]]With1C [[.Name]]_update error: %s jsonStr: %s\n", err, jsonStr)
			continue
		}
		cnt++
	}
	utils.HttpSuccess(c, fmt.Sprintf("синхронизировано записей: %v. Время %v", cnt, time.Since(start)))
	return
}
