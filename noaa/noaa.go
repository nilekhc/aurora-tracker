package noaa

import (
	"aurora-tracker/utils"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	noaaGeomagneticForecastEndpoint = "https://services.swpc.noaa.gov/text/3-day-geomag-forecast.txt"
)

// CheckForAuroraProbability ...
func CheckForAuroraProbability(kpIndexThreshold int) {
	forecast := getGeomagneticForecast()
	fmt.Println(parseGeomagneticForecast(forecast, kpIndexThreshold))
}

func parseGeomagneticForecast(forecast string, kpIndexThreshold int) string {
	// parse issued date
	issuedDateRegex := `Issued: (\d+) (\w+) (\d+) (\d{2})(\d{2}) UTC`
	issuedDate := regexp.MustCompile(issuedDateRegex).FindStringSubmatch(forecast)
	year := issuedDate[1]

	// parse forecast dates
	forecastDatesRegex := `([A-Z][a-z]{2}) ([0-9]{2})\s+([A-Z][a-z]{2}) ([0-9]{2})\s+([A-Z][a-z]{2}) ([0-9]{2})`
	forecastDates := regexp.MustCompile(forecastDatesRegex).FindStringSubmatch(forecast)

	// parse kp index
	kpIndexDataRegex := `(\d{2})-(\d{2})UT\s+(\d)\s+(\d)\s+(\d)`
	kpIndexData := regexp.MustCompile(kpIndexDataRegex).FindAllString(forecast, 8)

	// parse kp index threshold
	var auroraProbabilities []string
	for _, kp := range kpIndexData {
		kpIndexRegex := `(\d{2})-(\d{2})UT\s+(\d)\s+(\d)\s+(\d)`
		kpIndex := regexp.MustCompile(kpIndexRegex).FindStringSubmatch(kp)

		if utils.ConvertToInt(kpIndex[3]) >= kpIndexThreshold {
			dateTime, err := time.Parse("2006-Jan-02", year+"-"+forecastDates[1]+"-"+forecastDates[2])
			if err != nil {
				panic(fmt.Errorf("Error parsing date: %s", err))
			}

			auroraProbabilities = append(auroraProbabilities, fmt.Sprintf("%s during %s-%sUT with KP Index %d!\n", fmt.Sprintf(dateTime.Format("2006-Jan-02")), kpIndex[1], kpIndex[2], utils.ConvertToInt(kpIndex[3])))
		}

		if utils.ConvertToInt(kpIndex[4]) >= kpIndexThreshold {
			dateTime, err := time.Parse("2006-Jan-02", year+"-"+forecastDates[3]+"-"+forecastDates[4])
			if err != nil {
				panic(fmt.Errorf("Error parsing date: %s", err))
			}

			auroraProbabilities = append(auroraProbabilities, fmt.Sprintf("%s during %s-%sUT with KP Index %d!\n", fmt.Sprintf(dateTime.Format("2006-Jan-02")), kpIndex[1], kpIndex[2], utils.ConvertToInt(kpIndex[4])))
		}

		if utils.ConvertToInt(kpIndex[5]) >= kpIndexThreshold {
			dateTime, err := time.Parse("2006-Jan-02", year+"-"+forecastDates[5]+"-"+forecastDates[6])
			if err != nil {
				panic(fmt.Errorf("Error parsing date: %s", err))
			}

			auroraProbabilities = append(auroraProbabilities, fmt.Sprintf("%s during %s-%sUT with KP Index %d!\n", fmt.Sprintf(dateTime.Format("2006-Jan-02")), kpIndex[1], kpIndex[2], utils.ConvertToInt(kpIndex[5])))
		}
	}

	if len(auroraProbabilities) == 0 {
		return ""
	}

	return fmt.Sprintf("Date %s\n\nAurora/e can be possible on: \n%s", issuedDate[0], strings.Join(auroraProbabilities, ""))
}

// GetGeomagneticForecast gets the 3 day geomagnetic forecast from NOAA.
func getGeomagneticForecast() string {
	response, err := http.Get(noaaGeomagneticForecastEndpoint)
	if err != nil {
		panic(fmt.Errorf("Error getting geomagnetic forecast data from NOAA: %s", err))
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(fmt.Errorf("Error reading geomagnetic forecast data from NOAA: %s", err))
	}

	return string(body)
}
