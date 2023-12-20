package main

import (
	"database/sql"
	"errors"
	"time"
)

const saveTaxiInfoSQL = `
INSERT INTO
	work_home_taxi_price (datetime, from, to, waiting_time, duration, price, is_surge, temperature, rain_level, snow_level)
VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
`

func saveTaxiInfo(db *sql.DB, collectTime time.Time, fromLabel string, toLabel string, taxiInfo *TaxiInfo, weather *WeatherResponse) error {
	var err error

	if taxiInfo == nil {
		return errors.New("taxiInfo is nil")
	}

	if weather != nil {
		_, err = db.Exec(
			saveTaxiInfoSQL,
			collectTime,
			fromLabel,
			toLabel,
			taxiInfo.WaitingTime,
			taxiInfo.Time/60,
			taxiInfo.Price,
			taxiInfo.IsSurge,
			ConvertKelvinToCelsius(weather.Current.Temp),
			weather.Current.Rain.The1H,
			weather.Current.Snow.The1H,
		)
	} else {
		_, err = db.Exec(
			saveTaxiInfoSQL,
			collectTime,
			fromLabel,
			toLabel,
			taxiInfo.WaitingTime,
			taxiInfo.Time/60,
			taxiInfo.Price,
			taxiInfo.IsSurge,
			nil,
			nil,
			nil,
		)
	}

	return err
}
