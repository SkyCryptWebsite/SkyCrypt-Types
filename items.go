package skycrypttypes

type EncodedItems struct {
	Type int    `json:"type"`
	Data string `json:"data"`
}

type Item struct {
	Count         *int    `nbt:"Count" json:"Count,omitempty"`
	Damage        *int    `nbt:"Damage" json:"Damage,omitempty"`
	ID            *int    `nbt:"id" json:"id,omitempty"`
	Tag           *Tag    `nbt:"tag" json:"tag,omitempty"`
	ContainsItems []*Item `json:"containsItems,omitempty"`
	Price         float64 `json:"price,omitempty"`
}

type Tag struct {
	// HideFlags       int             `nbt:"HideFlags" json:"HideFlags,omitempty"`
	// Unbreakable     int             `nbt:"Unbreakable" json:"Unbreakable,omitempty"`
	// Enchantments    []Enchantment   `nbt:"ench" json:"ench,omitempty"`
	ExtraAttributes *ExtraAttributes `nbt:"ExtraAttributes" json:"ExtraAttributes,omitempty"`
	Display         Display          `nbt:"display" json:"display"`
	SkullOwner      *SkullOwner      `nbt:"SkullOwner" json:"SkullOwner,omitempty"`
	ItemModel       string           `nbt:"ItemModel" json:"ItemModel,omitempty"`
}

func (t *Tag) ToMap() TextureItemExtraAttributes {
	if t.ExtraAttributes == nil {
		t.ExtraAttributes = &ExtraAttributes{}
	}

	return TextureItemExtraAttributes{
		ExtraAttributes: t.ExtraAttributes.ToMap(),
		Display:         t.Display,
		SkullOwner:      t.SkullOwner,
		ItemModel:       t.ItemModel,
	}
}

type TextureItemExtraAttributes struct {
	ExtraAttributes map[string]any `nbt:"ExtraAttributes" json:"ExtraAttributes,omitempty"`
	Display         Display        `nbt:"display" json:"display"`
	SkullOwner      *SkullOwner    `nbt:"SkullOwner" json:"SkullOwner,omitempty"`
	ItemModel       string         `nbt:"ItemModel" json:"ItemModel,omitempty"`
}

