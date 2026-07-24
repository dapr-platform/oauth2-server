package service

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"oauth2-server/config"
	"oauth2-server/model"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dapr-platform/common"
	"github.com/pkg/errors"
)

// SSORestoreRequest 票据验证请求
type SSORestoreRequest struct {
	Ticket string `json:"ticket"`
}

// SSORestoreResponse 票据验证返回
type SSORestoreResponse struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    *struct {
		Content *SSOUserContent `json:"content"`
	} `json:"data"`
}

// SSOUserContent 中台返回的用户信息
type SSOUserContent struct {
	LoginName        string `json:"loginName"`
	Code             string `json:"code"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	OrgID            string `json:"orgId"`
	OrgName          string `json:"orgName"`
	OutsideOrgID     string `json:"outsideOrgId"`
	OutsideOrgName   string `json:"outsideOrgName"`
	SocialCreditCode string `json:"socialCreditCode"`
}

// SSORevokeRequest 登出请求
type SSORevokeRequest struct {
	Code string `json:"code"`
}

// SSORevokeResponse 登出返回
type SSORevokeResponse struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// SSOSyncRequest 人员全量同步请求
type SSOSyncRequest struct {
	RequestID string         `json:"requestId"`
	Timestamp int64          `json:"timestamp"`
	NotifyURL string         `json:"notifyUrl"`
	Params    *SSOSyncParams `json:"params"`
	Sort      *SSOSort       `json:"sort"`
	PageInfo  *SSOPageInfo   `json:"pageInfo"`
}

type SSOSyncParams struct {
	Code         string `json:"code,omitempty"`
	IncludeChild bool   `json:"includeChild"`
	StartTime    string `json:"startTime,omitempty"`
	EndTime      string `json:"endTime,omitempty"`
}

type SSOSort struct {
	Orders []SSOSortOrder `json:"orders"`
}

type SSOSortOrder struct {
	Property  string `json:"property"`
	Direction string `json:"direction"`
}

type SSOPageInfo struct {
	PageNumber int  `json:"pageNumber"`
	PageSize   int  `json:"pageSize"`
	Total      int  `json:"total"`
	Pages      int  `json:"pages"`
	NeedTotal  bool `json:"needTotal"`
}

// SSOSyncResponse 人员同步返回
type SSOSyncResponse struct {
	Status  interface{} `json:"status"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    *struct {
		PageInfo *SSOPageInfo    `json:"pageInfo"`
		Content  []SSOSyncMember `json:"content"`
	} `json:"data"`
}

type SSOSyncMember struct {
	ThirdID     string              `json:"thirdId"`
	Name        string              `json:"name"`
	Code        string              `json:"code"`
	Username    string              `json:"username"`
	Gender      string              `json:"gender"`
	PhoneNumber string              `json:"phoneNumber"`
	Email       string              `json:"email"`
	IsEnable    string              `json:"isEnable"`
	MemberType  string              `json:"memberType"`
	MemberPosts []SSOSyncMemberPost `json:"memberPosts"`
}

type SSOSyncMemberPost struct {
	Main     string `json:"main"`
	UnitCode string `json:"unitCode"`
}

// StartSSOSyncScheduler 启动每日定时同步任务
func StartSSOSyncScheduler() {
	syncHour := 2 // 默认凌晨2点同步
	if val := os.Getenv("SSO_SYNC_HOUR"); val != "" {
		if h, err := strconv.Atoi(val); err == nil && h >= 0 && h < 24 {
			syncHour = h
		}
	}

	go func() {
		for {
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day(), syncHour, 0, 0, 0, now.Location())
			if now.After(next) {
				next = next.Add(24 * time.Hour)
			}
			waitDuration := next.Sub(now)
			common.Logger.Infof("SSO定时同步: 下次执行时间 %s (等待 %v)", next.Format("2006-01-02 15:04:05"), waitDuration)

			timer := time.NewTimer(waitDuration)
			<-timer.C

			common.Logger.Info("SSO定时同步: 开始执行")
			ctx := context.Background()
			count, _, err := SSOSyncMembers(ctx)
			if err != nil {
				common.Logger.Error("SSO定时同步失败: ", err)
			} else {
				common.Logger.Infof("SSO定时同步完成, 同步 %d 个用户", count)
			}
		}
	}()
}

