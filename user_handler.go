package main

import (
	"net/http"
	"oauth2-server/model"
	"oauth2-server/service"

	"github.com/dapr-platform/common"
	"github.com/dchest/captcha"
	"github.com/go-oauth2/oauth2/v4"
)

// 用户注册
// @Description 用户注册
// @Tags Oauth2
// @Param sms_code query string false "短信验证码,如果系统配置为不需要验证码，则不传"
// @Param data body model.User true "{}"
// @Produce  json
// @Success 200 {object} common.Response{data=model.User} "用户信息"
// @Failure 500 {object} common.Response "错误code和错误信息"
// @Router /users/register [post]
func userRegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := common.ReadRequestBody(r, &user)
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg("user register error"))
		return
	}

	count, err := common.DbGetCount(r.Context(), common.GetDaprClient(), model.UserTableInfo.Name, "name", "name="+user.Name)
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg("系统错误:" + err.Error()))
		return
	}
	if count > 0 {
		common.HttpResult(w, common.ErrParam.AppendMsg("用户名已存在"))
		return
	}
	smsCode := r.URL.Query().Get("sms_code")
	if smsCode == "" {
		if REGISTER_SMS_CODE {
			common.HttpResult(w, common.ErrParam.AppendMsg("sms code blank"))
			return
		}
	} else {
		valid, err := service.CheckMobileSmsCode(r.Context(), user.Mobile, smsCode)
		if err != nil {
			common.HttpResult(w, common.ErrParam.AppendMsg(err.Error()))
			return
		}
		if !valid {
			common.HttpResult(w, common.ErrParam.AppendMsg("验证码错误"))
			return
		}
	}

	err = service.CreateUser(r.Context(), &user)
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg(err.Error()))
		return
	}
	common.HttpSuccess(w, common.OK.WithData(user))
}

// @Summary 用户登录
// @Description 用户登录,简单方式
// @Tags Oauth2
// @Param data body model.UserLogin true "{}"
// @Produce  json
// @Success 200 {object} model.TokenInfo "token信息"
// @Failure 500 {object} string "错误code和错误信息"
// @Router /users/login [post]
func userLoginHandler(w http.ResponseWriter, r *http.Request) {
	var userLogin model.UserLogin
	err := common.ReadRequestBody(r, &userLogin)
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg("user login error"))
		return
	}
	if VERIFY_CAPTCHA {
		vKey := userLogin.VerifyKey
		vVal := userLogin.VerifyValue
		if vKey == "" || vVal == "" {
			common.HttpResult(w, common.ErrParam.AppendMsg("verify value blank"))
			return
		}
		if !captcha.VerifyString(vKey, vVal) {
			common.HttpResult(w, common.ErrParam.AppendMsg("verify  error"))
			return
		}

	}
	user, err := service.GetUserByFieldName(r.Context(), "name", userLogin.UserName, false)
	if err != nil {
		http.Error(w, "系统错误:"+err.Error(), http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "用户不存在", 499)
		return
	}
	if user.Status != 1 {
		http.Error(w, "用户已停用", 498)
		return
	}

	gt := oauth2.GrantType("password")
	tgr := &oauth2.TokenGenerateRequest{
		ClientID:     "dapr-client",
		ClientSecret: "123456",
		Request:      r,
	}
	userID, err := oauthServer.PasswordAuthorizationHandler(r.Context(), "", user.ID, userLogin.Password)
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg(err.Error()))
		return
	} else if userID == "" {
		common.HttpResult(w, common.ErrParam.AppendMsg("ErrInvalidGrant"))
		return
	}
	tgr.UserID = userID
	ti, err := oauthServer.GetAccessToken(r.Context(), gt, tgr)
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg(err.Error()))
		return
	}
	common.HttpSuccess(w, common.OK.WithData(oauthServer.GetTokenData(ti)))
	return

}

// @Summary 发送短信验证码
// @Description 发送短信验证码
// @Tags Oauth2
// @Param data body model.SmsCodeGet true "{}"
// @Produce  json
// @Success 200 {object} model.SmsCodeGet "短信验证码"
// @Failure 500 {object} string "错误code和错误信息"
// @Router /sms-code/send [post]
func smsCodeSendHandler(w http.ResponseWriter, r *http.Request) {
	var smsCodeGet model.SmsCodeGet
	err := common.ReadRequestBody(r, &smsCodeGet)
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg("sms code get error"))
		return
	}
	code, err := service.SendSmsCode(r.Context(), smsCodeGet.Phone)
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg(err.Error()))
		return
	}
	common.HttpSuccess(w, common.OK.WithData(code))
	return
}
