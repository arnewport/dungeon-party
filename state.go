package main

import "fmt"

// PARTY

var State = Party{
	Characters: []Character{},
}

type Party struct {
	Characters []Character
}

// REGISTRIES

var ItemsByID = map[int]*Item{}
var WeaponsByID = map[int]*Weapon{}
var ArmorByID = map[int]*Armor{}
var ShieldsByID = map[int]*Shield{}
var JewelryByID = map[int]*Jewelry{}
var RodsWandsStavesByID = map[int]*RodWandStaff{}

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

var SpellsByID = map[int]Spell{}

var ArcaneSpellIDs = []int{}
var DivineSpellIDs = []int{}

type SpellValidationConfig struct {
	EnforceKnown bool // cannot know more spells than you can cast
	EnforceLevel bool // cannot learn spells of a higher level than you can cast (learning != memorizing)
}

var SpellRules = SpellValidationConfig{true, true}

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

func AddKnownSpell(c *Character, spellID int) error {
	if _, ok := SpellsByID[spellID]; !ok {
		return fmt.Errorf("spell %d not found", spellID)
	}
	for _, id := range c.KnownSpells {
		if id == spellID {
			return fmt.Errorf("spell %d already known", spellID)
		}
	}
	c.KnownSpells = append(c.KnownSpells, spellID)
	return nil
}

func RemoveKnownSpell(c *Character, spellID int) {
	for i, id := range c.KnownSpells {
		if id == spellID {
			c.KnownSpells = append(c.KnownSpells[:i], c.KnownSpells[i+1:]...)
			return
		}
	}
}

func CanLearnSpell(c *Character, spellID int) error {
	spell, ok := SpellsByID[spellID]
	if !ok {
		return fmt.Errorf("spell %d does not exist", spellID)
	}

	// Check if already known
	for _, id := range c.KnownSpells {
		if id == spellID {
			return fmt.Errorf("spell %s already known", spell.Name)
		}
	}

	// Determine how many spells of this level are already known
	spellLevel := spell.Level
	knownAtLevel := 0
	for _, id := range c.KnownSpells {
		s, ok := SpellsByID[id]
		if !ok {
			continue
		}
		if s.Level == spellLevel {
			knownAtLevel++
		}
	}

	// Determine how many spells can be known at this level
	allowed := GetSpellSlots(c.Class, c.Level, spellLevel)
	if knownAtLevel >= allowed {
		return fmt.Errorf("too many known spells at level %d (limit %d)", spellLevel, allowed)
	}

	return nil
}

func GetSpellSlots(class CharacterClass, level int, spellLevel int) int {
	byClass, ok := spellSlotTables[class]
	if !ok {
		return 0
	}
	byLevel, ok := byClass[level]
	if !ok {
		return 0
	}
	return byLevel[spellLevel]
}

func MaxSpellLevelAvailable(class CharacterClass, level int) int {
	classTable, ok := spellSlotTables[class]
	if !ok {
		return 0
	}
	levelTable, ok := classTable[level]
	if !ok {
		return 0
	}

	maxLevel := 0
	for spellLevel, count := range levelTable {
		if count > 0 && spellLevel > maxLevel {
			maxLevel = spellLevel
		}
	}

	return maxLevel
}

// spellSlotTables[class][charLevel][spellLevel] = number of slots
var spellSlotTables = map[CharacterClass]map[int]map[int]int{
	ClassCleric: {
		1:  {},
		2:  {1: 1},
		3:  {1: 2},
		4:  {1: 2, 2: 1},
		5:  {1: 2, 2: 2},
		6:  {1: 2, 2: 2, 3: 1, 4: 1},
		7:  {1: 2, 2: 2, 3: 2, 4: 1, 5: 1},
		8:  {1: 3, 2: 3, 3: 2, 4: 2, 5: 1},
		9:  {1: 3, 2: 3, 3: 3, 4: 2, 5: 2},
		10: {1: 4, 2: 4, 3: 3, 4: 3, 5: 2},
		11: {1: 4, 2: 4, 3: 4, 4: 3, 5: 3},
		12: {1: 5, 2: 5, 3: 4, 4: 4, 5: 3},
		13: {1: 5, 2: 5, 3: 5, 4: 4, 5: 4},
		14: {1: 6, 2: 5, 3: 5, 4: 5, 5: 4},
	},
	ClassMagicUser: {
		1:  {1: 1},
		2:  {1: 2},
		3:  {1: 2, 2: 1},
		4:  {1: 2, 2: 2},
		5:  {1: 2, 2: 2, 3: 1},
		6:  {1: 2, 2: 2, 3: 2},
		7:  {1: 3, 2: 2, 3: 2, 4: 1},
		8:  {1: 3, 2: 3, 3: 2, 4: 2},
		9:  {1: 3, 2: 3, 3: 3, 4: 2, 5: 1},
		10: {1: 3, 2: 3, 3: 3, 4: 3, 5: 2},
		11: {1: 4, 2: 3, 3: 3, 4: 3, 5: 2, 6: 1},
		12: {1: 4, 2: 4, 3: 3, 4: 3, 5: 3, 6: 2},
		13: {1: 4, 2: 4, 3: 4, 4: 3, 5: 3, 6: 3},
		14: {1: 4, 2: 4, 3: 4, 4: 4, 5: 3, 6: 3},
	},
}

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

func RegisterRodWandStaff(r RodWandStaff) error {
	if _, exists := RodsWandsStavesByID[r.ID]; exists {
		return fmt.Errorf("rod, wand, or staff with ID %d already registered", r.ID)
	}
	RodsWandsStavesByID[r.ID] = &r
	return nil
}

func UnregisterItem(id int) {
	delete(ItemsByID, id)
	delete(WeaponsByID, id)
	delete(ArmorByID, id)
	delete(ShieldsByID, id)
	delete(JewelryByID, id)
	delete(RodsWandsStavesByID, id)
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

func GetRodWandStaffByID(id int) *RodWandStaff {
	return RodsWandsStavesByID[id]
}
