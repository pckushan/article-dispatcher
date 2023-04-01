package configs

import (
	"github.com/olekukonko/tablewriter"
	"gopkg.in/oleiade/reflections.v1"

	"fmt"
	"log"
	"os"
)

type Configure interface {
	Register() error
}
type Validator interface {
	Validate() error
}

type Printer interface {
	Print() interface{}
}

// Load register,validate and print all the given configurations
func Load(configs ...Configure) error {
	for _, c := range configs {
		err := c.Register()
		if err != nil {
			return err
		}

		v, ok := c.(Validator)
		if ok {
			err = v.Validate()
			if err != nil {
				return err
			}
		}

		p, ok := c.(Printer)
		if ok {
			printTable(p)
		}
	}
	return nil
}

// printTable log table writer
func printTable(p Printer) {
	table := tablewriter.NewWriter(os.Stdout)

	var data [][]string

	pr := p.Print()
	var fields []string
	fields, _ = reflections.Fields(pr)

	for _, field := range fields {
		value, err := reflections.GetField(pr, field)
		if err != nil {
			log.Printf("error printing the goconf table %s", err)
		}
		data = append(data, []string{field, fmt.Sprint(value)})
	}

	table.SetHeader([]string{"Config", "Value"})
	table.AppendBulk(data)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	table.Render()
}
