package model

const (
	UserTableName = "user"

	UserFieldID              = "id"
	UserFieldUserName        = "user_name"
	UserFieldUserPasswd      = "user_passwd"
	UserFieldAlias           = "alias"
	UserFieldPortTrait       = "portrait"
	UserFieldTel             = "tel"
	UserFieldEmail           = "email"
	UserFieldQQ              = "qq"
	UserFieldIsMember        = "is_member"
	UserFieldMemberStartTime = "member_start_time"
	UserFieldMemberDuration  = "member_duration"
	UserFieldRole            = "role"
	UserFieldTegTime         = "reg_time"
)

func dbMySQLExec(db db.MySQLExec, _SQL string, query ...interface{}) (int64, error) {
	result, err := db.Exec(_SQL, query...)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func sqlInsert(db db.MySQLExec, tableName string, datas map[string]interface{}) (int64, error) {
	var q []interface{}
	var u []string
	for key, value := range datas {
		u = append(u, fmt.Sprintf("%s = ?", key))
		q = append(q, value)
	}

	_SQL := fmt.Sprintf("INSERT INTO %s SET %s", tableName, strings.Join(u, ", "))
	return dbMySQLExec(db, _SQL, q...)
}

func sqlQueryRow(db db.MySQLExec, tableName string, selects, wheres map[string]interface{}) (bool, error) {
	var f []string
	var s []interface{}
	for key, value := range selects {
		f = append(f, key)
		s = append(s, value)
	}

	var w []string
	var q []interface{}
	for key, value := range wheres {
		w = append(w, fmt.Sprintf("%s = ?", key))
		q = append(q, value)
	}

	_SQL := fmt.Sprintf(`SELECT %s FROM %s WHERE %s`, strings.Join(f, ", "), tableName, strings.Join(w, " and "))
	err := db.QueryRow(_SQL, q...).Scan(s...)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func sqlUpdate(db db.MySQLExec, tableName string, updates, wheres map[string]interface{}) (int64, error) {
	var q []interface{}

	var u []string
	for key, value := range updates {
		u = append(u, fmt.Sprintf("%s = ?", key))
		q = append(q, value)
	}

	var w []string
	for key, value := range wheres {
		w = append(w, fmt.Sprintf("%s = ?", key))
		q = append(q, value)
	}

	_SQL := fmt.Sprintf("UPDATE %s SET %s WHERE %s", tableName, strings.Join(u, ", "), strings.Join(w, " and "))
	return dbMySQLExec(db, _SQL, q...)
}
