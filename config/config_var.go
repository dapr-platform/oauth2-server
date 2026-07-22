package config

import "os"

var ALI_SMS_ACCESS_ID = ""
var ALI_SMS_ACCESS_SECRET = ""
var ALI_SMS_REGION = ""
var ALI_SMS_SIGN_NAME = ""
var ALI_SMS_TEMPLATE_CODE = ""

var SSO_ENABLED = true
var SSO_BASE_URL = "http://10.6.160.102:82"
var SSO_TICKET_BASE_URL = "http://10.6.160.102:80"
var SSO_APP_KEY = "03754485bd7442a98205e8d90c088471"
var SSO_APP_SECRET = "5e88cabddb66461187c0516996a39d13"

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
	if val := os.Getenv("SSO_ENABLED"); val == "false" {
		SSO_ENABLED = false
	}
	if val := os.Getenv("SSO_BASE_URL"); val != "" {
		SSO_BASE_URL = val
	}
	if val := os.Getenv("SSO_APP_KEY"); val != "" {
		SSO_APP_KEY = val
	}
	if val := os.Getenv("SSO_APP_SECRET"); val != "" {
		SSO_APP_SECRET = val
	}
}
