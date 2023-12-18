package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"go.uber.org/zap"
)

func runApplication() int {
	var configFilePath string
	flag.StringVar(&configFilePath, "c", "config.yaml", "config file path")

	logger, err := zap.NewProduction()
	if err != nil {
		log.Printf("can't create logger: %s\n", err)
	}

	cfg, err := readConfig(configFilePath)
	if err != nil {
		logger.Error("can't read config file", zap.Error(err))
		return 1
	}

	db := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{cfg.Database.Address},
		Auth: clickhouse.Auth{
			Database: cfg.Database.Database,
			Username: cfg.Database.Username,
			Password: cfg.Database.Password,
		},
	})
	err = upMigrations(db)
	if err != nil {
		logger.Error("can't up migrations", zap.Error(err))
		return 1
	}

	proxies, err := getProxies(fmt.Sprintf("%s?t=%d", cfg.Proxies.URL, time.Now().Unix()), cfg.Proxies.Type)
	if err != nil {
		logger.Error("can't get proxies", zap.Error(err))
		return 1
	}

	cookie, err := getCookies(fmt.Sprintf("%s?t=%d", cfg.Cookies.URL, time.Now().Unix()))
	if err != nil {
		logger.Error("can't get cookies", zap.Error(err))
		return 1
	}

	collectTime := time.Now()

	wasError := false

	taxiInfoFromWorkToHome, err := getMoscowTaxiRouteWithProxies(proxies, cookie, cfg.Points.Work, cfg.Points.Home)
	if err != nil {
		logger.Error("can't get taxi info", zap.String("direction", "to_home"), zap.Error(err))
		wasError = true
	} else {
		_, err = db.Exec("INSERT INTO work_home_taxi_price (datetime, from, to, waiting_time, duration, price, is_surge) VALUES (?, ?, ?, ?, ?, ?, ?)",
			collectTime, "work", "home", taxiInfoFromWorkToHome.WaitingTime, taxiInfoFromWorkToHome.Time/60, taxiInfoFromWorkToHome.Price, taxiInfoFromWorkToHome.IsSurge)
		if err != nil {
			logger.Error("can't write route from work to home", zap.Error(err))
			wasError = true
		}
	}

	taxiInfoFromHomeToWork, err := getMoscowTaxiRouteWithProxies(proxies, cookie, cfg.Points.Home, cfg.Points.Work)
	if err != nil {
		logger.Error("can't get taxi info", zap.String("direction", "to_work"), zap.Error(err))
		wasError = true
	} else {
		_, err = db.Exec("INSERT INTO work_home_taxi_price (datetime, from, to, waiting_time, duration, price, is_surge) VALUES (?, ?, ?, ?, ?, ?, ?)",
			collectTime, "home", "work", taxiInfoFromHomeToWork.WaitingTime, taxiInfoFromHomeToWork.Time/60, taxiInfoFromHomeToWork.Price, taxiInfoFromHomeToWork.IsSurge)
		if err != nil {
			logger.Error("can't write route from home to work", zap.Error(err))
			wasError = true
		}
	}

	if wasError {
		return 1
	} else {
		return 0
	}
}

func main() {
	os.Exit(runApplication())
}
