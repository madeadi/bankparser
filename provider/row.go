package provider

import "strconv"

type Row struct {
	Account       string
	Date          string
	Description   string
	In            float64
	Out           float64
	ContraAccount string
}

func (r Row) toSlice() []string {
	return []string{
		r.Account,
		r.Date,
		r.Description,
		strconv.FormatFloat(r.In, 'f', 2, 64),
		strconv.FormatFloat(r.Out, 'f', 2, 64),
		r.ContraAccount,
	}
}
