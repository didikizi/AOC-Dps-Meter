package metrics

import "time"

// CombatStats представляет статистику боя
type CombatStats struct {
	StartTime        time.Time
	EndTime          time.Time
	Duration         time.Duration
	TotalDamage      int
	TotalHealing     int
	TotalKills       int
	CritHits         int
	TotalHits        int
	CritHealing      int
	TotalHealingHits int
}

// DPSStats представляет статистику DPS
type DPSStats struct {
	CurrentDPS  float64
	MaxDPS      float64
	AvgDPS      float64
	TotalDamage int
	Duration    time.Duration
}

// HPSStats представляет статистику HPS
type HPSStats struct {
	CurrentHPS   float64
	MaxHPS       float64
	AvgHPS       float64
	TotalHealing int
	Duration     time.Duration
}

// AbilityStats представляет статистику по способностям
type AbilityStats struct {
	Name        string
	Damage      int
	Healing     int
	Hits        int
	Crits       int
	Kills       int
	CritHealing int
	HealingHits int
	LastUsed    time.Time
}

// TargetStats представляет статистику по целям
type TargetStats struct {
	Name        string
	Damage      int
	Healing     int
	Hits        int
	Crits       int
	Kills       int
	CritHealing int
	HealingHits int
	LastHit     time.Time
}

// CombatEvent представляет событие боя с временной меткой
type CombatEvent struct {
	Timestamp time.Time
	Event     interface{}
}

// Combat представляет отдельный бой
type Combat struct {
	ID           string
	StartTime    time.Time
	EndTime      time.Time
	IsActive     bool
	TotalDamage  int
	TotalHealing int
	Duration     time.Duration
}

// CombatSession представляет сессию боя
type CombatSession struct {
	ID            string
	StartTime     time.Time
	EndTime       time.Time
	IsActive      bool
	Stats         CombatStats
	DPSStats      DPSStats
	HPSStats      HPSStats
	Abilities     map[string]*AbilityStats
	Targets       map[string]*TargetStats
	RecentEvents  []CombatEvent
	CurrentCombat *Combat
	LastActivity  time.Time
}
