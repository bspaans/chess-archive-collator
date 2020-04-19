package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/freeeve/pgn"
)

type MoveTree struct {
	Move       string
	Annotation string
	Replies    map[string]*MoveTree
	Parent     *MoveTree
	Games      []*pgn.Game
}

func NewMoveTree(move, annotation string) *MoveTree {
	return &MoveTree{
		Move:       move,
		Annotation: annotation,
		Replies:    map[string]*MoveTree{},
		Games:      []*pgn.Game{},
	}
}

func (m *MoveTree) ClassifyGame(game *pgn.Game) string {
	b := pgn.NewBoard()
	tree := m
	fmt.Println("Classifying new game")
	for _, move := range game.Moves {
		fmt.Println(move)
		// make the move on the board
		b.MakeMove(move)

		next, found := tree.Replies[move.String()]
		if found {
			fmt.Println("Book move", next.Annotation)
			tree = next
			tree.Games = append(tree.Games, game)
		} else {
			break
		}
	}
	annotation := tree.Annotation
	for annotation == "" {
		tree = tree.Parent
		annotation = tree.Annotation
	}
	return annotation
}

func (m *MoveTree) GetOrInsertMove(move string) *MoveTree {
	if t, ok := m.Replies[move]; ok {
		return t
	}
	fmt.Println("Inserting at", move)
	tree := NewMoveTree(move, "")
	m.Replies[move] = tree
	tree.Parent = m
	return tree
}

func (m *MoveTree) PruneGameLessBranches() *MoveTree {
	result := NewMoveTree(m.Move, m.Annotation)
	result.Games = m.Games
	for move, replyTree := range m.Replies {
		if len(replyTree.Games) > 0 {
			result.Replies[move] = replyTree.PruneGameLessBranches()
			result.Replies[move].Parent = result
		}
	}
	return result
}

func (m *MoveTree) String() string {
	indent := func(s string) string {
		lines := strings.Split(s, "\n")
		for i, line := range lines {
			if line != "" {
				lines[i] = "  " + line
			}
		}
		return strings.Join(lines, "\n")
	}
	result := fmt.Sprintf("%s [%s] (%d)\n", m.Move, m.Annotation, len(m.Games))
	for _, tree := range m.Replies {
		result += indent(tree.String())
	}
	return result
}

func (m *MoveTree) AddNodeForPGN(eco, annotation, pgnStr string) {
	f := bytes.NewReader([]byte(strings.TrimSpace(pgnStr)))
	ps := pgn.NewPGNScanner(f)
	for ps.Next() {
		game, err := ps.Scan()
		if err != nil {
			panic(err)
		}
		fmt.Println(game)
		tree := m
		for _, move := range game.Moves {
			tree = tree.GetOrInsertMove(move.String())
			fmt.Println(move)
		}
		if tree.Annotation == "" {
			tree.Annotation = annotation
		} else {
			//panic("Already annotated: " + tree.Annotation + ", " + annotation)
		}
	}
}

func ParseECOClassificationIntoTree() (*MoveTree, error) {

	root := NewMoveTree("", "Start position")
	content, err := ioutil.ReadFile("scid.eco")
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(content), "\n")
	inDefinition := false
	eco := ""
	annotation := ""
	pgn := ""
	for _, line := range lines {
		if strings.HasPrefix(line, "#") {
			continue
		}
		if line == "" {
			continue
		}
		if inDefinition && !strings.HasPrefix(line, " ") {
			inDefinition = false
			pgn = ""
			root.AddNodeForPGN(eco, annotation, pgn)
		}
		if !inDefinition {
			words := strings.Split(line, " ")
			eco = words[0]
			rest := strings.Join(words[1:], " ")
			inDefinition = true

			if strings.HasSuffix(rest, `"`) {
				annotation = rest
				pgn = ""
			} else {
				parts := strings.Split(rest, `"`)
				annotation = parts[1]
				pgn = parts[2]
				root.AddNodeForPGN(eco, annotation, pgn)
			}
		} else {
			pgn += line
		}
	}
	root.AddNodeForPGN(eco, annotation, pgn)
	return root, nil
}
