package main

import (
	"net/http"
	"oauth2-server/model"
	"oauth2-server/service"

	"github.com/dapr-platform/common"
	"github.com/dchest/captcha"
	"github.com/go-oauth2/oauth2/v4"
)

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
	code, err := service.GetSmsCode(r.Context(), smsCodeGet.Phone)
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg(err.Error()))
		return
	}
	common.HttpSuccess(w, common.OK.WithData(code))
	return
}
