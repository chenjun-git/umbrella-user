package model

import (
	"fmt"
	"database/sql"
	"strings"
	"time"

	"github.com/chenjun-git/umbrella-user/common"
	"github.com/chenjun-git/umbrella-user/db"
	//"business/user/utils"
)

const (
	UserTableName = "user"

	UserFieldID                 = "user_id"
	UserFieldUserName           = "user_name"
	UserFieldAlias              = "alias"
	UserFieldUserPasswd         = "user_passwd"
	UserFieldHashPasswd         = "hashed_passwd"
	UserFieldPasswdLevel        = "passwd_level"
	UserFieldLastUpdatePasswd   = "last_update_passwd"
	UserFieldRegisterSource     = "register_source"
	UserFieldPortTrait          = "portrait"
	UserFieldTel                = "tel"
	UserFieldEmail              = "email"
	UserFieldQQ                 = "qq"
	UserFieldIsMember           = "is_member"
	UserFieldMemberStartTime    = "member_start_time"
	UserFieldMemberDuration     = "member_duration"
	UserFieldRole               = "role"
	UserFieldCreateTime         = "create_time"
	UserFieldUpdateTime         = "update_time"
)

type User struct {
	ID               string
	UserName         string
	Alias            string
	UserPasswd       string
	HashPasswd       string
	PasswdLevel      int
	LastUpdatePasswd *time.Time
	RegisterSource   string
	PortRait         string
	Tel              sql.NullString
	Email            sql.NullString
	QQ               string
	IsMember         bool
	MemberStartTime  *time.Time
	MemberDuration   int
	role             bool
}

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

//////////////////////////////////////////////////////////////////////////////////

func Count(mysqlExex db.MySQLExec, fieldValues map[string]interface{}) (int, error) {
	_SQL := "select count(1) from user where %s"

	args := make([]interface{}, 0, len(fieldValues))
	conditions := make([]string, 0, len(fieldValues))
	for key, value := range fieldValues {
		conditions = append(conditions, fmt.Sprintf("%s = ?", key))
		args = append(args, value)
	}
	_SQL = fmt.Sprintf(_SQL, strings.Join(conditions, ","))

	var total int
	if err := mysqlExex.QueryRow(_SQL, args...).Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}

func getUserBy(mysqlExex db.MySQLExec, filter string, args []interface{}) (*User, error) {
	_SQL := "select user_id,user_name, alias, user_passwd, hashed_passwd, passwd_level, last_update_passwd, register_source, portrait, tel, email, qq, is_member, member_start_time, member_duration, role from user %s limit 1 offset 0"

	if filter != "" {
		_SQL = fmt.Sprintf(_SQL, " where "+filter)
	} else {
		_SQL = fmt.Sprintf(_SQL, "")
	}

	user := &User{}
	err := mysqlExex.QueryRow(_SQL, args...).Scan(
		&user.ID,
		&user.UserName,
		&user.Alias,
		&user.UserPasswd,
		&user.HashPasswd,
		&user.PasswdLevel,
		&user.LastUpdatePasswd,
		&user.RegisterSource,
		&user.PortRait,
		&user.Tel,
		&user.Email,
		&user.QQ,
		&user.IsMember,
		&user.MemberStartTime,
		&user.MemberDuration,
		&user.MemberDuration,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserByPhone(db db.MySQLExec, phone string) (*User, error) {
	return getUserBy(db, "tel = ?", []interface{}{phone})
}

func GetUserByEmail(db db.MySQLExec, email string) (*User, error) {
	return getUserBy(db, "email = ?", []interface{}{email})
}

func GetUserByPhoneEmail(db db.MySQLExec, mediaType, phone, email string) (*User, error) {
	if common.MediaTypeEmail == mediaType {
		return GetUserByEmail(db, email)
	} else {
		return GetUserByPhone(db, phone)
	}
}

func GetUserById(db db.MySQLExec, id string) (*User, error) {
	return getUserBy(db, "user_id = ?", []interface{}{id})
}

func CreateUserByPhone(db db.MySQLExec, id string, user User) (string, error) {
	_SQL := "insert into user(user_id,user_name, alias, user_passwd, hashed_passwd, passwd_level, last_update_passwd, register_source, portrait, tel, qq, role) value(?,?,?,?,?,?,?,?,?,?,?,?)"
	_, err := db.Exec(_SQL, id, user.UserName, user.Alias, user.UserPasswd, user.HashPasswd, user.PasswdLevel, time.Now().UTC(), user.RegisterSource, user.PortRait, user.Tel, user.QQ, 0)
	if err != nil {
		return "", err
	}

	return id, nil
}

func CreateUserByEmail(db db.MySQLExec, id string, user User) (string, error) {
	_SQL := "insert into user(user_id,user_name, alias, user_passwd, hashed_passwd, passwd_level, last_update_passwd, register_source, portrait, email, qq,role) value(?,?,?,?,?,?,?,?,?,?,?,?)"
	_, err := db.Exec(_SQL, id, user.UserName, user.Alias, user.UserPasswd, user.HashPasswd, user.PasswdLevel, time.Now().UTC(), user.RegisterSource, user.PortRait, user.Email, user.QQ, 0)
	if err != nil {
		return "", err
	}

	return id, nil
}

// 根据id更新用户信息
func updateUserById(db db.MySQLExec, id string, fieldValues map[string]interface{}) (int64, error) {
	_SQLTemp := "update user set %s where user_id = ?"

	sqlArgs := make([]interface{}, 0)
	sqlFields := make([]string, 0)
	for key, value := range fieldValues {
		sqlFields = append(sqlFields, fmt.Sprintf("`%s` = ?", key))
		sqlArgs = append(sqlArgs, value)
	}
	sqlArgs = append(sqlArgs, id)
	_SQL := fmt.Sprintf(_SQLTemp, strings.Join(sqlFields, ","))
	result, err := db.Exec(_SQL, sqlArgs...)
	if err != nil {
		return 0, err
	}

	rowCnt, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowCnt, nil
}

func UpdatePhone(db db.MySQLExec, userID string, phone string) (int64, error) {
	values := map[string]interface{}{
		"tel": phone,
	}
	return updateUserById(db, userID, values)
}

func UpdateEmail(db db.MySQLExec, userID string, email string) (int64, error) {
	values := map[string]interface{}{
		"email": email,
	}
	return updateUserById(db, userID, values)
}

func UpdatePassword(db db.MySQLExec, userID string, hashedPassword string, level int) (int64, error) {
	values := map[string]interface{}{
		UserFieldHashPasswd:       hashedPassword,
		UserFieldPasswdLevel:      level,
		UserFieldLastUpdatePasswd: time.Now().UTC(),
	}
	return updateUserById(db, userID, values)
}

// func UpdatePasswordLifetime(sqlExec db.MySQLExec, accountID string, enabled bool, maxAge int) (int64, error) {
// 	values := map[string]interface{}{}
// 	if enabled {
// 		values["enable_password_lifetime"] = 1
// 		values["max_password_lifetime"] = maxAge
// 	} else {
// 		values["enable_password_lifetime"] = 0
// 	}

// 	return updateAccountById(sqlExec, accountID, values)
// }

// 根据id删除用户
func DeleteUserById(mysqlExec db.MySQLExec, id string) error {
	_SQL := "delete from user where id = ?"

	_, err := mysqlExec.Exec(_SQL, id)
	if err != nil {
		return err
	}

	return nil
}
