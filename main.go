package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/fazarmitrais/atm-simulation-stage-3/config/postgre"
	"github.com/fazarmitrais/atm-simulation-stage-3/delivery/menu"
	accountCsv "github.com/fazarmitrais/atm-simulation-stage-3/domain/account/repository/csv"
	accountPostgreRepo "github.com/fazarmitrais/atm-simulation-stage-3/domain/account/repository/postgre"
	transactionMapdatastore "github.com/fazarmitrais/atm-simulation-stage-3/domain/transaction/repository/mapDatastore"
	"github.com/fazarmitrais/atm-simulation-stage-3/service"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	envInit()
	e := echo.New()
	context := e.AcquireContext()
	fmt.Println("Welcome to the atm simulation")
	db := postgre.Connection()

	var path string
	flag.StringVar(&path, "path", "", "Path to directory")
	flag.Parse()
	acctRepo := accountPostgreRepo.NewPostgre(db)
	trxMap := transactionMapdatastore.NewMapDatastore()
	acctCsv := accountCsv.NewCSV(path)
	svc := service.NewService(acctRepo, acctCsv, trxMap)
	importDataFromCSV(context, path, svc)
	m := menu.NewMenu(svc)
	m.Start(context)
}

func importDataFromCSV(ctx echo.Context, path string, svc *service.Service) {
	if strings.TrimSpace(path) == "" {
		log.Fatalln("Please provide correct csv file path to import data")
	}
	fmt.Println("Importing data from csv file...")
	err := svc.Import(ctx)
	if err != nil {
		log.Fatalf("Error importing: %v", err)
	}
	fmt.Println("Successfully imported data from csv: ", path)
}

func envInit() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Fatalln("No .env file found")
	}
}
