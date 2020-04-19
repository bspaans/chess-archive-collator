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
var Order = flag.String("order", "opening", "Order rows. One of: opening, played, played-white, played-black, won, lost, drawn, won-white, won-black, lost-white, lost-black, drawn-white, drawn-black")

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

func (r *Report) Count(openingTree *MoveTree, game *pgn.Game) {

	if game.Tags["White"] != *Player && game.Tags["Black"] != *Player {
		fmt.Printf("Skipping game, because player '%s' wasn't playing (NB. you can set the player username with --player)\n", *Player)
		return
	}

	playingWithWhitePieces := game.Tags["White"] == *Player
	gameResult := game.Tags["Result"]
	r.Statistic.Count(playingWithWhitePieces, gameResult)

	opening := openingTree.ClassifyGame(game)
	openingFound := opening != ""
	if openingFound {
		r.CountOpening(playingWithWhitePieces, gameResult, opening, game)
	}
	/*
		if game.Tags["ECO"] != "" {
			r.CountOpening(playingWithWhitePieces, gameResult, game.Tags["ECO"], game)
			openingFound = true
		}
	*/
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
	data := [][]string{}
	for opening, stats := range r.OpeningStats {
		row := append([]string{LookupECO(opening)}, stats.Data()...)
		data = append(data, row)

	}
	sort.Slice(data, func(i, j int) bool {
		if *Order == "opening" {
			return data[i][0] < data[j][0]
		} else if *Order == "played" && data[i][1] != data[j][1] {
			return data[i][1] < data[j][1]
		} else if *Order == "played-white" && data[i][2] != data[j][2] {
			return data[i][2] < data[j][2]
		} else if *Order == "played-black" && data[i][3] != data[j][3] {
			return data[i][3] < data[j][3]
		} else if *Order == "won" && data[i][4] != data[j][4] {
			return data[i][4] < data[j][4]
		} else if *Order == "lost" && data[i][5] != data[j][5] {
			return data[i][5] < data[j][5]
		} else if *Order == "drawn" && data[i][6] != data[j][6] {
			return data[i][6] < data[j][6]
		} else if *Order == "won-white" && data[i][7] != data[j][7] {
			return data[i][7] < data[j][7]
		} else if *Order == "won-black" && data[i][8] != data[j][8] {
			return data[i][8] < data[j][8]
		} else if *Order == "lost-white" && data[i][9] != data[j][9] {
			return data[i][9] < data[j][9]
		} else if *Order == "lost-black" && data[i][10] != data[j][10] {
			return data[i][10] < data[j][10]
		} else if *Order == "drawn-white" && data[i][11] != data[j][11] {
			return data[i][11] < data[j][11]
		} else if *Order == "drawn-black" && data[i][12] != data[j][12] {
			return data[i][12] < data[j][12]
		}
		return data[i][0] < data[j][0]
	})
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

	openingTree, err := ParseECOClassificationIntoTree()
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
			report.Count(openingTree, game)
		}
	}
	fmt.Println(report)
	fmt.Println(report.Statistic.Header())
	fmt.Println(report.Statistic)

	fmt.Println(openingTree.PruneGameLessBranches())
}
