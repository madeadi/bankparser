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
		text = strings.Trim(text, " ")
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

				// there's an edge case which is like this:
				// 31/01/2022 Interest Earned 1.58 38,335.19
				// where the last two tokens are the amount and balance.
				// if they'are token and balance, we need to pop the amount and balance
				// and pass it to the next process.
				if len(tokenArr) >= 4 {
					maybeAmount := tokenArr[len(tokenArr)-2]
					maybeBalance := tokenArr[len(tokenArr)-1]
					in, out, ok := parseInOut(maybeAmount, maybeBalance, &prevBalance)
					if ok {
						row.In = in
						row.Out = out
						row.Description = strings.Join(tokenArr[1:len(tokenArr)-2], " ")
						d.resetAndAppendRow(&row)
						continue
					}

				}

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
			var ok bool
			row.In, row.Out, ok = parseInOut(tokenArr[0], tokenArr[1], &prevBalance)
			if ok {
				// we are done with a single transaction. Therefore,
				// store the row into array and continue with next row
				prevBalance, _ = parseFloatMoney(tokenArr[1])
				d.resetAndAppendRow(&row)
				continue
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

/// A helper function to parse the amount and balance
/// and update the previous balance
func parseInOut(amountStr string, balanceStr string, prevBalance *float64) (in float64, out float64, ok bool) {
	// parsing amount
	amount, err := parseFloatMoney(amountStr)

	if err == nil {
		// parsing balance
		balance, err := parseFloatMoney(balanceStr)
		if err == nil {
			if *prevBalance < balance {
				// this means that the amount is positive
				in = amount
			} else {
				// this means that the amount is negative
				out = amount
			}
			prevBalance = &balance

			return in, out, true
		}
	}

	return 0, 0, false
}
