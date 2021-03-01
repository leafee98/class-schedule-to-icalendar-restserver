package dto

// ConfigCreateReq is used in the request to create a new global config
type ConfigCreateReq struct {
	Name    string `json:"name"`
	Content string `json:"content"`
	Format  int8   `json:"format"`
	Remark  string `json:"remark"`
}

// ConfigCreateRes is used as response of ConfigCreateReq
type ConfigCreateRes struct {
	ID int64 `json:"id"`
}
