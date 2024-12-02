package api

import (
	"github.com/dapr-platform/common"
	"github.com/go-chi/chi/v5"
	"net/http"
	"oauth2-server/model"

	"strings"
)

func InitClient_infoRoute(r chi.Router) {

	r.Get(common.BASE_CONTEXT+"/client-info/page", Client_infoPageListHandler)
	r.Get(common.BASE_CONTEXT+"/client-info", Client_infoListHandler)

	r.Post(common.BASE_CONTEXT+"/client-info", UpsertClient_infoHandler)

	r.Delete(common.BASE_CONTEXT+"/client-info/{id}", DeleteClient_infoHandler)

	r.Post(common.BASE_CONTEXT+"/client-info/batch-delete", batchDeleteClient_infoHandler)

	r.Post(common.BASE_CONTEXT+"/client-info/batch-upsert", batchUpsertClient_infoHandler)

	r.Get(common.BASE_CONTEXT+"/client-info/groupby", Client_infoGroupbyHandler)

}

// @Summary GroupBy
// @Description GroupBy, for example,  _select=level, then return  {level_val1:sum1,level_val2:sum2}, _where can input status=0
// @Tags Client_info
// @Param _select query string true "_select"
// @Param _where query string false "_where"
// @Produce  json
// @Success 200 {object} common.Response{data=[]map[string]any} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /client-info/groupby [get]
func Client_infoGroupbyHandler(w http.ResponseWriter, r *http.Request) {

	common.CommonGroupby(w, r, common.GetDaprClient(), "o_client_info")
}

// @Summary batch update
// @Description batch update
// @Tags Client_info
// @Accept  json
// @Param entities body []map[string]any true "objects array"
// @Produce  json
// @Success 200 {object} common.Response ""
// @Failure 500 {object} common.Response ""
// @Router /client-info/batch-upsert [post]
func batchUpsertClient_infoHandler(w http.ResponseWriter, r *http.Request) {

	var entities []model.Client_info
	err := common.ReadRequestBody(r, &entities)
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg(err.Error()))
		return
	}
	if len(entities) == 0 {
		common.HttpResult(w, common.ErrParam.AppendMsg("len of entities is 0"))
		return
	}

	beforeHook, exists := common.GetUpsertBeforeHook("Client_info")
	if exists {
		for _, v := range entities {
			_, err1 := beforeHook(r, v)
			if err1 != nil {
				common.HttpResult(w, common.ErrService.AppendMsg(err1.Error()))
				return
			}
		}

	}
	for _, v := range entities {
		if v.ID == "" {
			v.ID = common.NanoId()
		}
	}

	err = common.DbBatchUpsert[model.Client_info](r.Context(), common.GetDaprClient(), entities, model.Client_infoTableInfo.Name, model.Client_info_FIELD_NAME_id)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}

	common.HttpResult(w, common.OK)
}

// @Summary page query
// @Description page query, _page(from 1 begin), _page_size, _order, and others fields, status=1, name=$like.%CAM%
// @Tags Client_info
// @Param _page query int true "current page"
// @Param _page_size query int true "page size"
// @Param _order query string false "order"
// @Param id query string false "id"
// @Param password query string false "password"
// @Produce  json
// @Success 200 {object} common.Response{data=common.Page{items=[]model.Client_info}} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /client-info/page [get]
func Client_infoPageListHandler(w http.ResponseWriter, r *http.Request) {

	page := r.URL.Query().Get("_page")
	pageSize := r.URL.Query().Get("_page_size")
	if page == "" || pageSize == "" {
		common.HttpResult(w, common.ErrParam.AppendMsg("page or pageSize is empty"))
		return
	}
	common.CommonPageQuery[model.Client_info](w, r, common.GetDaprClient(), "o_client_info", "id")

}

// @Summary query objects
// @Description query objects
// @Tags Client_info
// @Param _select query string false "_select"
// @Param _order query string false "order"
// @Param id query string false "id"
// @Param password query string false "password"
// @Produce  json
// @Success 200 {object} common.Response{data=[]model.Client_info} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /client-info [get]
func Client_infoListHandler(w http.ResponseWriter, r *http.Request) {
	common.CommonQuery[model.Client_info](w, r, common.GetDaprClient(), "o_client_info", "id")
}

// @Summary save
// @Description save
// @Tags Client_info
// @Accept       json
// @Param item body model.Client_info true "object"
// @Produce  json
// @Success 200 {object} common.Response{data=model.Client_info} "object"
// @Failure 500 {object} common.Response ""
// @Router /client-info [post]
func UpsertClient_infoHandler(w http.ResponseWriter, r *http.Request) {
	var val model.Client_info
	err := common.ReadRequestBody(r, &val)
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg(err.Error()))
		return
	}

	beforeHook, exists := common.GetUpsertBeforeHook("Client_info")
	if exists {
		v, err1 := beforeHook(r, val)
		if err1 != nil {
			common.HttpResult(w, common.ErrService.AppendMsg(err1.Error()))
			return
		}
		val = v.(model.Client_info)
	}
	if val.ID == "" {
		val.ID = common.NanoId()
	}
	err = common.DbUpsert[model.Client_info](r.Context(), common.GetDaprClient(), val, model.Client_infoTableInfo.Name, "id")
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}
	common.HttpSuccess(w, common.OK.WithData(val))
}

// @Summary delete
// @Description delete
// @Tags Client_info
// @Param id  path string true "实例id"
// @Produce  json
// @Success 200 {object} common.Response{data=model.Client_info} "object"
// @Failure 500 {object} common.Response ""
// @Router /client-info/{id} [delete]
func DeleteClient_infoHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	beforeHook, exists := common.GetDeleteBeforeHook("Client_info")
	if exists {
		_, err1 := beforeHook(r, id)
		if err1 != nil {
			common.HttpResult(w, common.ErrService.AppendMsg(err1.Error()))
			return
		}
	}
	common.CommonDelete(w, r, common.GetDaprClient(), "o_client_info", "id", "id")
}

// @Summary batch delete
// @Description batch delete
// @Tags Client_info
// @Accept  json
// @Param ids body []string true "id array"
// @Produce  json
// @Success 200 {object} common.Response ""
// @Failure 500 {object} common.Response ""
// @Router /client-info/batch-delete [post]
func batchDeleteClient_infoHandler(w http.ResponseWriter, r *http.Request) {

	var ids []string
	err := common.ReadRequestBody(r, &ids)
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg(err.Error()))
		return
	}
	if len(ids) == 0 {
		common.HttpResult(w, common.ErrParam.AppendMsg("len of ids is 0"))
		return
	}
	beforeHook, exists := common.GetBatchDeleteBeforeHook("Client_info")
	if exists {
		_, err1 := beforeHook(r, ids)
		if err1 != nil {
			common.HttpResult(w, common.ErrService.AppendMsg(err1.Error()))
			return
		}
	}
	idstr := strings.Join(ids, ",")
	err = common.DbDeleteByOps(r.Context(), common.GetDaprClient(), "o_client_info", []string{"id"}, []string{"in"}, []any{idstr})
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}

	common.HttpResult(w, common.OK)
}
