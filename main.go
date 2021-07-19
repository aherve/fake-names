package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/urfave/cli"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {

	var genNumber int
	var filename string

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
					Value:       "",
				},
				cli.IntFlag{
					Name:        "number, n",
					Usage:       "set the number of fake names you want to generate",
					Destination: &genNumber,
					Value:       10,
				},
			},
			Action: func(c *cli.Context) error {
				var names []string
				if filename != "" {
					names = readNames(filename)
				} else {
					names = frenchTownNames
				}
				g := initializeGraph(names)
				for i := 0; i < genNumber; i++ {
					fmt.Println(g.generateName())
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

func readNames(csvName string) []string {
	f, err := os.Open(csvName)
	if err != nil {
		log.Fatal(err)
	}
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	names := make([]string, len(records), len(records))
	for i := 1; i < len(records); i++ {
		names[i] = records[i][0]
	}
	return names
}