type ExtraAttributes struct {
	Id                string   `nbt:"id" json:"id,omitempty"`
	Timestamp         any      `nbt:"timestamp" json:"timestamp,omitempty"`
	Uuid              string   `nbt:"uuid" json:"uuid,omitempty"`
	AbilityScroll     []string `nbt:"ability_scroll" json:"ability_scroll,omitempty"`
	AdditionalCoins   int64    `nbt:"additional_coins" json:"additional_coins,omitempty"`
	ArtOfPeaceApplied int      `nbt:"artOfPeaceApplied" json:"artOfPeaceApplied,omitempty"`
	ArtOfWarCount     int      `nbt:"art_of_war_count" json:"art_of_war_count,omitempty"`
	Boosters          []string `nbt:"boosters" json:"boosters,omitempty"`
	// BoosterTiers replaces Boosters on newer items: each applied booster and the tier it was upgraded to.
	BoosterTiers           map[string]int `nbt:"booster_tiers" json:"booster_tiers,omitempty"`
	DivanPowderCoating     int            `nbt:"divan_powder_coating" json:"divan_powder_coating,omitempty"`
	DungeonItemLevel       any            `nbt:"dungeon_item_level" json:"dungeon_item_level,omitempty"`
	DyeItem                string         `nbt:"dye_item" json:"dye_item,omitempty"`
	Enchantments           map[string]int `nbt:"enchantments" json:"enchantments,omitempty"`
	Ethermerge             int            `nbt:"ethermerge" json:"ethermerge,omitempty"`
	FarmingForDummies      int            `nbt:"farming_for_dummies_count" json:"farming_for_dummies_count,omitempty"`
	Gems                   map[string]any `nbt:"gems" json:"gems,omitempty"`
	HotPotatoCount         int            `nbt:"hot_potato_count" json:"hot_potato_count,omitempty"`
	JalapenoCount          int            `nbt:"jalapeno_count" json:"jalapeno_count,omitempty"`
	ManaDisintegrator      int            `nbt:"mana_disintegrator_count" json:"mana_disintegrator_count,omitempty"`
	Modifier               string         `nbt:"modifier" json:"modifier,omitempty"`
	NewYearCakeBagData     []byte         `nbt:"new_year_cake_bag_data" json:"new_year_cake_bag_data,omitempty"`
	NewYearCakeBagYears    []int          `nbt:"new_year_cake_bag_years" json:"new_year_cake_bag_years,omitempty"`
	NewYearsCake           int            `nbt:"new_years_cake" json:"new_years_cake,omitempty"`
	PartyHatEmoji          string         `nbt:"party_hat_emoji" json:"party_hat_emoji,omitempty"`
	PartyHatColor          string         `nbt:"party_hat_color" json:"party_hat_color,omitempty"`
	PetInfo                string         `nbt:"petInfo" json:"petInfo,omitempty"`
	PickonimbusDurability  int            `nbt:"pickonimbus_durability" json:"pickonimbus_durability,omitempty"`
	Polarvoid              int            `nbt:"polarvoid" json:"polarvoid,omitempty"`
	Price                  int64          `nbt:"price" json:"price,omitempty"`
	PowerAbilityScroll     string         `nbt:"power_ability_scroll" json:"power_ability_scroll,omitempty"`
	Recombobulated         int            `nbt:"rarity_upgrades" json:"rarity_upgrades,omitempty"`
	Runes                  map[string]int `nbt:"runes" json:"runes,omitempty"`
	SackPss                int            `nbt:"sack_pss" json:"sack_pss,omitempty"`
	Skin                   string         `nbt:"skin" json:"skin,omitempty"`
	TalismanEnrichment     string         `nbt:"talisman_enrichment" json:"talisman_enrichment,omitempty"`
	ThunderCharge          int            `nbt:"thunder_charge" json:"thunder_charge,omitempty"`
	TunedTransmission      int            `nbt:"tuned_transmission" json:"tuned_transmission,omitempty"`
	UpgradeLevel           any            `nbt:"upgrade_level" json:"upgrade_level,omitempty"`
	WinningBid             int64          `nbt:"winning_bid" json:"winning_bid,omitempty"`
	WoodSingularityCount   int            `nbt:"wood_singularity_count" json:"wood_singularity_count,omitempty"`
	DrillPartUpgradeModule string         `nbt:"drill_part_upgrade_module" json:"drill_part_upgrade_module,omitempty"`
	DrillPartFuelTank      string         `nbt:"drill_part_fuel_tank" json:"drill_part_fuel_tank,omitempty"`
	DrillPartEngine        string         `nbt:"drill_part_engine" json:"drill_part_engine,omitempty"`
	Line                   RodPart        `nbt:"line" json:"line"`
	Hook                   RodPart        `nbt:"hook" json:"hook"`
	Sinker                 RodPart        `nbt:"sinker" json:"sinker"`
	Soulbound              bool           `nbt:"donated_museum" json:"donated_museum,omitempty"`
	Attributes             map[string]int `nbt:"attributes" json:"attributes,omitempty"`
	Edition                int            `nbt:"edition" json:"edition,omitempty"`
	Shiny                  bool           `nbt:"is_shiny" json:"is_shiny,omitempty"`
	ItemTier               int            `nbt:"item_tier" json:"item_tier,omitempty"`
	HecatombSRuns          int            `nbt:"hecatomb_s_runs" json:"hecatomb_s_runs,omitempty"`
	ChampionCombatXP       float64        `nbt:"champion_combat_xp" json:"champion_combat_xp,omitempty"`
	FarmedCultivating      int            `nbt:"farmed_cultivating" json:"farmed_cultivating,omitempty"`
	ExpertiseKills         int            `nbt:"expertise_kills" json:"expertise_kills,omitempty"`
	CompactBlocks          int            `nbt:"compact_blocks" json:"compact_blocks,omitempty"`
	Auction                int64          `nbt:"auction" json:"auction,omitempty"`
	Bid                    int64          `nbt:"bid" json:"bid,omitempty"`
	Model                  string         `nbt:"model" json:"model,omitempty"`
	CollectedCoins         int64          `nbt:"collected_coins" json:"collected_coins,omitempty"`
	Overclocker3000        int            `nbt:"levelable_overclocks" json:"levelable_overclocks,omitempty"`
}
type RodPart struct {
	Part      string `nbt:"part" json:"part,omitempty"`
	Soulbound bool   `nbt:"donated_museum" json:"donated_museum,omitempty"`
}

func (t *ExtraAttributes) ToMap() map[string]any {
	return map[string]any{
		"id":                  t.Id,
		"uuid":                t.Uuid,
		"timestamp":           t.Timestamp,
		"recombobulated":      t.Recombobulated,
		"enchantments":        t.Enchantments,
		"gems":                t.Gems,
		"hecatomb_s_runs":     t.HecatombSRuns,
		"champion_combat_xp":  t.ChampionCombatXP,
		"farmed_cultivating":  t.FarmedCultivating,
		"expertise_kills":     t.ExpertiseKills,
		"compact_blocks":      t.CompactBlocks,
		"modifier":            t.Modifier,
		"model":               t.Model,
		"talisman_enrichment": t.TalismanEnrichment,
	}
}

type Display struct {
	Name  string   `nbt:"Name" json:"Name,omitempty"`
	Lore  []string `nbt:"Lore" json:"Lore,omitempty"`
	Color int      `nbt:"color" json:"color,omitempty"`
}

type Enchantment struct {
	ID    int `nbt:"id" json:"id,omitempty"`
	Level int `nbt:"lvl" json:"lvl,omitempty"`
}

type SkullOwner struct {
	ID         string     `nbt:"Id" json:"Id,omitempty"`
	Properties Properties `nbt:"Properties" json:"Properties"`
}

type Properties struct {
	Textures []Texture `nbt:"textures" json:"textures,omitempty"`
}

type Texture struct {
	Value     string `nbt:"Value" json:"Value,omitempty"`
	Signature string `nbt:"Signature" json:"Signature,omitempty"`
}
