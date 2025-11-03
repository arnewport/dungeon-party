package main

import (
	"fmt"
)

// FUNCTIONS

// CHARACTER

var nextCharacterID = 1

func generateUniqueCharacterID() int {
	id := nextCharacterID
	nextCharacterID++
	return id
}

func AddCharacter(p *Party, name string) Character {
	newID := generateUniqueCharacterID()
	char := Character{
		ID:               newID,
		Name:             name,
		Class:            ClassNone,
		Level:            1,
		Alignment:        AlignmentNone,
		ArmorBonus:       0,
		RolledHitPoints:  0,
		CurrentHitPoints: 0,
		MaximumHitPoints: 0,
		Strength:         3,
		Intelligence:     3,
		Wisdom:           3,
		Dexterity:        3,
		Constitution:     3,
		Charisma:         3,
		Items:            []int{},
		ArmorID:          0,
		ShieldID:         0,
		Spellcasting:     []SpellType{},
		KnownSpells:      []int{},
		MemorizedSpells:  []MemorizedSpell{},
	}
	p.Characters = append(p.Characters, char)
	return char
}

// DeleteCharacter moves all of the character's items to Limbo, then removes the character.
func DeleteCharacter(charID int, p *Party) error {
	_, err := FindChar(p, charID)
	if err != nil {
		return err
	}

	// Re-home all items first (clears equips as a side effect).
	if err := MoveAllFromCharacter(charID, LocationLimbo, p); err != nil {
		return err
	}

	// Remove character from party slice.
	for i := range p.Characters {
		if p.Characters[i].ID == charID {
			p.Characters = append(p.Characters[:i], p.Characters[i+1:]...)
			return nil
		}
	}
	// Shouldn't reach here since FindChar succeeded, but be defensive.
	return fmt.Errorf("character %d disappeared during deletion", charID)
}

func ApplyCharacterPatch(c *Character, patch CharacterPatch) {
	if patch.Name != nil {
		c.Name = *patch.Name
	}
	if patch.Class != nil {
		c.Class = *patch.Class
	}
	if patch.Level != nil {
		c.Level = *patch.Level
	}
	if patch.Alignment != nil {
		c.Alignment = *patch.Alignment
	}
	if patch.ArmorBonus != nil {
		c.ArmorBonus = *patch.ArmorBonus
	}
	if patch.RolledHitPoints != nil {
		c.RolledHitPoints = *patch.RolledHitPoints
	}
	if patch.CurrentHitPoints != nil {
		c.CurrentHitPoints = *patch.CurrentHitPoints
	}
	if patch.MaximumHitPoints != nil {
		c.MaximumHitPoints = *patch.MaximumHitPoints
	}
	if patch.Strength != nil {
		c.Strength = *patch.Strength
	}
	if patch.Intelligence != nil {
		c.Intelligence = *patch.Intelligence
	}
	if patch.Wisdom != nil {
		c.Wisdom = *patch.Wisdom
	}
	if patch.Dexterity != nil {
		c.Dexterity = *patch.Dexterity
	}
	if patch.Constitution != nil {
		c.Constitution = *patch.Constitution
	}
	if patch.Charisma != nil {
		c.Charisma = *patch.Charisma
	}
	if patch.Items != nil {
		c.Items = *patch.Items
	}
	if patch.ArmorID != nil {
		c.ArmorID = *patch.ArmorID
	}
	if patch.ShieldID != nil {
		c.ShieldID = *patch.ShieldID
	}
	if patch.Spellcasting != nil {
		c.Spellcasting = *patch.Spellcasting
	}
	if patch.KnownSpells != nil {
		c.KnownSpells = *patch.KnownSpells
	}
	if patch.MemorizedSpells != nil {
		c.MemorizedSpells = *patch.MemorizedSpells
	}
}

