package main

import (
	"fmt"
	"os"

	"github.com/sbreitf1/money-graph/internal/moneydb"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

func main() {
	formatter := &nested.Formatter{
		HideKeys:        true,
		TimestampFormat: "2006-01-02 15:04:05",
		NoColors:        true,
	}
	logrus.SetFormatter(formatter)
	logrus.Infof("Startup")

	logrus.SetLevel(logrus.DebugLevel)

	db, err := moneydb.OpenOrCreateInFolder(".runtime/db", "TESTING")
	if err != nil {
		fmt.Println("ERR", err.Error())
		os.Exit(1)
	}

	fmt.Println(db.ImportCSV("testing/csv-examples/private-20220215-126303197-umsatz.CSV"))
}
