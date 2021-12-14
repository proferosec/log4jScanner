package cmd

import (
	"encoding/csv"
	"fmt"
	"github.com/pterm/pterm"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

var csvPath string

func checkCsvPath() {
	if csvPath == "" {
		csvPath = "log4jScanner-results.csv"
	}

	if !strings.HasSuffix(strings.ToLower(csvPath), ".csv") {
		fmt.Println("csv-output path is not a CSV file. Output will be saved to running folder")
		csvPath = "log4jScanner-results.csv"
	}
}

func createCsvRecords(resChan chan string) [][]string {
	csvRecords := [][]string{
		{"type", "ip", "port", "status_code"},
	}
	for res := range resChan {
		csvRes := strings.Split(res, ",")
		msg := fmt.Sprintf("Summary: %s:%s ==> %s", csvRes[1], csvRes[2], csvRes[3])
		pterm.Info.Println(msg)
		log.Info(msg)
		csvRecords = append(csvRecords, csvRes)
	}

	close(TCPServer.sChan)

	for suc := range TCPServer.sChan {
		csvSuc := strings.Split(suc, ",")
		msg := fmt.Sprintf("Summary: Callback from %s:%s", csvSuc[1], csvSuc[2])
		pterm.Info.Println(msg)
		log.Info(msg)
		csvRecords = append(csvRecords, csvSuc)
	}

	return csvRecords
}

func saveCSV(csvRecords [][]string) {
	checkCsvPath()

	f, err := os.Create(csvPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	for _, record := range csvRecords {
		if err = w.Write(record); err != nil {
			log.Fatal(err)
		}
	}
}
