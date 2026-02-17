package models

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/net/html/charset"
)

var (
	ValuteNotFound = errors.New("Valute not found")
)

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Date    string   `xml:"Date,attr"`
	Name    string   `xml:"name,attr"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	ID        string `xml:"ID,attr"`
	NumCode   string `xml:"NumCode"`
	CharCode  string `xml:"CharCode"`
	Nominal   int    `xml:"Nominal"`
	Name      string `xml:"Name"`
	Value     string `xml:"Value"`
	VunitRate string `xml:"VunitRate"`
}

func (v *Valute) GetNumericValue() (float64, error) {
	// В ответе ЦБ используется запятая как десятичный разделитель
	str := strings.Replace(v.Value, ",", ".", 1)
	return strconv.ParseFloat(str, 64)
}

// GetNumericVunitRate преобразует VunitRate в float64
func (v *Valute) GetNumericVunitRate() (float64, error) {
	str := strings.Replace(v.VunitRate, ",", ".", 1)
	return strconv.ParseFloat(str, 64)
}

// должны написать функцию получения данных о валютах на сегодняшний день
func GetVaultExchange() ([]byte, error) {
	date := time.Now().Format("02.01.2006")
	URI := fmt.Sprintf("http://www.cbr.ru/scripts/XML_daily.asp?date_req=%s", date)

	req, err := http.NewRequest("GET", URI, nil)
	if err != nil {
		return nil, fmt.Errorf("GetVaultExchange: error creating request: %v", err)
	}

	// Добавляем User-Agent, чтобы не казаться ботом
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GetVaultExchange: error while fetching data: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GetVaultExchange: unexpected status code: %d", resp.StatusCode)
	}

	var body []byte
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("GetVaultExchange: error while reading response body: %v", err)
	}

	if len(body) == 0 {
		return nil, errors.New("GetVaultExchange: empty response body")
	}

	// Проверим, не получили ли мы HTML вместо XML
	if bytes.Contains(body, []byte("<!DOCTYPE html")) || bytes.Contains(body, []byte("<html")) {
		return nil, errors.New("GetVaultExchange: received HTML instead of XML, possibly due to access restriction")
	}

	return body, nil
}

func ProceedExchangeVaults(body []byte) ([]Valute, error) {
	// Проверяем, валиден ли body как UTF-8
	if !utf8.Valid(body) {
		// Тогда пробуем декодировать из windows-1251
		reader, err := charset.NewReaderLabel("windows-1251", bytes.NewReader(body))
		if err != nil {
			return nil, fmt.Errorf("GetVaultExchange:\n\t\tError creating decoder: %v", err)
		}
		body, err = io.ReadAll(reader)
		if err != nil {
			return nil, fmt.Errorf("GetVaultExchange:\n\t\tError while reading UTF-8: %v", err)
		}
	}

	// Теперь body гарантированно валиден как UTF-8
	utf8Str := string(body)

	// Удаляем декларацию encoding из XML заголовка
	// Ищем <?xml ... encoding="..." ... ?>
	re := regexp.MustCompile(`(<\?xml[^>]*?)\s+encoding\s*=\s*["'][^"']*["']([^>]*?\?>)`)
	utf8StrCleaned := re.ReplaceAllString(utf8Str, "${1}${2}")

	// Преобразуем обратно в []byte
	cleanedBody := []byte(utf8StrCleaned)

	var valCurs ValCurs
	decoder := xml.NewDecoder(bytes.NewReader(cleanedBody))
	// Теперь CharsetReader можно не устанавливать, так как декларации нет
	// decoder.CharsetReader = charset.NewReaderLabel // Не нужно
	err := decoder.Decode(&valCurs)
	if err != nil {
		return nil, fmt.Errorf("GetVaultExchange:\n\t\tError while parsing from XML: %v", err)
	}

	return valCurs.Valutes, nil
}

func PrintValutes(data []Valute) error {
	now := string(time.Now().Format("02/01/2006"))
	fmt.Printf("	Курс следующих валют на сегодняшний день (%s):\n", now)
	for _, valute := range data {
		vUnit, err := valute.GetNumericVunitRate()
		if err != nil {
			return fmt.Errorf("PrintValutes: Failed to get numericVunitRate, error: %v", err)
		}
		fmt.Printf("Курс %s(%s) составляет %.3f руб.\n", valute.CharCode, valute.Name, vUnit)
	}
	return nil
}

func GetValuteByName(vaults []Valute, name string) (Valute, error) {
	var result Valute
	for _, vault := range vaults {
		if vault.Name == name {
			result = vault
			return result, nil
		}
	}
	return Valute{}, ValuteNotFound
}
