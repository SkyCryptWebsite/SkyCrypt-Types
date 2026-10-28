package skycrypttypes

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Profile struct {
	ProfileID         string             `json:"profile_id"`
	CuteName          string             `json:"cute_name"`
	Selected          bool               `json:"selected"`
	Members           map[string]Member  `json:"members"`
	GameMode          string             `json:"game_mode,omitempty"`
	Banking           *Banking           `json:"banking,omitempty"`
	CommunityUpgrades *CommunityUpgrades `json:"community_upgrades,omitempty"`
}

type Member struct {
	PlayerData          *PlayerData            `json:"player_data"`
	Profile             *ProfileData           `json:"profile"`
	JacobsContest       JacobsContest          `json:"jacobs_contest"`
	Pets                *Pets                  `json:"pets_data,omitempty"`
	Leveling            Leveling               `json:"leveling"`
	Currencies          Currencies             `json:"currencies"`
	FairySouls          *FairySouls            `json:"fairy_soul,omitempty"`
	Inventory           *Inventory             `json:"inventory,omitempty"`
	SharedInventory     SharedInventory        `json:"shared_inventory"`
	Rift                Rift                   `json:"rift"`
	AccessoryBagStorage AccessoryBagStorage    `json:"accessory_bag_storage"`
	CrimsonIsle         CrimsonIsleData        `json:"nether_island_player_data"`
	Mining              Mining                 `json:"mining_core"`
	Objectives          *Objectives            `json:"objectives,omitempty"`
	GlaciteTunnels      *GlaciteData           `json:"glacite_player_data,omitempty"`
	Forge               Forge                  `json:"forge"`
	Quests              Quests                 `json:"quests"`
	Garden              GardenProfileData      `json:"garden_player_data"`
	PlayerStats         PlayerStats            `json:"player_stats"`
	TrophyFish          MemberTrophyFish       `json:"trophy_fish"`
	Experimentation     ExperimentationData    `json:"experimentation"`
	Dungeons            Dungeons               `json:"dungeons"`
	Slayer              Slayer                 `json:"slayer"`
	Bestiary            *Bestiary              `json:"bestiary,omitempty"`
	Collections         map[string]int         `json:"collection,omitempty"`
	ItemData            ItemData               `json:"item_data"`
	WinterPlayerData    WinterPlayerIslandData `json:"winter_player_data"`
	SackCounts          map[string]int         `json:"sack_counts"`
	Foraging            Foraging               `json:"foraging"`
	SkillTree           SkillTree              `json:"skill_tree"`
	ForagingCore        ForagingCore           `json:"foraging_core"`
	Shards              Shards                 `json:"shards"`
	Temples             Temples                `json:"temples"`
	Attributes          Attributes             `json:"attributes"`
	Events              Events                 `json:"events"`
	Loadout             Loadouts               `json:"loadout"`
}

type Loadouts struct {
	Armor     ArmorLoadoutSets     `json:"armor"`
	Equipment EquipmentLoadoutSets `json:"equipment"`
	Loadouts  SavedLoadoutSets     `json:"loadouts"`
}

type ArmorLoadoutSets struct {
	EquippedSet int
	Sets        map[int]ArmorLoadout
}

func (l *ArmorLoadoutSets) UnmarshalJSON(data []byte) error {
	sets, equippedSet, err := unmarshalIndexedSets[ArmorLoadout](data)
	if err != nil {
		return err
	}

	l.Sets = sets
	l.EquippedSet = equippedSet

	return nil
}

func (l ArmorLoadoutSets) MarshalJSON() ([]byte, error) {
	return marshalIndexedSets(l.Sets, l.EquippedSet)
}

type EquipmentLoadoutSets struct {
	EquippedSet int
	Sets        map[int]EquipmentLoadout
}

func (l *EquipmentLoadoutSets) UnmarshalJSON(data []byte) error {
	sets, equippedSet, err := unmarshalIndexedSets[EquipmentLoadout](data)
	if err != nil {
		return err
	}

	l.Sets = sets
	l.EquippedSet = equippedSet

	return nil
}

func (l EquipmentLoadoutSets) MarshalJSON() ([]byte, error) {
	return marshalIndexedSets(l.Sets, l.EquippedSet)
}

type SavedLoadoutSets struct {
	Sets map[int]SavedLoadout
}

func (l *SavedLoadoutSets) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	l.Sets = make(map[int]SavedLoadout)
	for key, value := range raw {
		setID, err := strconv.Atoi(key)
		if err != nil {
			continue
		}

		var loadout SavedLoadout
		if err := json.Unmarshal(value, &loadout); err != nil {
			return err
		}

		l.Sets[setID] = loadout
	}

	return nil
}

func (l SavedLoadoutSets) MarshalJSON() ([]byte, error) {
	raw := make(map[string]any, len(l.Sets))
	for setID, loadout := range l.Sets {
		raw[strconv.Itoa(setID)] = loadout
	}

	return json.Marshal(raw)
}

