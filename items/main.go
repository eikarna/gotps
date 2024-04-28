package items

import (
	"encoding/binary"
	"fmt"
	"os"

	log "github.com/codecat/go-libs/log"
)

type Item struct {
	Name               string
	TexturePath        string
	ExtraFilePath      string
	PetName            string
	PetPrefix          string
	PetSuffix          string
	PetAbility         string
	ExtraOptions       string
	TexturePath2       string
	ExtraOptions2      string
	PunchOption        string
	StrData11          string
	StrData15          string
	StrData16          string
	ItemID             int32
	TextureHash        int32
	Val1               int32
	DropChance         int32
	ExtrafileHash      int32
	AudioVolume        int32
	WeatherID          int32
	SeedColor          int32
	SeedOverlayColor   int32
	GrowTime           int32
	IntData13          int32
	IntData14          int32
	Rarity             int16
	Val2               int16
	IsRayman           int16
	EditableType       int8
	ItemCategory       int8
	ActionType         int8
	HitsoundType       int8
	ItemKind           int8
	TextureX           int8
	TextureY           int8
	SpreadType         int8
	CollisionType      int8
	BreakHits          int8
	ClothingType       int8
	MaxAmount          int8
	SeedBase           int8
	SeedOverlay        int8
	TreeBase           int8
	TreeLeaves         int8
	IsStripeyWallpaper int8
}

type ItemInfo struct {
	ItemVersion int16
	ItemCount   int32

	Items []Item

	//items.dat packet
	FileBufferPacket []byte
	FileSize         int32
	FileHash         uint32
}

func getHash(str []byte, length int) uint32 {
	n := str
	acc := uint32(0x55555555)
	if length == 0 {
		for _, c := range n {
			acc = (acc >> 27) + (acc << 5) + uint32(c)
		}
	} else {
		for i := 0; i < length; i++ {
			acc = (acc >> 27) + (acc << 5) + uint32(n[i])
		}
	}
	return acc
}

func getFileHash(pathFile string) uint32 {
	file, err := os.Open(pathFile)
	if err != nil {
		fmt.Errorf("error opening file: %v", err)
	}

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Errorf("error getting file information: %v", err)
	}

	size := fileInfo.Size()
	if size == -1 {
		fmt.Errorf("error: File size is -1")
	}

	data := make([]byte, size)
	_, err = file.Read(data)
	if err != nil {
		fmt.Errorf("error reading file: %v", err)
	}

	return getHash(data, int(size))
}

func (Info *ItemInfo) GetItemHash() uint32 {
	return uint32(Info.FileHash)
}

func SerializeItemsDat(pathFile string) (*ItemInfo, error) {
	itemInfo := &ItemInfo{}
	file, err := os.Open(pathFile)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("error getting file information: %v", err)
	}

	size := fileInfo.Size()
	if size == -1 {
		return nil, fmt.Errorf("error: File size is -1")
	}
	itemInfo.FileSize = int32(size)
	itemInfo.FileHash = getFileHash(pathFile)
	data := make([]byte, size)
	_, err = file.Read(data)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	itemInfo.FileBufferPacket = make([]byte, 60+size)
	binary.LittleEndian.PutUint32(itemInfo.FileBufferPacket[0:], 4)
	binary.LittleEndian.PutUint32(itemInfo.FileBufferPacket[4:], 16)
	binary.LittleEndian.PutUint32(itemInfo.FileBufferPacket[8:], ^uint32(0))
	binary.LittleEndian.PutUint32(itemInfo.FileBufferPacket[16:], 8)
	binary.LittleEndian.PutUint32(itemInfo.FileBufferPacket[56:], uint32(size))
	copy(itemInfo.FileBufferPacket[60:], data)

	log.Info("Items.dat serialized with itemcount: %d, itemversion: %d, itemhash: %v", itemInfo.ItemCount, itemInfo.ItemVersion, itemInfo.FileHash)
	return itemInfo, nil
}