// ssoSign 计算 MD5 签名: MD5(AppSecret + body + AppSecret)
func ssoSign(body []byte) string {
	raw := config.SSO_APP_SECRET + string(body) + config.SSO_APP_SECRET
	hash := md5.Sum([]byte(raw))
	return fmt.Sprintf("%x", hash[:])
}

// ssoDoRequest 执行带签名的 SSO 请求
func ssoDoRequest(url string, reqBody interface{}) ([]byte, error) {
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, errors.Wrap(err, "json marshal")
	}
	common.Logger.Infof("SSO请求: %s", string(bodyBytes))
	sign := ssoSign(bodyBytes)

	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, errors.Wrap(err, "new request")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("app-key", config.SSO_APP_KEY)
	req.Header.Set("sign-type", "MD5")
	req.Header.Set("sign", sign)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "http do")
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read response")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("SSO请求失败, HTTP %d: %s", resp.StatusCode, string(respBody))
	}
	return respBody, nil
}

// SSORestoreTicket 调用中台验证 ticket，返回用户信息
func SSORestoreTicket(ticket string) (*SSOUserContent, error) {
	url := strings.TrimRight(config.SSO_TICKET_BASE_URL, "/") + "/service/ctp-user/auth/restore"
	reqBody := &SSORestoreRequest{Ticket: ticket}

	respBody, err := ssoDoRequest(url, reqBody)
	if err != nil {
		return nil, errors.Wrap(err, "调用SSO ticket验证接口")
	}

	var result SSORestoreResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, errors.Wrap(err, "解析SSO返回")
	}

	if result.Code != "BOOT_0000" || result.Data == nil || result.Data.Content == nil {
		return nil, fmt.Errorf("SSO ticket验证失败: code=%s, message=%s", result.Code, result.Message)
	}
	return result.Data.Content, nil
}

// SSORevokeByCode 调用中台登出接口
func SSORevokeByCode(code string) error {
	url := strings.TrimRight(config.SSO_TICKET_BASE_URL, "/") + "/service/ctp-user/auth/token/revoke-by-code"
	reqBody := &SSORevokeRequest{Code: code}

	respBody, err := ssoDoRequest(url, reqBody)
	if err != nil {
		return errors.Wrap(err, "调用SSO登出接口")
	}

	var result SSORevokeResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return errors.Wrap(err, "解析SSO登出返回")
	}

	if result.Code != "BOOT_0000" {
		return fmt.Errorf("SSO登出失败: code=%s, message=%s", result.Code, result.Message)
	}
	return nil
}

// SSOSyncMembers 从中台全量同步人员到本地 o_user 表
func SSOSyncMembers(ctx context.Context) (syncCount int, reqBody string, err error) {
	pageNumber := 1
	pageSize := 100
	totalSynced := 0

	for {
		members, hasMore, reqBodyJson, err := ssoFetchMembersPage(pageNumber, pageSize)
		if err != nil {
			return totalSynced, reqBodyJson, errors.Wrapf(err, "获取第%d页人员数据", pageNumber)
		}
		reqBody = reqBodyJson

		for _, member := range members {
			if err := ssoUpsertLocalUser(ctx, &member); err != nil {
				common.Logger.Error("同步用户失败: ", member.Code, " error: ", err)
				continue
			}
			memberJson, _ := json.Marshal(member)
			common.Logger.Infof("同步用户成功: %s", string(memberJson))
			totalSynced++
		}

		if !hasMore {
			break
		}
		pageNumber++
	}

	common.Logger.Infof("SSO用户同步完成, 共同步 %d 个用户", totalSynced)
	return totalSynced, reqBody, nil
}

