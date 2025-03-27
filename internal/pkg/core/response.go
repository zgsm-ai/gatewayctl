// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
	errv1 "github.com/zgsm-ai/gatewayctl/internal/pkg/error"
)

type response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func HandleSuccess(ctx *gin.Context, data interface{}) {
	if data == nil {
		data = map[string]interface{}{}
	}
	resp := response{Code: errv1.Success.Code, Msg: errv1.Success.Error(), Data: data}
	ctx.JSON(http.StatusOK, resp)
}

func HandleError(ctx *gin.Context, httpCode int, err errv1.CustomError, data interface{}) {
	if data == nil {
		data = map[string]string{}
	}
	resp := response{Code: err.Code, Msg: err.Error(), Data: data}
	ctx.JSON(httpCode, resp)
}
