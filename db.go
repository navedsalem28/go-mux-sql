package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	//_ "odbc/driver"
	//_ "github.com/alexbrainman/odbc"
	_ "github.com/go-sql-driver/mysql"
)

var dbConnection map[string]*sql.DB

// DisconnectDB close down db connection and clear memory
func DisconnectDB() {

	for key, _ := range dbConnection {
		dbConnection[key].Close()
	}
	//dbConnection["default"].Close()
}

var (
	DBUser     string
	DBPassword string
	DBHost     string
)

// ConnectDB make a connection to database using provide DSN as envoirent variable ConnectDB()
func ConnectDB() bool {
	var err error
	dbConnection = make(map[string]*sql.DB)
	var value *sql.DB

	dbType, ok := GetConfigData("DbType")
	if !ok {
		fmt.Println("------")
		dbType = os.Getenv("DbType")
	}

	if dbType == "" {
		Warning("can't get DbType value from environment variable")
		return false
	}
	switch dbType {
	case "ODBC":
		DSN := os.Getenv("DbDNS")
		if DSN == "" {
			Warning("can't get DbDNS value from environment variable")
			return false
		}
		value, err = sql.Open("odbc", "DSN="+DSN)
		dbConnection["default"] = value
		break
	case "MYSQL":
		DBNAME, ok := GetConfigData("DbName")
		if !ok {
			DBNAME = os.Getenv("DbName")
		}

		if DBNAME == "" {
			Warning("can't get DbName value from environment variable")
			return false
		}

		DBUser, ok = GetConfigData("DbUser")
		if !ok {
			DBUser = os.Getenv("DbUser")
		}

		if DBUser == "" {
			Warning("can't get DbUser value from environment variable")
			return false
		}

		DBPassword, ok = GetConfigData("DbPassword")
		if !ok {
			DBPassword = os.Getenv("DbPassword")
		}

		if DBPassword == "" {
			Warning("can't get DbPassword value from environment variable")
			return false
		}
		strConnection := ""
		DBHost, ok = GetConfigData("DbHost")
		if !ok {
			DBHost = os.Getenv("DbHost")
		}

		if DBHost == "" {
			strConnection = fmt.Sprintf("%s:%s@/%s", DBUser, DBPassword, DBNAME)
		} else {
			strConnection = fmt.Sprintf("%s:%s@%s/%s", DBUser, DBPassword, DBHost, DBNAME)
		}
		value, err = sql.Open("mysql", strConnection)
		dbConnection["default"] = value

		break
	default:
		Warning("Invalid DbType, currently ODBC / MYSQL supported")
		break
	}
	if err != nil {
		log.Println("DB Connection Error " + err.Error())
		return false
	}

	dbConnection["default"].SetMaxOpenConns(100)
	dbConnection["default"].SetMaxIdleConns(90)
	dbConnection["default"].SetConnMaxLifetime(time.Hour)
	keepAliveDB()
	return true
}

func keepAliveDB() {
	timer := time.NewTimer(time.Second * 60)
	go func() {
		<-timer.C
		log.Println("Refresh DB Connection")
		params := make([]interface{}, 0)
		Row, ok := GetSingleRow("SELECT NOW() AS 'Time'", params, "default")
		if ok {
			log.Println("DB Time " + Row["Time"])
		}
		keepAliveDB()
	}()
}

// GetSingleRow (Query as string, QueyParameters as []interface{}) return (Row as map[string]string, staus as bool)
func GetSingleRow(Query string, params []interface{}, DBConnectionName string) (map[string]string, bool) {
	log.Println("GetSingleRow : Executing Query :: " + Query)
	var rows *sql.Rows
	var err error
	if len(params) > 0 {
		paramPointers := make([]interface{}, len(params))
		for i := range params {
			paramPointers[i] = &params[i]
		}
		rows, err = dbConnection[DBConnectionName].Query(Query, paramPointers...)
	} else {
		rows, err = dbConnection[DBConnectionName].Query(Query)
	}
	if err != nil {
		log.Println("GetSingleRow : Error:" + err.Error() + "\tExecuting Query: " + Query)
		return nil, false
	}
	defer rows.Close()
	columnNames, err := rows.Columns()
	if err != nil {
		log.Println("GetSingleRow : Error:" + err.Error() + "\tFetching Column Names Query: " + Query)
		return nil, false
	}
	if len(columnNames) < 1 {
		log.Println("GetSingleRow : Error:" + err.Error() + "\tNo Column in Resultset Query: " + Query)
		return nil, false
	}
	columns := make([]interface{}, len(columnNames))
	for i := range columnNames {
		columns[i] = new(sql.RawBytes)
	}
	log.Println("GetSingleRow :\tscaning rows for Query: " + Query)
	ThisRow := make(map[string]string)
	ReturnStatus := false
	for rows.Next() {
		err := rows.Scan(columns...)
		if err != nil {
			log.Println("GetSingleRow : Error:" + err.Error() + "\tScannig rows for Query: " + Query)
			return nil, false
		}
		for i, colName := range columnNames {
			if rb, ok := columns[i].(*sql.RawBytes); ok {
				colValue := string(*rb)
				ThisRow[colName] = colValue
				*rb = nil // reset pointer to discard current value to avoid a bug
				log.Println("Column : " + colName + "\tValue : " + colValue)
			} else {
				log.Println("GetSingleRow : Column " + colName + " contains nil value for Query: " + Query)
				ThisRow[colName] = ""
			}
		}
		ReturnStatus = true
		return ThisRow, true
	}
	return ThisRow, ReturnStatus
}

