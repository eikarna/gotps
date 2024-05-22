package worlds

import (
	"errors"
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
	OwnerName  string
	PlayersIn  int16
	SizeX      int32
	SizeY      int32
	TotalTiles int32
	Admins     []Admins
	Tiles      []Tiles
	PosDoor    int32
}

// Flags constants
const (
	FlagsTileExtra   = 0x0001
	FlagsLocked      = 0x0002
	FlagsSeed        = 0x0010
	FlagsTree        = 0x0019
	FlagsFlipped     = 0x0020
	FlagsRotatedLeft = 0x0030
	FlagsOpen        = 0x0040
	FlagsPublic      = 0x0080
	FlagsSilenced    = 0x0200
	FlagsWater       = 0x0400
	FlagsFire        = 0x1000
	FlagsRed         = 0x2000
	FlagsBlue        = 0x8000
	FlagsGreen       = 0x4000
	FlagsYellow      = 0x6000
	FlagsPurple      = 0xa000
)

// ItemCollisionType enum
type ItemCollisionType int

const (
	ItemCollisionNone ItemCollisionType = iota
	ItemCollisionNormal
	ItemCollisionJumpThrough
	ItemCollisionGateway
	ItemCollisionIfOff
	ItemCollisionOneWay
	ItemCollisionVIP
	ItemCollisionWaterfall
	ItemCollisionAdventure
	ItemCollisionIfOn
	ItemCollisionTeamEntrance
	ItemCollisionGuild
	ItemCollisionCloud
	ItemCollisionFriendEntrance
)

// ActionTypes enum
type ActionTypes int

const (
	Fist ActionTypes = iota
	Wrench
	Door
	Lock
	Gems
	Treasure
	DeadlyBlock
	Trampoline
	Consumable
	Gateway
	Sign
	SfxWithExtraFrame
	BoomBox
	MainDoor
	Platform
	Bedrock
	Lava
	Foreground
	Background
	Seed
	Clothes
	ForegroundWithExtraFrame
	BackgdSfxExtraFrame
	BackBoomBox
	Bouncy
	Pointy
	Portal
	Checkpoint
	SheetMusic
	Ice
	Switcheroo
	Chest
	Mailbox
	Bulletin
	Pinata
	Dice
	Chemical
	Provider
	Lab
	Achievement
	WeatherMachine
	ScoreBoard
	Sungate
	Profile
	DeadlyIfOn
	HeartMonitor
	DonationBox
	Toybox
	Mannequin
	SecurityCamera
	MagicEgg
	GameResources
	GameGenerator
	Xenonite
	Dressup
	Crystal
	Burglar
	Compactor
	Spotlight
	Wind
	DisplayBlock
	VendingMachine
	Fishtank
	Petfish
	Solar
	Forge
	GivingTree
	GivingTreeStump
	Steampunk
	SteamLavaIfOn
	SteamOrgan
	Tamagotchi
	Swing
	Flag
	LobsterTrap
	ArtCanvas
	BattleCage
	PetTrainer
	SteamEngine
	Lockbot
	WeatherSpecial
	SpiritStorage
	DisplayShelf
	VipEntrance
	ChallengeTimer
	ChallengeFlag
	FishMount
	Portrait
	WeatherSpecial2
	Fossil
	FossilPrep
	DnaMachine
	Blaster
	Valhowla
	Chemsynth
	Chemtank
	Storage
	Oven
	SuperMusic
	GeigerCharger
	AdventureReset
	TombRobber
	Faction
	RedFaction
	GreenFaction
	BlueFaction
	Ances
	FishgotchiTank
	FishingBlock
	ItemSucker
	ItemPlanter
	Robot
	Command
	Ticket
	StatsBlock
	FieldNode
	OuijaBoard
	ArchitectMachine
	Starship
	Autodelete
	GreenFountain
	AutoActionBreak
	AutoActionHarvest
	AutoActionHarvestSuck
	LightningIfOn
	PhasedBlock
	Mud
	RootCutting
	PasswordStorage
	PhasedBlock2
	Bomb
	WeatherInfinity
	Slime
	Unk1
	Completionist
	Unk3
	FeedingBlock
	KrankenBlock
	FriendsEntrance
)

