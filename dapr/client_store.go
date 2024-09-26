package dapr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dapr-platform/common"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
	"log"
	"oauth2-server/config"
	"oauth2-server/model"
)

type ClientStore struct {
}

func NewClientStore() (*ClientStore, error) {

	return &ClientStore{}, nil
}
func (cs *ClientStore) refreshClientInfoFromDb(ctx context.Context, id string) error {
	data, err := common.GetDaprClient().InvokeMethod(ctx, "db-service", "/"+config.DBNAME+"/public/"+config.CLIENT_INFO_TABLE_NAME+"?id="+id, "get")
	if err != nil {
		log.Printf("refreshClientInfoFromDb error.%s", err.Error())
		return err
	}
	var infolist []model.Client_info
	err = json.Unmarshal(data, &infolist)
	if err != nil {
		log.Printf("refreshClientInfoFromDb Unmarshal data error.%s", err.Error())
		return err
	}
	for _, info := range infolist {
		err = cs.Set(info.ID, &models.Client{
			ID:     info.ID,
			Secret: info.Password,
		})
		if err != nil {
			log.Printf("clientStore.Set error.%s", err.Error())
			return err
		}
	}
	return nil

}

// GetByID according to the ID for the alpr-client information
func (cs *ClientStore) GetByID(ctx context.Context, id string) (oauth2.ClientInfo, error) {
	refreshed := false
	for {
		item, err := common.GetDaprClient().GetState(ctx, common.DAPR_STATESTORE_NAME, id, nil)
		if err != nil {
			fmt.Println("clientStore GetState error:", err)
			return nil, err
		}
		if len(item.Value) == 0 {
			fmt.Println("clientStore GetByID not found. id=" + id)
			if refreshed {
				return nil, errors.New("not found, refresh error")
			}
			err = cs.refreshClientInfoFromDb(ctx, id)
			if err != nil {
				return nil, errors.New("not found, refresh error, " + err.Error())
			}
			refreshed = true

		} else {
			var cli oauth2.ClientInfo
			client := models.Client{}
			err = json.Unmarshal(item.Value, &client)
			if err != nil {
				fmt.Println("clientStore GetByID unmarshal value error. ")
				return nil, err
			}
			cli = &client
			return cli, nil
		}
	}

}

// Set set alpr-client information
func (cs *ClientStore) Set(id string, cli oauth2.ClientInfo) (err error) {
	jv, err := json.Marshal(cli)
	if err != nil {
		fmt.Println("marshal clientInfo error.id=" + id)
		return err
	}

	return common.GetDaprClient().SaveState(context.Background(), common.DAPR_STATESTORE_NAME, id, jv, nil)
}