// GetAllRows (Query as string, QueyParameters as []interface{}) return (Row as []map[string]string, staus as bool)
func GetAllRows(Query string, params []interface{}, DBConnectionName string) ([]map[string]string, bool) {
	log.Println("GetAllRows : Executing Query :: " + Query)
	var rows *sql.Rows
	var err error
	if len(params) > 0 {
		paramPointers := make([]interface{}, len(params))
		for i := range params {
			paramPointers[i] = &params[i]
		}
		rows, err = dbConnection[DBConnectionName].Query(Query, paramPointers...)
	} else {
		rows, err = dbConnection[DBConnectionName].Query(Query)
	}
	if err != nil {
		log.Println("GetAllRows : Error:" + err.Error() + "\tExecuting Query: " + Query)
		return nil, false
	}
	defer rows.Close()
	columnNames, err := rows.Columns()
	if err != nil {
		log.Println("GetAllRows : Error:" + err.Error() + "\tFetching Column Names Query: " + Query)
		return nil, false
	}
	if len(columnNames) < 1 {
		log.Println("GetAllRows : Error:" + err.Error() + "\tNo Column in Resultset Query: " + Query)
		return nil, false
	}
	columns := make([]interface{}, len(columnNames))
	for i := range columnNames {
		columns[i] = new(sql.RawBytes)
	}
	log.Println("GetAllRows :\tscaning rows for Query: " + Query)
	var RecordsSet []map[string]string

	RowCount := 0
	for rows.Next() {
		RowCount++
		ThisRow := make(map[string]string)
		err := rows.Scan(columns...)
		if err != nil {
			log.Println("GetSingleRow : Error:" + err.Error() + "\tScannig rows for Query: " + Query)
			return nil, false
		}
		for i, colName := range columnNames {
			if rb, ok := columns[i].(*sql.RawBytes); ok {
				colValue := string(*rb)
				ThisRow[colName] = colValue
				*rb = nil // reset pointer to discard current value to avoid a bug
				log.Println("Column : " + colName + "\tValue : " + colValue)
			} else {
				log.Println("GetSingleRow : Column " + colName + " contains nil value for Query: " + Query)
				ThisRow[colName] = ""
			}
		}
		RecordsSet = append(RecordsSet, ThisRow)
	}
	if RowCount > 0 {
		return RecordsSet, true
	}
	return nil, false
}

// UpdateDB (Query string, params []interface{}) and return (LastInsertID for INSERT OR Updated rows as int, status as bool)
func UpdateDB(Query string, params []interface{}, DBConnectionName string) (int64, bool) {
	log.Println("UpdateDB : Executing Query :: " + Query)
	var Res sql.Result
	var err error
	if len(params) > 0 {
		paramPointers := make([]interface{}, len(params))
		for i := range params {
			paramPointers[i] = &params[i]
		}
		Res, err = dbConnection[DBConnectionName].Exec(Query, paramPointers...)
	} else {
		Res, err = dbConnection[DBConnectionName].Exec(Query)
	}
	if err != nil {
		log.Println("UpdateDB : Error: " + err.Error() + " Executing Query: " + Query)
		return 0, false
	}
	var Num int64
	Num = 0
	switch strings.ToUpper(Query[0:6]) {
	case "UPDATE":
		count, err := Res.RowsAffected()
		if err == nil {
			Num = count
			log.Println(fmt.Sprintf("UpdateDB : %d Rows(s) Updated for Query: %s", Num, Query))
		} else {
			log.Println(fmt.Sprintf("UpdateDB: can't get updated rows, Error: %s", err.Error()))
		}
		break
	case "INSERT":
		count, err := Res.LastInsertId()
		if err == nil {
			Num = count
			log.Println(fmt.Sprintf("UpdateDB : Last Insert ID is %d for Query: %s", Num, Query))
		} else {
			log.Println(fmt.Sprintf("UpdateDB: can't get LastInsertID, Error: %s", err.Error()))
		}
		break
	case "DELETE":
		count, err := Res.RowsAffected()
		if err == nil {
			Num = count
			log.Println(fmt.Sprintf("UpdateDB : %d Rows(s) Updated for Query: %s", Num, Query))
		} else {
			log.Println(fmt.Sprintf("UpdateDB: can't get updated rows, Error: %s", err.Error()))
		}
		break
	default:
		break
	}
	return Num, true
}

// IsSQLSafe (strValue string) check if the string is safe to put in SQL
func IsSQLSafe(strValue string) bool {
	if (strings.Index(strValue, "'") > -1) || (strings.Index(strValue, ";") > -1) || (strings.Index(strValue, "`") > -1) {
		return false
	}
	return true
}
