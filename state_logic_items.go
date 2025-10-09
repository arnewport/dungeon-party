package main

import (
	"fmt"
	"log"
	"strings"
)

// Inventory utilities

const NoItemEquipped = 0

var nextItemID = 1

func generateUniqueItemID() int {
	id := nextItemID
	nextItemID++
	return id
}

// Inventory helpers

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
	if lu, ok := LimitedUseItemsByID[id]; ok {
		return &lu.Item, nil
	}
	return nil, fmt.Errorf("item %d not found", id)
}

func DetachItemFromCharacter(p *Party, itemID int) {
	it, err := FindItemByID(itemID)
	if err != nil {
		log.Printf("warning: attempted to detach unknown item %d", itemID)
		return
	}

	if it.Location != LocationCharacter {
		return // not held by a character — nothing to do
	}

	ch, err := FindChar(p, it.HolderID)
	if err != nil {
		log.Printf("warning: item %d claimed by missing character %d", itemID, it.HolderID)
		return
	}

	// Remove from inventory
	ch.Items = removeID(ch.Items, itemID)

	// Unequip if necessary
	if ch.ArmorID == itemID {
		ch.ArmorID = NoItemEquipped
	}
	if ch.ShieldID == itemID {
		ch.ShieldID = NoItemEquipped
	}

	// Item is now unbound — you'll want to update .Location and .HolderID externally
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
	if c.ArmorID != NoItemEquipped {
		it := ItemsByID[c.ArmorID]
		if !hasID(c.Items, c.ArmorID) || it.Type != ItemArmor {
			return fmt.Errorf("armor id invalid or not armor")
		}
	}
	if c.ShieldID != NoItemEquipped {
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
		return nil // already there
	}
	if len(ch.Items) >= 10 {
		return fmt.Errorf("inventory full")
	}

	// Detach if coming from another character
	if it.Location == LocationCharacter && it.HolderID != charID {
		DetachItemFromCharacter(p, itemID)
	}

	// Mutate
	it.Location = LocationCharacter
	it.HolderID = charID
	ch.Items = append(ch.Items, itemID)

	if err := ValidateCharacterInventory(ch); err != nil {
		return fmt.Errorf("post-move validation failed: %w", err)
	}
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

	// Detach if held by a character
	if it.Location == LocationCharacter {
		DetachItemFromCharacter(p, itemID)
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
	// if a class has no armor restrictions, the following check is skipped
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
	if ch.ArmorID == NoItemEquipped {
		return nil // nothing equipped; no-op
	}
	ch.ArmorID = NoItemEquipped
	return nil
}

// UnequipShield clears a character's equipped shield, leaving the item in inventory.
func UnequipShield(charID int, p *Party) error {
	ch, _ := FindChar(p, charID)
	if ch == nil {
		return fmt.Errorf("character %d not found", charID)
	}
	if ch.ShieldID == NoItemEquipped {
		return nil // nothing equipped; no-op
	}
	ch.ShieldID = NoItemEquipped
	return nil
}

// MoveAllFromCharacter moves every item from the given character to the specified bucket location.
// Allowed targets: LocationParty, LocationStorage, LocationLimbo.
// It clears ArmorID/ShieldID if those items are moved.
// TODO: Make MoveAllFromCharacter atomic
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
	ch.ArmorID = NoItemEquipped
	ch.ShieldID = NoItemEquipped

	return nil
}

// DeleteItem fully removes an item from the system.
// It detaches the item from any character and deletes it from all registries.
// This is a permanent deletion — use with care.
func DeleteItem(itemID int, p *Party) error {
	it, err := FindItemByID(itemID)
	if err != nil {
		return fmt.Errorf("cannot delete item %d: %w", itemID, err)
	}

	// Detach from character if needed
	if it.Location == LocationCharacter {
		DetachItemFromCharacter(p, itemID)
	}

	// Remove from all global registries
	delete(ItemsByID, itemID)
	delete(WeaponsByID, itemID)
	delete(ArmorByID, itemID)
	delete(ShieldsByID, itemID)
	delete(JewelryByID, itemID)
	delete(LimitedUseItemsByID, itemID)

	return nil
}

// Create Items

// newBaseItem creates a base Item with a unique ID.
// It is used internally by all item subtype constructors.
// The returned Item is not registered anywhere — it must be wrapped and inserted by the caller.
func newBaseItem(name string, itemType ItemType, loc ItemLocation) Item {
	return Item{
		ID:       generateUniqueItemID(),
		Name:     name,
		Type:     itemType,
		Location: loc,
		HolderID: 0,
	}
}

func NewGenericItem(name string, loc ItemLocation) *Item {
	it := &Item{
		ID:       generateUniqueItemID(),
		Name:     name,
		Type:     ItemGeneric,
		Location: loc,
		HolderID: 0,
	}
	ItemsByID[it.ID] = it
	return it
}

func NewWeapon(
	name string,
	damage int,
	bonus int,
	isMelee bool,
	isRanged bool,
	isTwoHanded bool,
	isBlunt bool,
	loc ItemLocation,
) *Weapon {
	w := &Weapon{
		Item:        newBaseItem(name, ItemWeapon, loc),
		Damage:      damage,
		Bonus:       bonus,
		IsMelee:     isMelee,
		IsRanged:    isRanged,
		IsTwoHanded: isTwoHanded,
		IsBlunt:     isBlunt,
	}
	WeaponsByID[w.ID] = w
	return w
}

