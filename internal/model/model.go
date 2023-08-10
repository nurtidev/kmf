package model

import "time"

type XMLRate struct {
	FullName    string `xml:"fullname"`
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Quant       int    `xml:"quant"`
	Index       string `xml:"index"`
	Change      string `xml:"change"`
}

type XMLRates struct {
	Items []XMLRate `xml:"item"`
}

type Rate struct {
	FullName    string
	Title       string
	Description float64
	Quant       int
	Index       string
	Change      string
}

type Currency struct {
	ID    int
	Title string
	Code  string
	Value float64
	Date  time.Time
}
