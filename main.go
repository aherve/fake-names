package main

import (
	"encoding/csv"
	"fake-names/graph"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {

	var genNumber int
	var filename string
	var col int
	var delimiterStr string

	app := cli.NewApp()
	app.Name = "Fake-names"
	app.Usage = "Generate fake french town names"
	app.Description = "Read a list of names, and generate fake new ones that resemble the existing names"
	app.Commands = []cli.Command{
		{
			Name:    "generate",
			Aliases: []string{"g"},
			Usage:   "generate fake names",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "input, i",
					Usage:       "provide an input csv of names to build the graph with",
					Destination: &filename,
					Value:       "names.csv",
				},
				cli.IntFlag{
					Name:        "number, n",
					Usage:       "set the number of fake names you want to generate",
					Destination: &genNumber,
					Value:       10,
				},
				cli.IntFlag{
					Name:        "col, c",
					Usage:       "Column of interest in CSV",
					Destination: &col,
					Value:       0,
				},
				cli.StringFlag{
					Name:        "delimiter, d",
					Usage:       "custom CSV comma",
					Destination: &delimiterStr,
					Value:       ",",
				},
			},
			Action: func(c *cli.Context) error {
				var names []string
				if filename != "" {
					names = readNames(filename, col, []rune(delimiterStr)[0])
				} else {
					log.Fatalf("Please provide a filename to build a graph")
				}
				g := graph.InitializeGraph(names)
				for i := 0; i < genNumber; i++ {
					fmt.Println(g.GenerateName())
				}
				return nil
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func readNames(csvName string, col int, delimiter rune) []string {
	f, err := os.Open(csvName)
	if err != nil {
		log.Fatal(err)
	}
	r := csv.NewReader(f)
	r.Comma = delimiter

	names := []string{}

	for {
		line, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if line == nil {
			break
		}
		names = append(names, line[col])
	}

	return names
}
