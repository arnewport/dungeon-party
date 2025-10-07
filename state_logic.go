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
	if r, ok := LimitedUseItemsByID[id]; ok {
		return &r.Item, nil
	}
	return nil, fmt.Errorf("item %d not found", id)
}

// Inventory utilities
var nextItemID = 1

func generateUniqueItemID() int {
	id := nextItemID
	nextItemID++
	return id
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

// ITEM LOGIC

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
	ch, err := FindChar(p, charID)
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
		if prev, _ := FindChar(p, it.HolderID); prev != nil {
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

func MoveItemToParty(itemID int, p *Party) error { return MoveItemToBucket(itemID, LocationParty, p) }
func MoveItemToStorage(itemID int, p *Party) error {
	return MoveItemToBucket(itemID, LocationStorage, p)
}
func MoveItemToLimbo(itemID int, p *Party) error { return MoveItemToBucket(itemID, LocationLimbo, p) }

func MoveItemToBucket(itemID int, loc ItemLocation, p *Party) error {
	it, err := FindItemByID(itemID)
	if err != nil {
		return err
	}
	// detach from character if needed
	if it.Location == LocationCharacter {
		if ch, _ := FindChar(p, it.HolderID); ch != nil {
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
	ch, _ := FindChar(p, charID)
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
	ch, _ := FindChar(p, charID)
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
	ch, _ := FindChar(p, charID)
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
	ch, _ := FindChar(p, charID)
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

	ch, _ := FindChar(p, charID)
	if ch == nil {
		return fmt.Errorf("character %d not found", charID)
	}

	// Copy list so we can mutate ch.Items safely during moves.
	itemIDs := make([]int, len(ch.Items))
	copy(itemIDs, ch.Items)

	for _, id := range itemIDs {
		if err := MoveItemToBucket(id, to, p); err != nil {
			return err
		}
	}

	// After successful moves, ch.Items should already be cleared by MoveItemToBucket detaching each item.
	// Belt-and-suspenders: ensure it’s empty.
	ch.Items = ch.Items[:0]
	ch.ArmorID = 0
	ch.ShieldID = 0

	return nil
}

// DumpInventory returns a human-readable string listing all items a character holds.
// Marks equipped armor and shield.
func DumpInventory(charID int, p *Party) (string, error) {
	ch, _ := FindChar(p, charID)
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
