package Module

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var dbinfo = fmt.Sprintf("%s:%s@/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci", Config().DB.DBuser, Config().DB.DBpasswd, Config().DB.DBname)

func ConnectDB() *sql.DB {
	db, err := sql.Open("mysql", dbinfo)
	checkErr(err)
	db.SetMaxOpenConns(5000) //用于设置最大打开的连接数，默认值为0表示不限制。
	db.SetMaxIdleConns(10)   //用于设置闲置的连接数。
	err = db.Ping()
	checkErr(err)
	fmt.Println("Connected!")
	return db
}

//增
func InserttoDB(db *sql.DB, cmd string) (bool, int64) {
	res, err := db.Exec(cmd)
	checkErr(err)
	id, _ := res.LastInsertId()
	fmt.Printf("[Database] Insert ID: %d successd \n", id)
	return true, id
}

//删
func DeleteFromDB(db *sql.DB, cmd string) bool {

	stmt, err := db.Prepare(cmd)
	checkErr(err)

	res, err := stmt.Exec()
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

	return true
}

//查
func SelectUserInfo(db *sql.DB, cmd string) UserInfo {
	rows, err := db.Query(cmd)
	checkErr(err)
	var userInfo UserInfo
	for rows.Next() {
		err = rows.Scan(&userInfo.Uid, &userInfo.TelegramId, &userInfo.Age, &userInfo.Role, &userInfo.Height,
			&userInfo.Bodytype, &userInfo.Size)
		checkErr(err)
	}
	return userInfo
}

func SelectAdminInfo(db *sql.DB, cmd string) AdminInfo {
	rows, err := db.Query(cmd)
	checkErr(err)
	var adminInfo AdminInfo
	for rows.Next() {
		err = rows.Scan(&adminInfo.Admin_id, &adminInfo.Admin_uid)
		checkErr(err)
	}
	return adminInfo
}

func SelectWelcome(db *sql.DB, cmd string) WelcomeMessage {
	rows, err := db.Query(cmd)
	checkErr(err)
	var welcomeMessage WelcomeMessage
	for rows.Next() {
		err = rows.Scan(&welcomeMessage.Group_username, &welcomeMessage.Group_welcome, &welcomeMessage.Ask_role)
		checkErr(err)
	}
	return welcomeMessage
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
