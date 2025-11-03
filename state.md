# State

**State**

`State` holds a single `Party`.

This `Party` is initially empty.

**Party**

A `Party` holds zero or more `Characters`.

**Character**

A `Character` represents a Dungeons & Dragons character from the 1981 version of the game.

```go
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
```

`CharacterClass` is effectively an `enum` representing one of the D&D classes.

`Alignment` is effectively an `enum` representing one of the D&D alignments.

`ArmorBonus` represents any hidden, innate bonuses to armor class. Most characters do not have one of these, but the Barbarian from Old School Essentials does. It is included to allow that class to be added in the future.

`RolledHitPoints` is the value of the raw, unmodified dice roll determining character hit points before it is modified by `Constitution`.
`MaximumHitPoints` is `RolledHitPoints` after it has been modified by `Constitution`.
`CurrentHitPoints` is the current value of character hit points, which ranges from 0 to `MaximumHitPoints`.

`Items[]` contains items held by the character.

`ArmorID` and `ShieldID` contain the ID of the any actively equipped armor or shield.
A character may only have one suit of armor or shield equipped; this designs enforces that constraint.

`Spellcasting` contains the types of spells a character is able to cast. Most character classes are only able to cast one type of spell, but it's designed this way to allow support for classes that can cast multiple types (such as the original Ranger class from The Strategic Review).

`KnownSpells` contains spells known by the character if they are a Magic User or Elf. These characters only know a limited number of Magic User spells. Clerics know all of Cleric spells.
`MemorizedSpells` contains spells currently memorized (or "prepared" for Clerics) by Magic Users, Elves, and Clerics.
If you "Know" a spell, you have learned how to cast that spell.
If you have "Memorized" a spell, you made all of the preparations that spell at the beginning of the day and are able to cast it at will.

**Item**

An `Item` represents an item held by a `Character`.

```go
type Item struct {
    ID       int
    HolderID int
    Name     string
    Type     ItemType
    Location ItemLocation
    URL      string
}
```

`HolderID` is the ID of the character carrying the item.

`ItemType` is effectively an `enum` representing the subtype of an item. All item types excluding `ItemGeneric` have additional properties.

`ItemLocation` is effectively an `enum` representing the approximate location of an item.
`none`: no location type. The default empty type. It's also a way to hide an item from the players, as all of the categories below will be visible to the players.
`character`: held by a character
`party`: held by an unspecified character within the party
`storage`: owned by the party, but not currently held by anyone in the party. Stored safely somewhere in the game world
`limbo`: not owned by the party, but able to be currently picked up by the party. It exists somewhere in the game world and can be picked up.

`url` represents a link to an item description page somewhere on https://oldschoolessentials.necroticgnome.com/srd/index.php/Main_Page
Many magical items have complex abilities and this allows for one to link their description to the system reference document.
For example, here is a link to a Crystal Ball: https://oldschoolessentials.necroticgnome.com/srd/index.php/Crystal_Ball

**Item Subtypes**

```go
type Weapon struct {
	Item
	Damage      int
	Bonus       int // 0 for non-magical, 1 to 3 for magical
	IsMelee     bool
	IsRanged    bool
	IsTwoHanded bool
	IsBlunt     bool
}
```

`Damage` represents the maximum value of the damage die (4, 6, 8, or 10)

`Bonus` represents the magical bonus to hit and damage
TODO: Complex weapons such as "Sword +1, +3 vs Dragons" will be supported later

```go
type Armor struct {
	Item
	Type  ArmorType
	Bonus int
}
```

`ArmorType` is effectively an `enum` representing the type of armor. Robes or clothing, leather, chainmail, and plate mail are the four types.

`Bonus` represents the magical bonus to armor class.

```go
type Shield struct {
	Item
	Bonus int
}
```

`Bonus` represents the magical bonus to armor class.

```go
type Jewelry struct {
	Item
	ArmorBonus int
	SaveBonus  int
	/* Many pieces of magical jewelry and defensive items do not follow a clear pattern
	I'll account for this later. For now, this accounts for the Ring of Protection */
}
```

`ArmorBonus` represents the magical bonus to armor class.
`SaveBonus` represents the magical bonus to saving throws.
TODO: More complex magical jewelry will be supported later.

```go
type LimitedUseItem struct {
    Item
    Charges         int
    ArcaneAllowed   bool
    DivineAllowed   bool
    NonMagicAllowed bool
}
```

`LimitedUseItem` represents a wide variety of items with limited uses. Potions, scrolls, and magical wands as well as adventuring items like bundles of torches, rations, and iron spikes.

`Charges` lists the number of uses an item has before being exhausted.
By the book, all limited use items are effectively worthless once all charges are expended, but I may add another flag to prevent automatic deletion of an item when it reaches 0 charges.
It makes sense for torches to disappear from your inventory after using the final torch, but it doesn't make sense for a former magical staff to disappear.

**Spell**

A `Spell` represents a spell that is known, but not memorized or prepared, by a spellcaster.

```
type Spell struct {
	ID    int
	Name  string
	Level int
	Type  SpellType
}
```

`SpellType` is effectively an `enum` representing the type of spell, Arcane (Magic User and Elf) or Divine (Cleric).

```
type MemorizedSpell struct {
	SpellID int
	Cast    bool
}
```

`SpellID` is the ID of the spell.

`Cast` represents whether the spell has been cast or not.

TODO: Add URL support to spells

**Registries & Registration**

```go
var ItemsByID = map[int]*Item{}
var WeaponsByID = map[int]*Weapon{}
var ArmorByID = map[int]*Armor{}
var ShieldsByID = map[int]*Shield{}
var JewelryByID = map[int]*Jewelry{}
var LimitedUseItemsByID = map[int]*LimitedUseItem{}
```

Maps are used over arrays and slices as pointer values allow us to modify items directly in Go.

I have functions to register each type of item as well as an unregister function to ensure an item is removed from all registries.