// Options constants
const (
	MusicBlocksDisabled = 0x10
	MusicBlocksInvis    = 0x20
)

// ExtraTypes enum
type ExtraTypes int

const (
	ExtraNone ExtraTypes = iota
	ExtraDoor
	ExtraMainDoor
	ExtraSign
	ExtraLock
	ExtraSeed
	ExtraMailbox
	ExtraBulletin
	ExtraDice
	ExtraProvider
	ExtraAchievement
	ExtraHeartMonitor
	ExtraDonationBlock
	ExtraToyBox
	ExtraMannequin
	ExtraMagicEgg
	ExtraGameResources
	ExtraGameGenerator
	ExtraXenonite
	ExtraDressUp
	ExtraCrystal
	ExtraBurglar
	ExtraSpotlight
	ExtraDisplayBlock
	ExtraVendingMachine
	ExtraFishTank
	ExtraSolar
	ExtraForge
	ExtraGivingTree
	ExtraGivingTreeStump
	ExtraSteamOrgan
	ExtraTamagotchi
	ExtraSwing
	ExtraFlag
	ExtraLobsterTrap
	ExtraArtCanvas
	ExtraBattleCage
	ExtraPetTrainer
	ExtraSteamEngine
	ExtraLockBot
	ExtraWeatherSpecial
	ExtraSpiritStorage
	ExtraUnknown1
	ExtraDisplayShelf
	ExtraVipEntrance
	ExtraChallengeTimer
	ExtraChallengeFlag
	ExtraFishMount
	ExtraPortrait
	ExtraWeatherSpecial2
	ExtraFossilPrep
	ExtraDnaMachine
	ExtraMagPlant             = 0x3e
	ExtraGrowScan             = 66
	ExtraTesseractManipulator = 0x45
	ExtraGaiaHeart            = 0x46
	ExtraTechnoOrganicEngine  = 0x47
	ExtraKrankenGalactic      = 0x50
	ExtraWeatherInfinity      = 0x4d
)

var Worlds = make(map[string]*World)

