package sms

import (
	"errors"
	"fmt"
	"sync"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v4/client"
	console "github.com/alibabacloud-go/tea-console/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

var client *dysmsapi20170525.Client
var clientOnce sync.Once

func GetClient(region, accessId, accessSecret string) (*dysmsapi20170525.Client, error) {
	clientOnce.Do(func() {
		config := &openapi.Config{
			// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_ID。
			AccessKeyId: tea.String(accessId),
			// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_SECRET。
			AccessKeySecret: tea.String(accessSecret),
		}
		// Endpoint 请参考 https://api.aliyun.com/product/Dysmsapi
		config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
		_result, _err := dysmsapi20170525.NewClient(config)
		if _err != nil {
			panic(_err)
		}
		client = _result
	})
	return client, nil
}

func SendSmsCode(region, accessId, accessSecret, signName, templateCode, phone, code string) error {
	// 参数校验
	if region == "" || accessId == "" || accessSecret == "" || signName == "" || templateCode == "" || phone == "" || code == "" {
		return errors.New("required parameters cannot be empty")
	}

	client, err := GetClient(region, accessId, accessSecret)
	if err != nil {
		return fmt.Errorf("create client failed: %v", err)
	}

	request := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:   tea.String(phone),
		SignName:       tea.String(signName),
		TemplateCode:   tea.String(templateCode),
		TemplateParam:  tea.String("{\"code\":\"" + code + "\"}"),
	}
	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		resp, _err := client.SendSmsWithOptions(request, runtime)
		if _err != nil {
			return _err
		}

		console.Log(util.ToJSONString(resp))

		return nil
	}()

	return tryErr
}
