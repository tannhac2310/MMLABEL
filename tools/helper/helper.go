package helper

import (
	"encoding/json"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

func ParseResponse(resp *transportutil.BaseResponse, structEmpty interface{}) {
	data, err := json.Marshal(resp.Data)
	if err != nil {
		panic("err parseResponse: " + err.Error())
	}

	err = json.Unmarshal(data, structEmpty)
	if err != nil {
		panic(err)
	}
}