// ssoFetchMembersPage 分页获取中台人员数据
func ssoFetchMembersPage(pageNumber, pageSize int) ([]SSOSyncMember, bool, string, error) {
	url := strings.TrimRight(config.SSO_BASE_URL, "/") + "/organization/unit/members"
	endTimeStr := time.Now().Format("2006-01-02")
	reqBody := &SSOSyncRequest{
		RequestID: fmt.Sprintf("%d", time.Now().UnixMilli()),
		Timestamp: time.Now().UnixMilli(),
		NotifyURL: "",
		Params: &SSOSyncParams{
			Code:         "group",
			StartTime:    "2023-04-17",
			EndTime:      endTimeStr,
			IncludeChild: true,
		},
		Sort: &SSOSort{
			Orders: []SSOSortOrder{{Property: "createTime", Direction: "ASC"}},
		},
		PageInfo: &SSOPageInfo{
			PageNumber: pageNumber,
			PageSize:   pageSize,
			NeedTotal:  true,
		},
	}

	respBody, err := ssoDoRequest(url, reqBody)
	if err != nil {
		return nil, false, "", errors.Wrap(err, "请求中台人员列表")
	}

	var result SSOSyncResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, false, "", errors.Wrap(err, "解析人员列表返回")
	}

	if result.Data == nil {
		return nil, false, "", fmt.Errorf("中台人员列表返回为空: code=%s, message=%s", result.Code, result.Message)
	}

	hasMore := false
	if result.Data.PageInfo != nil && result.Data.PageInfo.Pages > pageNumber {
		hasMore = true
	}
	reqBodyJson, _ := json.Marshal(reqBody)
	return result.Data.Content, hasMore, string(reqBodyJson), nil
}

// ssoUpsertLocalUser 将中台人员数据写入/更新到本地 o_user 表
// 中台 code（员工编码）作为本地 identity 的唯一映射
func ssoUpsertLocalUser(ctx context.Context, member *SSOSyncMember) error {
	if member.Code == "" {
		return fmt.Errorf("中台人员编号(code)为空，跳过")
	}

	existingUser, _ := GetUserByFieldName(ctx, "identity", member.Code, false)

	now := common.LocalTime(time.Now())
	gender := ssoMapGender(member.Gender)
	status := int32(1)
	if member.IsEnable != "true" {
		status = 2
	}

	orgID := ""
	if len(member.MemberPosts) > 0 {
		for _, post := range member.MemberPosts {
			if post.Main == "true" {
				orgID = post.UnitCode
				break
			}
		}
		if orgID == "" {
			orgID = member.MemberPosts[0].UnitCode
		}
	}

	if existingUser != nil {
		info := make(map[string]any)
		info[model.User_FIELD_NAME_id] = existingUser.ID
		info[model.User_FIELD_NAME_name] = member.Name
		info[model.User_FIELD_NAME_zh_name] = member.Name
		info[model.User_FIELD_NAME_mobile] = member.PhoneNumber
		info[model.User_FIELD_NAME_email] = member.Email
		info[model.User_FIELD_NAME_gender] = gender
		info[model.User_FIELD_NAME_org_id] = orgID
		info[model.User_FIELD_NAME_work_number] = member.Username
		info[model.User_FIELD_NAME_status] = status
		info[model.User_FIELD_NAME_update_at] = now
		return common.DbUpsert[map[string]any](ctx, common.GetDaprClient(), info, model.UserTableInfo.Name, model.User_FIELD_NAME_id)
	}

	newUser := model.User{
		ID:         common.NanoId(),
		TenantID:   "default",
		Identity:   member.Code,
		Name:       member.Name,
		ZhName:     member.Name,
		Mobile:     member.PhoneNumber,
		Email:      member.Email,
		Gender:     gender,
		Password:   common.NanoId(),
		Type:       2,
		OrgID:      orgID,
		WorkNumber: member.Username,
		Status:     status,
		CreateAt:   now,
		UpdateAt:   now,
	}
	return common.DbUpsert[model.User](ctx, common.GetDaprClient(), newUser, model.UserTableInfo.Name, model.User_FIELD_NAME_id)
}

func ssoMapGender(gender string) int32 {
	switch gender {
	case "MALE":
		return 1
	case "FEMALE":
		return 2
	default:
		return 0
	}
}
