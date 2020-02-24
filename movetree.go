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
}

func NewMoveTree(move, annotation string) *MoveTree {
	return &MoveTree{
		Move:       move,
		Annotation: annotation,
		Replies:    map[string]*MoveTree{},
	}
}

func (m *MoveTree) GetOrInsertMove(move string) *MoveTree {
	if t, ok := m.Replies[move]; ok {
		return t
	}
	fmt.Println("Inserting at", move)
	tree := NewMoveTree(move, "")
	m.Replies[move] = tree
	return tree
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