func NewArmor(
	name string,
	armorType ArmorType,
	bonus int,
	loc ItemLocation,
) *Armor {
	a := &Armor{
		Item:  newBaseItem(name, ItemArmor, loc),
		Type:  armorType,
		Bonus: bonus,
	}
	ArmorByID[a.ID] = a
	return a
}

func NewShield(name string, bonus int, loc ItemLocation) *Shield {
	s := &Shield{
		Item:  newBaseItem(name, ItemShield, loc),
		Bonus: bonus,
	}
	ShieldsByID[s.ID] = s
	return s
}

func NewJewelry(name string, armorBonus int, saveBonus int, loc ItemLocation) *Jewelry {
	j := &Jewelry{
		Item:       newBaseItem(name, ItemJewelry, loc),
		ArmorBonus: armorBonus,
		SaveBonus:  saveBonus,
	}
	JewelryByID[j.ID] = j
	return j
}

func NewLimitedUseItem(
	name string,
	charges int,
	arcaneAllowed bool,
	divineAllowed bool,
	nonMagicAllowed bool,
	loc ItemLocation,
) *LimitedUseItem {
	lu := &LimitedUseItem{
		Item:            newBaseItem(name, ItemLimitedUse, loc),
		Charges:         charges,
		ArcaneAllowed:   arcaneAllowed,
		DivineAllowed:   divineAllowed,
		NonMagicAllowed: nonMagicAllowed,
	}
	LimitedUseItemsByID[lu.ID] = lu
	return lu
}

// Edit Items

// EditGenericItem updates an existing generic item by ID.
func EditGenericItem(id int, name, url string) error {
	it, ok := ItemsByID[id]
	if !ok {
		return fmt.Errorf("item %d not found", id)
	}
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("name cannot be empty")
	}
	it.Name = name
	it.URL = url
	return nil
}

// EditWeapon updates an existing weapon item by ID.
// Only editable fields are modified.
func EditWeapon(id int, newDamage, newBonus int, isMelee, isRanged, isTwoHanded, isBlunt bool) error {
	weapon, ok := WeaponsByID[id]
	if !ok {
		return fmt.Errorf("weapon %d not found", id)
	}

	if newDamage < 1 || newDamage > 10 {
		return fmt.Errorf("damage must be between 1 and 10")
	}

	if newBonus < 0 || newBonus > 3 {
		return fmt.Errorf("bonus must be between 0 and 3")
	}

	weapon.Damage = newDamage
	weapon.Bonus = newBonus
	weapon.IsMelee = isMelee
	weapon.IsRanged = isRanged
	weapon.IsTwoHanded = isTwoHanded
	weapon.IsBlunt = isBlunt

	return nil
}

// EditArmor updates an existing armor item by ID.
func EditArmor(id int, armorType ArmorType, bonus int, p *Party) error {
	armor, ok := ArmorByID[id]
	if !ok {
		return fmt.Errorf("armor %d not found", id)
	}

	// Optional: Validate armor type
	validTypes := map[ArmorType]bool{
		Robes: true, Leather: true, Chain: true, Plate: true,
	}
	if !validTypes[armorType] {
		return fmt.Errorf("invalid armor type: %s", armorType)
	}

	if bonus < 0 || bonus > 3 {
		return fmt.Errorf("bonus must be between 0 and 3")
	}

	// Check if the armor is currently equipped by a character
	if armor.Location == LocationCharacter && armor.HolderID != 0 {
		ch, err := FindChar(p, armor.HolderID)
		if err != nil {
			return fmt.Errorf("armor %d is equipped by missing character %d", id, armor.HolderID)
		}

		// Determine if this class is still allowed to wear the edited armor
		allowedArmorTypes := map[CharacterClass]map[ArmorType]bool{
			ClassMagicUser: {Robes: true},
			ClassThief:     {Robes: true, Leather: true},
		}

		if allowed, ok := allowedArmorTypes[ch.Class]; ok {
			if !allowed[armorType] {
				return fmt.Errorf(
					"cannot change armor: %s is wearing this item and cannot wear %s",
					ch.Class, armorType,
				)
			}
		}
	}

	armor.Type = armorType
	armor.Bonus = bonus

	return nil
}

// EditShield updates an existing shield item by ID.
func EditShield(id int, bonus int) error {
	shield, ok := ShieldsByID[id]
	if !ok {
		return fmt.Errorf("shield %d not found", id)
	}
	if bonus < 0 || bonus > 3 {
		return fmt.Errorf("bonus must be between 0 and 3")
	}
	shield.Bonus = bonus
	return nil
}

// EditJewelry updates an existing jewelry item by ID.
func EditJewelry(id, armorBonus, saveBonus int) error {
	j, ok := JewelryByID[id]
	if !ok {
		return fmt.Errorf("jewelry %d not found", id)
	}
	if armorBonus < 0 || armorBonus > 3 {
		return fmt.Errorf("armor bonus must be between 0 and 3")
	}
	if saveBonus < 0 || saveBonus > 3 {
		return fmt.Errorf("save bonus must be between 0 and 3")
	}
	j.ArmorBonus = armorBonus
	j.SaveBonus = saveBonus
	return nil
}

// EditLimitedUseItem updates an existing limited use item by ID.
func EditLimitedUseItem(id int, charges int, arcane, divine, nonmagic bool) error {
	lu, ok := LimitedUseItemsByID[id]
	if !ok {
		return fmt.Errorf("limited-use item %d not found", id)
	}
	if charges < 0 {
		return fmt.Errorf("charges cannot be negative")
	}
	lu.Charges = charges
	lu.ArcaneAllowed = arcane
	lu.DivineAllowed = divine
	lu.NonMagicAllowed = nonmagic
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
