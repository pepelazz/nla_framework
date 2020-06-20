package bitrix

import (
	"[[LocalProjectPath]]/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"encoding/json"
	"[[LocalProjectPath]]/pg"
	"github.com/spf13/cast"
)

type (
	[[DocNameCamel]]FromBtx struct {
	[[- range .Flds]]
	[[- if IsBtxFld .]]
		[[ToCamel .Name]] 		[[GetBtxFldType .]]  		`json:"[[GetBtxFldName .]]"`
	[[- end]]
	[[- end]]
	}

	[[DocNameCamel]] struct {
		Id                  int           `json:"id"`
	[[- range .Flds]]
	[[- if IsBtxFld .]]
		[[ToCamel .Name]] 		[[.GoType]]  		`json:"[[.Name]]"`
	[[- end]]
	[[- end]]
	}
)

func Get[[DocNameCamel]]History(c *gin.Context) {
	utils.HttpSuccess(c, "ok")

	var err error

	userId, _ := utils.ExtractUserIdString(c)

	go func() {
		nextId := 0
		lastProcessedId := 0
		for {
			lastId := 0
			fmt.Printf("getAll[[DocNameCamel]]HistoryAndSave nextId: %v\n", nextId)
			nextId, lastId, err = getAll[[DocNameCamel]]HistoryAndSave(nextId, nil)
			if err != nil {
				fmt.Printf("getAll[[DocNameCamel]]HistoryAndSave err %s\n", err)
				return
			}
			[[if .Integrations.Bitrix.IsNoPagination]]
			fmt.Printf("getAll[[DocNameCamel]]HistoryAndSave finished")
			saveResultMsgToPg(userId, "[[DocNameCamel]] импортированы из Битрикс")
			fmt.Printf("lastProcessedId: %v lastId: %v %v\n", lastProcessedId, lastId, time.Second)
			return
			[[else]]
			// прерываем процесс когда id'шники пошли на второй круг. Определяем это по тому что новый lastId меньше последнего обработанного id'шника
			if lastProcessedId > 0 && lastId < lastProcessedId {
				fmt.Printf("getAll[[DocNameCamel]]HistoryAndSave finished")
				saveResultMsgToPg(userId, "[[DocNameCamel]] импортированы из Битрикс")
				return
			}
			lastProcessedId = lastId
			time.Sleep(100 * time.Millisecond)
			[[end]]
		}
	}()
}

