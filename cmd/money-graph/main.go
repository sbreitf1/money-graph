package main

import (
	"os"
	"path/filepath"

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

	var dbDir string
	if len(os.Args) > 1 {
		dbDir = os.Args[1]
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			logrus.Fatalf("failed to detect user home dir: %s", err.Error())
		}
		dbDir = filepath.Join(homeDir, ".moneydb/db")
	}
	logrus.Infof("use %q as database dir", dbDir)

	db, err := moneydb.OpenOrCreateInFolder(".runtime/db", "New Money Database")
	if err != nil {
		logrus.Fatalf("could not open or create db: %s", err.Error())
	}

	//fmt.Println(db.ImportCSV("testing/csv-examples/private-20220215-126303197-umsatz.CSV"))

	if err := runServer("localhost:8081", db); err != nil {
		logrus.Fatalf("failed to run server: %s", err.Error())
	}

	logrus.Info("application shut down")
}