func ActionType(value int16) ActionTypes {
	switch value {
	case 0:
		return Fist
	case 1:
		return Wrench
	case 2:
		return Door
	case 3:
		return Lock
	case 4:
		return Gems
	case 5:
		return Treasure
	case 6:
		return DeadlyBlock
	case 7:
		return Trampoline
	case 8:
		return Consumable
	case 9:
		return Gateway
	case 10:
		return Sign
	case 11:
		return SfxWithExtraFrame
	case 12:
		return BoomBox
	case 13:
		return MainDoor
	case 14:
		return Platform
	case 15:
		return Bedrock
	case 16:
		return Lava
	case 17:
		return Foreground
	case 18:
		return Background
	case 19:
		return Seed
	case 20:
		return Clothes
	case 21:
		return ForegroundWithExtraFrame
	case 22:
		return BackgdSfxExtraFrame
	case 23:
		return BackBoomBox
	case 24:
		return Bouncy
	case 25:
		return Pointy
	case 26:
		return Portal
	case 27:
		return Checkpoint
	case 28:
		return SheetMusic
	case 29:
		return Ice
	case 31:
		return Switcheroo
	case 32:
		return Chest
	case 33:
		return Mailbox
	case 34:
		return Bulletin
	case 35:
		return Pinata
	case 36:
		return Dice
	case 37:
		return Chemical
	case 38:
		return Provider
	case 39:
		return Lab
	case 40:
		return Achievement
	case 41:
		return WeatherMachine
	case 42:
		return ScoreBoard
	case 43:
		return Sungate
	case 44:
		return Profile
	case 45:
		return DeadlyIfOn
	case 46:
		return HeartMonitor
	case 47:
		return DonationBox
	case 48:
		return Toybox
	case 49:
		return Mannequin
	case 50:
		return SecurityCamera
	case 51:
		return MagicEgg
	case 52:
		return GameResources
	case 53:
		return GameGenerator
	case 54:
		return Xenonite
	case 55:
		return Dressup
	case 56:
		return Crystal
	case 57:
		return Burglar
	case 58:
		return Compactor
	case 59:
		return Spotlight
	case 60:
		return Wind
	case 61:
		return DisplayBlock
	case 62:
		return VendingMachine
	case 63:
		return Fishtank
	case 64:
		return Petfish
	case 65:
		return Solar
	case 66:
		return Forge
	case 67:
		return GivingTree
	case 68:
		return GivingTreeStump
	case 69:
		return Steampunk
	case 70:
		return SteamLavaIfOn
	case 71:
		return SteamOrgan
	case 72:
		return Tamagotchi
	case 73:
		return Swing
	case 74:
		return Flag
	case 75:
		return LobsterTrap
	case 76:
		return ArtCanvas
	case 77:
		return BattleCage
	case 78:
		return PetTrainer
	case 79:
		return SteamEngine
	case 80:
		return Lockbot
	case 81:
		return WeatherSpecial
	case 82:
		return SpiritStorage
	case 83:
		return DisplayShelf
	case 84:
		return VipEntrance
	case 85:
		return ChallengeTimer
	case 86:
		return ChallengeFlag
	case 87:
		return FishMount
	case 88:
		return Portrait
	case 89:
		return WeatherSpecial2
	case 90:
		return Fossil
	case 91:
		return FossilPrep
	case 92:
		return DnaMachine
	case 93:
		return Blaster
	case 94:
		return Valhowla
	case 95:
		return Chemsynth
	case 96:
		return Chemtank
	case 97:
		return Storage
	case 98:
		return Oven
	case 99:
		return SuperMusic
	case 100:
		return GeigerCharger
	case 101:
		return AdventureReset
	case 102:
		return TombRobber
	case 103:
		return Faction
	case 104:
		return RedFaction
	case 105:
		return GreenFaction
	case 106:
		return BlueFaction
	case 107:
		return Ances
	case 109:
		return FishgotchiTank
	case 110:
		return FishingBlock
	case 111:
		return ItemSucker
	case 112:
		return ItemPlanter
	case 113:
		return Robot
	case 114:
		return Command
	case 115:
		return Ticket
	case 116:
		return StatsBlock
	case 117:
		return FieldNode
	case 118:
		return OuijaBoard
	case 119:
		return ArchitectMachine
	case 120:
		return Starship
	case 121:
		return Autodelete
	case 122:
		return GreenFountain
	case 123:
		return AutoActionBreak
	case 124:
		return AutoActionHarvest
	case 125:
		return AutoActionHarvestSuck
	case 126:
		return LightningIfOn
	case 127:
		return PhasedBlock
	case 128:
		return Mud
	case 129:
		return RootCutting
	case 130:
		return PasswordStorage
	case 131:
		return PhasedBlock2
	case 132:
		return Bomb
	case 134:
		return WeatherInfinity
	case 135:
		return Slime
	case 136:
		return Unk1
	case 137:
		return Completionist
	case 138:
		return Unk3
	case 140:
		return FeedingBlock
	case 141:
		return KrankenBlock
	case 142:
		return FriendsEntrance
	default:
		// Return a default value or handle the error accordingly
		return -1 // For simplicity, return an invalid value
	}
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
			tile.Label = "`cEXIT``"
			tile.Fg = 6
			world.PosDoor = int32(i + randomPosDoor)
		}
		if i == 2500+randomPosDoor {
			tile.Fg = 8
		}
		if i >= 2500 {
			tile.Bg = 14
		}
		world.Tiles = append(world.Tiles, tile)
	}
	Worlds[name] = world
	return world
}

func GetWorld(name string) (*World, error) {
	world, ok := Worlds[name]
	if ok {
		return world, nil
	} else {
		return nil, errors.New("World with name " + name + " is not exist in our database!")
	}
}
