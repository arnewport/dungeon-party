package main

// PARTY

// ItemsByID All items that exist in the world by ID.
var ItemsByID = make(map[int]*Item)

var State = Party{
	Characters: []Character{},
}

type Party struct {
	Characters []Character
}

// CHARACTER

type Character struct {
	ID               int
	Name             string
	Class            CharacterClass
	Level            int
	Alignment        Alignment
	ArmorBonus       int
	RolledHitPoints  int
	CurrentHitPoints int
	Strength         int
	Intelligence     int
	Wisdom           int
	Dexterity        int
	Constitution     int
	Charisma         int
	Items            []int
	ArmorID          int
	ShieldID         int
	ArcaneSpells     bool
	DivineSpells     bool
	KnownSpells      []int // Known spells for Magic Users and Elves. Left empty for Clerics
	MemorizedSpells  []MemorizedSpell
}

type CharacterClass string

const (
	ClassNone      CharacterClass = ""
	ClassCleric    CharacterClass = "cleric"
	ClassFighter   CharacterClass = "fighter"
	ClassMagicUser CharacterClass = "magicuser"
	ClassThief     CharacterClass = "thief"
	ClassDwarf     CharacterClass = "dwarf"
	ClassElf       CharacterClass = "elf"
	ClassHalfling  CharacterClass = "halfling"
)

type Alignment string

const (
	AlignmentNone    Alignment = ""
	AlignmentLawful  Alignment = "lawful"
	AlignmentNeutral Alignment = "neutral"
	AlignmentChaotic Alignment = "chaotic"
)

// ITEMS

type Item struct {
	ID       int
	HolderID int
	Name     string
	Type     ItemType
	Location ItemLocation
	URL      string
}

type ItemType string

const (
	ItemGeneric      ItemType = "item"
	ItemWeapon       ItemType = "weapon"
	ItemArmor        ItemType = "armor"
	ItemShield       ItemType = "shield"
	ItemJewelry      ItemType = "jewelry"
	ItemRodWandStaff ItemType = "rodwandstaff"
)

type ItemLocation string

const (
	LocationNone      ItemLocation = "none"
	LocationCharacter ItemLocation = "character" // Items carried by a specified member of the party
	LocationParty     ItemLocation = "party"     // Items carried by an unspecified member of the party; acts as an intra-party item sorting space
	LocationStorage   ItemLocation = "storage"   // Items stored safely in the world
	LocationLimbo     ItemLocation = "limbo"     // Items available to be picked up in the world
)

type Weapon struct {
	Item
	Damage      int
	Bonus       int // 0 for non-magical, 1 to 3 for magical
	IsMelee     bool
	IsRanged    bool
	IsTwoHanded bool
	IsBlunt     bool
}

type Armor struct {
	Item
	Type  ArmorType
	Bonus int
}

type ArmorType string

const (
	Robes   ArmorType = "robes"
	Leather ArmorType = "leather"
	Chain   ArmorType = "chain"
	Plate   ArmorType = "plate"
)

type Shield struct {
	Item
	Bonus int
}

type Jewelry struct {
	Item
	ArmorBonus int
	SaveBonus  int
	/* Many pieces of magical jewelry and defensive items do not follow a clear pattern
	I'll account for this later. For now, this accounts for the Ring of Protection */
}

type RodWandStaff struct {
	Item
	Charges         int
	ArcaneAllowed   bool
	DivineAllowed   bool
	NonMagicAllowed bool
}

// SPELLS

type Spell struct {
	ID    int
	Name  string
	Level int
}

type MemorizedSpell struct {
	SpellID int
	Cast    bool
}
