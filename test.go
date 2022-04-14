package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"io/ioutil"
	"encoding/json"
	"runtime"
)
type Response struct {
	Meta struct {
		Code          string `json:"code"`
		ExecutionTime string `json:"execution_time"`
	} `json:"meta"`
	Data struct {
		IP          string      `json:"ip"`
		City        interface{} `json:"city"`
		Postal      interface{} `json:"postal"`
		State       interface{} `json:"state"`
		StateCode   interface{} `json:"state_code"`
		Country     string      `json:"country"`
		CountryCode string      `json:"country_code"`
		Location    string      `json:"location"`
		Timezone    struct {
			ID                  string `json:"id"`
			Location            string `json:"location"`
			CountryCode         string `json:"country_code"`
			CountryName         string `json:"country_name"`
			Iso31661Alpha2      string `json:"iso3166_1_alpha_2"`
			Iso31661Alpha3      string `json:"iso3166_1_alpha_3"`
			UnM49Code           string `json:"un_m49_code"`
			Itu                 string `json:"itu"`
			Marc                string `json:"marc"`
			Wmo                 string `json:"wmo"`
			Ds                  string `json:"ds"`
			PhonePrefix         string `json:"phone_prefix"`
			Fifa                string `json:"fifa"`
			Fips                string `json:"fips"`
			Gual                string `json:"gual"`
			Ioc                 string `json:"ioc"`
			CurrencyAlphaCode   string `json:"currency_alpha_code"`
			CurrencyCountryName string `json:"currency_country_name"`
			CurrencyMinorUnit   string `json:"currency_minor_unit"`
			CurrencyName        string `json:"currency_name"`
			CurrencyCode        string `json:"currency_code"`
			Independent         string `json:"independent"`
			Capital             string `json:"capital"`
			Continent           string `json:"continent"`
			Tld                 string `json:"tld"`
			Languages           string `json:"languages"`
			GeonameID           string `json:"geoname_id"`
			Edgar               string `json:"edgar"`
		} `json:"timezone"`
		Datetime struct {
			Date          string `json:"date"`
			DateTime      string `json:"date_time"`
			DateTimeTxt   string `json:"date_time_txt"`
			DateTimeWti   string `json:"date_time_wti"`
			DateTimeYmd   string `json:"date_time_ymd"`
			Time          string `json:"time"`
			Month         string `json:"month"`
			MonthWilz     string `json:"month_wilz"`
			MonthAbbr     string `json:"month_abbr"`
			MonthFull     string `json:"month_full"`
			MonthDays     string `json:"month_days"`
			Day           string `json:"day"`
			DayWilz       string `json:"day_wilz"`
			DayAbbr       string `json:"day_abbr"`
			DayFull       string `json:"day_full"`
			Year          string `json:"year"`
			YearAbbr      string `json:"year_abbr"`
			Hour12Wolz    string `json:"hour_12_wolz"`
			Hour12Wilz    string `json:"hour_12_wilz"`
			Hour24Wolz    string `json:"hour_24_wolz"`
			Hour24Wilz    string `json:"hour_24_wilz"`
			HourAmPm      string `json:"hour_am_pm"`
			Minutes       string `json:"minutes"`
			Seconds       string `json:"seconds"`
			Week          string `json:"week"`
			OffsetSeconds string `json:"offset_seconds"`
			OffsetMinutes string `json:"offset_minutes"`
			OffsetHours   string `json:"offset_hours"`
			OffsetGmt     string `json:"offset_gmt"`
			OffsetTzid    string `json:"offset_tzid"`
			OffsetTzab    string `json:"offset_tzab"`
			OffsetTzfull  string `json:"offset_tzfull"`
			TzString      string `json:"tz_string"`
			Dst           string `json:"dst"`
			DstObserves   string `json:"dst_observes"`
			TimedaySpe    string `json:"timeday_spe"`
			TimedayGen    string `json:"timeday_gen"`
		} `json:"datetime"`
	} `json:"data"`
}

func main() {
	fmt.Println("Hello [%s]", runtime.GOARCH)
	LOG_FILE := "./app.log"
	PORT := "8082"
	logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE,0644)
	if err != nil {
        log.Panic(err)
    }
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	
	http.HandleFunc("/", test)
	http.HandleFunc("/log", logShow)
	log.Fatal(http.ListenAndServe(":" + PORT, nil))
	log.Println("Running on port " + PORT)
}

func test(w http.ResponseWriter, r *http.Request) {
	IPAddress := r.Header.Get("X-Real-Ip")
	msg := ""
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	IPAddress = "66.220.144.0"
	tmpStr := strings.Split(IPAddress, ":")
	IPAddress = tmpStr[0]
	msg = IPAddress
	if net.ParseIP(IPAddress).IsPrivate(){
		msg += " <- adres prywatny"
	}
	if net.ParseIP(IPAddress).IsLoopback(){
		msg += " <- adres loopback"
	}
	
	url := "https://timezoneapi.io/api/ip/?" + IPAddress + "&token=abQGXSYEczrgJCKoQmuA"
	response, err := http.Get(url)
	if err != nil {
        fmt.Print(err.Error())
        os.Exit(1)
    }
	responseData, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Fatal(err)
    }
	var result Response
	json.Unmarshal(responseData, &result)
	msg += "\n" + result.Data.Datetime.Date + " " + result.Data.Datetime.Time
	fmt.Fprint(w, msg)
	log.Println("NEW IP: " + IPAddress)
}

func logShow(w http.ResponseWriter, r *http.Request){
	content, err := ioutil.ReadFile("app.log")
    if err != nil {
        log.Fatal(err)
    }
	fmt.Fprint(w, string(content))
}
