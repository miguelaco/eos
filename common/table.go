package common

import (
	"io"

	"github.com/olekukonko/tablewriter"
)

func NewTable(writer io.Writer, header []string) *tablewriter.Table {
	table := tablewriter.NewWriter(writer)

	table.SetHeader(header)
	table.SetHeaderLine(false)
	table.SetBorder(false)
	//	table.SetRowSeparator(" ")
	table.SetColumnSeparator(" ")
	//	table.SetCenterSeparator(" ")
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoWrapText(false)

	return table
}
