/*
Copyright 2019, Bart Spaans

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/freeeve/pgn"
	"github.com/olekukonko/tablewriter"
)

var Player = flag.String("player", "bartspaans", "The player's name.")

// TODO: classify openings
// TODO: build move tree for white
// TODO: build move tree for black

type Report struct {
	Openings     map[string][]*pgn.Game
	OpeningStats map[string]*Statistic
	Statistic    *Statistic
}

func NewReport() *Report {
	return &Report{
		Openings:     map[string][]*pgn.Game{},
		OpeningStats: map[string]*Statistic{},
		Statistic:    NewStatistic(),
	}
}

func (r *Report) Count(game *pgn.Game) {

	if game.Tags["White"] != *Player && game.Tags["Black"] != *Player {
		fmt.Printf("Skipping game, because player '%s' wasn't playing (NB. you can set the player username with --player)\n", *Player)
		return
	}

	playingWithWhitePieces := game.Tags["White"] == *Player
	gameResult := game.Tags["Result"]
	r.Statistic.Count(playingWithWhitePieces, gameResult)

	openingFound := false
	if game.Tags["ECO"] != "" {
		r.CountOpening(playingWithWhitePieces, gameResult, game.Tags["ECO"], game)
		openingFound = true
	}
	if !openingFound {
		r.CountOpening(playingWithWhitePieces, gameResult, "Unknown opening", game)

		fmt.Println("Unknown opening: ")
		b := pgn.NewBoard()
		for _, move := range game.Moves {
			// make the move on the board
			fmt.Println(move)
			b.MakeMove(move)
			// print out FEN for each move in the game
			fmt.Println(b)
		}
	}
}

func (r *Report) CountOpening(white bool, gameResult, opening string, game *pgn.Game) {
	if _, ok := r.Openings[opening]; !ok {
		r.Openings[opening] = []*pgn.Game{}
		r.OpeningStats[opening] = NewStatistic()
	}
	r.Openings[opening] = append(r.Openings[opening], game)
	r.OpeningStats[opening].Count(white, gameResult)
}

func (r *Report) String() string {
	keys := []string{}
	for key, _ := range r.Openings {
		keys = append(keys, key)

	}
	data := [][]string{}
	sort.Strings(keys)
	for _, opening := range keys {
		stats, ok := r.OpeningStats[opening]
		if ok {
			row := append([]string{LookupECO(opening)}, stats.Data()...)
			data = append(data, row)
		}
	}
	b := bytes.NewBuffer([]byte{})
	table := tablewriter.NewWriter(b)
	table.SetHeader(append([]string{"Opening"}, r.Statistic.Headers()...))
	table.AppendBulk(data)
	//table.SetAutoWrapText(false)
	table.SetRowLine(true)
	colors := []tablewriter.Colors{
		tablewriter.Colors{tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.FgWhiteColor},
		tablewriter.Colors{tablewriter.FgMagentaColor},

		tablewriter.Colors{tablewriter.FgHiGreenColor},
		tablewriter.Colors{tablewriter.FgRedColor},
		tablewriter.Colors{tablewriter.FgBlueColor},

		tablewriter.Colors{tablewriter.FgHiGreenColor},
		tablewriter.Colors{tablewriter.FgHiGreenColor},

		tablewriter.Colors{tablewriter.FgRedColor},
		tablewriter.Colors{tablewriter.FgRedColor},

		tablewriter.Colors{tablewriter.FgBlueColor},
		tablewriter.Colors{tablewriter.FgBlueColor},
	}
	table.SetHeaderColor(colors...)
	table.SetColumnColor(colors...)
	table.Render()
	return string(b.Bytes())
}

var Openings = map[string]string{
	"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR":          "King's Pawn Opening",
	"rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR":        "King's Pawn game",
	"rnbqkbnr/pppp1ppp/8/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R":      "King's Knight Opening",
	"rnbqkbnr/pppp1ppp/8/4p3/2B1P3/8/PPPP1PPP/RNBQK1NR":      "Bishop's Opening",
	"r1bqkbnr/pppp1ppp/2n5/4p3/2B1P3/5N2/PPPP1PPP/RNBQK2R":   "Italian Game",
	"r1bqkbnr/pppp1ppp/2n5/1b2p3/2B1P3/5N2/PPPP1PPP/RNBQK2R": "Ruy Lopez",
	"rnbqkbnr/pppp1ppp/8/4p3/4PP2/8/PPPP2PP/RNBQKBNR":        "King's Gambit",
	"rnbqkbnr/pppp1ppp/8/8/4Pp2/8/PPPP2PP/RNBQKBNR":          "King's Gambit Accepted",
	"rnbqkb1r/pppp1ppp/5n2/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R":    "Petrov",
	"rnbqkbnr/pppp1ppp/4p3/8/4P3/8/PPPP1PPP/RNBQKBNR":        "French",
	"r1bqkbnr/pppppppp/2n5/8/4P3/8/PPPP1PPP/RNBQKBNR":        "Nimzowitsch Defence",
	"rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR":        "Sicilian",
	"rnbqkbnr/pp1ppppp/2p5/8/4P3/8/PPPP1PPP/RNBQKBNR":        "Caro Kann",
	"rnbqkbnr/pppppp1p/6p1/8/4P3/8/PPPP1PPP/RNBQKBNR":        "Modern Defense",
	"rnbqkbnr/pppppppp/8/8/8/1P6/P1PPPPPP/RNBQKBNR":          "Larsen's Opening",
	"rnbqkbnr/pppppppp/8/8/6P1/8/PPPPPP1P/RNBQKBNR":          "Grob's Attack",
	"rnbqkbnr/pppppppp/8/8/1P6/8/P1PPPPPP/RNBQKBNR":          "Polish Opening",
	"rnbqkbnr/pppppppp/8/8/3P4/8/PPP1PPPP/RNBQKBNR":          "Queen's Pawn Opening",
}

func main() {
	flag.Parse()

	_, err := ParseECOClassificationIntoTree()
	if err != nil {
		panic(err)
	}

	report := NewReport()
	for _, arg := range flag.Args() {
		fmt.Println("Processing", arg)
		f, err := os.Open(arg)
		if err != nil {
			log.Fatal(err)
		}
		ps := pgn.NewPGNScanner(f)
		// while there's more to read in the file
		for ps.Next() {
			// scan the next game
			game, err := ps.Scan()
			if err != nil {
				log.Fatal(err)
			}
			report.Count(game)
		}
	}
	fmt.Println(report)
	fmt.Println(report.Statistic.Header())
	fmt.Println(report.Statistic)
}
