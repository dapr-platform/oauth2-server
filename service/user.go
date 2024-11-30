package service

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"net/url"
	"oauth2-server/model"
	"strconv"
	"time"

	"github.com/dapr-platform/common"
	"github.com/pkg/errors"
	"golang.org/x/exp/rand"
)

var smsVerfyCodeKeyPrefix = "SmsCode:"

func CheckMobileSmsCode(ctx context.Context, mobile, code string) (valid bool, err error) {
	exists, err := common.GetInStateStore(ctx, common.GetDaprClient(), common.GLOBAL_STATESTOR_NAME, smsVerfyCodeKeyPrefix+mobile)
	if err != nil {
		err = errors.Wrap(err, "GetInStateStore")
		return
	}
	if len(exists) == 0 {
		err = errors.New("验证码不存在")
		return
	}
	var x uint32
	err = binary.Read(bytes.NewBuffer(exists), binary.BigEndian, &x)
	if err != nil {
		err = errors.Wrap(err, "验证码处理错误")
		return
	}
	codi, err := strconv.Atoi(code)
	if err != nil {
		err = errors.Wrap(err, "验证码不是数字")
		return
	}
	if x != uint32(codi) {
		err = errors.New("验证码错误")
		return
	}
	valid = true
	return
}
func SendSmsCode(ctx context.Context, mobile string) (code string, err error) {
	exists, err := common.GetInStateStore(ctx, common.GetDaprClient(), common.GLOBAL_STATESTOR_NAME, smsVerfyCodeKeyPrefix+mobile)
	if err != nil {
		err = errors.Wrap(err, "GetInStateStore")
		return
	}
	if exists != nil {
		return
	}
	code, err = GenerateSmsCode(ctx, mobile)
	if err != nil {
		return
	}
	err = common.SaveInStateStore(ctx, common.GetDaprClient(), common.GLOBAL_STATESTOR_NAME, smsVerfyCodeKeyPrefix+mobile, []byte(code), true, time.Minute*1)
	if err != nil {
		return
	}
	common.Logger.Info("发送短信验证码", "mobile", mobile, "code", code)
	//TODO 发送短信
	return
}
func GenerateSmsCode(ctx context.Context, mobile string) (code string, err error) {
	code = strconv.Itoa(rand.Intn(10000))
	return
}
func CreateUser(ctx context.Context, user *model.User) (err error) {
	user.ID = common.NanoId()
	user.CreateAt = common.LocalTime(time.Now())
	user.UpdateAt = common.LocalTime(time.Now())
	exists, err := GetUserByFieldName(ctx, model.User_FIELD_NAME_name, user.Name, false)
	if err != nil {
		return errors.Wrap(err, "get user by name error")
	}
	if exists != nil {
		return errors.New("用户名已存在")
	}
	_, err = common.DbInsert[model.User](ctx, common.GetDaprClient(), *user, model.UserTableInfo.Name)
	if err != nil {
		return errors.Wrap(err, "db insert error")
	}
	return
}

func GetUserByFieldName(ctx context.Context, field, value string, isTravel bool) (user *model.User, err error) {
	value = url.QueryEscape(value)
	qstr := ""
	if isTravel {
		qstr = "&type=9"
	}

	datas, err := common.DbQuery[model.User](ctx, common.GetDaprClient(), model.UserTableInfo.Name, field+"="+value+qstr)
	if err != nil {
		err = errors.WithMessage(err, "db query error field="+field+" value="+value)
		return
	}
	if len(datas) == 0 {
		common.Logger.Error("not found field=" + field + " value=" + value)
		return
	}
	user = &datas[0]
	return
}
func GetUserPasswordByField(ctx context.Context, field, value string, isTravel bool) (password string, err error) {
	value = url.QueryEscape(value)
	qstr := ""
	if isTravel {
		qstr = "&type=9"
	}

	datas, err := common.DbQuery[model.User](ctx, common.GetDaprClient(), model.UserTableInfo.Name, field+"="+value+qstr)
	if err != nil {
		err = errors.WithMessage(err, "db query error field="+field+" value="+value)
		return
	}
	if len(datas) == 0 {
		common.Logger.Error("not found field=" + field + " value=" + value)
		return
	}
	password = datas[0].Password
	return
}
func GetUserByIdAndPassword(ctx context.Context, id, password string) (user *model.User, err error) {

	users, err := common.DbQuery[model.User](ctx, common.GetDaprClient(), model.UserTableInfo.Name, "id="+id+"&password="+password)
	if err != nil {
		common.Logger.Error("db query error ", err)
		return nil, nil
	}
	if len(users) == 0 {
		common.Logger.Error("user not found ", id)
		return nil, nil
	}
	user = &users[0]
	return
}

func SaveUserInfoInStore(ctx context.Context, id string) (err error) {
	users, err := common.DbQuery[model.User](ctx, common.GetDaprClient(), model.UserTableInfo.Name, "id="+id)
	if err != nil {
		common.Logger.Error("db query error", err)
		return errors.Wrap(err, "db query error")
	}
	if len(users) == 0 {
		common.Logger.Error("user not found", id)
		return errors.Wrap(err, "user not found "+id)
	}
	user := &users[0]
	buf, _ := json.Marshal(user)
	return common.SaveInStateStore(ctx, common.GetDaprClient(), common.GLOBAL_STATESTOR_NAME, common.USER_STATESTORE_KEY_PREFIX+user.ID, buf, true, time.Second*time.Duration(common.USER_EXPIRED_SECONDS))
}
