package dto

import "time"

// Response is the outer json object used in response.
// The response data should be added in the Data field
type Response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
	Time   int64       `json:"time"`
}

// NewResponseFine will create a new Response with a "ok" as Status
func NewResponseFine(data interface{}) Response {
	return Response{
		Status: "ok",
		Time:   time.Now().UnixNano() / int64(time.Millisecond) / int64(time.Nanosecond),
		Data:   data,
	}
}

// NewResponseBad will create a new Response with a "bad" as Status,
// mostly there should be a string work as a message in data field
func NewResponseBad(data interface{}) Response {
	return Response{
		Status: "bad",
		Time:   time.Now().UnixNano() / int64(time.Millisecond) / int64(time.Nanosecond),
		Data:   data,
	}
}
