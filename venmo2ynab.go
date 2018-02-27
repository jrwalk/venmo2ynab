package main

import (
	"fmt"
	"flag"
	"time"
	"os"
	"encoding/csv"
	"io"
    "path/filepath"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Record struct {
	// ` ID`    string
	Date           string
	Type           string
	// Status      string
	Note           string
	From           string
	To             string
	Amount         string
	// `Amount (fees)`     string
	// Source         interface{}
	// Destination    interface{}
}

type CleanRecord struct {
	Date     string
	Payee    string
	Category string
	Memo     string
	Outflow  string
	Inflow   string
}

func buildRecord(row []string) (rec Record) {
	rec = Record{
		Date:        ripDate(row[1]),
        Type:        string(row[2]),
		Note:        string(row[4]),
		From:        string(row[5]),
		To:          string(row[6]),
		Amount:      string(row[7]),
	}
	return
}

func ripDate(dt string) string {
	const refTime = "2006-01-02T15:04:05"
	timestamp, _ := time.Parse(refTime, dt)
	return fmt.Sprintf("%04d-%02d-%02d", 
                       timestamp.Year(), 
                       timestamp.Month(), 
                       timestamp.Day())
}

func mungeTransactions(rec Record) CleanRecord {

    isPayment := (rec.Type == "Payment")
    isOutflow := (fmt.Sprintf("%c", rec.Amount[0]) == "-")

    payee := ""
    outflow := ""
    inflow := ""
    if isPayment && isOutflow {
        payee = rec.To
        outflow = rec.Amount[2:]
    } else if isPayment && !isOutflow {
        payee = rec.From
        inflow = rec.Amount[2:]
    } else if !isPayment && isOutflow {
        payee = rec.From
        outflow = rec.Amount[2:]
    } else {
        payee = rec.To
        inflow = rec.Amount[2:]
    }


    return CleanRecord{
        Date:       rec.Date,
        Payee:      payee,
        Category:   "",
        Memo:       rec.Note,
        Outflow:    outflow,
        Inflow:     inflow,
    }

}

func recToList(rec CleanRecord) []string {
    vals := make([]string, 6)
    vals[0] = rec.Date
    vals[1] = rec.Payee
    vals[2] = rec.Category
    vals[3] = rec.Memo
    vals[4] = rec.Outflow
    vals[5] = rec.Inflow
    return vals
}

func main() {

    dirPtr := flag.String("dir", "./", "working directory")
    inPtr := flag.String("inFile", "", "input file")
    outPtr := flag.String("outFile", "", "output file")

    flag.Parse()

	inPath := filepath.Join(*dirPtr, *inPtr)
    outPath := filepath.Join(*dirPtr, *outPtr)
    file, err := os.Open(inPath)
	check(err)
    defer file.Close()
    
	r := csv.NewReader(file)
	_, err = r.Read()  // clear header row
	check(err)

    var records [][]string
    headers := []string{
        "Date",
        "Payee",
        "Category",
        "Memo",
        "Outflow",
        "Inflow",
    }
    records = append(records, headers)

	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		rec := buildRecord(row)
		cleanRec := mungeTransactions(rec)
        output := recToList(cleanRec)

        records = append(records, output)
	}

    outFile, err := os.Create(outPath)
    check(err)

    w := csv.NewWriter(outFile)
    for _, record := range records {
        if err := w.Write(record); err != nil {
            panic(err)
        }
    }
    w.Flush()
}
