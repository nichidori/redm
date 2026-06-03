package main

import (
	"fmt"
	"strings"
)

type Column struct {
	Header string
	Width  int
	Value  func(int) string
}

func NewColumn(header string, width int, value func(int) string) Column {
	return Column{Header: header, Width: width, Value: value}
}

func PrintTable(cols []Column, n int) {
	for i, col := range cols {
		if i > 0 {
			fmt.Print(" | ")
		}
		fmt.Print(FixLength(col.Header, col.Width))
	}
	fmt.Println()

	for i, col := range cols {
		if i > 0 {
			fmt.Print("-+-")
		}
		fmt.Print(strings.Repeat("-", col.Width))
	}
	fmt.Println()

	for row := range n {
		for i, col := range cols {
			if i > 0 {
				fmt.Print(" | ")
			}
			fmt.Print(FixLength(col.Value(row), col.Width))
		}
		fmt.Println()
	}
}
