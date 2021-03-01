package db

// ConfigCreate will create a new config in database,
// format: 0=json, other value is unimplemented
func ConfigCreate(name string, content string, format int8, ownerID int64, remark string) (int64, error) {
	var configType int8 = 0
	res, err := DB.Exec(
		"insert into t_config (c_type, c_name, c_content, c_format, c_owner_id, c_remark)"+
			" values (?, ?, ?, ?, ?, ?)",
		configType, name, content, format, ownerID, remark)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
