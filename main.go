package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/fazarmitrais/atm-simulation-stage-3/delivery/menu"
	accountCsv "github.com/fazarmitrais/atm-simulation-stage-3/domain/account/repository/csv"
	accountMapdatastore "github.com/fazarmitrais/atm-simulation-stage-3/domain/account/repository/mapDatastore"
	transactionMapdatastore "github.com/fazarmitrais/atm-simulation-stage-3/domain/transaction/repository/mapDatastore"
	"github.com/fazarmitrais/atm-simulation-stage-3/service"
	"github.com/joho/godotenv"
)

func main() {
	envInit()
	c := context.Background()
	fmt.Println("Welcome to the atm simulation")
	var path string
	flag.StringVar(&path, "path", "", "Path to directory")
	flag.Parse()
	acctMap := accountMapdatastore.NewMapDatastore()
	trxMap := transactionMapdatastore.NewMapDatastore()
	acctCsv := accountCsv.NewCSV(path)
	svc := service.NewService(acctMap, acctCsv, trxMap)
	importDataFromCSV(c, path, svc)
	m := menu.NewMenu(svc)
	m.Start()
}

func importDataFromCSV(c context.Context, path string, svc *service.Service) {
	if strings.TrimSpace(path) == "" {
		log.Fatalln("Please provide correct csv file path to import data")
	}
	fmt.Println("Importing data from csv file...")
	err := svc.Import(c)
	if err != nil {
		log.Fatalf("Error importing : %v", err)
	}
	fmt.Println("Successfully imported data from csv: ", path)
}

func envInit() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Fatalln("No .env file found")
	}
}
