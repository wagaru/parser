package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Printf("Usage:\n\n\t parser.exe [targetDirectory]\n")
		return
	}
	targetDir := args[0]
	fmt.Printf("Parse folder %s\n", targetDir)

	files, err := ioutil.ReadDir(targetDir)
	if err != nil {
		fmt.Printf("Parse folder err %v", err)
		return
	}

	resultCsvStrs := [][]string{}
	for _, f := range files {
		if f.IsDir() && strings.HasPrefix(f.Name(), "PAH") {
			innerFiles, err := ioutil.ReadDir(filepath.Join(targetDir, f.Name()))
			if err != nil {
				fmt.Printf("Parse %s failed. Skipped...\n", f.Name())
				continue
			}
			fmt.Printf("Check folder %s...", f.Name())
			for _, innerFile := range innerFiles {
				if strings.HasPrefix(innerFile.Name(), "Cap_") {
					parseResult := parseCsv(filepath.Join(targetDir, f.Name(), innerFile.Name()))
					parseResult = append([]string{f.Name()}, parseResult...)
					resultCsvStrs = append(resultCsvStrs, parseResult)
					break
				}
			}
			fmt.Printf("finished\n")
		}
	}

	fmt.Printf("Write to %s...", filepath.Join(targetDir, "result.csv"))
	resultCsvFile, err := os.Create(filepath.Join(targetDir, "result.csv"))
	if err != nil {
		fmt.Printf("Generate result.csv failed %v", err)
		return
	}
	defer resultCsvFile.Close()
	w := csv.NewWriter(resultCsvFile)
	defer w.Flush()
	for _, record := range resultCsvStrs {
		if err := w.Write(record); err != nil {
			fmt.Printf("Write result csv failed.")
			break
		}
	}
	fmt.Printf("finished\n")
}

func parseCsv(filepath string) []string {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("Open %s failed\n", filepath)
		return []string{}
	}
	defer file.Close()

	lineCount := 0
	targetLine := 400
	Index4Sum, Index5Sum, Index6Sum, Index7Sum := 0, 0, 0, 0

	r := csv.NewReader(file)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Read file %s failed", filepath)
		}
		lineCount += 1
		if lineCount > targetLine {
			break
		}

		Index4Sum += toInt(record[4])
		Index5Sum += toInt(record[5])
		Index6Sum += toInt(record[6])
		Index7Sum += toInt(record[7])
	}

	if lineCount < targetLine {
		return []string{"less than 400 lines"}
	}
	return []string{
		strconv.FormatFloat(float64(Index4Sum)/float64(targetLine), 'f', 2, 32),
		strconv.FormatFloat(float64(Index5Sum)/float64(targetLine), 'f', 2, 32),
		strconv.FormatFloat(float64(Index6Sum)/float64(targetLine), 'f', 2, 32),
		strconv.FormatFloat(float64(Index7Sum)/float64(targetLine), 'f', 2, 32),
	}
}

func toInt(str string) int {
	v, err := strconv.Atoi(str)
	if err != nil {
		fmt.Printf("Convert %s to int failed\n", str)
	}
	return v
}
