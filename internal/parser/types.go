package parser

import "time"

// CombatEvent представляет событие боя из лога
type CombatEvent struct {
	Timestamp string `json:"timestamp"`
	Frame     int    `json:"frame"`
	Category  string `json:"category"`
	Message   string `json:"message"`
}

// DamageEvent представляет событие урона
type DamageEvent struct {
	Timestamp time.Time
	Amount    int
	IsCrit    bool
	IsLethal  bool
	Target    string
	Source    string
	Ability   string
	IsDealt   bool // true если урон нанесен, false если получен
}

// HealEvent представляет событие исцеления
type HealEvent struct {
	Timestamp time.Time
	Amount    int
	IsCrit    bool
	Target    string
	Source    string
	Ability   string
	IsDealt   bool // true если исцеление нанесено, false если получено
}

// KillEvent представляет событие убийства
type KillEvent struct {
	Timestamp time.Time
	Target    string
	Source    string
	Ability   string
	Damage    int
	IsCrit    bool
}

// BuffEvent представляет событие баффа/дебаффа
type BuffEvent struct {
	Timestamp time.Time
	Type      string // "Received", "Applied", "Removed"
	BuffName  string
	Target    string
	Source    string
}

// CombatStateEvent представляет событие изменения состояния боя
type CombatStateEvent struct {
	Timestamp time.Time
	State     string // "Entered", "Exited", "Started", "Ended"
	Target    string
	Source    string
}
