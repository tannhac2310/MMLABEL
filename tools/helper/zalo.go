package helper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

func CallZaloMethod(ctx context.Context, client *http.Client, url string, reqData interface{}, respData *transportutil.BaseResponse, token string) error {
	data, _ := json.Marshal(reqData)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.WithContext(ctx)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("s.client.Do: %w", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, respData)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %w: %v", err, string(body))
	}

	if respData.Code != 0 {
		return fmt.Errorf("err: %v, debugMesage: %v", respData.Message, respData.DebugMessage)
	}

	return nil
}
