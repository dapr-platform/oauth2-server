package model

import (
	"database/sql"
	"github.com/dapr-platform/common"
	"github.com/satori/go.uuid"
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


Table: p_user
[ 0] id                                             VARCHAR(25)          null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 25      default: []
[ 1] mobile                                         VARCHAR(15)          null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 15      default: []
[ 2] email                                          VARCHAR(255)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 255     default: []
[ 3] name                                           VARCHAR(255)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 255     default: []
[ 4] zh_name                                        VARCHAR(255)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 255     default: []
[ 5] nick_name                                      VARCHAR(255)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 255     default: []
[ 6] password                                       VARCHAR(255)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 255     default: []
[ 7] type                                           INT4                 null: false  primary: false  isArray: false  auto: false  col: INT4            len: -1      default: []
[ 8] org_id                                         VARCHAR(255)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 255     default: []
[ 9] dep_id                                         VARCHAR(255)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 255     default: []
[10] avatar_url                                     VARCHAR(255)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 255     default: []
[11] intialized                                     INT4                 null: false  primary: false  isArray: false  auto: false  col: INT4            len: -1      default: []
[12] create_at                                      TIMESTAMP            null: false  primary: false  isArray: false  auto: false  col: TIMESTAMP       len: -1      default: []
[13] update_at                                      TIMESTAMP            null: false  primary: false  isArray: false  auto: false  col: TIMESTAMP       len: -1      default: []
[14] feishu_id                                      VARCHAR(255)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 255     default: []
[15] status                                         INT4                 null: false  primary: false  isArray: false  auto: false  col: INT4            len: -1      default: []
[16] sex                                            INT4                 null: false  primary: false  isArray: false  auto: false  col: INT4            len: -1      default: [0]
[17] address                                        VARCHAR(1024)        null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 1024    default: []
[18] province_id                                    VARCHAR(255)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 255     default: []
[19] city_id                                        VARCHAR(255)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 255     default: []


JSON Sample
-------------------------------------
{    "id": "ftcsOHkBonFdNkEpxnmMdlBEs",    "mobile": "CeovGZCKxjDkEgIuyLfOmXdID",    "email": "VKySJHaIWfnDFwIZEcEgVQngp",    "name": "MFmYXgeBhHcPADneDmienBdIL",    "zh_name": "dPMCEDcikYVIIAbuJmVFVTMNM",    "nick_name": "MtFxGpCKInKIdnrqfOvOrBNkF",    "password": "NSeropIVPiFvAjNkgWXTXaIws",    "type": 69,    "org_id": "jQJnMGCKEreOXIxnlYbbCYcbb",    "dep_id": "CgDmjYebKhmaJPXaExuZAXAIS",    "avatar_url": "ipNqDkVIAOAeIyZhLiAPITbFU",    "intialized": 23,    "create_at": 99,    "update_at": 86,    "feishu_id": "gKQsvAPjeKonNMIcLTAsluFit",    "status": 28,    "sex": 88,    "address": "CbsyUfDIFcgtIoTGtkaHwXexB",    "province_id": "YkCjByvnSawEFmWVJESrelBHO",    "city_id": "dbWXIXUDvGVdZkXufoOnYKheU"}



*/

var (
	User_FIELD_NAME_id = "id"

	User_FIELD_NAME_mobile = "mobile"

	User_FIELD_NAME_email = "email"

	User_FIELD_NAME_name = "name"

	User_FIELD_NAME_zh_name = "zh_name"

	User_FIELD_NAME_nick_name = "nick_name"

	User_FIELD_NAME_password = "password"

	User_FIELD_NAME_type = "type"

	User_FIELD_NAME_org_id = "org_id"

	User_FIELD_NAME_dep_id = "dep_id"

	User_FIELD_NAME_avatar_url = "avatar_url"

	User_FIELD_NAME_intialized = "intialized"

	User_FIELD_NAME_create_at = "create_at"

	User_FIELD_NAME_update_at = "update_at"

	User_FIELD_NAME_feishu_id = "feishu_id"

	User_FIELD_NAME_status = "status"

	User_FIELD_NAME_sex = "sex"

	User_FIELD_NAME_address = "address"

	User_FIELD_NAME_province_id = "province_id"

	User_FIELD_NAME_city_id = "city_id"
)

// User struct is a row record of the p_user table in the thingsdb database
type User struct {
	//[ 0] id                                             VARCHAR(25)          null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 25      default: []
	ID string `json:"id"`
	//[ 1] mobile                                         VARCHAR(15)          null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 15      default: []
	Mobile string `json:"mobile"`
	//[ 2] email                                          VARCHAR(255)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 255     default: []
	Email string `json:"email"`
	//[ 3] name                                           VARCHAR(255)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 255     default: []
	Name string `json:"name"`
	//[ 4] zh_name                                        VARCHAR(255)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 255     default: []
	ZhName string `json:"zh_name"`
	//[ 5] nick_name                                      VARCHAR(255)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 255     default: []
	NickName string `json:"nick_name"`
	//[ 6] password                                       VARCHAR(255)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 255     default: []
	Password string `json:"password"`
	//[ 7] type                                           INT4                 null: false  primary: false  isArray: false  auto: false  col: INT4            len: -1      default: []
	Type int32 `json:"type"`
	//[ 8] org_id                                         VARCHAR(255)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 255     default: []
	OrgID string `json:"org_id"`
	//[ 9] dep_id                                         VARCHAR(255)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 255     default: []
	DepID string `json:"dep_id"`
	//[10] avatar_url                                     VARCHAR(255)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 255     default: []
	AvatarURL string `json:"avatar_url"`
	//[11] intialized                                     INT4                 null: false  primary: false  isArray: false  auto: false  col: INT4            len: -1      default: []
	Intialized int32 `json:"intialized"`
	//[12] create_at                                      TIMESTAMP            null: false  primary: false  isArray: false  auto: false  col: TIMESTAMP       len: -1      default: []
	CreateAt common.LocalTime `json:"create_at"`
	//[13] update_at                                      TIMESTAMP            null: false  primary: false  isArray: false  auto: false  col: TIMESTAMP       len: -1      default: []
	UpdateAt common.LocalTime `json:"update_at"`
	//[14] feishu_id                                      VARCHAR(255)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 255     default: []
	FeishuID string `json:"feishu_id"`
	//[15] status                                         INT4                 null: false  primary: false  isArray: false  auto: false  col: INT4            len: -1      default: []
	Status int32 `json:"status"`
	//[16] sex                                            INT4                 null: false  primary: false  isArray: false  auto: false  col: INT4            len: -1      default: [0]
	Sex int32 `json:"sex"`
	//[17] address                                        VARCHAR(1024)        null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 1024    default: []
	Address string `json:"address"`
	//[18] province_id                                    VARCHAR(255)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 255     default: []
	ProvinceID string `json:"province_id"`
	//[19] city_id                                        VARCHAR(255)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 255     default: []
	CityID string `json:"city_id"`
}

var UserTableInfo = &TableInfo{
	Name: "p_user",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "id",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(25)",
			IsPrimaryKey:       true,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       25,
			GoFieldName:        "ID",
			GoFieldType:        "string",
			JSONFieldName:      "id",
			ProtobufFieldName:  "id",
			ProtobufType:       "string",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "mobile",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(15)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       15,
			GoFieldName:        "Mobile",
			GoFieldType:        "string",
			JSONFieldName:      "mobile",
			ProtobufFieldName:  "mobile",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "email",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(255)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       255,
			GoFieldName:        "Email",
			GoFieldType:        "string",
			JSONFieldName:      "email",
			ProtobufFieldName:  "email",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "name",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(255)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       255,
			GoFieldName:        "Name",
			GoFieldType:        "string",
			JSONFieldName:      "name",
			ProtobufFieldName:  "name",
			ProtobufType:       "string",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "zh_name",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(255)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       255,
			GoFieldName:        "ZhName",
			GoFieldType:        "string",
			JSONFieldName:      "zh_name",
			ProtobufFieldName:  "zh_name",
			ProtobufType:       "string",
			ProtobufPos:        5,
		},

		&ColumnInfo{
			Index:              5,
			Name:               "nick_name",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(255)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       255,
			GoFieldName:        "NickName",
			GoFieldType:        "string",
			JSONFieldName:      "nick_name",
			ProtobufFieldName:  "nick_name",
			ProtobufType:       "string",
			ProtobufPos:        6,
		},

		&ColumnInfo{
			Index:              6,
			Name:               "password",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(255)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       255,
			GoFieldName:        "Password",
			GoFieldType:        "string",
			JSONFieldName:      "password",
			ProtobufFieldName:  "password",
			ProtobufType:       "string",
			ProtobufPos:        7,
		},

		&ColumnInfo{
			Index:              7,
			Name:               "type",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "INT4",
			DatabaseTypePretty: "INT4",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "INT4",
			ColumnLength:       -1,
			GoFieldName:        "Type",
			GoFieldType:        "int32",
			JSONFieldName:      "type",
			ProtobufFieldName:  "type",
			ProtobufType:       "int32",
			ProtobufPos:        8,
		},

		&ColumnInfo{
			Index:              8,
			Name:               "org_id",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(255)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       255,
			GoFieldName:        "OrgID",
			GoFieldType:        "string",
			JSONFieldName:      "org_id",
			ProtobufFieldName:  "org_id",
			ProtobufType:       "string",
			ProtobufPos:        9,
		},

		&ColumnInfo{
			Index:              9,
			Name:               "dep_id",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(255)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       255,
			GoFieldName:        "DepID",
			GoFieldType:        "string",
			JSONFieldName:      "dep_id",
			ProtobufFieldName:  "dep_id",
			ProtobufType:       "string",
			ProtobufPos:        10,
		},

		&ColumnInfo{
			Index:              10,
			Name:               "avatar_url",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(255)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       255,
			GoFieldName:        "AvatarURL",
			GoFieldType:        "string",
			JSONFieldName:      "avatar_url",
			ProtobufFieldName:  "avatar_url",
			ProtobufType:       "string",
			ProtobufPos:        11,
		},

		&ColumnInfo{
			Index:              11,
			Name:               "intialized",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "INT4",
			DatabaseTypePretty: "INT4",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "INT4",
			ColumnLength:       -1,
			GoFieldName:        "Intialized",
			GoFieldType:        "int32",
			JSONFieldName:      "intialized",
			ProtobufFieldName:  "intialized",
			ProtobufType:       "int32",
			ProtobufPos:        12,
		},

		&ColumnInfo{
			Index:              12,
			Name:               "create_at",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "TIMESTAMP",
			DatabaseTypePretty: "TIMESTAMP",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "TIMESTAMP",
			ColumnLength:       -1,
			GoFieldName:        "CreateAt",
			GoFieldType:        "common.LocalTime",
			JSONFieldName:      "create_at",
			ProtobufFieldName:  "create_at",
			ProtobufType:       "uint64",
			ProtobufPos:        13,
		},

		&ColumnInfo{
			Index:              13,
			Name:               "update_at",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "TIMESTAMP",
			DatabaseTypePretty: "TIMESTAMP",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "TIMESTAMP",
			ColumnLength:       -1,
			GoFieldName:        "UpdateAt",
			GoFieldType:        "common.LocalTime",
			JSONFieldName:      "update_at",
			ProtobufFieldName:  "update_at",
			ProtobufType:       "uint64",
			ProtobufPos:        14,
		},

		&ColumnInfo{
			Index:              14,
			Name:               "feishu_id",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(255)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       255,
			GoFieldName:        "FeishuID",
			GoFieldType:        "string",
			JSONFieldName:      "feishu_id",
			ProtobufFieldName:  "feishu_id",
			ProtobufType:       "string",
			ProtobufPos:        15,
		},

		&ColumnInfo{
			Index:              15,
			Name:               "status",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "INT4",
			DatabaseTypePretty: "INT4",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "INT4",
			ColumnLength:       -1,
			GoFieldName:        "Status",
			GoFieldType:        "int32",
			JSONFieldName:      "status",
			ProtobufFieldName:  "status",
			ProtobufType:       "int32",
			ProtobufPos:        16,
		},

		&ColumnInfo{
			Index:              16,
			Name:               "sex",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "INT4",
			DatabaseTypePretty: "INT4",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "INT4",
			ColumnLength:       -1,
			GoFieldName:        "Sex",
			GoFieldType:        "int32",
			JSONFieldName:      "sex",
			ProtobufFieldName:  "sex",
			ProtobufType:       "int32",
			ProtobufPos:        17,
		},

		&ColumnInfo{
			Index:              17,
			Name:               "address",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(1024)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       1024,
			GoFieldName:        "Address",
			GoFieldType:        "string",
			JSONFieldName:      "address",
			ProtobufFieldName:  "address",
			ProtobufType:       "string",
			ProtobufPos:        18,
		},

		&ColumnInfo{
			Index:              18,
			Name:               "province_id",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(255)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       255,
			GoFieldName:        "ProvinceID",
			GoFieldType:        "string",
			JSONFieldName:      "province_id",
			ProtobufFieldName:  "province_id",
			ProtobufType:       "string",
			ProtobufPos:        19,
		},

		&ColumnInfo{
			Index:              19,
			Name:               "city_id",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(255)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       255,
			GoFieldName:        "CityID",
			GoFieldType:        "string",
			JSONFieldName:      "city_id",
			ProtobufFieldName:  "city_id",
			ProtobufType:       "string",
			ProtobufPos:        20,
		},
	},
}

// TableName sets the insert table name for this struct type
func (u *User) TableName() string {
	return "p_user"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (u *User) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (u *User) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (u *User) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (u *User) TableInfo() *TableInfo {
	return UserTableInfo
}
