package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/cakramediadata2022/chs_cloud_general/pkg/global_var"
	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
)

func GetMD5Hash(text string) string {
	text = text + global_var.PasswordKeyString
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func EncryptString(Key []byte, Text string) (string, error) {
	TextByte := []byte(Text)
	block, err := aes.NewCipher(Key)
	if err != nil {
		return "", err
	}
	b := base64.StdEncoding.EncodeToString(TextByte)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	Result := hex.EncodeToString(ciphertext[:])
	return Result, nil
}

func DecryptString(Key []byte, Text string) (string, error) {
	TextByte, err := hex.DecodeString(Text)
	block, err := aes.NewCipher(Key)
	if err != nil {
		return "", err
	}
	if len(TextByte) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := TextByte[:aes.BlockSize]
	TextByte = TextByte[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(TextByte, TextByte)
	data, err := base64.StdEncoding.DecodeString(string(TextByte))
	if err != nil {
		return "", err
	}
	Result := string(data)
	return Result, nil
}

func generateIV(iv []byte, salt string) []byte {
	ivString := string(iv[:])
	// fmt.Println("iv", ivString)
	var newIV string
	if salt != "" {
		ivLen := len(ivString)
		saltLen := len(salt)
		if saltLen > ivLen {
			newIV = salt[0:ivLen]
		} else {
			// keyR := ivString[0:saltLen]
			newIV = salt + ivString[saltLen:]
			// fmt.Println("r", keyR)
		}
	} else {
		newIV = ivString
	}
	// fmt.Println("ne", newIV)
	IV := []byte(newIV)
	return IV
}

func OpensslDecrypt(encryptedText string, salt string) (string, error) {
	key := global_var.AESSecretKey
	iv := generateIV(global_var.AESiv, salt)
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)

	plaintext = PKCS7UnPadding(plaintext)
	return fmt.Sprintf("%s\n", plaintext), nil
}

func OpensslEncrypt(text string, salt string) (string, error) {
	plaintext := []byte(text)
	key := global_var.AESSecretKey
	iv := generateIV(global_var.AESiv, salt)

	plaintext = PKCS7Padding(plaintext)
	ciphertext := make([]byte, len(plaintext))
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintext)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func PKCS7Padding(ciphertext []byte) []byte {
	padding := aes.BlockSize - len(ciphertext)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	unpadding := int(plantText[length-1])
	return plantText[:(length - unpadding)]
}

func DateOf(ADateTime time.Time) time.Time {
	Result := time.Date(ADateTime.Year(), ADateTime.Month(), ADateTime.Day(), 0, 0, 0, 0, time.UTC)
	return Result
}

func TimeOf(ADateTime time.Time) time.Time {
	Result := time.Date(0, 0, 0, ADateTime.Hour(), ADateTime.Minute(), ADateTime.Second(), 0, time.UTC)
	return Result
}

func ReplaceTime(ADate time.Time, ATime time.Time) time.Time {
	Result := time.Date(ADate.Year(), ADate.Month(), ADate.Day(), ATime.Hour(), ATime.Minute(), ATime.Second(), 0, time.UTC)
	return Result
}

func ReplaceTimeLocation(ADate time.Time, ATime time.Time, Loc *time.Location) time.Time {
	Result := time.Date(ADate.Year(), ADate.Month(), ADate.Day(), ATime.Hour(), ATime.Minute(), ATime.Second(), 0, Loc)
	return Result
}

func StrTimeToDateTime(timeString string) (time.Time, error) {
	// Get the current date
	currentDate := time.Now().Format("2006-01-02")

	// Combine the current date with the input time string
	fullTimeString := fmt.Sprintf("%s %s", currentDate, timeString)

	// Parse the combined string into a time.Time value
	parsedTime, err := time.Parse("2006-01-02 15:04", fullTimeString)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return time.Time{}, err
	}
	return parsedTime, nil
}

