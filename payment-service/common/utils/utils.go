package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"html/template"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/dustin/go-humanize"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type PaginationParam struct {
	Count int64       `json:"count"`
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
	Data  interface{} `json:"data"`
}

type PaginationResult struct {
	TotalPage    int         `json:"totalPage"`
	TotalData    int64       `json:"totalData"`
	NextPage     *int        `json:"nextPage"`
	PreviousPage *int        `json:"previousPage"`
	Page         int         `json:"page"`
	Limit        int         `json:"limit"`
	Data         interface{} `json:"data"`
}

func GeneratePagination(params PaginationParam) PaginationResult {
	totalPage := int(math.Ceil(float64(params.Count) / float64(params.Limit)))

	var (
		nextPage     int
		previousPage int
	)

	if params.Page < totalPage {
		nextPage = params.Page + 1
	}
	if params.Page > 1 {
		previousPage = params.Page - 1
	}

	result := PaginationResult{
		TotalPage:    totalPage,
		TotalData:    params.Count,
		NextPage:     &nextPage,
		PreviousPage: &previousPage,
		Page:         params.Page,
		Limit:        params.Limit,
		Data:         params.Data,
	}
	return result
}

func GenerateSha256(inputString string) string {
	hash := sha256.New()
	hash.Write([]byte(inputString))
	hashBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)
	return hashString
}

func RupiahFormat(amount *float64) string {
	stringValue := "0"
	if amount != nil {
		humanizeValue := humanize.CommafWithDigits(*amount, 0)
		stringValue = strings.ReplaceAll(humanizeValue, ",", ",")
	}
	return fmt.Sprintf("Rp. %s", stringValue)
}

// Read and Unmarshal file
func BindFromJSON(dest any, filename, path string) error {

	v := viper.New() //use viper as JSON parser
	v.SetConfigType("json")
	v.AddConfigPath(path)
	v.SetConfigName(filename)

	err := v.ReadInConfig()
	if err != nil {
		return err
	}

	// Convert []byte of JSON to a struct
	err = v.Unmarshal(&dest)
	if err != nil {
		logrus.Errorf("failed to unmarshal %s", err)
	}

	return nil
}

func SetEnvFromConsulKV(v *viper.Viper) error {
	env := make(map[string]any)

	err := v.Unmarshal(&env)
	if err != nil {
		logrus.Errorf("failed to unmarshal %s", err)
		return err
	}

	for k, v := range env {
		var (
			valOf = reflect.ValueOf(v)
			val   string
		)

		switch valOf.Kind() {
		case reflect.String:
			val = valOf.String()
		case reflect.Int:
			val = strconv.Itoa(int(valOf.Int()))
		case reflect.Uint:
			val = strconv.Itoa(int(valOf.Uint()))
		case reflect.Float32:
			val = strconv.Itoa(int(valOf.Float()))
		case reflect.Float64:
			val = strconv.Itoa(int(valOf.Float()))
		case reflect.Bool:
			val = strconv.FormatBool(valOf.Bool())
		default:
			panic("Unsupported Type")
		}

		err = os.Setenv(k, val)
		if err != nil {
			logrus.Errorf("failed to set env: %v", err)
			return err
		}
	}

	return nil

}

func BindFromConsul(dest any, endPoint, path string) error {
	v := viper.New()
	v.SetConfigType("json")

	err := v.AddRemoteProvider("consul", endPoint, path)
	if err != nil {
		logrus.Errorf("failed to add remote provider %v", err)
		return err

	}

	err = v.ReadRemoteConfig()
	if err != nil {
		logrus.Errorf("failed to read remote provider %v", err)
		return err

	}

	err = v.Unmarshal(&dest)
	if err != nil {
		logrus.Errorf("failed to unmarshal %v", err)
		return err

	}

	err = SetEnvFromConsulKV(v)
	if err != nil {
		logrus.Errorf("failed to set env from cosul kv: %v", err)
		return err
	}

	return nil
}

func Add1(a int) int {
	return a + 1
}

// Create PDF File
func GeneratePDFFromHTML(htmlTemplate string, data any) ([]byte, error) {
	funcMap := template.FuncMap{
		"add1": Add1,
	}

	tmpl, err := template.New("htmlTemplate").Funcs(funcMap).Parse(htmlTemplate)
	if err != nil {
		return nil, err
	}

	var filledTemplate bytes.Buffer
	if err := tmpl.Execute(&filledTemplate, data); err != nil {
		return nil, err
	}
	htmlContent := filledTemplate.String()

	pdfGenerator, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		logrus.Errorf("failed to create pdf generator: %v", err)
		return nil, err
	}

	pdfGenerator.Dpi.Set(600)
	pdfGenerator.NoCollate.Set(false)
	pdfGenerator.Orientation.Set(wkhtmltopdf.OrientationPortrait)
	pdfGenerator.PageSize.Set(wkhtmltopdf.PageSizeA6)
	pdfGenerator.Grayscale.Set(false)
	pdfGenerator.AddPage(wkhtmltopdf.NewPageReader(strings.NewReader(htmlContent)))

	err = pdfGenerator.Create()
	if err != nil {
		logrus.Errorf("failed to create pdf generator: %v", err)
		return nil, err
	}

	return pdfGenerator.Bytes(), err
}