func unmarshalIndexedSets[T any](data []byte) (map[int]T, int, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, 0, err
	}

	sets := make(map[int]T)
	equippedSet := 0
	for key, value := range raw {
		switch key {
		case "equipped_set":
			if err := json.Unmarshal(value, &equippedSet); err != nil {
				return nil, 0, err
			}
		default:
			setID, err := strconv.Atoi(key)
			if err != nil {
				continue
			}

			var loadout T
			if err := json.Unmarshal(value, &loadout); err != nil {
				return nil, 0, err
			}
			sets[setID] = loadout
		}
	}

	return sets, equippedSet, nil
}

func marshalIndexedSets[T any](sets map[int]T, equippedSet int) ([]byte, error) {
	if sets == nil {
		sets = map[int]T{}
	}

	raw := make(map[string]any, len(sets)+1)
	raw["equipped_set"] = equippedSet
	for setID, loadout := range sets {
		raw[strconv.Itoa(setID)] = loadout
	}

	return json.Marshal(raw)
}

type ArmorLoadout struct {
	Id         int          `json:"id"`
	Helmet     EncodedItems `json:"HELMET"`
	Chestplate EncodedItems `json:"CHESTPLATE"`
	Leggings   EncodedItems `json:"LEGGINGS"`
	Boots      EncodedItems `json:"BOOTS"`
}

type EquipmentLoadout struct {
	Id             int          `json:"id"`
	EquipmentSlot1 EncodedItems `json:"EQUIPMENT_SLOT_1"`
	EquipmentSlot2 EncodedItems `json:"EQUIPMENT_SLOT_2"`
	EquipmentSlot3 EncodedItems `json:"EQUIPMENT_SLOT_3"`
	EquipmentSlot4 EncodedItems `json:"EQUIPMENT_SLOT_4"`
}

type SavedLoadout struct {
	ID                       int    `json:"id"`
	Name                     string `json:"name"`
	ArmorSetID               int    `json:"armor_set_id,omitempty"`
	EquipmentSetID           int    `json:"equipment_set_id,omitempty"`
	PowerStoneID             string `json:"power_stone,omitempty"`
	PetUUID                  string `json:"pet,omitempty"`
	TuningPointsSlot         int    `json:"tuning_points_slot,omitempty"`
	MiningCoreSelectedSlot   int    `json:"mining_core_selected_slot,omitempty"`
	ForagingCoreSelectedSlot int    `json:"foraging_core_selected_slot,omitempty"`
}

type Events struct {
	Easter EasterEvent `json:"easter"`
}

type EasterEvent struct {
	RefinedDarkCacaoTruffles int `json:"refined_dark_cacao_truffles,omitempty"`
}

type WinterPlayerIslandData struct {
	RefinedJyrreUses int `json:"refined_jyrre_uses,omitempty"`
}

type ItemData struct {
	Soulflow               float64 `json:"soulflow,omitempty"`
	TeleporterPillConsumed bool    `json:"teleporter_pill_consumed,omitempty"`
}

type Races struct {
	ForagingRaceBestTime float64            `json:"foraging_race_best_time"`
	EndRaceBestTime      float64            `json:"end_race_best_time"`
	ChickenRaceBestTime2 float64            `json:"chicken_race_best_time_2"`
	DungeonHub           map[string]float64 `json:"dungeon_hub"`
	RiftRaceBestTime     float64            `json:"rift_race_best_time"`
}

type CoopInvitation struct {
	Confirmed bool `json:"confirmed,omitempty"`
}

type PlayerData struct {
	Experience         *Experience    `json:"experience"`
	Minions            []string       `json:"crafted_generators"`
	ReaperPeppersEaten int            `json:"reaper_peppers_eaten,omitempty"`
	GardenChips        map[string]int `json:"garden_chips,omitempty"`
}

type Experience struct {
	SkillFishing      float64 `json:"SKILL_FISHING"`
	SkillAlchemy      float64 `json:"SKILL_ALCHEMY"`
	SkillMining       float64 `json:"SKILL_MINING"`
	SkillFarming      float64 `json:"SKILL_FARMING"`
	SkillEnchanting   float64 `json:"SKILL_ENCHANTING"`
	SkillTaming       float64 `json:"SKILL_TAMING"`
	SkillForaging     float64 `json:"SKILL_FORAGING"`
	SkillSocial       float64 `json:"SKILL_SOCIAL"`
	SkillCarpentry    float64 `json:"SKILL_CARPENTRY"`
	SkillCombat       float64 `json:"SKILL_COMBAT"`
	SkillRunecrafting float64 `json:"SKILL_RUNECRAFTING"`
	SkillHunting      float64 `json:"SKILL_HUNTING"`
}

type ProfileData struct {
	DeletionNotice *DeletionNotice `json:"deletion_notice"`
	// CoopInvitation is set for members invited to the coop; Confirmed is false until they accept.
	CoopInvitation      *CoopInvitation `json:"coop_invitation,omitempty"`
	FirstJoin           int64           `json:"first_join,omitempty"`
	BankAccount         float64         `json:"bank_account,omitempty"`
	PersonalBankUpgrade int             `json:"personal_bank_upgrade,omitempty"`
}

