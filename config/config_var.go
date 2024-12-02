package config

import "os"

var ALI_SMS_ACCESS_ID = ""
var ALI_SMS_ACCESS_SECRET = ""
var ALI_SMS_REGION = ""
var ALI_SMS_SIGN_NAME = ""
var ALI_SMS_TEMPLATE_CODE = ""

func init() {
	if os.Getenv("ALI_SMS_ACCESS_ID") != "" {
		ALI_SMS_ACCESS_ID = os.Getenv("ALI_SMS_ACCESS_ID")
	}
	if os.Getenv("ALI_SMS_ACCESS_SECRET") != "" {
		ALI_SMS_ACCESS_SECRET = os.Getenv("ALI_SMS_ACCESS_SECRET")
	}
	if os.Getenv("ALI_SMS_REGION") != "" {
		ALI_SMS_REGION = os.Getenv("ALI_SMS_REGION")
	}
	if os.Getenv("ALI_SMS_SIGN_NAME") != "" {
		ALI_SMS_SIGN_NAME = os.Getenv("ALI_SMS_SIGN_NAME")
	}
	if os.Getenv("ALI_SMS_TEMPLATE_CODE") != "" {
		ALI_SMS_TEMPLATE_CODE = os.Getenv("ALI_SMS_TEMPLATE_CODE")
	}
}