func DaysBetween(StartDate, EndDate time.Time) int {
	StartDate = DateOf(StartDate)
	EndDate = DateOf(EndDate)
	Result := int(EndDate.Sub(StartDate).Hours() / 24)
	return Result
}

func IncDay(ADate time.Time, Count int) time.Time {
	ADate = DateOf(ADate)
	Result := ADate.AddDate(0, 0, Count)
	return Result
}

func StrToDate(ADateStr string) time.Time {
	Date, err := time.ParseInLocation("2006-01-02", ADateStr, time.UTC)
	if err != nil {
		return time.Time{}
	} else {
		return Date
	}
}

func StrToDateTime(ADateStr string) time.Time {
	Date, err := time.ParseInLocation("2006-01-02 15:04:05", ADateStr, time.UTC)
	if err != nil {
		return time.Time{}
	} else {
		return Date
	}
}

func StrZToDate(ADateStr string) time.Time {
	Date, err := time.Parse(time.RFC3339, ADateStr)
	if err != nil {
		fmt.Println("Failed to parse date:", err)
		return time.Time{}
	} else {
		return Date
	}
}

func StrToTime(ATimeStr string) time.Time {
	Time, err := time.ParseInLocation("03:04:05", ATimeStr, time.UTC)
	if err != nil {
		return time.Time{}
	} else {
		return Time
	}
}

func BoolToUint8(Value bool) uint8 {
	if Value {
		return 1
	}
	return 0
}

func BoolToUint8String(Value bool) string {
	if Value {
		return "1"
	}
	return "0"
}

func BoolToString(Value bool) string {
	if Value {
		return "true"
	}
	return "false"
}

func StartDateOfTheMonth(ADate time.Time) time.Time {
	ADate = DateOf(ADate)
	Result := time.Date(ADate.Year(), ADate.Month(), 1, 0, 0, 0, 0, time.UTC)
	return Result
}

func FormatDate1(ADate time.Time) string {
	//yyyy-dd-mm
	ADate = DateOf(ADate)
	Result := ADate.Format("2006-01-02")
	return Result
}

func FormatTime1(ADate time.Time) string {
	//yyyy-dd-mm
	Result := ADate.Format("15:04:05")
	return Result
}

func FormatTimeHour(Hour int) string {
	formattedHour := strconv.Itoa(Hour)
	if len(formattedHour) == 1 {
		formattedHour = "0" + formattedHour
	}
	formattedHour += ":00"

	formattedHour = strconv.Itoa(Hour)
	formattedHour += ":00"
	return formattedHour
}

func FormatDate2(Format string, ADate time.Time) string {
	//yyyy-dd-mm
	ADate = DateOf(ADate)
	Result := ADate.Format(Format)
	return Result
}

func FormatDatePrefix(ADate time.Time) string {
	//yy
	ADate = DateOf(ADate)
	Result := ADate.Format("06")
	return Result
}

func Div(AValue, ADivider int64) int64 {
	return int64(math.Trunc(float64(AValue) / float64(ADivider)))
}

func InterfaceToUint64(v interface{}) uint64 {
	switch t := v.(type) {
	case int:
		return uint64(t)
	case int8:
		return uint64(t)
	case int16:
		return uint64(t)
	case int32:
		return uint64(t)
	case int64:
		return uint64(t)
	case float64:
		return uint64(t)
	default:
		return 0
	}
}

func InterfaceToInt64(v interface{}) int64 {
	switch t := v.(type) {
	case int:
		return int64(t)
	case int8:
		return int64(t)
	case int16:
		return int64(t)
	case int32:
		return int64(t)
	case int64:
		return int64(t)
	case float64:
		return int64(t)
	default:
		return 0
	}
}

func InterfaceToInt8(v interface{}) int8 {
	switch t := v.(type) {
	case int:
		return int8(t)
	case int8:
		return int8(t)
	case int16:
		return int8(t)
	case int32:
		return int8(t)
	case int64:
		return int8(t)
	case float64:
		return int8(t)
	default:
		return 0
	}
}

