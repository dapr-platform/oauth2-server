package sms

import (
	"errors"
	"fmt"
	"sync"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v4/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

var (
	client *dysmsapi20170525.Client
	clientOnce sync.Once
	clientErr error
)

func GetClient(region, accessId, accessSecret string) (*dysmsapi20170525.Client, error) {
	clientOnce.Do(func() {
		config := &openapi.Config{
			AccessKeyId:     tea.String(accessId),
			AccessKeySecret: tea.String(accessSecret),
		}
		config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
		_result, _err := dysmsapi20170525.NewClient(config)
		if _err != nil {
			clientErr = fmt.Errorf("create aliyun sms client failed: %v", _err)
			return
		}
		client = _result
	})

	if clientErr != nil {
		return nil, clientErr
	}
	
	if client == nil {
		return nil, fmt.Errorf("aliyun sms client is nil")
	}

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
		PhoneNumbers:  tea.String(phone),
		SignName:      tea.String(signName),
		TemplateCode:  tea.String(templateCode),
		TemplateParam: tea.String("{\"code\":\"" + code + "\"}"),
	}
	runtime := &util.RuntimeOptions{}
	
	resp, err := client.SendSmsWithOptions(request, runtime)
	if err != nil {
		return fmt.Errorf("send sms failed: %v", err)
	}

	if resp == nil {
		return fmt.Errorf("send sms response is nil")
	}

	if resp.StatusCode == nil || *resp.StatusCode != 200 {
		return fmt.Errorf("send sms code failed with status code: %v", resp.StatusCode)
	}

	if resp.Body == nil || resp.Body.Code == nil || *resp.Body.Code != "OK" {
		var msg string
		if resp.Body != nil && resp.Body.Message != nil {
			msg = *resp.Body.Message
		}
		return fmt.Errorf("send sms code failed: %v", msg)
	}

	return nil
}
