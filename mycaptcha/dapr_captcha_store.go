package mycaptcha

import (
	"context"
	"github.com/dapr-platform/common"
	daprc "github.com/dapr/go-sdk/client"
	"oauth2-server/dapr"
	"strconv"
	"time"
)

type DaprCaptchaStore struct {
	Expiration time.Duration
}

func (impl *DaprCaptchaStore) Set(id string, digits []byte) {
	ttlstr := strconv.FormatFloat(float64(impl.Expiration), 'f', 0, 64)

	item := &daprc.SetStateItem{
		Key: id,
		Etag: &daprc.ETag{
			Value: "1",
		},
		Metadata: map[string]string{
			"ttlInSeconds": ttlstr,
		},
		Value: digits,
		Options: &daprc.StateOptions{
			Concurrency: daprc.StateConcurrencyLastWrite,
			Consistency: daprc.StateConsistencyStrong,
		},
	}

	err := common.GetDaprClient().SaveBulkState(context.Background(), dapr.StateStoreName, item)
	if err != nil {
		common.Logger.Error("save mycaptcha error", err)
		panic(err)
	}
}

func (impl *DaprCaptchaStore) Get(id string, clear bool) (digits []byte) {
	item, err := common.GetDaprClient().GetState(context.Background(), dapr.StateStoreName, id, nil)
	if err != nil {
		common.Logger.Error("mycaptcha GetState error:", err)
		panic(err)
	}
	if clear {
		err = common.GetDaprClient().DeleteState(context.Background(), dapr.StateStoreName, id, make(map[string]string))
		if err != nil {
			common.Logger.Error("mycaptcha delete error:", err)
			panic(err)
		}
	}
	return item.Value
}
