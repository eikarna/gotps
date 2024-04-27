package worlds

import (
	"math/rand"
	"time"
)

type Tiles struct {
	Fg      int16
	Bg      int16
	Flags   int32
	Label   string
	IntData int32
}

type Admins struct {
	AdminUid int32
	Name     string
}

type World struct {
	Name       string
	OwnerUid   int32
	PlayersIn  int32
	SizeX      int32
	SizeY      int32
	TotalTiles int32
	Admins     []Admins
	Tiles      []Tiles
}

var (
	Worlds []World
)

func GenerateWorld(name string, sizeX int32, sizeY int32) *World {
	rand.Seed(time.Now().UnixNano())
	world := &World{}
	world.Name = name
	world.SizeX = sizeX
	world.SizeY = sizeY
	world.TotalTiles = sizeX * sizeY
	world.Tiles = make([]Tiles, 0)
	randomPosDoor := rand.Intn(int(world.TotalTiles)/(int(world.TotalTiles/100)-4) + 2)
	for i := 0; i < int(world.TotalTiles); i++ {
		tile := Tiles{}
		if i >= 2500 && i < 5400 && rand.Intn(50) == 0 {
			tile.Fg = 10
		} else if i >= 2500 && i < 5400 {
			if i > 5000 {
				if rand.Intn(8) < 3 {
					tile.Fg = 4
				} else {
					tile.Fg = 2
				}
			} else {
				tile.Fg = 2
			}
		} else if i >= 5400 {
			tile.Fg = 8
		}
		if i == 2400+randomPosDoor {
			tile.Label = "EXIT"
			tile.Fg = 6
		}
		if i == 2500+randomPosDoor {
			tile.Fg = 8
		}
		if i >= 2500 {
			tile.Bg = 14
		}
		world.Tiles = append(world.Tiles, tile)
	}
	Worlds = append(Worlds, *world)
	return world
}

func GetWorld(name string) (*World, error) {
	for _, world := range Worlds {
		if world.Name == name {
			return &world, nil
		}
	}
	return GenerateWorld(name, 100, 60), nil
}
