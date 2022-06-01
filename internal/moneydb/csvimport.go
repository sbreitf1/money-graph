package moneydb

import (
	"crypto/sha256"
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	headersSparkasseCSV    = []string{"Auftragskonto", "Buchungstag", "Valutadatum", "Buchungstext", "Verwendungszweck", "Glaeubiger ID", "Mandatsreferenz", "Kundenreferenz (End-to-End)", "Sammlerreferenz", "Lastschrift Ursprungsbetrag", "Auslagenersatz Ruecklastschrift", "Beguenstigter/Zahlungspflichtiger", "Kontonummer/IBAN", "BIC (SWIFT-Code)", "Betrag", "Waehrung", "Info"}
	patternSparkasseCSVDay = regexp.MustCompile(`^\s*(\d+)\.(\d+)\.(\d+)\s*$`)
)

func (db *Database) ImportCSV(path string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	tOpenFile := time.Now()

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	csvReader.Comma = ';'
	records, err := csvReader.ReadAll()
	if err != nil {
		return err
	}

	var entryParser func([][]string) ([]Entry, error)
	switch {
	case db.isSparkasseCSV(records[0]):
		entryParser = db.parseSparkasseCSVEntries
	default:
		return fmt.Errorf("unrecognized headers: %v", records[0])
	}

	entries, err := entryParser(records[1:])
	if err != nil {
		return err
	}

	logrus.Debugf("reading csv file took %v", time.Since(tOpenFile))

	return db.addEntries(entries)
}

func (db *Database) isSparkasseCSV(header []string) bool {
	if len(header) != len(headersSparkasseCSV) {
		return false
	}
	for i := range header {
		if header[i] != headersSparkasseCSV[i] {
			return false
		}
	}
	return true
}

func (db *Database) parseSparkasseCSVEntries(records [][]string) ([]Entry, error) {
	parseDay := func(str string) (time.Time, error) {
		m := patternSparkasseCSVDay.FindStringSubmatch(str)
		if len(m) != 4 {
			return time.Time{}, fmt.Errorf("invalid format")
		}
		day, _ := strconv.Atoi(m[1])
		month, _ := strconv.Atoi(m[2])
		shortYear, _ := strconv.Atoi(m[3])
		return time.Date(2000+shortYear, time.Month(month), day, 0, 0, 0, 0, time.Local), nil
	}

	// actual payload
	indexIBAN := 0
	indexDate := 1
	indexType := 3
	indexMessage := 4
	indexAmount := 14
	indexOtherName := 11
	indexOtherIBAN := 12
	// additional values for hash
	indexMandateReference := 6
	indexCustomerReference := 7
	indexCollectionReference := 8
	entries := make([]Entry, len(records))
	for i := range records {
		// compute entry hash from raw data
		strRepresentation := fmt.Sprintf("sparkasse-csv|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s",
			records[i][indexIBAN],
			records[i][indexDate],
			records[i][indexType],
			records[i][indexMessage],
			records[i][indexAmount],
			records[i][indexOtherName],
			records[i][indexOtherIBAN],
			records[i][indexMandateReference],
			records[i][indexCustomerReference],
			records[i][indexCollectionReference])
		hasher := sha256.New()
		hasher.Write([]byte(strRepresentation))
		hash := fmt.Sprintf("%x", hasher.Sum(nil))

		// extract actual payload
		iban, err := ParseIBAN(records[i][indexIBAN])
		if err != nil {
			return nil, fmt.Errorf("invalid iban %q in line %d: %s", records[i][indexIBAN], i+2, err.Error())
		}
		date, err := parseDay(records[i][indexDate])
		if err != nil {
			return nil, fmt.Errorf("invalid date %q in line %d: %s", records[i][indexDate], i+2, err.Error())
		}
		entryType := records[i][indexType]
		message := records[i][indexMessage]
		amount, err := ParseEuro(records[i][indexAmount])
		if err != nil {
			return nil, fmt.Errorf("invalid amount %q in line %d: %s", records[i][indexAmount], i+2, err.Error())
		}
		var otherName string
		var otherIBAN IBAN
		if records[i][indexOtherName] != "" && records[i][indexOtherIBAN] != "0000000000" {
			otherName = records[i][indexOtherName]
			otherIBAN, err = ParseIBAN(records[i][indexOtherIBAN])
			if err != nil {
				return nil, fmt.Errorf("invalid iban %q in line %d: %s", records[i][indexOtherIBAN], i+2, err.Error())
			}
		}

		// assemble final entry
		entries[i] = Entry{
			Hash:      hash,
			IBAN:      iban,
			Date:      date,
			Type:      entryType,
			Message:   message,
			Amount:    amount,
			OtherName: otherName,
			OtherIBAN: otherIBAN,
		}
	}

	return entries, nil
}

func (db *Database) addEntries(entries []Entry) error {
	tBegin := time.Now()

	chunks := make(map[chunkDescriptor]*chunk)
	addedCount := 0

	for i, e := range entries {
		desc := e.ChunkDescriptor()

		chunk, ok := chunks[desc]
		if !ok {
			var err error
			chunk, err = db.loadChunk(desc)
			if err != nil {
				return fmt.Errorf("load chunk %v: %s", desc, err.Error())
			}
			chunks[desc] = chunk
		}

		added, err := chunk.AddEntry(e)
		if err != nil {
			return fmt.Errorf("add entry %d to chunk %v: %s", i, desc, err.Error())
		}
		if added {
			addedCount++
		}
	}

	logrus.Infof("added %d of %d new entries", addedCount, len(entries))

	for _, c := range chunks {
		if err := db.saveChunk(c); err != nil {
			return fmt.Errorf("save chunk %v: %s", c.desc, err.Error())
		}
	}

	logrus.Debugf("adding and writing entries took %v", time.Since(tBegin))
	return nil
}
