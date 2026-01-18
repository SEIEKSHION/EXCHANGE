package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

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

func ProceedExchangeVaults(body []byte) (map[string]float64, error) {

	var reader io.Reader = nil
	reader, err := charset.NewReaderLabel("windows-1251", bytes.NewReader(body)) // декодирование из windows1251 в utf8
	if err != nil {
		return nil, fmt.Errorf("GetVaultExchange:\n\t\tError while decoding: %v", err)
	}

	var utf8Body []byte
	utf8Body, err = io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("GetVaultExchange:\n\t\tError while reading UTF-8: %v", err)
	}

	// Декодируем XML
	var valCurs ValCurs
	decoder := xml.NewDecoder(bytes.NewReader(utf8Body))
	decoder.CharsetReader = charset.NewReaderLabel // Указываем обработчик кодировок

	err = decoder.Decode(&valCurs)
	if err != nil {
		return nil, fmt.Errorf("GetVaultExchange:\n\t\tError while parsing from XML: %v", err)
	}
	for _, valute := range valCurs.Valutes {
		if valute.ID == "R01030" {
			fmt.Println(valute.Name)
			fmt.Println(valute.GetNumericVunitRate())
		}
		// ID        string `xml:"ID,attr"`
		// NumCode   string `xml:"NumCode"`
		// CharCode  string `xml:"CharCode"`
		// Nominal   int    `xml:"Nominal"`
		// Name      string `xml:"Name"`
		// Value     string `xml:"Value"`
		// VunitRate string `xml:"VunitRate"`
	}
	return nil, nil
}

func main() {
	// получение данных
	body, err := GetVaultExchange()
	if err != nil {
		panic(fmt.Errorf("Main: \n\t%v", err))
	}
	var data map[string]float64
	data, err = ProceedExchangeVaults(body)
	if err != nil {
		panic(fmt.Errorf("Main:\n\t%v", err))
	}
	fmt.Println(data)

}
