package model

import (
	"database/sql"
	"github.com/dapr-platform/common"
	"time"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = common.LocalTime{}
)

/*
DB Table Details
-------------------------------------


Table: o_user
[ 0] id                                             VARCHAR(25)          null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 25      default: []
[ 1] tenant_id                                      VARCHAR(25)          null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 25      default: []
[ 2] mobile                                         VARCHAR(15)          null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 15      default: []
[ 3] email                                          VARCHAR(255)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 255     default: []
[ 4] identity                                       VARCHAR(255)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 255     default: []
[ 5] name                                           VARCHAR(255)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 255     default: []
[ 6] zh_name                                        VARCHAR(255)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 255     default: []
[ 7] gender                                         INT4                 null: false  primary: false  isArray: false  auto: false  col: INT4            len: -1      default: [0]
[ 8] address                                        VARCHAR(1024)        null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 1024    default: []
[ 9] password                                       VARCHAR(255)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 255     default: []
[10] type                                           INT4                 null: false  primary: false  isArray: false  auto: false  col: INT4            len: -1      default: [0]
[11] org_id                                         VARCHAR(255)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 255     default: []
[12] id_card                                        VARCHAR(255)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 255     default: []
[13] work_number                                    VARCHAR(255)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 255     default: []
[14] avatar_url                                     VARCHAR(255)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 255     default: []
[15] create_at                                      TIMESTAMP            null: false  primary: false  isArray: false  auto: false  col: TIMESTAMP       len: -1      default: []
[16] update_at                                      TIMESTAMP            null: false  primary: false  isArray: false  auto: false  col: TIMESTAMP       len: -1      default: []
[17] remark                                         VARCHAR(255)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 255     default: []
[18] status                                         INT4                 null: false  primary: false  isArray: false  auto: false  col: INT4            len: -1      default: [0]


JSON Sample
-------------------------------------
{    "id": "dcJIudoVkEIubLiEjZxgafume",    "tenant_id": "RWxSPovCbAOlFVAQZfvEoVHpO",    "mobile": "reLycKRukNIxdLXxaPbmiPVdN",    "email": "LDbpnMXrmfVklsAiGBZahNZaX",    "identity": "UaLnsQhnrAsOshWBBGnWGdKTv",    "name": "xagguPJlhKZCUTsZptIIgFWyx",    "zh_name": "sIgTWomEODsYoFyBXRhyXQSEu",    "gender": 0,    "address": "eKamQZpRYATabJxtnFYIkPilt",    "password": "TORNeXxDKEdyMCSYcnHakHBrT",    "type": 63,    "org_id": "qdiYQuvMvtckjJGUhLsFKbrmH",    "id_card": "lkNruRvKymKMcUhDvxcPAcsuC",    "work_number": "WsCEZLTnQCsvVXZhuLFYcnrsN",    "avatar_url": "tPqDmAWywVAVxyuxhrhoBMkUg",    "create_at": 55,    "update_at": 25,    "remark": "RsTtbTvAEHprOfwxjhkeMQLuZ",    "status": 37}



*/

var (
	User_FIELD_NAME_id = "id"

	User_FIELD_NAME_tenant_id = "tenant_id"

	User_FIELD_NAME_mobile = "mobile"

	User_FIELD_NAME_email = "email"

	User_FIELD_NAME_identity = "identity"

	User_FIELD_NAME_name = "name"

	User_FIELD_NAME_zh_name = "zh_name"

	User_FIELD_NAME_gender = "gender"

	User_FIELD_NAME_address = "address"

	User_FIELD_NAME_password = "password"

	User_FIELD_NAME_type = "type"

	User_FIELD_NAME_org_id = "org_id"

	User_FIELD_NAME_id_card = "id_card"

	User_FIELD_NAME_work_number = "work_number"

	User_FIELD_NAME_avatar_url = "avatar_url"

	User_FIELD_NAME_create_at = "create_at"

	User_FIELD_NAME_update_at = "update_at"

	User_FIELD_NAME_remark = "remark"

	User_FIELD_NAME_status = "status"
)

