// This script converts the Freebase database of notable people from
// Excel format to CSV, json, and gob formats.  The original data are
// available from here:
//
// https://science.sciencemag.org/content/345/6196/558/tab-figures-data
//
// The relevant file is under the link reading "Data S1".
//
// This script will store the rows of the dataset as arrays of strings.
// See convert_structs.go and convert_structs_cols.go to convert to
// alternative formats.
//
// To obtain the dependencies for this script, run the following:
//  go get github.com/kshedden/godata_workshop/notable/notable
//  go get github.com/tealeg/xlsx
package main

import (
	"fmt"

	"github.com/kshedden/godata_workshop/notable/notable"
	"github.com/tealeg/xlsx"
)

// saveSheetCSV saves the contents of an Excel sheet to the
// named file, in gzip-compressed text/csv format.
func saveSheetCSV(sheet *xlsx.Sheet, fname string) {

	f, g, cout := notable.GetCSVWriter(fname)
	defer f.Close()
	defer g.Close()
	defer cout.Flush()

	// The dimensions of the sheet
	r, c := sheet.MaxRow, sheet.MaxCol

	// Storage for the data in one row
	trow := make([]string, c)

	// Loop over the rows of the Excel sheet
	for i := 0; i < r; i++ {

		// Read one row from the Excel sheet
		row, err := sheet.Row(i)
		if err != nil {
			panic(err)
		}

		// Copy the contents of the Excel cells into
		// an array of strings
		for j := 0; j < c; j++ {
			trow[j] = row.GetCell(j).Value
		}

		// Save one row to the output file in CSV format.
		if err := cout.Write(trow); err != nil {
			panic(err)
		}
	}
}

// saveSheetJSON saves the contents of an Excel sheet to the
// named file, in gzip-compressed json format.
func saveSheetJSON(sheet *xlsx.Sheet, fname string) {

	f, g, enc := notable.GetJSONEncoder(fname)
	defer f.Close()
	defer g.Close()

	// The dimensions of the sheet
	r, c := sheet.MaxRow, sheet.MaxCol

	// Storage for the data in one row
	trow := make([]string, c)

	// Loop over the rows of the Excel sheet
	for i := 0; i < r; i++ {

		// Read one row from the Excel sheet
		row, err := sheet.Row(i)
		if err != nil {
			panic(err)
		}

		// Copy the contents of the Excel cells into
		// an array of strings
		for j := 0; j < c; j++ {
			trow[j] = row.GetCell(j).Value
		}

		// Save one row to the output file in CSV format.
		if err := enc.Encode(trow); err != nil {
			panic(err)
		}
	}
}

// saveSheetGob saves the contents of an Excel sheet to the
// named file, in gzip-compressed json format.
func saveSheetGob(sheet *xlsx.Sheet, fname string) {

	f, g, enc := notable.GetGobEncoder(fname)
	defer f.Close()
	defer g.Close()

	// The dimensions of the sheet
	r, c := sheet.MaxRow, sheet.MaxCol

	// Storage for the data in one row
	trow := make([]string, c)

	// Loop over the rows of the Excel sheet
	for i := 0; i < r; i++ {

		// Read one row from the Excel sheet
		row, err := sheet.Row(i)
		if err != nil {
			panic(err)
		}

		// Copy the contents of the Excel cells into
		// an array of strings
		for j := 0; j < c; j++ {
			trow[j] = row.GetCell(j).Value
		}

		// Save one row to the output file in CSV format.
		if err := enc.Encode(trow); err != nil {
			panic(err)
		}
	}
}

func main() {

	// Read the Excel file
	f, err := xlsx.OpenFile("SchichDataS1_FB.xlsx")
	if err != nil {
		panic(err)
	}

	// Print out the names of all work sheets
	fmt.Printf("Sheets:\n")
	for j, sheet := range f.Sheets {
		fmt.Printf("%4d %s\n", j, sheet.Name)
	}

	// Get the only sheet in the workbook
	fb := f.Sheet["FB"]

	// Save it in compressed csv format
	saveSheetCSV(fb, "fb.csv.gz")

	// Save it in compressed json format
	saveSheetJSON(fb, "fb.json.gz")

	// Save it in compressed gob format
	saveSheetGob(fb, "fb.gob.gz")
}
