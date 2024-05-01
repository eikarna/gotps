package worlds

import (
	"errors"
	//	"fmt"
	"github.com/bvinc/go-sqlite-lite/sqlite3"
	"github.com/codecat/go-libs/log"
	"github.com/vmihailenco/msgpack/v5"
	"math/rand"
	//	"reflect"
	//	"strings"
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
	OwnerName  string
	PlayersIn  int16
	SizeX      int32
	SizeY      int32
	TotalTiles int32
	Admins     []Admins
	Tiles      []Tiles
}

var Worlds = make(map[string]World)

/* autoTagMsgpackStruct automatically adds msgpack struct tags to fields of a struct
func autoTagMsgpackStruct(s interface{}) interface{} {
	typ := reflect.TypeOf(s)
	if typ.Kind() != reflect.Struct {
		return nil
	}

	resultTyp := reflect.StructOf([]reflect.StructField{})
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("msgpack")
		if tag == "" {
			fieldName := field.Name
			tagValue := fieldName
			tagValue = cleanTag(tagValue)
			tagValue = strings.ToLower(tagValue)
			tagValue = strings.ReplaceAll(tagValue, " ", "")
			field.Tag = reflect.StructTag(fmt.Sprintf(`msgpack:"%s"`, tagValue))
		}
		resultTyp = reflect.Append(resultTyp, field)
	}

	return reflect.New(resultTyp).Elem().Interface()
}*/

// autoTagMsgpackStruct automatically adds msgpack struct tags to fields of a struct
/*func AutoTagMsgpackStruct(s World) World {
	typ := reflect.TypeOf(s)
	if typ.Kind() != reflect.Struct {
		return World{}
	}

	resultFields := make([]reflect.StructField, 0)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("msgpack")
		if tag == "" {
			fieldName := field.Name
			tagValue := fieldName
			tagValue = CleanTag(tagValue)
			tagValue = strings.ToLower(tagValue)
			tagValue = strings.ReplaceAll(tagValue, " ", "")
			field.Tag = reflect.StructTag(fmt.Sprintf(`msgpack:"%s"`, tagValue))
		}
		resultFields = append(resultFields, field)
	}

	resultTyp := reflect.StructOf(resultFields)
	resultSlice := reflect.MakeSlice(reflect.SliceOf(resultTyp), 0, 0)
	return reflect.Append(resultSlice, reflect.New(resultTyp).Elem())
}

// cleanTag cleans up the tag value by removing special characters
func CleanTag(tag string) string {
	return strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			return r
		default:
			return -1
		}
	}, tag)
}*/

// SaveWorld saves a single World struct to the database with the given name
func SaveWorld(db *sqlite3.Conn, name string, world World) error {
	// Serialize World struct to MessagePack binary format
	worldBytes, err := msgpack.Marshal(world)
	//log.Warn("WorldBytes (msgpack) to be saved: %s", worldBytes)
	if err != nil {
		return err
	}

	// Prepare statement for insertion
	stmt, err := db.Prepare("INSERT INTO worlds (name, data) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Insert name and binary data into the database
	stmt.Exec(name, worldBytes)
	return nil
}

// LoadWorld loads a single World struct from the database by its name
func LoadWorld(db *sqlite3.Conn, name string) (*World, error) {
	var world World

	// Query the data
	row, err := db.Prepare("SELECT data FROM worlds WHERE name = ?", name)

	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}
	defer row.Close()

	for {
		hasRow, err := row.Step()
		if err != nil {
			log.Fatal(err.Error())

			// panic("Error Parsing World: " + name)
		}
		if !hasRow {
			// The query is finished
			log.Fatal("World %s Not Found in our database!", name)
			return nil, errors.New("World " + name + " Not Found in our database!")
			break
		}

		// Deserialize MessagePack binary data into World struct
		var worldBytes []byte
		if err := row.Scan(&worldBytes); err != nil {
			log.Fatal(err.Error())
			return nil, err
			break
		}
		if err := msgpack.Unmarshal(worldBytes, &world); err != nil {
			log.Fatal(err.Error())
			return nil, err
			break
		}
		//log.Warn("[Decoded] WorldBytes (msgpack) to be loaded: %v", world)
		if world.Name == name {
			log.Error("Found World named: %s", world.Name)
			return &world, nil
			break
		}
	}
	return &world, nil
}

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
	Worlds[name] = *world
	return world
}

func GetWorld(name string) (*World, error) {
	world, ok := Worlds[name]
	if ok {
		return &world, nil
	} else {
		return nil, errors.New("World with name " + name + " is not exist in our database!")
	}
}