// User struct is a row record of the o_user table in the  database
type User struct {
	ID string `json:"id"` //Primary Key

	TenantID string `json:"tenant_id"` //租户ID

	Mobile string `json:"mobile"` //手机号

	Email string `json:"email"` //邮箱

	Identity string `json:"identity"` //用户标识

	Name string `json:"name"` //姓名

	ZhName string `json:"zh_name"` //中文名

	Gender int32 `json:"gender"` //性别(0:未知,1:男,2:女)

	Address string `json:"address"` //地址

	Password string `json:"password"` //密码

	Type int32 `json:"type"` //用户类型,1:管理员,2:普通用户,3:访客

	OrgID string `json:"org_id"` //组织ID

	IDCard string `json:"id_card"` //身份证

	WorkNumber string `json:"work_number"` //工号

	AvatarURL string `json:"avatar_url"` //头像

	CreateAt common.LocalTime `json:"create_at"` //创建时间

	UpdateAt common.LocalTime `json:"update_at"` //更新时间

	Remark string `json:"remark"` //备注

	Status int32 `json:"status"` //状态(1正常，2:禁止登陆，3:删除

}

var UserTableInfo = &TableInfo{
	Name: "o_user",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "id",
			Comment:            `Primary Key`,
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
			Name:               "tenant_id",
			Comment:            `租户ID`,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(25)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       25,
			GoFieldName:        "TenantID",
			GoFieldType:        "string",
			JSONFieldName:      "tenant_id",
			ProtobufFieldName:  "tenant_id",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "mobile",
			Comment:            `手机号`,
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
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "email",
			Comment:            `邮箱`,
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
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "identity",
			Comment:            `用户标识`,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(255)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       255,
			GoFieldName:        "Identity",
			GoFieldType:        "string",
			JSONFieldName:      "identity",
			ProtobufFieldName:  "identity",
			ProtobufType:       "string",
			ProtobufPos:        5,
		},

		&ColumnInfo{
			Index:              5,
			Name:               "name",
			Comment:            `姓名`,
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
			ProtobufPos:        6,
		},

		&ColumnInfo{
			Index:              6,
			Name:               "zh_name",
			Comment:            `中文名`,
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
			ProtobufPos:        7,
		},

		&ColumnInfo{
			Index:              7,
			Name:               "gender",
			Comment:            `性别(0:未知,1:男,2:女)`,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "INT4",
			DatabaseTypePretty: "INT4",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "INT4",
			ColumnLength:       -1,
			GoFieldName:        "Gender",
			GoFieldType:        "int32",
			JSONFieldName:      "gender",
			ProtobufFieldName:  "gender",
			ProtobufType:       "int32",
			ProtobufPos:        8,
		},

		&ColumnInfo{
			Index:              8,
			Name:               "address",
			Comment:            `地址`,
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
			ProtobufPos:        9,
		},

		&ColumnInfo{
			Index:              9,
			Name:               "password",
			Comment:            `密码`,
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
			ProtobufPos:        10,
		},

		&ColumnInfo{
			Index:              10,
			Name:               "type",
			Comment:            `用户类型,1:管理员,2:普通用户,3:访客`,
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
			ProtobufPos:        11,
		},

		&ColumnInfo{
			Index:              11,
			Name:               "org_id",
			Comment:            `组织ID`,
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
			ProtobufPos:        12,
		},

		&ColumnInfo{
			Index:              12,
			Name:               "id_card",
			Comment:            `身份证`,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(255)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       255,
			GoFieldName:        "IDCard",
			GoFieldType:        "string",
			JSONFieldName:      "id_card",
			ProtobufFieldName:  "id_card",
			ProtobufType:       "string",
			ProtobufPos:        13,
		},

		&ColumnInfo{
			Index:              13,
			Name:               "work_number",
			Comment:            `工号`,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(255)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       255,
			GoFieldName:        "WorkNumber",
			GoFieldType:        "string",
			JSONFieldName:      "work_number",
			ProtobufFieldName:  "work_number",
			ProtobufType:       "string",
			ProtobufPos:        14,
		},

		&ColumnInfo{
			Index:              14,
			Name:               "avatar_url",
			Comment:            `头像`,
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
			ProtobufPos:        15,
		},

		&ColumnInfo{
			Index:              15,
			Name:               "create_at",
			Comment:            `创建时间`,
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
			ProtobufPos:        16,
		},

		&ColumnInfo{
			Index:              16,
			Name:               "update_at",
			Comment:            `更新时间`,
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
			ProtobufPos:        17,
		},

		&ColumnInfo{
			Index:              17,
			Name:               "remark",
			Comment:            `备注`,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(255)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       255,
			GoFieldName:        "Remark",
			GoFieldType:        "string",
			JSONFieldName:      "remark",
			ProtobufFieldName:  "remark",
			ProtobufType:       "string",
			ProtobufPos:        18,
		},

		&ColumnInfo{
			Index:              18,
			Name:               "status",
			Comment:            `状态(1正常，2:禁止登陆，3:删除`,
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
			ProtobufPos:        19,
		},
	},
}

// TableName sets the insert table name for this struct type
func (u *User) TableName() string {
	return "o_user"
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
