package provider

type Provider interface {
	Parse(inputCsv string) *Csv
	GenerateKeywords(inputCsv string) []string
}
