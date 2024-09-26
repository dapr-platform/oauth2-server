package model

import (
	"database/sql"
	"github.com/dapr-platform/common"
	"github.com/guregu/null"
	"time"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
	_ = uuid.UUID{}
	_ = common.LocalTime{}
)

/*
DB Table Details
-------------------------------------


Table: p_client_info
[ 0] id                                             VARCHAR(32)          null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 32      default: []
[ 1] password                                       VARCHAR(32)          null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 32      default: []


JSON Sample
-------------------------------------
{    "id": "wPpvfbKQfByLjRQVEaRiXrrOS",    "password": "QSPQGDdYbLBMHkYPmCHPxPPLL"}



*/

var (
	Client_info_FIELD_NAME_id = "id"

	Client_info_FIELD_NAME_password = "password"
)

// Client_info struct is a row record of the p_client_info table in the thingsdb database
type Client_info struct {
	//[ 0] id                                             VARCHAR(32)          null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 32      default: []
	ID string `json:"id"`
	//[ 1] password                                       VARCHAR(32)          null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 32      default: []
	Password string `json:"password"`
}

var Client_infoTableInfo = &TableInfo{
	Name: "p_client_info",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "id",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(32)",
			IsPrimaryKey:       true,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       32,
			GoFieldName:        "ID",
			GoFieldType:        "string",
			JSONFieldName:      "id",
			ProtobufFieldName:  "id",
			ProtobufType:       "string",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "password",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(32)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       32,
			GoFieldName:        "Password",
			GoFieldType:        "string",
			JSONFieldName:      "password",
			ProtobufFieldName:  "password",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},
	},
}

// TableName sets the insert table name for this struct type
func (c *Client_info) TableName() string {
	return "p_client_info"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (c *Client_info) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (c *Client_info) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (c *Client_info) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (c *Client_info) TableInfo() *TableInfo {
	return Client_infoTableInfo
}
