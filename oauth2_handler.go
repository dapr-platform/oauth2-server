package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"oauth2-server/event"
	"oauth2-server/service"
	"os"
	"strings"
	"time"

	"github.com/dapr-platform/common"
	"github.com/dchest/captcha"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-session/session"
	"github.com/pkg/errors"
)

func passwordAuthHandler(ctx context.Context, clientId, id, password string) (userID string, err error) {

	user, err := service.GetUserByIdAndPassword(ctx, id, password)
	if err != nil {
		common.Logger.Error("GetUserByIdAndPassword error", err)
		return "", errors.Wrap(err, "getuserByid")
	}
	if user == nil {
		common.Logger.Error("GetUserByIdAndPassword not found" + id)
		return "", nil
	}
	msg := &common.InternalMessage{
		common.INTERNAL_MESSAGE_KEY_TYPE:    common.INTERNAL_MESSAGE_TYPE_USER_LOGIN,
		common.INTERNAL_MESSAGE_KEY_USER_ID: id,
	}
	err = event.PublishInternalMessage(ctx, msg)
	if err != nil {
		common.Logger.Error(err)
	}
	return id, nil

}

// @Summary login
// @Description 登录
// @Tags Oauth2
// @Router /login [get]
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if dumpvar {
		_ = dumpRequest(os.Stdout, "login", r) // Ignore the error
	}
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == "POST" {
		if r.Form == nil {
			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		store.Set("LoggedInUserID", r.Form.Get("username"))
		store.Save()

		w.Header().Set("Location", "/auth")
		w.WriteHeader(http.StatusFound)
		return
	}
	outputHTML(w, r, "static/login.html")
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	if dumpvar {
		_ = dumpRequest(os.Stdout, "auth", r) // Ignore the error
	}
	store, err := session.Start(nil, w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, ok := store.Get("LoggedInUserID"); !ok {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}

	outputHTML(w, r, "static/auth.html")
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Form == nil {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	redirectURI := r.Form.Get("redirect_uri")
	if _, err := url.Parse(redirectURI); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	store.Delete("LoggedInUserID")
	store.Save()

	w.Header().Set("Location", redirectURI)
	w.WriteHeader(http.StatusFound)
}

func authorizeHandler(w http.ResponseWriter, r *http.Request) {
	if dumpvar {
		dumpRequest(os.Stdout, "authorize", r)
	}

	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var form url.Values
	if v, ok := store.Get("ReturnUri"); ok {
		form = v.(url.Values)
	}
	r.Form = form

	store.Delete("ReturnUri")
	store.Save()

	err = oauthServer.HandleAuthorizeRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

// @Summary 获取token
// @Description 获取token
// @Tags Oauth2
// @Produce  json
// @Success 200 {object} model.TokenInfo "token信息"
// @Failure 500 {object} string "错误code和错误信息"
// @Router /oauth/token [post]
func tokenHandler(w http.ResponseWriter, r *http.Request) {
	if dumpvar {
		_ = dumpRequest(os.Stdout, "token", r) // Ignore the error
	}

	grantType := r.FormValue("grant_type")

	if grantType != "refresh_token" { //refresh token
		username := r.FormValue("username")
		var field string
		field = "name"

		value := username

		isTravelStr := r.FormValue("is_travel")
		isTravel := false
		if isTravelStr == "1" {
			isTravel = true
		}
		sms_code := r.FormValue("sms_code")
		if sms_code != "" { //如果是验证码登录，那么就先校验验证码，成功后，获取密码，后面走oauth流程
			valid, err := service.CheckMobileSmsCode(r.Context(), value, sms_code)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if !valid {
				http.Error(w, "短信验证码错误", http.StatusNotAcceptable)
				return
			}
			passwd, err := service.GetUserPasswordByField(r.Context(), field, value, isTravel)
			if err != nil {
				http.Error(w, "获取用户错误 "+err.Error(), http.StatusInternalServerError)
				return
			}
			r.Form.Set("password", passwd)

		}
		if VERIFY_CAPTCHA && !isTravel && sms_code == "" {
			vKey := r.FormValue("verify_key")
			vVal := r.FormValue("verify_value")
			if vKey == "" || vVal == "" {
				http.Error(w, "验证码为空", http.StatusBadRequest)
				return
			}
			if !captcha.VerifyString(vKey, vVal) {
				http.Error(w, "验证码错误", http.StatusNotAcceptable)
				return
			}

		}
		user, err := service.GetUserByFieldName(r.Context(), field, value, isTravel)
		if err != nil {
			http.Error(w, "GetUserByFieldName "+err.Error(), http.StatusBadRequest)
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
		r.Form.Set("username", user.ID)

	}

	err := oauthServer.HandleTokenRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}

}

// @Summary 获取Captcha
// @Description 获取Captcha
// @Tags Oauth2
// @Produce  json
// @Success 200 {object} common.Response{data=string} "token信息"
// @Failure 500 {object} string "错误code和错误信息"
// @Router /captcha-gen [get]
func captchaGen(w http.ResponseWriter, r *http.Request) {
	d := struct {
		CaptchaId string
	}{
		captcha.New(),
	}
	ret := common.OK.WithData(d.CaptchaId)
	bytes, _ := json.Marshal(ret)
	w.Write(bytes)
	return
}

// @Summary 获取token
// @Description 获取token
// @Tags Oauth2
// @Param username formData string false "username"
// @Param password formData string false "password"
// @Param grant_type formData string false "grant_type"
// @Param scope formData string false "scope"
// @Param client_id formData string false "client_id"
// @Param client_secret formData string false "client_secret"
// @Param verify_key formData string false "verify_key"
// @Param verify_value formData string false "verify_value"
// @Param sms_code formData string false "sms_code"
// @Param is_travel formData string false "is_travel"
// @Produce  json
// @Success 200 {object} model.TokenInfo "token信息"
// @Failure 500 {object} string "错误code和错误信息"
// @Router /oauth/token-by-field [post]
func tokenByFieldHandler(w http.ResponseWriter, r *http.Request) {
	if dumpvar {
		_ = dumpRequest(os.Stdout, "token", r) // Ignore the error
	}
	grantType := r.FormValue("grant_type")
	if grantType != "refresh_token" {
		username := r.FormValue("username")
		if username == "" {
			http.Error(w, "用户名不能为空", http.StatusBadRequest)
			return
		}
		var field string
		field = "name"

		value := username
		isTravelStr := r.FormValue("is_travel")
		isTravel := false
		if isTravelStr == "1" {
			isTravel = true
		}
		user, err := service.GetUserByFieldName(r.Context(), field, value, isTravel)
		if err != nil {
			http.Error(w, "系统错误:"+err.Error(), http.StatusInternalServerError)
			return
		}
		if user == nil {
			http.Error(w, "用户不存在", 498)
			return
		}

		if user.Status != 1 {
			http.Error(w, "用户已停用", 497)
			return
		}
		sms_code := r.FormValue("sms_code")
		if sms_code != "" { //如果是验证码登录，那么就先校验验证码，成功后，获取密码，后面走oauth流程
			valid, err := service.CheckMobileSmsCode(r.Context(), value, sms_code)
			if err != nil {
				http.Error(w, err.Error(), 499)
				return
			}
			if !valid {
				http.Error(w, "短信验证码错误", 499)
				return
			}

		}

		if VERIFY_CAPTCHA && !isTravel && sms_code == "" {
			vKey := r.FormValue("verify_key")
			vVal := r.FormValue("verify_value")
			if vKey == "" || vVal == "" {
				http.Error(w, "验证码为空", http.StatusBadRequest)
				return
			}
			if !captcha.VerifyString(vKey, vVal) {
				http.Error(w, "验证码错误", http.StatusNotAcceptable)
				return
			}

		}

		passwd, err := service.GetUserPasswordByField(r.Context(), field, value, isTravel)
		if err != nil {
			http.Error(w, "获取用户错误 "+err.Error(), http.StatusInternalServerError)
			return
		}

		if passwd != r.FormValue("password") {
			http.Error(w, "密码错误", 496)
			return
		}
		r.Form.Set("password", passwd)
		r.Form.Set("username", user.ID)
	}
	err := oauthServer.HandleTokenRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	if dumpvar {
		_ = dumpRequest(os.Stdout, "test", r) // Ignore the error
	}
	token, err := oauthServer.ValidationBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := map[string]interface{}{
		"expires_in": int64(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
		"client_id":  token.GetClientID(),
		"user_id":    token.GetUserID(),
	}
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	e.Encode(data)
}

func testUsePostHandler(w http.ResponseWriter, r *http.Request) {
	if dumpvar {
		_ = dumpRequest(os.Stdout, "test_use_post", r) // Ignore the error
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		common.Logger.Error("read body error ." + err.Error())
		http.Error(w, "read body error ."+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	s := string(body)

	params := map[string]string{}
	arr := strings.Split(s, "&")
	for _, item := range arr {
		itemarr := strings.Split(item, "=")
		if len(itemarr) == 2 {
			params[itemarr[0]] = itemarr[1]
		}
	}
	access_token, exists := params["token"]
	if exists {

		token, err := oauthServer.Manager.LoadAccessToken(r.Context(), access_token)
		if err != nil {
			common.Logger.Error(params)
			common.Logger.Error("LoadAccessToken error . " + err.Error())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		log.Printf("client_id=%s user_id=%s expires_in=%d", token.GetClientID(), token.GetUserID(), int64(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()).Sub(time.Now()).Seconds()))
		//将权限信息加到这里？ 通过X-Userinfo 进一步处理？
		data := map[string]interface{}{
			"expires_in": int64(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
			"client_id":  token.GetClientID(),
			"user_id":    token.GetUserID(),
			"active":     true,
		}

		e := json.NewEncoder(w)
		w.WriteHeader(http.StatusOK)
		e.SetIndent("", "  ")
		e.Encode(data)
	} else {
		common.Logger.Error("can't find token in request body")
		http.Error(w, "can't find token in request body", http.StatusBadRequest)
		return
	}

}

func initTestClientStoreHandler(w http.ResponseWriter, r *http.Request) {
	clientStore.Set("dapr-client", &models.Client{
		ID:     "dapr-client",
		Secret: "123456",
	})

	w.Write([]byte("ok"))
	return

}

func refreshClientFromDbHandler(w http.ResponseWriter, r *http.Request) {

	err := refreshClientInfoFromDb(clientStore)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("ok"))
	return

}
