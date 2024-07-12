package utils

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cakramediadata2022/chs_cloud_general/pkg/global_var"
	"github.com/cakramediadata2022/chs_cloud_general/pkg/httpErrors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

type ReqIDCtxKey struct{}

// Get config path for local or docker
func GetConfigPath(configPath string) string {
	if configPath == "docker" {
		return "./config/config-docker"
	}

	if configPath == "staging" {
		return "./config/config-staging"
	}
	return "./config/config-local"
}

// Get user ip address
func GetIPAddress(c *gin.Context) string {
	return c.Request.RemoteAddr
}

func GetRequestCtx(c *gin.Context) context.Context {
	return context.WithValue(c.Request.Context(), ReqIDCtxKey{}, GetRequestID(c))
}

// Get request id from gi context
func GetRequestID(c *gin.Context) string {
	reqID := c.GetHeader("X-Request-ID")
	if reqID == "" {
		reqID = c.GetString("RequestID")
	}
	fmt.Println(reqID)
	return reqID
}

func GetUnitCode(c *gin.Context) string {
	return c.GetString("UnitCode")
}

// Error response with logging error for gin context
func ErrResponseWithLog(ctx *gin.Context, logger *otelzap.Logger, err error) {
	logger.Ctx(ctx).Error("ErrResponseWithLog", zap.Error(err),
		zap.String("UnitCode", GetUnitCode(ctx)),
		zap.String("RequestID", GetRequestID(ctx)),
		zap.String("IPAddress", GetIPAddress(ctx)))

	ctx.JSON(httpErrors.ErrorResponse(err))
}

// func ErrResponseWithLog(ctx *gin.Context, logger logger.Logger, err error) {
// 	logger.Errorf(
// 		"ErrResponseWithLog, UnitCode: %s, RequestID: %s, IPAddress: %s, Error: %s",
// 		GetUnitCode(ctx),
// 		GetRequestID(ctx),
// 		GetIPAddress(ctx),
// 		err,
// 	)
// 	ctx.JSON(httpErrors.ErrorResponse(err))

// }

// Error response with logging error for echo context
//
//	func LogResponseError(ctx *gin.Context, logger logger.Logger, err error) {
//		logger.Errorf(
//			"ErrResponseWithLog, UnitCode: %s, RequestID: %s, IPAddress: %s, Error: %s",
//			GetUnitCode(ctx),
//			GetRequestID(ctx),
//			GetIPAddress(ctx),
//			err,
//		)
//	}
//
// Error response with logging error for gin context
func LogResponseError(c *gin.Context, ctx context.Context, logger *otelzap.Logger, err error) {
	logger.Ctx(ctx).Error("ErrResponseWithLog", zap.Error(err),
		zap.String("UnitCode", GetUnitCode(c)),
		zap.String("RequestID", GetRequestID(c)),
		zap.String("IPAddress", GetIPAddress(c)))
}

func SendResponse(StatusCode uint, Message interface{}, Result interface{}, c *gin.Context) {
	var RequestResponse = global_var.TRequestResponse{
		StatusCode: StatusCode,
		Message:    Message,
		Result:     Result}

	if RequestResponse.StatusCode == global_var.ResponseCode.Successfully || RequestResponse.StatusCode == global_var.ResponseCode.SuccessfullyWithStatus {
		c.JSON(http.StatusOK, &RequestResponse)
	} else if (RequestResponse.StatusCode == global_var.ResponseCode.NotAuthorized) || (RequestResponse.StatusCode == global_var.ResponseCode.ErrorCreateToken) {
		c.JSON(http.StatusUnauthorized, &RequestResponse)
	} else if (RequestResponse.StatusCode == global_var.ResponseCode.InvalidDataFormat) || (RequestResponse.StatusCode == global_var.ResponseCode.DataNotFound) ||
		(RequestResponse.StatusCode == global_var.ResponseCode.InvalidDataValue) || (RequestResponse.StatusCode == global_var.ResponseCode.DatabaseValueChanged) ||
		(RequestResponse.StatusCode == global_var.ResponseCode.DatabaseError) || (RequestResponse.StatusCode == global_var.ResponseCode.DuplicateEntry) ||
		(RequestResponse.StatusCode == global_var.ResponseCode.OtherResult) || (RequestResponse.StatusCode == global_var.ResponseCode.Unregistered) || (RequestResponse.StatusCode == global_var.ResponseCode.SubscriptionExpired) {
		c.JSON(http.StatusBadRequest, &RequestResponse)
	} else {
		c.JSON(http.StatusBadRequest, Message)
	}
}

func SendWebsocketResponse(StatusCode uint, Message interface{}, Result interface{}, con *websocket.Conn) {
	var RequestResponse = global_var.TRequestResponse{
		StatusCode: StatusCode,
		Message:    Message,
		Result:     Result}

	if RequestResponse.StatusCode == global_var.ResponseCode.Successfully || RequestResponse.StatusCode == global_var.ResponseCode.SuccessfullyWithStatus {
		con.WriteJSON(&RequestResponse)
	} else if (RequestResponse.StatusCode == global_var.ResponseCode.NotAuthorized) || (RequestResponse.StatusCode == global_var.ResponseCode.ErrorCreateToken) {
		con.WriteJSON(&RequestResponse)
	} else if (RequestResponse.StatusCode == global_var.ResponseCode.InvalidDataFormat) || (RequestResponse.StatusCode == global_var.ResponseCode.DataNotFound) ||
		(RequestResponse.StatusCode == global_var.ResponseCode.InvalidDataValue) || (RequestResponse.StatusCode == global_var.ResponseCode.DatabaseValueChanged) ||
		(RequestResponse.StatusCode == global_var.ResponseCode.DatabaseError) || (RequestResponse.StatusCode == global_var.ResponseCode.DuplicateEntry) ||
		(RequestResponse.StatusCode == global_var.ResponseCode.OtherResult) {
		con.WriteJSON(&RequestResponse)
	} else {
		con.WriteJSON(Message)
	}
}
