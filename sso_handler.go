package main

import (
	"net/http"
	"oauth2-server/config"
	"oauth2-server/model"
	"oauth2-server/service"

	"github.com/dapr-platform/common"
	"github.com/go-oauth2/oauth2/v4"
)

type SSOTokenRequest struct {
	Ticket string `json:"ticket"`
}

// @Summary SSO登录
// @Description 通过中台SSO ticket换取本地OAuth2 token
// @Tags SSO
// @Accept json
// @Produce json
// @Param data body SSOTokenRequest true "SSO ticket"
// @Success 200 {object} common.Response "token信息"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 500 {object} common.Response "系统错误"
// @Router /sso/token [post]
func ssoTokenHandler(w http.ResponseWriter, r *http.Request) {
	if !config.SSO_ENABLED {
		common.HttpResult(w, common.ErrParam.AppendMsg("SSO功能未启用"))
		return
	}

	var req SSOTokenRequest
	if err := common.ReadRequestBody(r, &req); err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg("请求参数错误"))
		return
	}
	if req.Ticket == "" {
		common.HttpResult(w, common.ErrParam.AppendMsg("ticket不能为空"))
		return
	}

	ssoUser, err := service.SSORestoreTicket(req.Ticket)
	if err != nil {
		common.Logger.Error("SSO ticket验证失败: " + err.Error())
		common.HttpResult(w, common.ErrService.AppendMsg("SSO ticket验证失败: "+err.Error()))
		return
	}

	user, err := service.GetUserByFieldName(r.Context(), "identity", ssoUser.Code, false)
	if err != nil {
		common.Logger.Error("查询本地用户失败: " + err.Error())
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}
	if user == nil {
		common.Logger.Error("SSO用户在本地不存在, code: " + ssoUser.Code)
		common.HttpResult(w, common.ErrParam.AppendMsg("用户不存在，请先同步用户数据"))
		return
	}
	if user.Status != 1 {
		common.HttpResult(w, common.ErrParam.AppendMsg("用户已停用"))
		return
	}

	gt := oauth2.GrantType("password")
	tgr := &oauth2.TokenGenerateRequest{
		ClientID:     "dapr-client",
		ClientSecret: "123456",
		UserID:       user.ID,
		Request:      r,
	}
	ti, err := oauthServer.GetAccessToken(r.Context(), gt, tgr)
	if err != nil {
		common.Logger.Error("生成token失败: " + err.Error())
		common.HttpResult(w, common.ErrService.AppendMsg("生成token失败"))
		return
	}

	common.HttpSuccess(w, common.OK.WithData(oauthServer.GetTokenData(ti)))
}

// @Summary SSO登出
// @Description 登出本地会话并通知中台登出
// @Tags SSO
// @Produce json
// @Success 200 {object} common.Response "登出成功"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 500 {object} common.Response "系统错误"
// @Router /sso/logout [post]
func ssoLogoutHandler(w http.ResponseWriter, r *http.Request) {
	if !config.SSO_ENABLED {
		common.HttpResult(w, common.ErrParam.AppendMsg("SSO功能未启用"))
		return
	}

	sub, err := common.ExtractUserSub(r)
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg("无法获取当前用户: "+err.Error()))
		return
	}

	user, err := service.GetUserByFieldName(r.Context(), "id", sub, false)
	if err != nil || user == nil {
		common.HttpResult(w, common.ErrParam.AppendMsg("用户不存在"))
		return
	}

	if err := service.SSORevokeByCode(user.Identity); err != nil {
		common.Logger.Error("SSO登出失败: " + err.Error())
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}

	common.HttpSuccess(w, common.OK)
}

// @Summary SSO用户同步
// @Description 从中台全量同步用户到本地
// @Tags SSO
// @Produce json
// @Success 200 {object} common.Response "同步结果"
// @Failure 400 {object} common.Response "参数错误"
// @Failure 500 {object} common.Response "系统错误"
// @Router /sso/sync-members [post]
func ssoSyncMembersHandler(w http.ResponseWriter, r *http.Request) {
	if !config.SSO_ENABLED {
		common.HttpResult(w, common.ErrParam.AppendMsg("SSO功能未启用"))
		return
	}

	syncCount, err := service.SSOSyncMembers(r.Context())
	if err != nil {
		common.Logger.Error("SSO用户同步失败: " + err.Error())
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}

	common.HttpSuccess(w, common.OK.WithData(map[string]interface{}{
		"synced_count": syncCount,
	}))
}

// SSOUserInfo 不含密码的用户信息（调试用）
type SSOUserInfo struct {
	ID         string           `json:"id"`
	Name       string           `json:"name"`
	Identity   string           `json:"identity"`
	Mobile     string           `json:"mobile"`
	Email      string           `json:"email"`
	Gender     int32            `json:"gender"`
	OrgID      string           `json:"org_id"`
	WorkNumber string           `json:"work_number"`
	Status     int32            `json:"status"`
	Type       int32            `json:"type"`
	CreateAt   common.LocalTime `json:"create_at"`
	UpdateAt   common.LocalTime `json:"update_at"`
}

// @Summary SSO状态
// @Description 查看SSO配置和同步状态
// @Tags SSO
// @Produce json
// @Success 200 {object} common.Response "SSO状态"
// @Router /sso/status [get]
func ssoStatusHandler(w http.ResponseWriter, r *http.Request) {
	common.HttpSuccess(w, common.OK.WithData(map[string]interface{}{
		"sso_enabled":  config.SSO_ENABLED,
		"sso_base_url": config.SSO_BASE_URL,
		"sso_app_key":  config.SSO_APP_KEY,
	}))
}

// @Summary SSO同步用户列表
// @Description 查看通过SSO同步到本地的用户列表（identity不为空的用户）
// @Tags SSO
// @Produce json
// @Success 200 {object} common.Response "用户列表"
// @Failure 500 {object} common.Response "系统错误"
// @Router /sso/users [get]
func ssoUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := common.DbQuery[model.User](r.Context(), common.GetDaprClient(),
		model.UserTableInfo.Name, "identity_ne=&_select=id,name,identity,mobile,email,gender,org_id,work_number,status,type,create_at,update_at")
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}

	result := make([]SSOUserInfo, 0, len(users))
	for _, u := range users {
		result = append(result, SSOUserInfo{
			ID: u.ID, Name: u.Name, Identity: u.Identity,
			Mobile: u.Mobile, Email: u.Email, Gender: u.Gender,
			OrgID: u.OrgID, WorkNumber: u.WorkNumber,
			Status: u.Status, Type: u.Type,
			CreateAt: u.CreateAt, UpdateAt: u.UpdateAt,
		})
	}

	common.HttpSuccess(w, common.OK.WithData(map[string]interface{}{
		"total": len(result),
		"users": result,
	}))
}
