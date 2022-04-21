package model

import "encoding/json"

func UnmarshalResponse(data []byte) (Response, error) {
	var r Response
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Response) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Response struct {
	Status int64          `json:"Status"`
	Text   *string        `json:"text"`
	Flex   *[]interface{} `json:"Flex"`
}
