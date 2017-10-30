package customisor

import (
	"log"
	"encoding/json"
	"github.com/go-ini/ini"
	"strings"
)

type MysqlJsonData struct {
	Mysqld struct {
		ReadOnly	bool	`json:"read_only"`
		TransactionReadOnly	bool	`json:"transaction_read_only"`
		AutoCommit	bool	`json:"autocommit"`
		BigTables	bool	`json:"big_tables"`
		LowPriorityUpdates	bool	`json:"low_priority_updates"`
		EventScheduler	string	`json:"event_scheduler"`
		CompletionType	string	`json:"completion_type"`
		ConcurrentInsert	string	`json:"concurrent_insert"`
		DivPrecisionIncrement	int64	`json:"div_precision_increment"`
		SqlMode []string	`json:"sql_mode"`
		TransactionIsolation	string	`json:"transaction_isolation"`
		DefaultTimeZone	string	`json:"default_time_zone"`
		DefaultWeekFormat	string	`json:"default_week_format"`
		ConnectTimeout	int64	`json:"connect_timeout"`
		LockWaitTimeout	int64	`json:"lock_wait_timeout"`
		InteractiveTimeout	int64	`json:"interactive_timeout"`
		WaitTimeout	int64	`json:"wait_timeout"`
		NetReadTimeout	int64	`json:"net_read_timeout"`
		NetWriteTimeout	int64	`json:"net_write_timeout"`
		NetRetryCount	int64	`json:"net_retry_count"`
		MaxConnections	int64	`json:"max_connections"`
		MaxUserConnections	int64	`json:"max_user_connections"`
		MaxConnectErrors	int64	`json:"max_connect_errors"`
		MaxErrorCount	int64	`json:"max_error_count"`
		MaxSpRecursionDepth	int64	`json:"max_sp_recursion_depth"`
		QueryCacheLimit	int64	`json:"query_cache_limit"`
		QueryCacheSize int64	`json:"query_cache_size"`
		QueryCacheType string	`json:"query_cache_type"`
	}
}

type MysqlParserData struct {
	JsonData MysqlJsonData
	Section ini.Section
}

func (mysql *MysqlParserData) MysqlJsonLoader(data []byte) {
	err := json.Unmarshal(data, &mysql.JsonData)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	//log.Printf("%s\n", mysql.JsonData.Mysqld.TransactionIsolation)
}

func (mysql *MysqlParserData) ApplyCustomisations() {
	_, iniFile, iniFilePath := unpack_etc_ini(mysql.Section, true)
	if iniFile != nil {
		mysqld, err := iniFile.GetSection("mysqld")
		if err == nil {
			UpdateBoolKey("MYSQL", mysqld, "read_only", mysql.JsonData.Mysqld.ReadOnly)
			UpdateBoolKey("MYSQL", mysqld, "transaction_read_only", mysql.JsonData.Mysqld.TransactionReadOnly)
			UpdateBoolKey("MYSQL", mysqld, "autocommit", mysql.JsonData.Mysqld.AutoCommit)
			UpdateBoolKey("MYSQL", mysqld, "big_tables", mysql.JsonData.Mysqld.BigTables)
			UpdateBoolKey("MYSQL", mysqld, "low_priority_updates", mysql.JsonData.Mysqld.LowPriorityUpdates)

			UpdateStringKey("MYSQL", mysqld, "event_scheduler", mysql.JsonData.Mysqld.EventScheduler)
			UpdateStringKey("MYSQL", mysqld, "completion_type", mysql.JsonData.Mysqld.CompletionType)
			UpdateStringKey("MYSQL", mysqld, "concurrent_insert", mysql.JsonData.Mysqld.ConcurrentInsert)

			UpdateInt64Key("MYSQL", mysqld, "div_precision_increment", mysql.JsonData.Mysqld.DivPrecisionIncrement)

			UpdateStringKey("MYSQL", mysqld, "sql_mode", strings.Join(mysql.JsonData.Mysqld.SqlMode, ","))
			UpdateStringKey("MYSQL", mysqld, "transaction_isolation", mysql.JsonData.Mysqld.TransactionIsolation)
			UpdateStringKey("MYSQL", mysqld, "default_time_zone", mysql.JsonData.Mysqld.DefaultTimeZone)
			UpdateStringKey("MYSQL", mysqld, "default_week_format", mysql.JsonData.Mysqld.DefaultWeekFormat)

			UpdateInt64Key("MYSQL", mysqld, "connect_timeout", mysql.JsonData.Mysqld.ConnectTimeout)
			UpdateInt64Key("MYSQL", mysqld, "lock_wait_timeout", mysql.JsonData.Mysqld.LockWaitTimeout)
			UpdateInt64Key("MYSQL", mysqld, "interactive_timeout", mysql.JsonData.Mysqld.InteractiveTimeout)
			UpdateInt64Key("MYSQL", mysqld, "wait_timeout", mysql.JsonData.Mysqld.WaitTimeout)
			UpdateInt64Key("MYSQL", mysqld, "net_read_timeout", mysql.JsonData.Mysqld.NetReadTimeout)
			UpdateInt64Key("MYSQL", mysqld, "net_write_timeout", mysql.JsonData.Mysqld.NetWriteTimeout)
			UpdateInt64Key("MYSQL", mysqld, "net_retry_count", mysql.JsonData.Mysqld.NetRetryCount)
			UpdateInt64Key("MYSQL", mysqld, "max_connections", mysql.JsonData.Mysqld.MaxConnections)
			UpdateInt64Key("MYSQL", mysqld, "max_user_connections", mysql.JsonData.Mysqld.MaxUserConnections)
			UpdateInt64Key("MYSQL", mysqld, "max_connect_errors", mysql.JsonData.Mysqld.MaxConnectErrors)
			UpdateInt64Key("MYSQL", mysqld, "max_error_count", mysql.JsonData.Mysqld.MaxErrorCount)
			UpdateInt64Key("MYSQL", mysqld, "max_sp_recursion_depth", mysql.JsonData.Mysqld.MaxSpRecursionDepth)
			UpdateInt64Key("MYSQL", mysqld, "query_cache_limit", mysql.JsonData.Mysqld.QueryCacheLimit)
			UpdateInt64Key("MYSQL", mysqld, "query_cache_size", mysql.JsonData.Mysqld.QueryCacheSize)

			UpdateStringKey("MYSQL", mysqld, "query_cache_type", mysql.JsonData.Mysqld.QueryCacheType)

			SaveIniFile(*iniFile, iniFilePath, "mysql.ini")
		}
	} else {
		log.Print("No ini file")
	}
}
