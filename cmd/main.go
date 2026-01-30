package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/SEIEKSHION/Exchanger/internal/server"
	"golang.org/x/net/html/charset"
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
func GetVaultExchange() ([]byte, error) { // заменить первый вывод на map[string]float64
	date := string(time.Now().Format("02/01/2006"))
	URI := fmt.Sprintf("http://www.cbr.ru/scripts/XML_daily.asp?date_req=%s", date)
	resp, err := http.Get(URI)
	if err != nil {
		return nil, fmt.Errorf("GetVaultExchange:\n\tError while fetching data: %v", err)
	}

	defer resp.Body.Close() // гарантия закрытия тела ответа

	var body []byte
	body, err = io.ReadAll(resp.Body) // чтение тела ответа
	if err != nil {
		return nil, fmt.Errorf("GetVaultExchange:\n\t\tError while parsing data: %v", err)
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

func main() {
	body, err := GetVaultExchange()
	if err != nil {
		panic(fmt.Errorf("Main: \n\t%v", err))
	}
	var valutes []Valute
	valutes, err = ProceedExchangeVaults(body)
	if err != nil {
		panic(fmt.Errorf("Main:\n\t%v", err))
	}
	err = PrintValutes(valutes)
	if err != nil {
		panic(fmt.Errorf("Main: \n\t%v", err))
	}
	server.StartServer()
}
