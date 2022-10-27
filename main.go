package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// Record
type Record struct {
	Name   string
	Price  int
	Rating int
}

func (r *Record) ToString() string {
	return r.Name + "," + strconv.Itoa(r.Price) + "," + strconv.Itoa(r.Rating)
}

// search: читаем построчно, запоминаем и обновляем нужное
func search(inputPath string) (Record, Record, error) {
	var mPriceRecord, mRatingRecord Record

	file, err := os.Open(inputPath)
	if err != nil {
		return Record{}, Record{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		record, err := readRecord(scanner.Text())
		if err != nil {
			continue
		}

		if mostPrice(mPriceRecord, record) {
			mPriceRecord = record
		}

		if mostRating(mPriceRecord, record) {
			mRatingRecord = record
		}
	}

	if err := scanner.Err(); err != nil {
		return Record{}, Record{}, err
	}

	return mPriceRecord, mRatingRecord, nil
}

func readRecord(rec string) (Record, error) {
	rawFields := strings.Split(rec, ",")
	if len(rawFields) < 3 {
		return Record{}, fmt.Errorf("Record " + rec + " must have 3 fields")
	}

	price, err := strconv.Atoi(rawFields[1])
	if err != nil {
		return Record{}, fmt.Errorf("Failed to read Price field: %w", err)
	}
	rating, err := strconv.Atoi(rawFields[2])
	if err != nil {
		return Record{}, fmt.Errorf("Failed to read Rating field: %w", err)
	}

	return Record{rawFields[0], price, rating}, nil
}

func mostPrice(old Record, new Record) bool { return old.Price < new.Price }

func mostRating(old Record, new Record) bool { return old.Rating < new.Rating }

func write(outputPath string, mPriceRecord, mRatingRecord Record) error {

	file, err := json.MarshalIndent([]Record{mPriceRecord, mRatingRecord}, "", " ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(outputPath, file, 0644)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	mPriceRecord, mRatingRecord, err := search("./data/db.csv")
	if err != nil {
		fmt.Println("Failed to read input file: %w", err)
	}

	err = write("./data/db.json", mPriceRecord, mRatingRecord)
	if err != nil {
		fmt.Println("Failed to write output file: %w", err)
	}
}
