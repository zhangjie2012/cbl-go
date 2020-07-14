package cbl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const yyServer = "https://www.feifeicloud.cn/yinyang"

type YangDate struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}

type YinDate struct {
	YearNum    int    `json:"year_num"`
	YearTian   string `json:"year_tian"`
	YearDi     string `json:"year_di"`
	YearZodiac string `json:"year_zodiac"`
	MonthNum   int    `json:"month_num"`
	MonthName  string `json:"month_name"`
	MonthLeap  bool   `json:"month_leap"`
	DayNum     int    `json:"day_num"`
	DayName    string `json:"day_name"`
	WeekDay    string `json:"weekday"`
	SolarTerm  string `json:"solarterm"`
}

func (y *YinDate) ToString1() string {
	month := y.MonthName
	if y.MonthLeap {
		month = "闰" + month
	}
	return fmt.Sprintf("%d年%s%s %s", y.YearNum, y.MonthName, y.DayName, y.WeekDay)
}

func (y *YinDate) ToString2() string {
	month := y.MonthName
	if y.MonthLeap {
		month = "闰" + month
	}
	full := fmt.Sprintf("%s%s年【%s年】%s%s %s", y.YearTian, y.YearDi, y.YearZodiac,
		y.MonthName, y.DayName, y.WeekDay)
	if y.SolarTerm != "" {
		full = full + " " + y.SolarTerm
	}
	return full
}

// ConvYinYang 农历转公历
func ConvYinYang(year int, month int, leap int, day int) (time.Time, error) {
	url := fmt.Sprintf("%s/api/v1/conv/yin-yang/%d/%d/%d/%d", yyServer, year, month, leap, day)
	resp, err := http.Get(url)
	if err != nil {
		return time.Time{}, err
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return time.Time{}, err
	}

	data := struct {
		Code  int       `json:"code"`
		Data  *YangDate `json:"data"`
		Error string    `json:"error"`
	}{}
	if err := json.Unmarshal(bs, &data); err != nil {
		return time.Time{}, err
	}

	if data.Code != 0 {
		return time.Time{}, fmt.Errorf(data.Error)
	}

	t := time.Date(data.Data.Year, time.Month(data.Data.Month), data.Data.Day, 0, 0, 0, 0, time.Now().Location())
	return t, nil
}

// ConvYangYin 公历转农历
func ConvYangYin(year int, month int, day int) (*YinDate, error) {
	url := fmt.Sprintf("%s/api/v1/conv/yang-yin/%d/%d/%d", yyServer, year, month, day)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	data := struct {
		Code  int      `json:"code"`
		Data  *YinDate `json:"data"`
		Error string   `json:"error"`
	}{}
	if err := json.Unmarshal(bs, &data); err != nil {
		return nil, err
	}

	if data.Code != 0 {
		return nil, fmt.Errorf(data.Error)
	}

	return data.Data, nil
}
