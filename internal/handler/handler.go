package handler

import (
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"

	"github.com/zgsm-ai/gatewayctl/internal/pkg/core"
	"github.com/zgsm-ai/gatewayctl/internal/pkg/error"
	"github.com/zgsm-ai/gatewayctl/internal/service"
	"github.com/zgsm-ai/gatewayctl/internal/store"
)

type registerBody struct {
	Id    string `json:"id"`
	Extra string `json:"extra"`
}

type unregisterBody struct {
	Id string `json:"id"`
}

type queryParam struct {
	Id string `uri:"id" binding:"required"`
}

func RegisterPlugins(c *gin.Context) {
	var body registerBody
	if err := c.BindJSON(&body); err != nil {
		core.HandleError(c, http.StatusBadRequest, error.ErrBadRequest, nil)
		return
	}

	partUri := base64.RawURLEncoding.EncodeToString([]byte(body.Id))

	uri, err := service.AddRouterToGateway(partUri, body.Id)
	if err != nil {
		core.HandleError(c, http.StatusOK, error.NewError(error.ErrCreateRoute, err.Error()), nil)
		return
	}

	plugin := store.Plugin{
		ID:      body.Id,
		URL:     uri,
		Extra:   datatypes.JSON([]byte(body.Extra)),
		Deleted: false,
	}
	if err := store.PluginModel.Create(&plugin); err != nil {
		core.HandleError(c, http.StatusOK, error.NewError(error.ErrDatabase, err.Error()), nil)
		return
	}

	core.HandleSuccess(c, map[string]string{"uri": uri})
}

func UnregisterPlugins(c *gin.Context) {
	var body unregisterBody
	if err := c.BindJSON(&body); err != nil {
		core.HandleError(c, http.StatusBadRequest, error.ErrBadRequest, nil)
		return
	}

	partUri := base64.RawURLEncoding.EncodeToString([]byte(body.Id))
	err := service.RemoveRouterFromGateway(partUri)
	if err != nil {
		core.HandleError(c, http.StatusOK, error.NewError(error.ErrCreateRoute, err.Error()), nil)
		return
	}

	if err := store.PluginModel.Delete(body.Id); err != nil {
		core.HandleError(c, http.StatusOK, error.NewError(error.ErrDatabase, err.Error()), nil)
		return
	}

	core.HandleSuccess(c, nil)
}

func ListPlugins(c *gin.Context) {
	data, err := store.PluginModel.List()
	if err != nil {
		core.HandleError(c, http.StatusOK, error.NewError(error.ErrDatabase, err.Error()), nil)
		return
	}

	core.HandleSuccess(c, data)
}

func GetPlugin(c *gin.Context) {
	var param queryParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
		return
	}

	data, err := store.PluginModel.Get(param.Id)
	if err != nil {
		core.HandleError(c, http.StatusOK, error.NewError(error.ErrDatabase, err.Error()), nil)
		return
	}

	core.HandleSuccess(c, data)
}
