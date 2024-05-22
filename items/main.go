package items

import (
	"encoding/binary"
	"fmt"
	"os"
	// "strconv"

	log "github.com/codecat/go-libs/log"
	"time"
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
	ActionType         int16
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

func byteArrayToInt(byteSlice []byte) (int, error) {
	var result int
	for _, b := range byteSlice {
		if b < '0' || b > '9' {
			return 0, fmt.Errorf("Invalid byte: %c", b)
		}
		result = result*10 + int(b-'0')
	}
	return result, nil
}

func SerializeItemsDat(pathFile string, timestamp time.Time) (*ItemInfo, error) {
	itemInfo := &ItemInfo{}
	key := "PBG892FXX982ABC*"

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
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	itemInfo.FileBufferPacket = make([]byte, 60+size)
	binary.LittleEndian.PutUint32(itemInfo.FileBufferPacket[0:], 4)
	binary.LittleEndian.PutUint32(itemInfo.FileBufferPacket[4:], 16)
	binary.LittleEndian.PutUint32(itemInfo.FileBufferPacket[8:], ^uint32(0))
	binary.LittleEndian.PutUint32(itemInfo.FileBufferPacket[16:], 8)
	binary.LittleEndian.PutUint32(itemInfo.FileBufferPacket[56:], uint32(size))
	copy(itemInfo.FileBufferPacket[60:], data)
	// log.Info("FileBufferPacket Length: %d, data Length: %d", len(itemInfo.FileBufferPacket))
	memPos := 0
	itemInfo.ItemVersion = int16(binary.LittleEndian.Uint16(data[memPos:]))
	memPos += 2
	itemInfo.ItemCount = int32(binary.LittleEndian.Uint32(data[memPos:]))
	memPos += 4
	itemInfo.Items = make([]Item, itemInfo.ItemCount)
	for i := 0; int32(i) < itemInfo.ItemCount; i++ {
		// Items Dat Info start from 66
		// if memPos <= len(itemInfo.FileBufferPacket[66:]) {
		// if memPos < int(size) {
		// log.Info("MemPos: %d", memPos)
		itemId := binary.LittleEndian.Uint32(data[memPos:])
		// log.Info("Got Items ID: %d", int(binary.LittleEndian.Uint16(data[memPos:])))
		itemInfo.Items[i].ItemID = int32(itemId)
		memPos += 4
		if int32(itemId) < itemInfo.ItemCount {
			// log.Info("Got Items EditableType: %d", int(data[memPos]))
			itemInfo.Items[i].EditableType = int8(data[memPos])
			memPos += 1
			// log.Info("Got Items ItemCategory: %d", int(data[memPos]))
			itemInfo.Items[i].ItemCategory = int8(data[memPos])
			memPos += 1
			// log.Info("Got Items ActionType: %d", int(data[memPos]))
			itemInfo.Items[i].ActionType = int16(data[memPos])
			memPos += 1
			// log.Info("Got Items hitSoundType: %d", int(data[memPos]))
			itemInfo.Items[i].HitsoundType = int8(data[memPos])
			memPos += 1
			// Read first strLen
			strLen := int(binary.LittleEndian.Uint16(data[memPos:]))
			memPos += 2
			for j := 0; j < strLen; j++ {
				itemInfo.Items[i].Name += string(data[memPos] ^ key[(int32(j)+itemInfo.Items[i].ItemID)%int32(len(key))])
				memPos++
			}

			// Read second strLen
			strLen = int(binary.LittleEndian.Uint16(data[memPos:]))
			memPos += 2
			itemInfo.Items[i].TexturePath = string(data[memPos : memPos+strLen])
			memPos += strLen
			// log.Info("Got Items Name: %s, and Texture: %s", itemInfo.Items[i].Name, itemInfo.Items[i].TexturePath)

			itemInfo.Items[i].TextureHash = int32(binary.LittleEndian.Uint32(data[memPos:]))
			memPos += 4
			itemInfo.Items[i].ItemKind = int8(data[memPos])
			memPos += 1
			itemInfo.Items[i].Val1 = int32(binary.LittleEndian.Uint32(data[memPos:]))
			memPos += 4
			itemInfo.Items[i].TextureX = int8(data[memPos])
			memPos += 1
			itemInfo.Items[i].TextureY = int8(data[memPos])
			memPos += 1
			itemInfo.Items[i].SpreadType = int8(data[memPos])
			memPos += 1
			itemInfo.Items[i].IsStripeyWallpaper = int8(data[memPos])
			memPos += 1
			itemInfo.Items[i].CollisionType = int8(data[memPos])
			memPos += 1
			itemInfo.Items[i].BreakHits = int8(data[memPos])
			memPos += 1
			itemInfo.Items[i].DropChance = int32(binary.LittleEndian.Uint32(data[memPos:]))
			memPos += 4
			itemInfo.Items[i].ClothingType = int8(data[memPos])
			memPos += 1
			itemInfo.Items[i].Rarity = int16(binary.LittleEndian.Uint16(data[memPos:]))
			memPos += 2
			itemInfo.Items[i].MaxAmount = int8(data[memPos])
			memPos += 1
			strLen = int(binary.LittleEndian.Uint16(data[memPos:]))
			memPos += 2
			itemInfo.Items[i].ExtraFilePath = string(data[memPos : memPos+strLen])
			// log.Info("Got Items ExtraFile: %s", itemInfo.Items[i].ExtraFilePath)
			memPos += strLen
			memPos += 8
			strLen = int(binary.LittleEndian.Uint16(data[memPos:]))
			memPos += 2
			itemInfo.Items[i].PetName = string(data[memPos : memPos+strLen])
			memPos += strLen
			strLen = int(binary.LittleEndian.Uint16(data[memPos:]))
			memPos += 2
			itemInfo.Items[i].PetPrefix = string(data[memPos : memPos+strLen])
			memPos += strLen
			strLen = int(binary.LittleEndian.Uint16(data[memPos:]))
			memPos += 2
			itemInfo.Items[i].PetSuffix = string(data[memPos : memPos+strLen])
			memPos += strLen
			strLen = int(binary.LittleEndian.Uint16(data[memPos:]))
			memPos += 2
			itemInfo.Items[i].PetAbility = string(data[memPos : memPos+strLen])
			memPos += strLen
			// log.Info("Got Items PetName: %s, PetPrefix: %s, PetSuffix: %s, PetAbility: %s", petName, petPrefix, petSuffix, petAbility)
			// TODO: Parse all byte
			itemInfo.Items[i].SeedBase = int8(data[memPos])
			memPos += 1
			itemInfo.Items[i].SeedOverlay = int8(data[memPos])
			memPos += 1
			itemInfo.Items[i].TreeBase = int8(data[memPos])
			memPos += 1
			itemInfo.Items[i].TreeLeaves = int8(data[memPos])
			memPos += 1
			itemInfo.Items[i].SeedColor = int32(binary.LittleEndian.Uint32(data[memPos:]))
			memPos += 4
			itemInfo.Items[i].SeedOverlayColor = int32(binary.LittleEndian.Uint32(data[memPos:]))
			memPos += 8
			itemInfo.Items[i].GrowTime = int32(binary.LittleEndian.Uint32(data[memPos:]))
			memPos += 4
			itemInfo.Items[i].Val2 = int16(binary.LittleEndian.Uint16(data[memPos:]))
			memPos += 2
			itemInfo.Items[i].IsRayman = int16(binary.LittleEndian.Uint16(data[memPos:]))
			memPos += 2

			strLen = int(binary.LittleEndian.Uint16(data[memPos:]))
			memPos += 2
			itemInfo.Items[i].ExtraOptions = string(data[memPos : memPos+int(strLen)])
			// log.Info("Got Items ExtraOptions: %s", itemInfo.Items[i].ExtraOptions)
			memPos += strLen

			strLen = int(binary.LittleEndian.Uint16(data[memPos:]))
			memPos += 2
			itemInfo.Items[i].TexturePath2 = string(data[memPos : memPos+strLen])
			memPos += strLen
			// log.Info("Got Items Texture2: %s", itemInfo.Items[i].TexturePath2)

			strLen = int(binary.LittleEndian.Uint16(data[memPos:]))
			memPos += 2
			itemInfo.Items[i].ExtraOptions2 = string(data[memPos : memPos+strLen])
			memPos += strLen
			// log.Info("Got Items ExtraOptions2: %s", itemInfo.Items[i].ExtraOptions2)
			memPos += 80
			if itemInfo.ItemVersion >= 11 {
				strLen := int(binary.LittleEndian.Uint16(data[memPos:]))
				memPos += 2
				itemInfo.Items[i].PunchOption = string(data[memPos : memPos+strLen])
				memPos += strLen
				// log.Info("Got Items PunchOptions: %s", itemInfo.Items[i].PunchOption)
			}
			if itemInfo.ItemVersion >= 12 {
				memPos += 13
			}
			if itemInfo.ItemVersion >= 13 {
				memPos += 4
			}
			if itemInfo.ItemVersion >= 14 {
				memPos += 4
			}
			if itemInfo.ItemVersion >= 15 {
				memPos += 25
				strLen := int(binary.LittleEndian.Uint16(data[memPos:]))
				memPos += 2 + strLen
			}
			if itemInfo.ItemVersion >= 16 {
				jLen := int(binary.LittleEndian.Uint16(data[memPos:]))
				memPos += 2 + jLen
			}
			if itemInfo.ItemVersion >= 17 {
				jLen := int(binary.LittleEndian.Uint16(data[memPos:]))
				memPos += 4 + jLen
			}
			if itemInfo.ItemVersion >= 18 {
				jLen := int(binary.LittleEndian.Uint16(data[memPos:]))
				memPos += 4 + jLen
			}

		} else {
			break
		}
	}

	log.Info("Items.dat serialized for %s. With Item Count: %d, ItemsDatVersion: %d, Item Hash: %v", time.Since(timestamp), itemInfo.ItemCount, itemInfo.ItemVersion, itemInfo.FileHash)
	return itemInfo, nil
}