func getAll[[DocNameCamel]]HistoryAndSave(startId int, errResultArr *[]errResult) (nextId, lastId int, err error) {

	res := struct {
		[[if .Integrations.Bitrix.Result.StructDesc]] [[.Integrations.Bitrix.Result.StructDesc]] [[- else -]]Result           [][[DocNameCamel]]FromBtx `json:"result"`[[end]]
		Error            interface{}  `json:"error"`
		ErrorDescription string       `json:"error_description"`
	}{}
	// https://crm.tian-trade.ru/rest/11161/cbwiqxom770hdgpm/crm.company.list.json?select[]=title&select[]=lead_id
	// https://crm.tian-trade.ru/rest/11161/cbwiqxom770hdgpm/crm.company.list.json?Filter[UF_CRM_1535355557]=12431
	selectFlds := []string{[[range .Flds]] [[- if IsBtxFld .]] "[[GetBtxFldName .]]", [[- end]] [[- end]]}
	url := fmt.Sprintf("%s/rest/%s/%s/[[.Integrations.Bitrix.UrlName]]?&start=%v&order[id]=asc[[.Integrations.Bitrix.UrlQuery]]", bitrixConfig.ApiUrl, bitrixConfig.UserId, bitrixConfig.WebhookToken, startId)
	for _, fld := range selectFlds {
		url = fmt.Sprintf("%s&select[]=%s", url, fld)
	}
	err = utils.GetJsonByUrl(url, &res)
	if err != nil {
		return
	}

	if len(res.ErrorDescription) > 0 {
		return 0, 0, errors.New(fmt.Sprintf("error: %s", res.ErrorDescription))
	}

	//fmt.Printf("result length: %v\n", len(res.Result))
	if [[.Integrations.Bitrix.PrintCheckResultIsEmpty]] {
		return 0, 0, nil
	} else {
		nextId = startId + 50
	}

	for _, v := range res.[[if .Integrations.Bitrix.Result.PathStr -]] [[.Integrations.Bitrix.Result.PathStr]] [[- else -]]Result[[end]] {
		//fmt.Printf("process %s %s %s %s\n", v.ID, v.TITLE, v.PHONE, v.EMAIL)
		doc, err := v.ConvertFromBitrix()
		if err != nil {
			fmt.Printf("ConvertFromBitrix err %s %s\n", err, v)
			if errResultArr != nil {
				*errResultArr = append(*errResultArr, errResult{
					JsonParams: doc,
					Message:    fmt.Sprintf("ConvertFromBitrix error: %s\n", err),
				})
			}
			continue
		}
		lastId = cast.ToInt(v.BtxId)
		if lastId == 0 {
			fmt.Printf("cast.ToInt err %s %s\n", err, v)
			continue
		}

		jsonData, _ := json.Marshal(doc)
		err = pg.CallPgFunc("[[.Name]]_update", jsonData, doc, nil)
		if err != nil {
			fmt.Printf("[[.Name]]_update error: %s %s\n", err, jsonData)
			if errResultArr != nil {
				*errResultArr = append(*errResultArr, errResult{
					JsonParams: doc,
					Message:    fmt.Sprintf("[[.Name]]_update error: %s\n", err),
				})
			}
			continue
		}
	}
	return
}

func (btxDoc *[[DocNameCamel]]FromBtx) ConvertFromBitrix() (res *[[DocNameCamel]], err error) {
	if btxDoc == nil {
		return nil, errors.New("[[DocNameCamel]]FromBtx is nil in [[DocNameCamel]]FromBtx.ConvertFromBitrix")
	}
	res = &[[DocNameCamel]]{}
	[[range .Flds]]
	[[- if IsBtxFld .]]
		[[CastToGoType .]]
	[[- end]]
	[[- end]]

	return res, nil
}

func Get[[DocNameCamel]]HistoryDebug(c *gin.Context) {

	var err error

	nextId := 0
	lastProcessedId := 0
	errResultArr := []errResult{}
	for {
		lastId := 0
		fmt.Printf("getAll[[DocNameCamel]]HistoryAndSave nextId: %v\n", nextId)
		nextId, lastId, err = getAll[[DocNameCamel]]HistoryAndSave(nextId, &errResultArr)
		if err != nil {
			fmt.Printf("getAll[[DocNameCamel]]HistoryAndSave err %s\n", err)
			break
		}

		if nextId > 100 {
			break
		}

		[[if .Integrations.Bitrix.IsNoPagination]]
		fmt.Printf("getAll[[DocNameCamel]]HistoryAndSave finished")
		fmt.Printf("lastProcessedId: %v lastId: %v %v\n", lastProcessedId, lastId, time.Second)
		break
		[[else]]
		// прерываем процесс когда id'шники пошли на второй круг. Определяем это по тому что новый lastId меньше последнего обработанного id'шника
		if lastProcessedId > 0 && lastId < lastProcessedId {
			fmt.Printf("getAll[[DocNameCamel]]HistoryAndSave finished")
			break
		}
		lastProcessedId = lastId
		time.Sleep(100 * time.Millisecond)
		[[end]]
	}
	if err != nil {
		utils.HttpError(c, http.StatusBadRequest, err.Error())
	} else if len(errResultArr) > 0  {
		c.JSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"message": errResultArr,
		})
	} else {
		utils.HttpSuccess(c, fmt.Sprintf("succesfully imported: %v", nextId))
	}

}