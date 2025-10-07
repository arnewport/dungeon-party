package main

import (
	"fmt"
)

// Spell utilities

var SpellsByID = map[int]Spell{}

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

// SPELL LOGIC

type SpellValidationConfig struct {
	EnforceKnown bool // cannot know more spells than you can cast
	EnforceLevel bool // cannot learn spells of a higher level than you can cast (learning != memorizing)
}

var SpellRules = SpellValidationConfig{true, true}

// Add & Remove Known Spells

func AddKnownSpell(c *Character, spellID int) error {
	if err := CanLearnSpell(c, spellID); err != nil {
		return err
	}
	c.KnownSpells = append(c.KnownSpells, spellID)
	return nil
}

func RemoveKnownSpell(c *Character, spellID int) error {
	for i, id := range c.KnownSpells {
		if id == spellID {
			c.KnownSpells = append(c.KnownSpells[:i], c.KnownSpells[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("spell %d not known", spellID)
}

func AddMemorizedSpell(c *Character, spellID int) error {
	// Check that the spell exists
	spell, ok := SpellsByID[spellID]
	if !ok {
		return fmt.Errorf("spell %d does not exist", spellID)
	}

	// Must be known
	if err := checkSpellIsKnown(c, spellID); err != nil {
		return err
	}

	// Must not exceed per-level memorized count
	if err := checkMemorizedSpellSlotLimit(c, spell.Level); err != nil {
		return err
	}

	c.MemorizedSpells = append(c.MemorizedSpells, MemorizedSpell{SpellID: spellID})
	return nil
}

func RemoveMemorizedSpell(c *Character, spellID int) error {
	for i, ms := range c.MemorizedSpells {
		if ms.SpellID == spellID {
			c.MemorizedSpells = append(c.MemorizedSpells[:i], c.MemorizedSpells[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("spell %d is not memorized", spellID)
}

// Known Spell Helpers

// Check if the spell exists
func getSpellIfExists(spellID int) (Spell, error) {
	spell, ok := SpellsByID[spellID]
	if !ok {
		return Spell{}, fmt.Errorf("spell %d does not exist", spellID)
	}
	return spell, nil
}

// Check if the character is a spell caster of the proper type
func checkIfCharacterMayCast(c *Character, spellType SpellType) error {
	for _, st := range c.Spellcasting {
		if st == spellType {
			return nil
		}
	}
	return fmt.Errorf("%s cannot learn %s spells", c.Class, spellType)
}

// Determine if the spell is already known (for the purpose of not learning duplicate spells)
func checkIfAlreadyKnown(c *Character, spellID int, spellName string) error {
	for _, id := range c.KnownSpells {
		if id == spellID {
			return fmt.Errorf("spell %s already known", spellName)
		}
	}
	return nil
}

// Determine how many spells of this level are already known
func countKnownSpellsAtLevel(c *Character, level int) int {
	count := 0
	for _, id := range c.KnownSpells {
		s, ok := SpellsByID[id]
		if !ok {
			continue
		}
		if s.Level == level {
			count++
		}
	}
	return count
}

// Determine if spell level is too high for the character at this level
func checkSpellLevelAvailable(c *Character, spellLevel int) error {
	maxLevel := MaxSpellLevelAvailable(c.Class, c.Level)
	if spellLevel > maxLevel {
		return fmt.Errorf("cannot learn spells at level %d (limit %d)", spellLevel, maxLevel)
	}
	return nil
}

// Determine how many spells can be known at this level
func checkSpellSlotLimit(c *Character, spellLevel int) error {
	known := countKnownSpellsAtLevel(c, spellLevel)
	allowed := GetSpellSlots(c.Class, c.Level, spellLevel)

	if known >= allowed {
		return fmt.Errorf("too many known spells at level %d (limit %d)", spellLevel, allowed)
	}
	return nil
}

// Memorized Spell Helpers

// Determine if the spell is known (for the purpose of memorization)
func checkSpellIsKnown(c *Character, spellID int) error {
	for _, id := range c.KnownSpells {
		if id == spellID {
			return nil
		}
	}
	return fmt.Errorf("spell %d is not known", spellID)
}

// Determine how many spells can be memorized at this level
func checkMemorizedSpellSlotLimit(c *Character, spellLevel int) error {
	count := 0
	for _, ms := range c.MemorizedSpells {
		s, ok := SpellsByID[ms.SpellID]
		if !ok {
			continue
		}
		if s.Level == spellLevel {
			count++
		}
	}

	allowed := GetSpellSlots(c.Class, c.Level, spellLevel)
	if count >= allowed {
		return fmt.Errorf("cannot memorize more level %d spells (limit: %d)", spellLevel, allowed)
	}
	return nil
}

func isSpellKnown(c *Character, spellID int) bool {
	for _, id := range c.KnownSpells {
		if id == spellID {
			return true
		}
	}
	return false
}

//

// CanLearnSpell Determine if a spell can be learned by applying several checks
func CanLearnSpell(c *Character, spellID int) error {

	// Check if the spell exists
	spell, err := getSpellIfExists(spellID)
	if err != nil {
		return err
	}

	// Check if the character is a spell caster of the proper type
	if err := checkIfCharacterMayCast(c, spell.Type); err != nil {
		return err
	}

	// Check if already known
	if err := checkIfAlreadyKnown(c, spellID, spell.Name); err != nil {
		return err
	}

	// Determine if spell level is too high for the character at this level
	if err := checkSpellLevelAvailable(c, spell.Level); err != nil {
		return err
	}

	// Determine how many spells can be known at this level
	if err := checkSpellSlotLimit(c, spell.Level); err != nil {
		return err
	}

	return nil
}

// ValidateSpells Validates known spells
func ValidateSpells(c *Character) []error {
	var errs []error
	seen := make(map[int]bool)        // track duplicates
	countByLevel := make(map[int]int) // for slot limits

	for _, id := range c.KnownSpells {
		spell, ok := SpellsByID[id]
		if !ok {
			errs = append(errs, fmt.Errorf("spell %d does not exist", id))
			continue
		}

		// Check for duplicates
		if seen[id] {
			errs = append(errs, fmt.Errorf("duplicate spell: %s", spell.Name))
		}
		seen[id] = true

		// Check if character can cast the type
		if err := checkIfCharacterMayCast(c, spell.Type); err != nil {
			errs = append(errs, fmt.Errorf("spell %s: %w", spell.Name, err))
		}

		// Check spell level availability
		if err := checkSpellLevelAvailable(c, spell.Level); err != nil {
			errs = append(errs, fmt.Errorf("spell %s: %w", spell.Name, err))
		}

		// Tally for slot limits
		countByLevel[spell.Level]++
	}

	// Check count limits per level
	for level, count := range countByLevel {
		allowed := GetSpellSlots(c.Class, c.Level, level)
		if count > allowed {
			errs = append(errs, fmt.Errorf("too many known spells at level %d: %d (limit: %d)", level, count, allowed))
		}
	}

	return errs
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

func CastMemorizedSpell(c *Character, spellID int) error {
	for i := range c.MemorizedSpells {
		if c.MemorizedSpells[i].SpellID == spellID && !c.MemorizedSpells[i].Cast {
			c.MemorizedSpells[i].Cast = true
			return nil
		}
	}
	return fmt.Errorf("no uncast memorized copy of spell %d found", spellID)
}

func UncastMemorizedSpell(c *Character, spellID int) error {
	for i := range c.MemorizedSpells {
		if c.MemorizedSpells[i].SpellID == spellID && c.MemorizedSpells[i].Cast {
			c.MemorizedSpells[i].Cast = false
			return nil
		}
	}
	return fmt.Errorf("no cast memorized copy of spell %d found", spellID)
}

func ResetAllMemorizedSpells(c *Character) {
	for i := range c.MemorizedSpells {
		c.MemorizedSpells[i].Cast = false
	}
}

func ValidateMemorizedSpells(c *Character) []error {
	var errs []error
	perLevelCount := map[int]int{}

	for _, ms := range c.MemorizedSpells {
		spell, ok := SpellsByID[ms.SpellID]
		if !ok {
			errs = append(errs, fmt.Errorf("memorized spell %d does not exist", ms.SpellID))
			continue
		}

		// Check it's known
		if !isSpellKnown(c, ms.SpellID) {
			errs = append(errs, fmt.Errorf("memorized spell %s is not known", spell.Name))
		}

		// Check level availability
		if spell.Level > MaxSpellLevelAvailable(c.Class, c.Level) {
			errs = append(errs, fmt.Errorf("memorized spell %s is too high level (level %d)", spell.Name, spell.Level))
		}

		// Tally per-level count
		perLevelCount[spell.Level]++
	}

	// Check slot limits
	for level, count := range perLevelCount {
		allowed := GetSpellSlots(c.Class, c.Level, level)
		if count > allowed {
			errs = append(errs, fmt.Errorf("too many memorized spells at level %d: %d (limit: %d)", level, count, allowed))
		}
	}

	return errs
}