type CharacterPatch struct {
	Name             *string
	Class            *CharacterClass
	Level            *int
	Alignment        *Alignment
	ArmorBonus       *int
	RolledHitPoints  *int
	CurrentHitPoints *int
	MaximumHitPoints *int
	Strength         *int
	Intelligence     *int
	Wisdom           *int
	Dexterity        *int
	Constitution     *int
	Charisma         *int
	Items            *[]int
	ArmorID          *int
	ShieldID         *int
	Spellcasting     *[]SpellType
	KnownSpells      *[]int
	MemorizedSpells  *[]MemorizedSpell
}

func ValidateCharacterPatch(p CharacterPatch) error {
	// validate identity
	validClassesAndLevels := map[CharacterClass]int{
		ClassCleric:    14,
		ClassFighter:   14,
		ClassMagicUser: 14,
		ClassThief:     14,
		ClassDwarf:     12,
		ClassElf:       10,
		ClassHalfling:  8,
	}

	if p.Name != nil && len(*p.Name) > 50 {
		return fmt.Errorf("name is too long")
	}

	if p.Class != nil {
		if _, ok := validClassesAndLevels[*p.Class]; !ok {
			return fmt.Errorf("invalid class: %q", *p.Class)
		}
	}

	if p.Level != nil && p.Class != nil {
		if *p.Level < 1 {
			return fmt.Errorf("level must be at least 1")
		}
		if maxLevel, ok := validClassesAndLevels[*p.Class]; ok && *p.Level > maxLevel {
			return fmt.Errorf("%s cannot exceed level %d", *p.Class, maxLevel)
		}
	}
	if p.Alignment != nil {
		switch *p.Alignment {
		case AlignmentLawful, AlignmentNeutral, AlignmentChaotic:
		default:
			return fmt.Errorf("invalid alignment: %q", *p.Alignment)
		}
	}
	// validate general statistics
	if p.RolledHitPoints != nil && *p.RolledHitPoints <= 0 {
		return fmt.Errorf("rolled hit points must be greater than 0")
	}
	if p.CurrentHitPoints != nil {
		if *p.CurrentHitPoints < 0 {
			return fmt.Errorf("current hit points cannot be negative")
		}
		if p.RolledHitPoints != nil && *p.CurrentHitPoints > *p.RolledHitPoints {
			return fmt.Errorf("current hit points cannot exceed rolled hit points")
		}
	}
	// validate ability scores
	if p.Strength != nil && (*p.Strength < 3 || *p.Strength > 18) {
		return fmt.Errorf("strength must be between 3 and 18")
	}
	if p.Intelligence != nil && (*p.Intelligence < 3 || *p.Intelligence > 18) {
		return fmt.Errorf("intelligence must be between 3 and 18")
	}
	if p.Wisdom != nil && (*p.Wisdom < 3 || *p.Wisdom > 18) {
		return fmt.Errorf("wisdom must be between 3 and 18")
	}
	if p.Dexterity != nil && (*p.Dexterity < 3 || *p.Dexterity > 18) {
		return fmt.Errorf("dexterity must be between 3 and 18")
	}
	if p.Constitution != nil && (*p.Constitution < 3 || *p.Constitution > 18) {
		return fmt.Errorf("constitution must be between 3 and 18")
	}
	if p.Charisma != nil && (*p.Charisma < 3 || *p.Charisma > 18) {
		return fmt.Errorf("charisma must be between 3 and 18")
	}
	// Optional: validate item IDs, spell IDs, etc.
	return nil
}

// Lookups

// FindChar Cycles through the party, looking for a character
func FindChar(p *Party, id int) (*Character, error) {
	for i := range p.Characters {
		if p.Characters[i].ID == id {
			return &p.Characters[i], nil
		}
	}
	return nil, fmt.Errorf("character %d not found", id)
}

func hasID(xs []int, id int) bool {
	for _, v := range xs {
		if v == id {
			return true
		}
	}
	return false
}
func removeID(xs []int, id int) []int {
	out := xs[:0]
	for _, v := range xs {
		if v != id {
			out = append(out, v)
		}
	}
	return out
}
