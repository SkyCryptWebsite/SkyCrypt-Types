# SkyCrypt-Types

[![Go Reference](https://pkg.go.dev/badge/github.com/SkyCryptWebsite/SkyCrypt-Types.svg)](https://pkg.go.dev/github.com/SkyCryptWebsite/SkyCrypt-Types)

Go types for the [Hypixel API](https://api.hypixel.net/)'s SkyBlock data: profiles and their members, museums,
gardens, players and inventory items. They are shared by
[SkyCrypt-Backend](https://github.com/SkyCryptWebsite/SkyCrypt-Backend) and
[SkyHelper-Networth-Go](https://github.com/SkyCryptWebsite/SkyHelper-Networth-Go), so data decoded once can be passed
between them.

The module has no dependencies. Fields are tagged for `encoding/json` (Hypixel's responses) and, for items, `nbt`
(the decoded inventory data), and work with any JSON library that honors those tags.

## Installation

```bash
go get github.com/SkyCryptWebsite/SkyCrypt-Types
```

## What's in it

| Type | Hypixel endpoint | Notes |
|---|---|---|
| `Profile` | `/v2/skyblock/profiles` (each of `profiles`) | `Members` is keyed by player UUID without dashes. |
| `Member` | a member of a profile | Skills, slayers, dungeons, inventories, pets, collections, Rift, garden progress, ... |
| `Museum` | `/v2/skyblock/museum` (each of `members`) | Donated items, keyed by item ID. |
| `Garden` | `/v2/skyblock/garden` (`garden`) | Plots, visitors, crops collected, crop upgrades, composter, garden XP. |
| `Player` | `/v2/player` (`player`) | Only what SkyBlock pages use: name, rank, social links, a few achievements. |
| `Item` | one decoded inventory slot | `Tag.ExtraAttributes` holds the SkyBlock data (ID, enchantments, gems, ...). |
| `EncodedItems` | any inventory field | Base64 of gzipped NBT, as Hypixel sends it. See [Decoding inventories](#decoding-inventories). |

The responses' outer objects (`success`, `profiles`, `members`, ...) aren't modelled, so you can keep only what you
read. Wrap the types in a struct of your own:

```go
type ProfilesResponse struct {
	Success  bool                    `json:"success"`
	Cause    string                  `json:"cause,omitempty"`
	Profiles []skycrypttypes.Profile `json:"profiles"`
}

type MuseumResponse struct {
	Success bool                             `json:"success"`
	Members map[string]*skycrypttypes.Museum `json:"members"`
}

type GardenResponse struct {
	Success bool                 `json:"success"`
	Garden  skycrypttypes.Garden `json:"garden"`
}
```

## Reading a profile

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	skycrypttypes "github.com/SkyCryptWebsite/SkyCrypt-Types"
)

func main() {
	// The body of https://api.hypixel.net/v2/skyblock/profiles?uuid=<uuid>
	body, err := os.ReadFile("profiles.json")
	if err != nil {
		log.Fatal(err)
	}
	var response struct {
		Profiles []skycrypttypes.Profile `json:"profiles"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatal(err)
	}

	const uuid = "4855c53ee4fb4100997600a92fc50984"
	for _, profile := range response.Profiles {
		member, ok := profile.Members[uuid]
		if !ok {
			continue
		}
		fmt.Printf("%s (%s, selected: %v)\n", profile.CuteName, gameMode(profile), profile.Selected)
		fmt.Printf("  SkyBlock level: %d\n", member.Leveling.Experience/100)
		fmt.Printf("  Purse: %.0f coins\n", member.Currencies.CoinPurse)
		if member.PlayerData != nil && member.PlayerData.Experience != nil {
			fmt.Printf("  Farming XP: %.0f\n", member.PlayerData.Experience.SkillFarming)
		}
		for boss, data := range member.Slayer.SlayerBosses {
			fmt.Printf("  %s slayer XP: %.0f\n", boss, data.Experience)
		}
		if member.Pets != nil {
			for _, pet := range member.Pets.Pets {
				if pet.Active {
					fmt.Printf("  Active pet: %s %s\n", pet.Rarity, pet.Type)
				}
			}
		}
	}
}

// Hypixel leaves game_mode out for normal profiles.
func gameMode(profile skycrypttypes.Profile) string {
	if profile.GameMode == "" {
		return "normal"
	}
	return profile.GameMode
}
```

Pointer fields (`PlayerData`, `Profile`, `Inventory`, `Pets`, `Bestiary`, ...) are nil when Hypixel leaves the
section out, which happens when the player turned that API setting off or never used the feature. Check them before
reading.

### Coop members

Every member of a coop is in `Profile.Members`, including players who left or were only invited. A member isn't (or
is no longer) part of the coop when:

```go
func isRemoved(member skycrypttypes.Member) bool {
	if member.Profile == nil {
		return false
	}
	// Invited, but hasn't accepted yet.
	if member.Profile.CoopInvitation != nil && !member.Profile.CoopInvitation.Confirmed {
		return true
	}
	// Left the coop or was kicked.
	return member.Profile.DeletionNotice != nil && member.Profile.DeletionNotice.Timestamp != 0
}
```

`Member.CoopInvitation` is deprecated and always nil: Hypixel sends the invitation under `profile`.

### Fields with their own decoding

Some sections don't map directly to Go types, so they have custom `UnmarshalJSON` (and `MarshalJSON`) methods:

- **`Loadout.Armor`, `Loadout.Equipment`, `Loadout.Loadouts`**: Hypixel sends `{"equipped_set": 2, "1": {...},
  "2": {...}}`. They decode into `EquippedSet` plus `Sets map[int]...` keyed by the set number.
- **`Bestiary`**: only the `kills` counts are kept, as `map[string]int` (mob ID and level, such as `"automaton_150"`).
- **`TuningData`** (accessory bag tuning): slots by number, with their stat points and purchase times.

## Decoding inventories

Inventories (`Member.Inventory.Inventory`, `Enderchest`, `Backpack`, `BagContents.TalismanBag`, wardrobe sets, museum
items, ...) are `EncodedItems`: base64 of gzipped [NBT](https://minecraft.wiki/w/NBT_format). Decode them with an NBT
library into `[]*Item`; the item list is under the `i` tag. With [go-mc](https://github.com/Tnze/go-mc):

```go
import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io"

	skycrypttypes "github.com/SkyCryptWebsite/SkyCrypt-Types"
	"github.com/Tnze/go-mc/nbt"
)

type decodedInventory struct {
	Items []*skycrypttypes.Item `nbt:"i"`
}

func decodeItems(encoded skycrypttypes.EncodedItems) ([]*skycrypttypes.Item, error) {
	if encoded.Data == "" {
		return nil, nil
	}
	compressed, err := base64.StdEncoding.DecodeString(encoded.Data)
	if err != nil {
		return nil, err
	}
	reader, err := gzip.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return nil, err
	}
	// Inflate first: NBT decoders read one byte at a time, which is slow straight from a gzip stream.
	raw, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	var inventory decodedInventory
	if _, err := nbt.NewDecoder(bytes.NewReader(raw)).Decode(&inventory); err != nil {
		return nil, err
	}
	return inventory.Items, nil
}
```

Empty slots decode as items with no `Tag`. For SkyBlock items:

```go
items, err := decodeItems(member.Inventory.Inventory)
if err != nil {
	log.Fatal(err)
}
for slot, item := range items {
	if item == nil || item.Tag == nil || item.Tag.ExtraAttributes == nil {
		continue // empty slot or a vanilla item
	}
	attributes := item.Tag.ExtraAttributes
	fmt.Printf("slot %d: %s (%s)\n", slot, item.Tag.Display.Name, attributes.Id)
	for enchantment, level := range attributes.Enchantments {
		fmt.Printf("  %s %d\n", enchantment, level)
	}
}
```

To price items too, [SkyHelper-Networth-Go](https://github.com/SkyCryptWebsite/SkyHelper-Networth-Go)'s
`CalculateFromSpecifiedInventories` decodes `EncodedItems` and returns each item with its price.

## Adding fields

Hypixel adds data often. To read something new:

1. Add the field to the matching struct with its exact Hypixel JSON key (and an `nbt` tag for item data).
2. Use a pointer or a map when Hypixel may leave it out, so "missing" and "zero" stay apart.
3. Keep it backwards compatible: other modules decode the same responses with these types. Deprecate a field
   instead of removing it.
