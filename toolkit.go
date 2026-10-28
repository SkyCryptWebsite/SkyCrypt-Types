package skycrypttypes

import (
	"encoding/json"
	"strconv"
)

// Toolkit is a farming or hunting toolkit: tools kept per category, like the garden's farming tools or the hunting
// scythe, net and traps. Hypixel sends
//
//	{"IS_UNLOCKED": true, "IN_USE": {"TRAP": {"0": true}}, "TRAP": [{"type": 0, "data": "<base64 NBT>"}, ...], ...}
//
// Each slot's data is a single item's ExtraAttributes, not an inventory: decode it into ExtraAttributes, not into a
// list of Items.
type Toolkit struct {
	IsUnlocked bool
	// InUse marks slots whose tool is currently taken out, by category and slot index. That tool is in one of the
	// player's inventories, so the slot here shouldn't be counted again.
	InUse map[string]map[int]bool
	// Slots holds each category's tools by slot index. Empty slots are left out.
	Slots map[string]map[int]EncodedItems
}

func (t *Toolkit) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*t = Toolkit{InUse: map[string]map[int]bool{}, Slots: map[string]map[int]EncodedItems{}}
	for key, value := range raw {
		switch key {
		case "IS_UNLOCKED":
			if err := json.Unmarshal(value, &t.IsUnlocked); err != nil {
				return err
			}
		case "IN_USE":
			var inUse map[string]map[string]bool
			if err := json.Unmarshal(value, &inUse); err != nil {
				return err
			}
			for category, slots := range inUse {
				t.InUse[category] = map[int]bool{}
				for index, used := range slots {
					if i, err := strconv.Atoi(index); err == nil && used {
						t.InUse[category][i] = true
					}
				}
			}
		default:
			slots, err := unmarshalToolkitSlots(value)
			if err != nil {
				return err
			}
			if len(slots) > 0 {
				t.Slots[key] = slots
			}
		}
	}
	return nil
}

// unmarshalToolkitSlots reads a category, which Hypixel sends as a list (index = position) or an object keyed by index.
func unmarshalToolkitSlots(data json.RawMessage) (map[int]EncodedItems, error) {
	slots := map[int]EncodedItems{}
	var list []*EncodedItems
	if err := json.Unmarshal(data, &list); err == nil {
		for i, slot := range list {
			if slot != nil && slot.Data != "" {
				slots[i] = *slot
			}
		}
		return slots, nil
	}

	var object map[string]*EncodedItems
	if err := json.Unmarshal(data, &object); err != nil {
		// Not a category of slots: a field Hypixel added later. Ignore it rather than fail the whole profile.
		return nil, nil
	}
	for index, slot := range object {
		if i, err := strconv.Atoi(index); err == nil && slot != nil && slot.Data != "" {
			slots[i] = *slot
		}
	}
	return slots, nil
}

func (t Toolkit) MarshalJSON() ([]byte, error) {
	out := map[string]any{"IS_UNLOCKED": t.IsUnlocked}
	if len(t.InUse) > 0 {
		inUse := map[string]map[string]bool{}
		for category, slots := range t.InUse {
			inUse[category] = map[string]bool{}
			for index, used := range slots {
				inUse[category][strconv.Itoa(index)] = used
			}
		}
		out["IN_USE"] = inUse
	}
	for category, slots := range t.Slots {
		size := 0
		for index := range slots {
			size = max(size, index+1)
		}
		list := make([]*EncodedItems, size)
		for index, slot := range slots {
			list[index] = &slot
		}
		out[category] = list
	}
	return json.Marshal(out)
}
