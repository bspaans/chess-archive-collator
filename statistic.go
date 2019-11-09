package main

import (
	"fmt"
	"strings"
)

type Statistic struct {
	TotalPlayed int
	TotalWon    int
	TotalLost   int
	TotalDrawn  int

	Played map[bool]int
	Won    map[bool]int
	Lost   map[bool]int
	Drawn  map[bool]int
}

func NewStatistic() *Statistic {
	return &Statistic{
		Played: map[bool]int{},
		Won:    map[bool]int{},
		Lost:   map[bool]int{},
		Drawn:  map[bool]int{},
	}
}

func (s *Statistic) Count(white bool, result string) {
	s.TotalPlayed += 1
	s.Played[white] += 1
	if result == "1-0" && white || result == "0-1" && !white {
		s.TotalWon += 1
		s.Won[white] += 1
	} else if result == "0-1" && white || result == "1-0" && !white {
		s.TotalLost += 1
		s.Lost[white] += 1
	} else {
		s.TotalDrawn += 1
		s.Drawn[white] += 1
	}
}

func (s *Statistic) Header() string {
	return "Games\tWhite\tBlack\tWon\tLost\tDrawn\tWon(W)\tWon(B)\tLost(W)\tLost(B)\tDraw(W)\tDraw(B)"
}

func (s *Statistic) Headers() []string {
	return []string{"Games", "White", "Black", "Won", "Lost", "Drawn", "Won(W)", "Won(B)", "Lost(W)", "Lost(B)", "Draw(W)", "Draw(B)"}
	return strings.Split(s.Header(), "\t")
}
func (s Statistic) Data() []string {
	percentage := func(c, d int) string {
		if d == 0 || c == 0 {
			return ""
		}
		return fmt.Sprintf("%3d (%0.f%%)", c, float64(c)/float64(d)*100)
	}
	return []string{
		fmt.Sprintf("%d", s.TotalPlayed),
		percentage(s.Played[true], s.TotalPlayed),
		percentage(s.Played[false], s.TotalPlayed),
		percentage(s.TotalWon, s.TotalPlayed),
		percentage(s.TotalLost, s.TotalPlayed),
		percentage(s.TotalDrawn, s.TotalPlayed),
		percentage(s.Won[true], s.Played[true]),
		percentage(s.Won[false], s.Played[false]),
		percentage(s.Lost[true], s.Played[true]),
		percentage(s.Lost[false], s.Played[false]),
		percentage(s.Drawn[true], s.Played[true]),
		percentage(s.Drawn[false], s.Played[false]),
	}

}