type DeletionNotice struct {
	Timestamp int64 `json:"timestamp,omitempty"`
}

type JacobsContest struct {
	Perks          *Perks                      `json:"perks,omitempty"`
	UniqueBrackets map[string][]string         `json:"unique_brackets,omitempty"`
	MedalsInv      map[string]int              `json:"medals_inv,omitempty"`
	Contests       map[string]JacobContestData `json:"contests,omitempty"`
}

type Perks struct {
	FarmingLevelCap int `json:"farming_level_cap,omitempty"`
}

type PetCare struct {
	PetTypesSacrificed []string `json:"pet_types_sacrificed,omitempty"`
}

type Leveling struct {
	Experience int `json:"experience,omitempty"`
}

type Currencies struct {
	CoinPurse  float64            `json:"coin_purse,omitempty"`
	MotesPurse float64            `json:"motes_purse,omitempty"`
	Essence    map[string]Essence `json:"essence,omitempty"`
}

type Essence struct {
	Current int `json:"current,omitempty"`
}

type Banking struct {
	Balance *float64 `json:"balance,omitempty"`
}

type FairySouls struct {
	TotalCollected int `json:"total_collected,omitempty"`
}

type Inventory struct {
	Inventory     EncodedItems            `json:"inv_contents"`
	Enderchest    EncodedItems            `json:"ender_chest_contents"`
	BackpackIcons map[string]EncodedItems `json:"backpack_icons"`
	Armor         EncodedItems            `json:"inv_armor"`
	Equipment     EncodedItems            `json:"equipment_contents"`
	PersonalVault EncodedItems            `json:"personal_vault_contents"`
	Backpack      map[string]EncodedItems `json:"backpack_contents"`
	// Wardrobe      EncodedItems            `json:"wardrobe_contents"`
	BagContents BagContents    `json:"bag_contents"`
	Sacks       map[string]int `json:"sacks_counts"`
}

type SharedInventory struct {
	CandyInventory        EncodedItems `json:"candy_inventory_contents"`
	CarnivalMaskInventory EncodedItems `json:"carnival_mask_inventory_contents"`
}

type BagContents struct {
	PotionBag   EncodedItems `json:"potion_bag"`
	TalismanBag EncodedItems `json:"talisman_bag"`
	FishingBag  EncodedItems `json:"fishing_bag"`
	SacksBag    EncodedItems `json:"sacks_bag"`
	Quiver      EncodedItems `json:"quiver"`
}

type Rift struct {
	Inventory  RiftInventory  `json:"inventory"`
	Access     RiftAccess     `json:"access"`
	DeadCats   DeadCats       `json:"dead_cats"`
	Enigma     RiftEnigma     `json:"enigma"`
	Castle     RiftCastle     `json:"castle"`
	Gallery    RiftGallery    `json:"gallery"`
	WitherCage RiftWitherCage `json:"wither_cage"`
}

type RiftWitherCage struct {
	KilledEyes []string `json:"killed_eyes,omitempty"`
}

type RiftGallery struct {
	SecuredTrophies []RiftSecuredTrophy `json:"secured_trophies,omitempty"`
}

