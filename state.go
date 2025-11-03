package main

import "fmt"

// PARTY

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
	MaximumHitPoints int
	Strength         int
	Intelligence     int
	Wisdom           int
	Dexterity        int
	Constitution     int
	Charisma         int
	Items            []int
	ArmorID          int
	ShieldID         int
	Spellcasting     []SpellType
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
	ItemGeneric    ItemType = "item"
	ItemWeapon     ItemType = "weapon"
	ItemArmor      ItemType = "armor"
	ItemShield     ItemType = "shield"
	ItemJewelry    ItemType = "jewelry"
	ItemLimitedUse ItemType = "limiteduseitem"
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

type LimitedUseItem struct {
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
	Type  SpellType
}

type SpellType string

const (
	SpellArcane SpellType = "arcane"
	SpellDivine SpellType = "divine"
)

type MemorizedSpell struct {
	SpellID int
	Cast    bool
}

// REGISTRIES

var ItemsByID = map[int]*Item{}
var WeaponsByID = map[int]*Weapon{}
var ArmorByID = map[int]*Armor{}
var ShieldsByID = map[int]*Shield{}
var JewelryByID = map[int]*Jewelry{}
var LimitedUseItemsByID = map[int]*LimitedUseItem{}

// REGISTRATION

func RegisterItem(item Item) error {
	if _, exists := ItemsByID[item.ID]; exists {
		return fmt.Errorf("item with ID %d already registered", item.ID)
	}
	ItemsByID[item.ID] = &item
	return nil
}

func RegisterWeapon(w Weapon) error {
	if _, exists := WeaponsByID[w.ID]; exists {
		return fmt.Errorf("weapon with ID %d already registered", w.ID)
	}
	WeaponsByID[w.ID] = &w
	return nil
}

func RegisterArmor(a Armor) error {
	if _, exists := ArmorByID[a.ID]; exists {
		return fmt.Errorf("armor with ID %d already registered", a.ID)
	}
	ArmorByID[a.ID] = &a
	return nil
}

func RegisterShield(s Shield) error {
	if _, exists := ShieldsByID[s.ID]; exists {
		return fmt.Errorf("shield with ID %d already registered", s.ID)
	}
	ShieldsByID[s.ID] = &s
	return nil
}

func RegisterJewelry(j Jewelry) error {
	if _, exists := JewelryByID[j.ID]; exists {
		return fmt.Errorf("jewelry with ID %d already registered", j.ID)
	}
	JewelryByID[j.ID] = &j
	return nil
}

func RegisterLimitedUseItem(lu LimitedUseItem) error {
	if _, exists := LimitedUseItemsByID[lu.ID]; exists {
		return fmt.Errorf("limited use item with ID %d already registered", lu.ID)
	}
	LimitedUseItemsByID[lu.ID] = &lu
	return nil
}

func UnregisterItem(id int) {
	delete(ItemsByID, id)
	delete(WeaponsByID, id)
	delete(ArmorByID, id)
	delete(ShieldsByID, id)
	delete(JewelryByID, id)
	delete(LimitedUseItemsByID, id)
}

// GETTERS

func GetItemByID(id int) *Item {
	return ItemsByID[id]
}

func GetWeaponByID(id int) *Weapon {
	return WeaponsByID[id]
}

func GetArmorByID(id int) *Armor {
	return ArmorByID[id]
}

func GetShieldByID(id int) *Shield {
	return ShieldsByID[id]
}

func GetJewelryByID(id int) *Jewelry {
	return JewelryByID[id]
}

func GetLimitedUseItemByID(id int) *LimitedUseItem {
	return LimitedUseItemsByID[id]
}
