package main

import (
	"fmt"
	"strings"
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
		Strength:         3,
		Intelligence:     3,
		Wisdom:           3,
		Dexterity:        3,
		Constitution:     3,
		Charisma:         3,
		Items:            []int{},
		ArmorID:          0,
		ShieldID:         0,
		ArcaneSpells:     false,
		DivineSpells:     false,
		KnownSpells:      []int{},
		MemorizedSpells:  []MemorizedSpell{},
	}
	p.Characters = append(p.Characters, char)
	return char
}

// DeleteCharacter moves all of the character's items to Limbo, then removes the character.
func DeleteCharacter(charID int, p *Party) error {
	_, err := findChar(p, charID)
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
	// Shouldn't reach here since findChar succeeded, but be defensive.
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
	if patch.ArcaneSpells != nil {
		c.ArcaneSpells = *patch.ArcaneSpells
	}
	if patch.DivineSpells != nil {
		c.DivineSpells = *patch.DivineSpells
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
	Strength         *int
	Intelligence     *int
	Wisdom           *int
	Dexterity        *int
	Constitution     *int
	Charisma         *int
	Items            *[]int
	ArmorID          *int
	ShieldID         *int
	ArcaneSpells     *bool
	DivineSpells     *bool
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
func findChar(p *Party, id int) (*Character, error) {
	for i := range p.Characters {
		if p.Characters[i].ID == id {
			return &p.Characters[i], nil
		}
	}
	return nil, fmt.Errorf("character %d not found", id)
}

// FindItemByID Cycles through all registries, looking for an item
func FindItemByID(id int) (*Item, error) {
	if it, ok := ItemsByID[id]; ok {
		return it, nil
	}
	if w, ok := WeaponsByID[id]; ok {
		return &w.Item, nil
	}
	if a, ok := ArmorByID[id]; ok {
		return &a.Item, nil
	}
	if s, ok := ShieldsByID[id]; ok {
		return &s.Item, nil
	}
	if j, ok := JewelryByID[id]; ok {
		return &j.Item, nil
	}
	if r, ok := RodsWandsStavesByID[id]; ok {
		return &r.Item, nil
	}
	return nil, fmt.Errorf("item %d not found", id)
}

// Inventory utilities
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

// ValidateCharacterInventory Core invariants for a single character (call after mutations or in tests)
func ValidateCharacterInventory(c *Character) error {
	seen := make(map[int]struct{}, len(c.Items))
	if len(c.Items) > 10 {
		return fmt.Errorf("inventory over capacity")
	}
	for _, id := range c.Items {
		if _, err := FindItemByID(id); err != nil {
			return fmt.Errorf("unknown item id %d", id)
		}
		if _, dup := seen[id]; dup {
			return fmt.Errorf("duplicate item id %d", id)
		}
		seen[id] = struct{}{}
	}
	// equip checks
	if c.ArmorID != 0 {
		it := ItemsByID[c.ArmorID]
		if !hasID(c.Items, c.ArmorID) || it.Type != ItemArmor {
			return fmt.Errorf("armor id invalid or not armor")
		}
	}
	if c.ShieldID != 0 {
		it := ItemsByID[c.ShieldID]
		if !hasID(c.Items, c.ShieldID) || it.Type != ItemShield {
			return fmt.Errorf("shield id invalid or not shield")
		}
	}
	return nil
}

// MoveItemToCharacter Moves (do all checks first, then mutate)
func MoveItemToCharacter(itemID, charID int, p *Party) error {
	it, err := FindItemByID(itemID)
	if err != nil {
		return err
	}
	ch, err := findChar(p, charID)
	if err != nil {
		return err
	}
	if hasID(ch.Items, itemID) {
		return nil
	} // already there
	if len(ch.Items) >= 10 {
		return fmt.Errorf("inventory full")
	}

	// if coming from another character, detach there
	if it.Location == LocationCharacter && it.HolderID != charID {
		if prev, _ := findChar(p, it.HolderID); prev != nil {
			prev.Items = removeID(prev.Items, itemID)
			if prev.ArmorID == itemID {
				prev.ArmorID = 0
			}
			if prev.ShieldID == itemID {
				prev.ShieldID = 0
			}
		}
	}

	// mutate
	it.Location = LocationCharacter
	it.HolderID = charID
	ch.Items = append(ch.Items, itemID)
	return nil
}

func MoveItemToParty(itemID int, p *Party) error { return moveItemToBucket(itemID, LocationParty, p) }
func MoveItemToStorage(itemID int, p *Party) error {
	return moveItemToBucket(itemID, LocationStorage, p)
}
func MoveItemToLimbo(itemID int, p *Party) error { return moveItemToBucket(itemID, LocationLimbo, p) }

func moveItemToBucket(itemID int, loc ItemLocation, p *Party) error {
	it, err := FindItemByID(itemID)
	if err != nil {
		return err
	}
	// detach from character if needed
	if it.Location == LocationCharacter {
		if ch, _ := findChar(p, it.HolderID); ch != nil {
			ch.Items = removeID(ch.Items, itemID)
			if ch.ArmorID == itemID {
				ch.ArmorID = 0
			}
			if ch.ShieldID == itemID {
				ch.ShieldID = 0
			}
		}
	}
	it.Location = loc
	it.HolderID = 0
	return nil
}

func EquipArmor(charID, itemID int, p *Party) error {
	ch, _ := findChar(p, charID)
	if ch == nil {
		return fmt.Errorf("character %d not found", charID)
	}
	it, err := FindItemByID(itemID)
	if err != nil {
		return err
	}
	if it.Location != LocationCharacter || it.HolderID != charID {
		return fmt.Errorf("item not owned by character")
	}
	if it.Type != ItemArmor {
		return fmt.Errorf("item %d is not armor", itemID)
	}
	if !hasID(ch.Items, itemID) {
		return fmt.Errorf("item not in inventory list")
	}

	var allowedArmorTypes = map[CharacterClass]map[ArmorType]bool{
		ClassMagicUser: {
			Robes: true,
		},
		ClassThief: {
			Robes:   true,
			Leather: true,
		},
	}

	armor, ok := ArmorByID[itemID]
	if !ok {
		return fmt.Errorf("armor data for item %d not found", itemID)
	}
	if allowed, ok := allowedArmorTypes[ch.Class]; ok {
		if !allowed[armor.Type] {
			return fmt.Errorf("%s cannot wear %s", ch.Class, armor.Type)
		}
	}
	ch.ArmorID = itemID
	return nil
}

func EquipShield(charID, itemID int, p *Party) error {
	ch, _ := findChar(p, charID)
	if ch == nil {
		return fmt.Errorf("character %d not found", charID)
	}
	it, err := FindItemByID(itemID)
	if err != nil {
		return err
	}
	if it.Location != LocationCharacter || it.HolderID != charID {
		return fmt.Errorf("item not owned by character")
	}
	if it.Type != ItemShield {
		return fmt.Errorf("item %d is not a shield", itemID)
	}
	if !hasID(ch.Items, itemID) {
		return fmt.Errorf("item not in inventory list")
	}
	var disallowedShieldClasses = map[CharacterClass]bool{
		ClassMagicUser: true,
		ClassThief:     true,
	}
	if disallowedShieldClasses[ch.Class] {
		return fmt.Errorf("%s cannot equip shield", ch.Class)
	}
	ch.ShieldID = itemID
	return nil
}

// UnequipArmor clears a character's equipped armor, leaving the item in inventory.
func UnequipArmor(charID int, p *Party) error {
	ch, _ := findChar(p, charID)
	if ch == nil {
		return fmt.Errorf("character %d not found", charID)
	}
	if ch.ArmorID == 0 {
		return nil // nothing equipped; no-op
	}
	ch.ArmorID = 0
	return nil
}

// UnequipShield clears a character's equipped shield, leaving the item in inventory.
func UnequipShield(charID int, p *Party) error {
	ch, _ := findChar(p, charID)
	if ch == nil {
		return fmt.Errorf("character %d not found", charID)
	}
	if ch.ShieldID == 0 {
		return nil // nothing equipped; no-op
	}
	ch.ShieldID = 0
	return nil
}

// MoveAllFromCharacter moves every item from the given character to the specified bucket location.
// Allowed targets: LocationParty, LocationStorage, LocationLimbo.
// It clears ArmorID/ShieldID if those items are moved.
func MoveAllFromCharacter(charID int, to ItemLocation, p *Party) error {
	if to == LocationCharacter || to == LocationNone {
		return fmt.Errorf("invalid target location for bulk move: %v", to)
	}
	if to != LocationParty && to != LocationStorage && to != LocationLimbo {
		return fmt.Errorf("unsupported target location: %v", to)
	}

	ch, _ := findChar(p, charID)
	if ch == nil {
		return fmt.Errorf("character %d not found", charID)
	}

	// Copy list so we can mutate ch.Items safely during moves.
	itemIDs := make([]int, len(ch.Items))
	copy(itemIDs, ch.Items)

	for _, id := range itemIDs {
		if err := moveItemToBucket(id, to, p); err != nil {
			return err
		}
	}

	// After successful moves, ch.Items should already be cleared by moveItemToBucket detaching each item.
	// Belt-and-suspenders: ensure it’s empty.
	ch.Items = ch.Items[:0]
	ch.ArmorID = 0
	ch.ShieldID = 0

	return nil
}

// DumpInventory returns a human-readable string listing all items a character holds.
// Marks equipped armor and shield.
func DumpInventory(charID int, p *Party) (string, error) {
	ch, _ := findChar(p, charID)
	if ch == nil {
		return "", fmt.Errorf("character %d not found", charID)
	}

	var b strings.Builder
	_, _ = fmt.Fprintf(&b, "Inventory for %s (ID %d):\n", ch.Name, ch.ID)

	if len(ch.Items) == 0 {
		b.WriteString("  (no items)\n")
		return b.String(), nil
	}

	for _, id := range ch.Items {
		it, err := FindItemByID(id)
		if err != nil {
			_, _ = fmt.Fprintf(&b, "  [Missing item ID %d]\n", id)
			continue
		}
		equipMarker := ""
		switch id {
		case ch.ArmorID:
			equipMarker = " (equipped armor)"
		case ch.ShieldID:
			equipMarker = " (equipped shield)"
		}
		_, _ = fmt.Fprintf(&b, "  - %s [%s]%s\n", it.Name, it.Type, equipMarker)
	}

	return b.String(), nil
}

// ITEM

var nextItemID = 1

func generateUniqueItemID() int {
	id := nextItemID
	nextItemID++
	return id
}
