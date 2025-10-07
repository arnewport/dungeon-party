package main

import (
	"fmt"
	"strings"
)

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
