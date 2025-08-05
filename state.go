package main

// PARTY

var State = Party{
	Characters: []Character{},
	Items:      []Item{}, // Items carried by an unspecified member of the party. Acts as an item sorting space
	Limbo:      []Item{}, // Items available to be picked up in the world
	Storage:    []Item{}, // Items stored safely in the world
}

type Party struct {
	Characters []Character
	Items      []Item
	Limbo      []Item
	Storage    []Item
}

// CHARACTER

type Character struct {
	ID               int
	Name             string
	Class            string
	Level            int
	Alignment        string
	ArmorBonus       int
	RolledHitPoints  int
	CurrentHitPoints int
	Strength         int
	Intelligence     int
	Wisdom           int
	Dexterity        int
	Constitution     int
	Charisma         int
	Items            []Item
	ArmorID          int
	ShieldID         int
	ArcaneSpells     bool
	DivineSpells     bool
	KnownSpells      []int // Known spells for Magic Users and Elves. Left empty for Clerics
	MemorizedSpells  []MemorizedSpell
}

// ITEMS

type Item struct {
	ID       int
	Name     string
	Type     string // "item", "weapon", "armor", "shield", "jewelry", "rodwandstaff"
	Location ItemLocation
	URL      string
}

type ItemLocation struct {
	Type        string // "character", "limbo", "party", "storage"
	CharacterID *int   // If Type == "character"
}

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
	Charges int
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
