package provider

import (
	"fmt"
	"strconv"
)

// Option represents a request option generator
type Option struct {
	Key   string // The key of the option
	Value string // The value of the option
}

const (
	OptionPageNumber = "pagenum"
	OptionPageSize   = "pagesize"
	OptionID         = "id"
)

// WithPageNumber builds a string reprenting a url option for the page number of the results
func WithPageNumber(n uint) Option {
	return Option{
		Key:   OptionPageNumber,
		Value: strconv.FormatInt(n, 10),
	}
	//return Option(fmt.Sprint("page[number]=", n))
}

// WithPageSize builds a string representing a url option for the page size of the results
func WithPageSize(s uint) Option {
	return Option(fmt.Sprint("page[size]=", s))
}
