package dbs

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/madeadi/bankparser/provider"
)

type dbs struct {
	rows []provider.Row
}

func NewDbs() provider.Provider {
	return dbs{
		rows: []provider.Row{},
	}
}

func (d dbs) Parse(input string) *provider.Csv {
	f, err := os.Open(input)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fs := bufio.NewScanner(f)
	fs.Split(bufio.ScanLines)

	row := provider.Row{}
	d.resetAndAppendRow(&row)
	var prevBalance float64

	for fs.Scan() {
		text := fs.Text()
		if strings.Contains(text, "Balance Brought Forward") {
			balance := strings.ReplaceAll(text, "Balance Brought Forward", "")
			prevBalance, err = parseFloatMoney(balance)
			if err != nil {
				panic(err)
			}
			d.resetAndAppendRow(&row)
			continue
		}
		if strings.Contains(text, "Date, Description") {
			d.resetAndAppendRow(&row)
			continue
		}

		tokenArr := strings.Split(text, " ")

		// parsing the date and description
		// 01/09/2022 Monthly Savings Amount for MySavings/POSB
		// parsing date, e.g. 01/01/2020
		if len(tokenArr) > 2 {
			dateArr := strings.Split(tokenArr[0], "/")
			if len(dateArr) == 3 {
				// that means, this is the date.
				// Convert it to Y-m-d format
				row.Date = fmt.Sprintf("%s-%s-%s", dateArr[2], dateArr[1], dateArr[0])

				// parsing description
				// e.g. Monthly Savings Amount for MySavings/POSB
				row.Description = strings.Join(tokenArr[1:], " ")
				continue
			}
		}

		// parsing the amount. The format is, "amount, balance". However, we don't know
		// if amount is positive or negative. We will have to check the balance.
		// e.g. 30.35 1,284.62
		if len(tokenArr) == 2 {
			// parsing amount
			amount, err := parseFloatMoney(tokenArr[0])

			if err == nil {
				// parsing balance
				balance, err := parseFloatMoney(tokenArr[1])
				if err == nil {
					if prevBalance < balance {
						// this means that the amount is positive
						row.In = amount
					} else {
						// this means that the amount is negative
						row.Out = amount
					}
					prevBalance = balance

					// we are done with a single transaction. Therefore,
					// store the row into array and continue with next row
					d.resetAndAppendRow(&row)
					continue
				}
			}
		}

		// else, this is a description
		row.Description += ", " + text
	}

	pc := provider.Csv{
		Rows: d.rows,
	}

	fmt.Println("Total rows:", len(pc.Rows))
	return &pc
}

func (d dbs) GenerateKeywords(inputCsv string) []string {
	return []string{}
}

func (d *dbs) resetAndAppendRow(row *provider.Row) {
	// don't append if it's empty
	if row.Date != "" && (row.In != 0 || row.Out != 0) {
		d.rows = append(d.rows, *row)
	}

	*row = provider.Row{Account: "dbs"}
}

func parseFloatMoney(val string) (float64, error) {
	val = strings.Trim(val, " ")
	val = strings.ReplaceAll(val, ",", "")
	amount, err := strconv.ParseFloat(val, 64)
	return amount, err
}
