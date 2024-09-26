package config

import (
	"os"
)

var DBNAME = "thingsdb"

var CLIENT_INFO_TABLE_NAME = "p_client_info"
var USER_INFO_TABLE_NAME = "p_user"

func init() {

	if val := os.Getenv("DB_NAME"); val != "" {
		DBNAME = val
	}

	if val := os.Getenv("CLIENT_INFO_TABLE_NAME"); val != "" {
		CLIENT_INFO_TABLE_NAME = val
	}
	if val := os.Getenv("USER_INFO_TABLE_NAME"); val != "" {
		USER_INFO_TABLE_NAME = val
	}
}