type RiftSecuredTrophy struct {
	Type      string `json:"type,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

type RiftCastle struct {
	GrubberStacks int `json:"grubber_stacks,omitempty"`
}

type RiftEnigma struct {
	FoundSouls []string `json:"found_souls,omitempty"`
}

type RiftInventory struct {
	Inventory  EncodedItems `json:"inv_contents"`
	Armor      EncodedItems `json:"inv_armor"`
	Enderchest EncodedItems `json:"ender_chest_contents"`
	Equipment  EncodedItems `json:"equipment_contents"`
}

type RiftAccess struct {
	ConsumedPrism bool `json:"consumed_prism,omitempty"`
}

type AccessoryBagStorage struct {
	SelectedPower        string     `json:"selected_power,omitempty"`
	HighestMagicalPower  int        `json:"highest_magical_power,omitempty"`
	UnlockedPowers       []string   `json:"unlocked_powers,omitempty"`
	BagUpgradesPurchased int        `json:"bag_upgrades_purchased,omitempty"`
	Tuning               TuningData `json:"tuning,omitempty"`
}

type TuningData struct {
	HighestUnlockedSlot int
	Slots               map[int]map[string]int
	PurchaseTS          map[int]int64
}

func (t *TuningData) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	t.Slots = make(map[int]map[string]int)
	t.PurchaseTS = make(map[int]int64)
	for key, value := range raw {
		switch {
		case key == "highest_unlocked_slot":
			if err := json.Unmarshal(value, &t.HighestUnlockedSlot); err != nil {
				return err
			}
		case len(key) > 5 && key[:5] == "slot_":
			slotID, err := strconv.Atoi(key[5:])
			if err != nil {
				continue
			}

			var slotRaw map[string]json.RawMessage
			if err := json.Unmarshal(value, &slotRaw); err != nil {
				return err
			}

			slotStats := make(map[string]int)
			for statKey, statValue := range slotRaw {
				if statKey == "purchase_ts" {
					var purchaseTS int64
					if err := json.Unmarshal(statValue, &purchaseTS); err != nil {
						return err
					}
					t.PurchaseTS[slotID] = purchaseTS
					continue
				}

				var stat int
				if err := json.Unmarshal(statValue, &stat); err != nil {
					return err
				}
				slotStats[statKey] = stat
			}

			t.Slots[slotID] = slotStats
		}
	}

	return nil
}

func (t TuningData) MarshalJSON() ([]byte, error) {
	if t.Slots == nil {
		t.Slots = map[int]map[string]int{}
	}

	if t.PurchaseTS == nil {
		t.PurchaseTS = map[int]int64{}
	}

	raw := make(map[string]any, len(t.Slots)+1)
	raw["highest_unlocked_slot"] = t.HighestUnlockedSlot
	for slotID, slotStats := range t.Slots {
		slotRaw := make(map[string]any, len(slotStats)+1)
		if purchaseTS, ok := t.PurchaseTS[slotID]; ok && purchaseTS != 0 {
			slotRaw["purchase_ts"] = purchaseTS
		}

		for statKey, statValue := range slotStats {
			slotRaw[statKey] = statValue
		}

		raw[fmt.Sprintf("slot_%d", slotID)] = slotRaw
	}

	for slotID, purchaseTS := range t.PurchaseTS {
		if _, exists := t.Slots[slotID]; exists || purchaseTS == 0 {
			continue
		}

		raw[fmt.Sprintf("slot_%d", slotID)] = map[string]any{
			"purchase_ts": purchaseTS,
		}
	}

	return json.Marshal(raw)
}

type CrimsonIsleData struct {
	Abiphone            Abiphone       `json:"abiphone"`
	Kuudra              map[string]int `json:"kuudra_completed_tiers,omitempty"`
	Dojo                map[string]int `json:"dojo,omitempty"`
	SelectedFaction     string         `json:"selected_faction,omitempty"`
	MagesReputation     float64        `json:"mages_reputation,omitempty"`
	BarbarianReputation float64        `json:"barbarians_reputation,omitempty"`
}

type Abiphone struct {
	ActiveContacts []string `json:"active_contacts,omitempty"`
}

type DeadCats struct {
	FoundCats []string `json:"found_cats,omitempty"`
	Montezuma Pet      `json:"montezuma"`
}

type Pet struct {
	UUID       string  `json:"uuid,omitempty"`
	UniqueId   string  `json:"uniqueId,omitempty"`
	Type       string  `json:"type,omitempty"`
	Experience float64 `json:"exp,omitempty"`
	Active     bool    `json:"active,omitempty"`
	Rarity     string  `json:"tier,omitempty"`
	HeldItem   string  `json:"heldItem,omitempty"`
	CandyUsed  int     `json:"candyUsed,omitempty"`
	Skin       string  `json:"skin,omitempty"`
}

type Pets struct {
	PetCare *PetCare `json:"pet_care,omitempty"`
	Pets    []Pet    `json:"pets,omitempty"`
}

type Mining struct {
	GreaterMinesLastAccess int64              `json:"greater_mines_last_access,omitempty"`
	PowderMithril          int                `json:"powder_mithril,omitempty"`
	PowderMithrilTotal     int                `json:"powder_mithril_total,omitempty"`
	PowderSpentMithril     int                `json:"powder_spent_mithril,omitempty"`
	PowderGemstone         int                `json:"powder_gemstone,omitempty"`
	PowderGemstoneTotal    int                `json:"powder_gemstone_total,omitempty"`
	PowderSpentGemstone    int                `json:"powder_spent_gemstone,omitempty"`
	PowderGlacite          int                `json:"powder_glacite,omitempty"`
	PowderGlaciteTotal     int                `json:"powder_glacite_total,omitempty"`
	PowderSpentGlacite     int                `json:"powder_spent_glacite,omitempty"`
	Crystals               map[string]Crystal `json:"crystals,omitempty"`
	Biomes                 Biomes             `json:"biomes"`
}

type Crystal struct {
	State       string `json:"state,omitempty"`
	TotalFound  int    `json:"total_found,omitempty"`
	TotalPlaced int    `json:"total_placed,omitempty"`
}

type Biomes struct {
	Precursor Precursor `json:"precursor"`
}

type Precursor struct {
	PartsDelivered []string `json:"parts_delivered,omitempty"`
}

type Objectives struct {
	Tutorial []string `json:"tutorial,omitempty"`
}

type GlaciteData struct {
	FossilsDonated    []string       `json:"fossils_donated,omitempty"`
	FossilDust        float64        `json:"fossil_dust,omitempty"`
	CorpsesLooted     map[string]int `json:"corpses_looted,omitempty"`
	MineshaftsEntered int            `json:"mineshafts_entered,omitempty"`
}

type Forge struct {
	ForgeProcesses ForgeProcesses `json:"forge_processes"`
}

type ForgeProcesses struct {
	Forge map[string]ForgeProcess `json:"forge_1"`
}

type ForgeProcess struct {
	Id        string `json:"id"`
	StartTime int64  `json:"startTime"`
	Slot      int    `json:"slot"`
}

type Quests struct {
	TrapperQuest TrapperQuest `json:"trapper_quest"`
}

type TrapperQuest struct {
	PeltCount    int            `json:"pelt_count,omitempty"`
	LastTaskTime int64          `json:"last_task_time,omitempty"`
	Kills        map[string]any `json:"kills,omitempty"`
}

func (t *TrapperQuest) UnmarshalJSON(data []byte) error {
	type Alias TrapperQuest
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	t.Kills = make(map[string]any)
	if killsRaw, ok := raw["kills"]; ok {
		if killsMap, ok := killsRaw.(map[string]any); ok {
			for k, v := range killsMap {
				switch val := v.(type) {
				case float64:
					t.Kills[k] = int(val)
				case int:
					t.Kills[k] = val
				case string:
					if i, err := strconv.Atoi(val); err == nil {
						t.Kills[k] = i
					} else {
						t.Kills[k] = val
					}
				default:
					t.Kills[k] = val
				}
			}
		}
	}

	return nil
}

type GardenProfileData struct {
	Copper                    int      `json:"copper,omitempty"`
	LarvaConsumed             int      `json:"larva_consumed,omitempty"`
	AnalyzedGreenhouseCrops   []string `json:"analyzed_greenhouse_crops,omitempty"`
	DiscoveredGreenhouseCrops []string `json:"discovered_greenhouse_crops,omitempty"`
	FarmingToolkit            *Toolkit `json:"farming_toolkit,omitempty"`
}

type JacobContestData struct {
	Collected           int    `json:"collected"`
	ClaimedPosition     *int   `json:"claimed_position,omitempty"`
	ClaimedParticipants *int   `json:"claimed_participants,omitempty"`
	ClaimedMedal        string `json:"claimed_medal"`
}

type PlayerStats struct {
	Kills       map[string]float64 `json:"kills,omitempty"`
	Deaths      map[string]float64 `json:"deaths,omitempty"`
	ItemsFished struct {
		Total         float64 `json:"total,omitempty"`
		Normal        float64 `json:"normal,omitempty"`
		Treasure      float64 `json:"treasure,omitempty"`
		LargeTreasure float64 `json:"large_treasure,omitempty"`
		TrophyFish    float64 `json:"trophy_fish,omitempty"`
	} `json:"items_fished"`
	ShredderRod struct {
		Fished float64 `json:"fished,omitempty"`
		Bait   float64 `json:"bait,omitempty"`
	} `json:"shredder_rod"`
	Pets struct {
		Milestone struct {
			SeaCreaturesKilled float64 `json:"sea_creatures_killed,omitempty"`
			OresMined          float64 `json:"ores_mined,omitempty"`
		} `json:"milestone"`
	} `json:"pets"`
	Rift                  RiftPlayerData   `json:"rift"`
	Races                 Races            `json:"races"`
	Gifts                 Gifts            `json:"gifts"`
	WinterIslandData      WinterIslandData `json:"winter"`
	EndIsland             EndIsland        `json:"end_island"`
	HighestCriticalDamage float64          `json:"highest_critical_damage"`
	Mythos                Mythos           `json:"mythos"`
	Auctions              Auctions         `json:"auctions"`
}

type Auctions struct {
	Bids        float64            `json:"bids"`
	HighestBid  float64            `json:"highest_bid"`
	Won         float64            `json:"won"`
	TotalBought map[string]float64 `json:"total_bought"`
	GoldSpent   float64            `json:"gold_spent"`
	Created     float64            `json:"created"`
	Fees        float64            `json:"fees"`
	TotalSold   map[string]float64 `json:"total_sold"`
	GoldEarned  float64            `json:"gold_earned"`
	NoBids      float64            `json:"no_bids"`
}

type Mythos struct {
	Kills                 float64            `json:"kills"`
	BurrowsDugNext        map[string]float64 `json:"burrows_dug_next"`
	BurrowsDugCombat      map[string]float64 `json:"burrows_dug_combat"`
	BurrowsDugTreasure    map[string]float64 `json:"burrows_dug_treasure"`
	BurrowsChainsComplete map[string]float64 `json:"burrows_chains_complete"`
}

type EndIsland struct {
	DragonFight DragonFight `json:"dragon_fight"`
}

type DragonFight struct {
	EnderCrystalsDestroyed float64            `json:"ender_crystals_destroyed"`
	MostDamage             map[string]float64 `json:"most_damage"`
	FastestKill            map[string]float64 `json:"fastest_kill"`
}

type WinterIslandData struct {
	MostSnowballsHit     float64 `json:"most_snowballs_hit"`
	MostDamageDealt      float64 `json:"most_damage_dealt"`
	MostMagmaDamageDealt float64 `json:"most_magma_damage_dealt"`
	MostCannonballsHit   float64 `json:"most_cannonballs_hit"`
}

type Gifts struct {
	Given    float64 `json:"total_given"`
	Received float64 `json:"total_received"`
}

type RiftPlayerData struct {
	Visits                 float64 `json:"visits,omitempty"`
	LifetimeMotesCollected float64 `json:"lifetime_motes_earned,omitempty"`
	MotesOrbPickup         float64 `json:"motes_orb_pickup,omitempty"`
}

type MemberTrophyFish struct {
	Rewards     []int          `json:"rewards"`
	TotalCaught int            `json:"total_caught"`
	Extra       map[string]int `json:"-"`
}

func (t *MemberTrophyFish) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	t.Extra = make(map[string]int)

	for key, value := range raw {
		switch key {
		case "rewards":
			if err := json.Unmarshal(value, &t.Rewards); err != nil {
				return err
			}
		case "total_caught":
			if err := json.Unmarshal(value, &t.TotalCaught); err != nil {
				return err
			}
		default:
			var num int
			if err := json.Unmarshal(value, &num); err == nil {
				t.Extra[key] = num
			}
		}
	}

	return nil
}

type ExperimentationGame struct {
	LastAttempt int64            `json:"last_attempt,omitempty"`
	LastClaimed int64            `json:"last_claimed,omitempty"`
	BonusClicks int              `json:"bonus_clicks,omitempty"`
	Claimed     bool             `json:"claimed,omitempty"`
	Attempts    map[int]int      `json:"-"`
	Claims      map[int]int      `json:"-"`
	BestScores  map[int]int      `json:"-"`
	Raw         map[string]int64 `json:"-"`
}

func (e *ExperimentationGame) UnmarshalJSON(data []byte) error {
	type Alias ExperimentationGame
	aux := &struct {
		*Alias
	}{Alias: (*Alias)(e)}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	e.Attempts = make(map[int]int)
	e.Claims = make(map[int]int)
	e.BestScores = make(map[int]int)
	e.Raw = make(map[string]int64)

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	for k, v := range raw {
		switch {
		case len(k) > 9 && k[:9] == "attempts_":
			var idx int
			fmt.Sscanf(k, "attempts_%d", &idx)
			e.Attempts[idx] = int(v.(float64))
		case len(k) > 7 && k[:7] == "claims_":
			var idx int
			fmt.Sscanf(k, "claims_%d", &idx)
			e.Claims[idx] = int(v.(float64))
		case len(k) > 11 && k[:11] == "best_score_":
			var idx int
			fmt.Sscanf(k, "best_score_%d", &idx)
			e.BestScores[idx] = int(v.(float64))
		}
	}
	return nil
}

type ExperimentationData struct {
	Simon                 *ExperimentationGame `json:"simon,omitempty"`
	Pairings              *ExperimentationGame `json:"pairings,omitempty"`
	Numbers               *ExperimentationGame `json:"numbers,omitempty"`
	ClaimsResets          *int64               `json:"claims_resets,omitempty"`
	ClaimsResetsTimestamp int64                `json:"claims_resets_timestamp,omitempty"`
	SerumsDrank           int                  `json:"serums_drank,omitempty"`
	ClaimedRetroactiveRng bool                 `json:"claimed_retroactive_rng,omitempty"`
	ChargeTrackTimestamp  int64                `json:"charge_track_timestamp,omitempty"`
}

type Dungeons struct {
	DungeonTypes         map[string]DungeonData `json:"dungeon_types,omitempty"`
	Classes              map[string]PlayerClass `json:"player_classes,omitempty"`
	SelectedDungeonClass string                 `json:"selected_dungeon_class,omitempty"`
	Secrets              float64                `json:"secrets,omitempty"`
}

type PlayerClass struct {
	Experience float64 `json:"experience,omitempty"`
}

type DungeonData struct {
	Experience float64 `json:"experience,omitempty"`

	HighestTierCompleted int                   `json:"highest_tier_completed,omitempty"`
	TimesPlayed          map[string]float64    `json:"times_played,omitempty"`
	TierCompletions      map[string]float64    `json:"tier_completions,omitempty"`
	MilestoneCompletions map[string]float64    `json:"milestone_completions,omitempty"`
	MobsKilled           map[string]float64    `json:"mobs_killed,omitempty"`
	MostMobsKilled       map[string]float64    `json:"most_mobs_killed,omitempty"`
	WatcherKills         map[string]float64    `json:"watcher_kills,omitempty"`
	MostDamageBerserk    map[string]float64    `json:"most_damage_berserk,omitempty"`
	MostDamageMage       map[string]float64    `json:"most_damage_mage,omitempty"`
	MostDamageHealer     map[string]float64    `json:"most_damage_healer,omitempty"`
	MostDamageArcher     map[string]float64    `json:"most_damage_archer,omitempty"`
	MostDamageTank       map[string]float64    `json:"most_damage_tank,omitempty"`
	MostHealing          map[string]float64    `json:"most_healing,omitempty"`
	FastestTime          map[string]float64    `json:"fastest_time,omitempty"`
	FastestTimeS         map[string]float64    `json:"fastest_time_s,omitempty"`
	FastestTimeSPlus     map[string]float64    `json:"fastest_time_s_plus,omitempty"`
	BestScore            map[string]float64    `json:"best_score,omitempty"`
	BestRuns             map[string]*[]BestRun `json:"best_runs,omitempty"`
}

type BestRun struct {
	Timestamp        int64   `json:"timestamp"`
	ScoreExploration int     `json:"score_exploration"`
	ScoreSpeed       int     `json:"score_speed"`
	ScoreSkill       int     `json:"score_skill"`
	ScoreBonus       int     `json:"score_bonus"`
	DungeonClass     string  `json:"dungeon_class"`
	ElapsedTime      int64   `json:"elapsed_time"`
	DamageDealt      float64 `json:"damage_dealt"`
	Deaths           int     `json:"deaths"`
	MobsKilled       int     `json:"mobs_killed"`
	SecretsFound     int     `json:"secrets_found"`
	DamageMitigated  float64 `json:"damage_mitigated"`
}

type Slayer struct {
	SlayerBosses map[string]SlayerBoss `json:"slayer_bosses,omitempty"`
}

type SlayerBoss struct {
	BossKillsTier0    int     `json:"boss_kills_tier_0,omitempty"`
	BossKillsTier1    int     `json:"boss_kills_tier_1,omitempty"`
	BossKillsTier2    int     `json:"boss_kills_tier_2,omitempty"`
	BossKillsTier3    int     `json:"boss_kills_tier_3,omitempty"`
	BossKillsTier4    int     `json:"boss_kills_tier_4,omitempty"`
	BossAttemptsTier0 int     `json:"boss_attempts_tier_0,omitempty"`
	BossAttemptsTier1 int     `json:"boss_attempts_tier_1,omitempty"`
	BossAttemptsTier2 int     `json:"boss_attempts_tier_2,omitempty"`
	BossAttemptsTier3 int     `json:"boss_attempts_tier_3,omitempty"`
	BossAttemptsTier4 int     `json:"boss_attempts_tier_4,omitempty"`
	Experience        float64 `json:"xp,omitempty"`
}

type CommunityUpgrades struct {
	UpgradeStates []CommunityUpgradeState `json:"upgrade_states,omitempty"`
}

type CommunityUpgradeState struct {
	Upgrade string `json:"upgrade,omitempty"`
	Tier    int    `json:"tier,omitempty"`
}
type Bestiary struct {
	Kills map[string]int `json:"kills,omitempty"`
}

func (b *Bestiary) UnmarshalJSON(data []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	b.Kills = make(map[string]int)
	if killsRaw, ok := raw["kills"]; ok {
		if killsMap, ok := killsRaw.(map[string]any); ok {
			for k, v := range killsMap {
				switch val := v.(type) {
				case float64:
					b.Kills[k] = int(val)
				case int:
					b.Kills[k] = val
				case string:
					if i, err := strconv.Atoi(val); err == nil {
						b.Kills[k] = i
					}
				}
			}
		}
	}

	return nil
}

type SkillTree struct {
	Nodes             map[string]SkillTreeNodeData `json:"-"`
	SelectedAbility   map[string]string            `json:"selected_ability,omitempty"`
	TokensSpent       map[string]int               `json:"tokens_spent,omitempty"`
	Experience        map[string]float64           `json:"experience,omitempty"`
	LastReset         map[string]int64             `json:"last_reset,omitempty"`
	RefundAbilityFree bool                         `json:"refund_ability_free,omitempty"`
}

type SkillTreeNodeData struct {
	Levels  map[string]int
	Toggles map[string]bool
}

func (s *SkillTree) UnmarshalJSON(data []byte) error {
	type Alias SkillTree
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	s.Nodes = make(map[string]SkillTreeNodeData)
	if nodesRaw, ok := raw["nodes"]; ok {
		if nodesMap, ok := nodesRaw.(map[string]any); ok {
			for skill, v := range nodesMap {
				if skillMap, ok := v.(map[string]any); ok {
					nodeData := SkillTreeNodeData{
						Levels:  make(map[string]int),
						Toggles: make(map[string]bool),
					}
					for k, val := range skillMap {
						switch typedVal := val.(type) {
						case float64:
							nodeData.Levels[k] = int(typedVal)
						case bool:
							nodeData.Toggles[k] = typedVal
						}
					}
					s.Nodes[skill] = nodeData
				}
			}
		}
	}

	return nil
}

type Foraging struct {
	Starlyn    *Starlyn   `json:"starlyn,omitempty"`
	FishFamily []string   `json:"fish_family,omitempty"`
	Hina       *Hina      `json:"hina,omitempty"`
	TreeGifts  *TreeGifts `json:"tree_gifts,omitempty"`
	Songs      *Songs     `json:"songs,omitempty"`
	// HuntingToolkit holds the hunting scythe, nets, traps and other hunting tools.
	HuntingToolkit *Toolkit `json:"hunting_toolkit,omitempty"`
}

type Starlyn struct {
	PersonalBests map[string]float64 `json:"personal_bests,omitempty"`
}

type Hina struct {
	Tasks *HinaTasks `json:"tasks,omitempty"`
}

type HinaTasks struct {
	TaskProgress   map[string]int `json:"task_progress,omitempty"`
	CompletedTasks []string       `json:"completed_tasks,omitempty"`
	ClaimedRewards []string       `json:"claimed_rewards,omitempty"`
	TierClaimed    int            `json:"tier_claimed,omitempty"`
}

type TreeGifts struct {
	Fig                  int            `json:"-"`
	Mangrove             int            `json:"-"`
	MilestoneTierClaimed map[string]int `json:"milestone_tier_claimed,omitempty"`
	Extra                map[string]int `json:"-"`
}

func (t *TreeGifts) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	t.Extra = make(map[string]int)

	for key, value := range raw {
		switch key {
		case "milestone_tier_claimed":
			if err := json.Unmarshal(value, &t.MilestoneTierClaimed); err != nil {
				return err
			}
		default:
			var num int
			if err := json.Unmarshal(value, &num); err == nil {
				t.Extra[key] = num
				switch key {
				case "FIG":
					t.Fig = num
				case "MANGROVE":
					t.Mangrove = num
				}
			}
		}
	}

	return nil
}

type Songs struct {
	Harp *HarpSongs `json:"harp,omitempty"`
}

type HarpSongs struct {
	ClaimedTalisman   bool               `json:"claimed_talisman,omitempty"`
	SelectedSong      string             `json:"selected_song,omitempty"`
	SelectedSongEpoch int64              `json:"selected_song_epoch,omitempty"`
	Extra             map[string]float64 `json:"-"`
}

func (h *HarpSongs) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	h.Extra = make(map[string]float64)

	for key, value := range raw {
		switch key {
		case "claimed_talisman":
			if err := json.Unmarshal(value, &h.ClaimedTalisman); err != nil {
				return err
			}
		case "selected_song":
			if err := json.Unmarshal(value, &h.SelectedSong); err != nil {
				return err
			}
		case "selected_song_epoch":
			if err := json.Unmarshal(value, &h.SelectedSongEpoch); err != nil {
				return err
			}
		default:
			var num float64
			if err := json.Unmarshal(value, &num); err == nil {
				h.Extra[key] = num
			}
		}
	}

	return nil
}

type ForagingCore struct {
	DailyTreesCutDay              int      `json:"daily_trees_cut_day,omitempty"`
	DailyTreesCut                 int      `json:"daily_trees_cut,omitempty"`
	DailyGifts                    int      `json:"daily_gifts,omitempty"`
	DailyLogCutDay                int      `json:"daily_log_cut_day,omitempty"`
	DailyLogCut                   []string `json:"daily_log_cut,omitempty"`
	ForestsWhispers               int      `json:"forests_whispers,omitempty"`
	ForestsWhispersSpent          int      `json:"forests_whispers_spent,omitempty"`
	CurrentDailyEffect            string   `json:"current_daily_effect,omitempty"`
	CurrentDailyEffectLastChanged int64    `json:"current_daily_effect_last_changed,omitempty"`
}

type Shards struct {
	Traps            ShardTraps   `json:"traps"`
	Owned            []OwnedShard `json:"owned,omitempty"`
	ShardSort        string       `json:"shard_sort,omitempty"`
	Fused            int          `json:"fused,omitempty"`
	FusionResultSort string       `json:"fusion_result_sort,omitempty"`
}

type ShardTraps struct {
	ActiveTraps []ShardTrap `json:"active_traps,omitempty"`
}

type ShardTrap struct {
	TrapItem            string `json:"trap_item,omitempty"`
	CaptureTime         int64  `json:"capture_time,omitempty"`
	Mode                string `json:"mode,omitempty"`
	Location            string `json:"location,omitempty"`
	PlacedAt            int64  `json:"placed_at,omitempty"`
	Shard               string `json:"shard,omitempty"`
	Captured            bool   `json:"captured,omitempty"`
	UUID                string `json:"uuid,omitempty"`
	HuntingToolkit      bool   `json:"hunting_toolkit,omitempty"`
	HuntingToolkitIndex int    `json:"hunting_toolkit_index,omitempty"`
}

type OwnedShard struct {
	Type        string `json:"type,omitempty"`
	AmountOwned int    `json:"amount_owned,omitempty"`
	Captured    int64  `json:"captured,omitempty"`
}

type Temples struct {
	UnlockedTemples []string `json:"unlocked_temples,omitempty"`
}

type Attributes struct {
	Stacks map[string]int `json:"stacks,omitempty"`
}