func InterfaceToFloat64(v interface{}) float64 {
	switch t := v.(type) {
	case int:
		return float64(t)
	case int8:
		return float64(t)
	case int16:
		return float64(t)
	case int32:
		return float64(t)
	case int64:
		return float64(t)
	case float32:
		return float64(t)
	case float64:
		return t
	default:
		return 0
	}
}

func InterfaceToBool(t interface{}) bool {
	switch t := t.(type) {
	case int:
		return t != 0
	case int8:
		return t != 0
	case int16:
		return t != 0
	case int32:
		return t != 0
	case int64:
		return t != 0
	case uint64:
		return t != 0
	default:
		return false
	}
}

func Uint8ToBool(i uint8) bool {
	return i != 0
}

func Uint64ToStr(i uint64) string {
	return strconv.FormatUint(i, 10)
}

func ByteToFloat64(bytes []byte) float64 {
	float, _ := strconv.ParseFloat(string(bytes), 64)
	return float
}

func ByteToUint64(bytes []byte) uint64 {
	Uint64, _ := strconv.ParseUint(string(bytes), 10, 64)
	return Uint64
}

// TODO Review result
func RoundTo(AFloat float64) float64 {
	return math.Trunc(AFloat)
}

// TODO Review result
func RoundToX2(AFloat float64) float64 {
	return math.Trunc(AFloat*100) / 100
}

// TODO Review result
func RoundToX3(AFloat float64) float64 {
	return math.Trunc(AFloat*1000) / 1000
}

func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func FloatToStrX3(AFloat float64) string {
	return fmt.Sprintf("%.3f", AFloat)
}

func FloatToStrX2(AFloat float64) string {
	return fmt.Sprintf("%.2f", AFloat)
}

func FormatNumber(number float64) string {
	// Convert number to string with desired format
	formatted := strconv.FormatFloat(number, 'f', 2, 64)

	// Split the number into integer and decimal parts
	integerPart := formatted[:len(formatted)-3]
	decimalPart := formatted[len(formatted)-2:]

	// Add comma as the thousands separator
	thousandsSeparator := ","
	formattedWithSeparator := ""
	for i, digit := range integerPart {
		if i > 0 && (len(integerPart)-i)%3 == 0 {
			formattedWithSeparator += thousandsSeparator
		}
		formattedWithSeparator += string(digit)
	}

	// Concatenate the formatted number with the decimal part
	return formattedWithSeparator + "." + decimalPart
}

func IsWeekend(ATime time.Time, Dataset *global_var.TDataset) bool {
	return (Dataset.ProgramConfiguration.FridayAsWeekend && (int(ATime.Weekday()) == 5)) || (Dataset.ProgramConfiguration.SaturdayAsWeekend && (int(ATime.Weekday()) == 6)) || (Dataset.ProgramConfiguration.SundayAsWeekend && (int(ATime.Weekday()) == 7))
}

func PtrString(i string) *string {
	return &i
}

func PtrInt(i int) *int {
	return &i
}

func PtrUint(i uint) *uint {
	return &i
}

func PtrUint8(i uint8) *uint8 {
	return &i
}

func PtrUint32(i uint32) *uint32 {
	return &i
}

func PtrUint64(i uint64) *uint64 {
	return &i
}

func PtrFloat64(i float64) *float64 {
	return &i
}

func PtrFloat32(i float32) *float32 {
	return &i
}

func PtrTime(i time.Time) *time.Time {
	return &i
}

func StrToInt(str string) int {
	number, err := strconv.Atoi(str)
	if err != nil {
		number = 0
	}
	return number
}

func StrToFloat64(str string) float64 {
	number, error := strconv.ParseFloat(str, 64)
	if error != nil {
		number = 0
	}
	return number
}

func StrToUint64(str string) uint64 {
	number, error := strconv.ParseUint(string(str), 10, 64)
	if error != nil {
		number = 0
	}
	return number
}

func StrToInt64(str string) int64 {
	number, error := strconv.ParseInt(string(str), 10, 64)
	if error != nil {
		number = 0
	}
	return number
}

