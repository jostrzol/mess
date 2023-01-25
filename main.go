package main

import (
	"flag"
	"log"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
)

type Config struct {
	Board        BoardConfig        `hcl:"board,block"`
	Pieces       PiecesConfig       `hcl:"pieces,block"`
	InitialState InitialStateConfig `hcl:"initial_state,block"`
}

type BoardConfig struct {
	Height uint `hcl:"height"`
	Width  uint `hcl:"width"`
}

type PiecesConfig struct {
	Pieces []PieceConfig `hcl:"piece,block"`
}

type PieceConfig struct {
	Name string `hcl:"piece_name,label"`
}

type InitialStateConfig struct {
	PiecePlacements []PiecePlacementConfig `hcl:"piece_placement,block"`
}

type PiecePlacementConfig struct {
	PlayerName string         `hcl:"player_name,label"`
	Placements hcl.Attributes `hcl:",remain"`
}

func main() {
	var configFilepath = flag.String("rules", "./rules.hcl", "path to a rules config file")
	flag.Parse()

	var config Config
	err := hclsimple.DecodeFile(*configFilepath, nil, &config)
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}
	log.Printf("Configuration is %#v", config)
}
