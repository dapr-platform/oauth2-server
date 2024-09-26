package dapr

import (
	"context"
	"encoding/json"
	"github.com/dapr-platform/common"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/google/uuid"
	"strconv"

	"time"

	dapr "github.com/dapr/go-sdk/client"
)

type TokenStore struct {
}

func NewDaprTokenStore() (oauth2.TokenStore, error) {

	return &TokenStore{}, nil
}

// Create create and store the new token information
func (ts *TokenStore) Create(ctx context.Context, info oauth2.TokenInfo) error {
	ct := time.Now()
	jv, err := json.Marshal(info)
	if err != nil {
		return err
	}
	if code := info.GetCode(); code != "" {
		return ts.saveState(ctx, code, jv, true, info.GetCodeExpiresIn())
	}
	basicID := uuid.Must(uuid.NewRandom()).String()
	aexp := info.GetAccessExpiresIn()
	rexp := aexp
	expires := true
	if refresh := info.GetRefresh(); refresh != "" {
		rexp = info.GetRefreshCreateAt().Add(info.GetRefreshExpiresIn()).Sub(ct)
		if aexp.Seconds() > rexp.Seconds() {
			aexp = rexp
		}
		expires = info.GetRefreshExpiresIn() != 0
		err := ts.saveState(ctx, refresh, []byte(basicID), expires, rexp)
		if err != nil {
			return err
		}

	}
	err = ts.saveState(ctx, basicID, jv, expires, rexp)
	if err != nil {
		return err
	}
	return ts.saveState(ctx, info.GetAccess(), []byte(basicID), expires, aexp)
}

func (ts *TokenStore) saveState(ctx context.Context, key string, value []byte, expires bool, ttl time.Duration) error {
	ttlstr := "0"
	if expires {
		ttlstr = strconv.FormatFloat(ttl.Seconds(), 'f', 0, 64)
	}
	item := &dapr.SetStateItem{
		Key: key,
		Etag: &dapr.ETag{
			Value: "1",
		},
		Metadata: map[string]string{
			"ttlInSeconds": ttlstr,
		},
		Value: value,
		Options: &dapr.StateOptions{
			Concurrency: dapr.StateConcurrencyLastWrite,
			Consistency: dapr.StateConsistencyStrong,
		},
	}

	return common.GetDaprClient().SaveBulkState(ctx, common.DAPR_STATESTORE_NAME, item)
}

// RemoveByCode use the authorization code to delete the token information
func (ts *TokenStore) RemoveByCode(ctx context.Context, code string) error {

	return common.GetDaprClient().DeleteState(ctx, common.DAPR_STATESTORE_NAME, code, nil)
}

// RemoveByAccess use the access token to delete the token information
func (ts *TokenStore) RemoveByAccess(ctx context.Context, access string) error {

	return common.GetDaprClient().DeleteState(ctx, common.DAPR_STATESTORE_NAME, access, nil)

}

// RemoveByRefresh use the refresh token to delete the token information
func (ts *TokenStore) RemoveByRefresh(ctx context.Context, refresh string) error {

	return common.GetDaprClient().DeleteState(ctx, common.DAPR_STATESTORE_NAME, refresh, nil)
}

func (ts *TokenStore) getData(ctx context.Context, key string) (oauth2.TokenInfo, error) {

	var ti oauth2.TokenInfo
	item, err := common.GetDaprClient().GetState(ctx, common.DAPR_STATESTORE_NAME, key, nil)
	if err != nil {
		return nil, err
	}
	var tm models.Token
	err = json.Unmarshal(item.Value, &tm)
	if err != nil {
		return nil, err
	}
	ti = &tm
	return ti, nil

}

func (ts *TokenStore) getBasicID(ctx context.Context, key string) (string, error) {

	var basicID string
	item, err := common.GetDaprClient().GetState(ctx, common.DAPR_STATESTORE_NAME, key, nil)
	if err != nil {
		return "", err
	}
	basicID = string(item.Value)
	return basicID, nil

}

// GetByCode use the authorization code for token information data
func (ts *TokenStore) GetByCode(ctx context.Context, code string) (oauth2.TokenInfo, error) {
	return ts.getData(ctx, code)
}

// GetByAccess use the access token for token information data
func (ts *TokenStore) GetByAccess(ctx context.Context, access string) (oauth2.TokenInfo, error) {
	basicID, err := ts.getBasicID(ctx, access)
	if err != nil {
		return nil, err
	}
	return ts.getData(ctx, basicID)
}

// GetByRefresh use the refresh token for token information data
func (ts *TokenStore) GetByRefresh(ctx context.Context, refresh string) (oauth2.TokenInfo, error) {
	basicID, err := ts.getBasicID(ctx, refresh)
	if err != nil {
		return nil, err
	}
	return ts.getData(ctx, basicID)
}