func StrToBool(str string) bool {
	return strings.ToUpper(str) == "TRUE" || str == "1"
}

func StrToUint8(str string) uint8 {
	number, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		number = 0
	}
	return uint8(number)
}

func GetMonthCountRange(FromMonth int, FromYear int, ToMonth int, ToYear int) (int, error) {
	if FromYear > ToYear {
		panic(errors.New("To year cannot lower than from year"))
	}
	fromDate := time.Date(FromYear, time.Month(FromMonth), 1, 0, 0, 0, 0, time.UTC)
	toDate := time.Date(ToYear, time.Month(ToMonth), 1, 0, 0, 0, 0, time.UTC)
	y, m, _, _, _, _ := TimeDiff(fromDate, toDate)
	month := (y * 12) + m

	return month + 1, nil
}

func GetFirstDateOfMonth(Date time.Time) time.Time {
	now := Date
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)

	return firstOfMonth
}

func GetLastDateOfMonth(Date time.Time) time.Time {
	now := Date
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	return lastOfMonth
}

func EndOfAMonth(currentYear int, currentMonth int) (time.Time, error) {
	if currentMonth > 12 {
		return time.Time{}, errors.New("Month cannot more than 12")
	}
	firstOfMonth := time.Date(currentYear, time.Month(currentMonth), 1, 0, 0, 0, 0, time.UTC)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	return lastOfMonth, nil
}

func StartOfAMonth(currentYear int, currentMonth int) (time.Time, error) {
	if currentMonth > 12 {
		return time.Time{}, errors.New("Month cannot more than 12")
	}
	firstOfMonth := time.Date(currentYear, time.Month(currentMonth), 1, 0, 0, 0, 0, time.UTC)

	return firstOfMonth, nil
}

func TimeDiff(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}

func GetStrBetween(str string, startS string, endS string) (result string, found bool) {
	s := strings.Index(str, startS)
	if s == -1 {
		return result, false
	}
	newS := str[s+len(startS):]
	e := strings.Index(newS, endS)
	if e == -1 {
		return result, false
	}
	result = newS[:e]
	return result, true
}

// enum validator
func Enum(
	fl validator.FieldLevel,
) bool {
	enumString := fl.Param()                    // get string
	value := fl.Field().String()                // the actual field
	enumSlice := strings.Split(enumString, "_") // convert to slice
	for _, v := range enumSlice {
		if value == v {
			return true
		}
	}
	return false
}

func IfThenElse(condition bool, a interface{}, b interface{}) interface{} {
	if condition {
		return a
	}
	return b
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "lte":
		return "Should be less than " + fe.Param()
	case "gte":
		return "Should be greater than " + fe.Param()
	case "min":
		return "Should be " + fe.Param() + " length at least\n"
	case "enum":
		replacer := *strings.NewReplacer("_", ",")
		return "Should be one of " + replacer.Replace(fe.Param())
	}
	return "Unknown error"
}

func GenerateValidateErrorMsg(c *gin.Context, err error) interface{} {

	type ErrorMsg struct {
		Field   string `json:"field"`
		Message string `json:"message"`
	}
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make([]ErrorMsg, len(ve))
		for i, fe := range ve {
			out[i] = ErrorMsg{fe.Field(), getErrorMsg(fe)}
		}
		// if out == nil {
		// 	return err
		// }
		return out
	}
	return err
}

func IntToStr(i int64) string {
	return strconv.FormatInt(i, 36)
}

func PtrToStr(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func PtrToUint8(s *uint8) uint8 {
	if s != nil {
		return *s
	}
	return 0
}

func PtrToInt(s *int) int {
	if s != nil {
		return *s
	}
	return 0
}

func PtrToFloat64(s *float64) float64 {
	if s != nil {
		return *s
	}
	return 0
}

// layout can be empty
func IsDate(dateString string, layout string) bool {
	if layout == "" {
		layout = "2006-01-02"
	}
	_, err := time.Parse(layout, dateString)
	return err == nil
}
