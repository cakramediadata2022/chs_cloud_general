package utils

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/cakramediadata2022/chs_cloud_general/internal/utils/cache"
	"github.com/cakramediadata2022/chs_cloud_general/pkg/global_var"
	"github.com/cakramediadata2022/chs_cloud_general/pkg/httpErrors"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

type ReqIDCtxKey struct{}

var Context = context.Background()

// Get config path for local or docker
func GetConfigPath(configPath string) string {
	if configPath == "docker" {
		return "./config/config-docker"
	}
	if configPath == "k8s" {
		return "./config/config-k8s"
	}
	if configPath == "staging" {
		return "./config/config-staging"
	}
	return "./config/config-local"
}

// Get user ip address
func GetIPAddress(c *gin.Context) string {
	ip := strings.Split(c.ClientIP(), ":")
	if ip[0] == "[" {
		ip = []string{"127.0.0.1"}
	}
	// fmt.Println("ip", ip)
	return ip[0]
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

func Post(URI string, Payload interface{}) (interface{}, int, error) {
	payload, err := json.Marshal(Payload)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return nil, http.StatusInternalServerError, err

	}

	req, err := http.NewRequest("POST", URI, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, http.StatusInternalServerError, err

	}

	Authorize := GetServiceAuthorized()
	req.Header.Set("Authorization", "Basic "+Authorize)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil, http.StatusInternalServerError, err

	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, http.StatusInternalServerError, err
	}
	// Check if the request was successful (status code 200)
	// if resp.StatusCode != http.StatusOK {
	// 	fmt.Println("Error: Unexpected status code:", resp.Status, URI, Payload)
	// 	return nil, errors.New(string(body))
	// }

	var jsonResponds interface{}
	if err := json.Unmarshal(body, &jsonResponds); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return "", http.StatusInternalServerError, err
	}
	// Print the HTTP status code and response body
	fmt.Println("Status Code:", resp.Status)
	fmt.Println("Response Body:", string(body))
	return jsonResponds, resp.StatusCode, nil
}

func Get(URI string, Payload interface{}) (interface{}, int, error) {
	req, err := http.NewRequest("GET", URI, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, http.StatusInternalServerError, err
	}

	Authorize := GetServiceAuthorized()
	req.Header.Set("Authorization", "Basic "+Authorize)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil, http.StatusInternalServerError, err

	}
	defer resp.Body.Close()
	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, http.StatusInternalServerError, err
	}
	// Check if the request was successful (status code 200)
	// if resp.StatusCode != http.StatusOK {
	// 	fmt.Println("Error: Unexpected status code:", resp.Status, URI, Payload)
	// 	return nil, errors.New(string(body))

	// }

	var jsonResponds interface{}
	if err := json.Unmarshal(body, &jsonResponds); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return "", http.StatusInternalServerError, err
	}
	// Print the HTTP status code and response body
	fmt.Println("Status Code:", resp.Status)
	fmt.Println("Response Body:", string(body))
	return jsonResponds, resp.StatusCode, nil
}

func GetToken(username string, password string) string {
	Token, err := cache.DataCache.GetString(Context, "TADA", username+password)
	if err != nil {
		return ""
	}
	return Token
}

func GetServiceAuthorized() string {
	API_Key := "CAKRASOFT_CLOUD"
	API_Secret := "ushdKUGHueKHBsdiyue324OUJNs"
	Authorization := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", API_Key, API_Secret)))

	return Authorization
}

// func RequestToken(username string, password string) (string, error) {
// 	API_Key := "vwNlJTehY24tArHWzpVWC0YVa"
// 	API_Secret := "4Ywhrkam5aOqZ7fKRyov8zD4H878ELsjqsZUGzkMDMD44odQCP"
// 	url := "https://api.gift.id/v1/pos/token"
// 	credential := map[string]interface{}{
// 		"username":   username,
// 		"password":   password,
// 		"grant_type": "password",
// 		"scope":      "offline_access",
// 	}
// 	payload, err := json.Marshal(credential)
// 	if err != nil {
// 		fmt.Println("Error marshaling JSON:", err)
// 		return "", err
// 	}

// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
// 	if err != nil {
// 		fmt.Println("Error creating request:", err)
// 		return "", err
// 	}

// 	Authorization := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", API_Key, API_Secret)))
// 	req.Header.Set("Authorization", "Basic "+Authorization)
// 	req.Header.Set("Content-Type", "application/json")
// 	client := &http.Client{}

// 	fmt.Println("request", Authorization)
// 	// Perform the request
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println("Error making request:", err)
// 		return "", err
// 	}
// 	defer resp.Body.Close()

// 	// Check if the request was successful (status code 200)
// 	if resp.StatusCode != http.StatusOK {
// 		fmt.Println("Error: Unexpected status code:", resp.Status, url, payload)
// 		return "", err
// 	}

// 	// Read the response body
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		fmt.Println("Error reading response body:", err)
// 		return "", err
// 	}
// 	var accessTokenResp map[string]interface{}
// 	if err := json.Unmarshal(body, &accessTokenResp); err != nil {
// 		fmt.Println("Error unmarshalling JSON:", err)
// 		return "", err
// 	}
// 	token := accessTokenResp["access_token"].(string)
// 	expire := accessTokenResp["expires_in"].(float64)

// 	cache.DataCache.Set(Context, "TADA", username+password, token, time.Second*time.Duration(expire))

// 	return token, nil
// }
