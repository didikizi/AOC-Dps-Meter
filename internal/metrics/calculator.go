package metrics

import (
	"fmt"
	"time"

	"aocdpsmetr/internal/parser"
)

// Calculator рассчитывает метрики боя
type Calculator struct {
	session *CombatSession
}

// NewCalculator создает новый калькулятор
func NewCalculator() *Calculator {
	return &Calculator{
		session: &CombatSession{
			ID:        generateSessionID(),
			StartTime: time.Now(),
			IsActive:  true,
			Abilities: make(map[string]*AbilityStats),
			Targets:   make(map[string]*TargetStats),
		},
	}
}

// ProcessEvent обрабатывает событие боя
func (c *Calculator) ProcessEvent(event interface{}) {
	now := time.Now()

	// Проверяем, нужно ли начать новый бой
	c.checkCombatStatus(now)

	switch e := event.(type) {
	case *parser.DamageEvent:
		fmt.Printf("Processing DamageEvent: %+v\n", e)
		c.processDamageEvent(e)
	case *parser.HealEvent:
		fmt.Printf("Processing HealEvent: %+v\n", e)
		c.processHealEvent(e)
	case *parser.KillEvent:
		fmt.Printf("Processing KillEvent: %+v\n", e)
		c.processKillEvent(e)
	case *parser.BuffEvent:
		fmt.Printf("Processing BuffEvent: %+v\n", e)
		c.processBuffEvent(e)
	default:
		fmt.Printf("Unknown event type: %T\n", event)
	}

	// Обновляем время последней активности
	c.session.LastActivity = now

	// Добавляем событие в список недавних событий
	c.addRecentEvent(event)
}

// processDamageEvent обрабатывает событие урона
func (c *Calculator) processDamageEvent(event *parser.DamageEvent) {
	if !c.session.IsActive {
		c.startNewSession()
	}

	// Обновляем общую статистику
	c.session.Stats.TotalDamage += event.Amount
	c.session.Stats.TotalHits++
	if event.IsCrit {
		c.session.Stats.CritHits++
	}

	// Обновляем статистику по способностям
	if ability, exists := c.session.Abilities[event.Ability]; exists {
		ability.Damage += event.Amount
		ability.Hits++
		if event.IsCrit {
			ability.Crits++
		}
		ability.LastUsed = event.Timestamp
	} else {
		c.session.Abilities[event.Ability] = &AbilityStats{
			Name:     event.Ability,
			Damage:   event.Amount,
			Hits:     1,
			Crits:    boolToInt(event.IsCrit),
			LastUsed: event.Timestamp,
		}
	}

	// Обновляем статистику по целям
	if target, exists := c.session.Targets[event.Target]; exists {
		target.Damage += event.Amount
		target.Hits++
		if event.IsCrit {
			target.Crits++
		}
		target.LastHit = event.Timestamp
	} else {
		c.session.Targets[event.Target] = &TargetStats{
			Name:    event.Target,
			Damage:  event.Amount,
			Hits:    1,
			Crits:   boolToInt(event.IsCrit),
			LastHit: event.Timestamp,
		}
	}

	// Пересчитываем DPS
	c.updateDPSStats()
}

// processHealEvent обрабатывает событие исцеления
func (c *Calculator) processHealEvent(event *parser.HealEvent) {
	if !c.session.IsActive {
		c.startNewSession()
	}

	// Обновляем общую статистику
	c.session.Stats.TotalHealing += event.Amount
	c.session.Stats.TotalHealingHits++
	if event.IsCrit {
		c.session.Stats.CritHealing++
	}

	// Обновляем статистику по способностям
	if ability, exists := c.session.Abilities[event.Ability]; exists {
		ability.Healing += event.Amount
		ability.HealingHits++
		if event.IsCrit {
			ability.CritHealing++
		}
		ability.LastUsed = event.Timestamp
	} else {
		c.session.Abilities[event.Ability] = &AbilityStats{
			Name:        event.Ability,
			Healing:     event.Amount,
			HealingHits: 1,
			CritHealing: boolToInt(event.IsCrit),
			LastUsed:    event.Timestamp,
		}
	}

	// Обновляем статистику по целям
	if target, exists := c.session.Targets[event.Target]; exists {
		target.Healing += event.Amount
		target.HealingHits++
		if event.IsCrit {
			target.CritHealing++
		}
		target.LastHit = event.Timestamp
	} else {
		c.session.Targets[event.Target] = &TargetStats{
			Name:        event.Target,
			Healing:     event.Amount,
			HealingHits: 1,
			CritHealing: boolToInt(event.IsCrit),
			LastHit:     event.Timestamp,
		}
	}

	// Пересчитываем HPS
	c.updateHPSStats()
}

// processKillEvent обрабатывает событие убийства
func (c *Calculator) processKillEvent(event *parser.KillEvent) {
	if !c.session.IsActive {
		c.startNewSession()
	}

	// Обновляем общую статистику
	c.session.Stats.TotalKills++

	// Обновляем статистику по способностям
	if ability, exists := c.session.Abilities[event.Ability]; exists {
		ability.Kills++
		ability.LastUsed = event.Timestamp
	} else {
		c.session.Abilities[event.Ability] = &AbilityStats{
			Name:     event.Ability,
			Kills:    1,
			LastUsed: event.Timestamp,
		}
	}

	// Обновляем статистику по целям
	if target, exists := c.session.Targets[event.Target]; exists {
		target.Kills++
		target.LastHit = event.Timestamp
	} else {
		c.session.Targets[event.Target] = &TargetStats{
			Name:    event.Target,
			Kills:   1,
			LastHit: event.Timestamp,
		}
	}
}

// processBuffEvent обрабатывает событие баффа/дебаффа
func (c *Calculator) processBuffEvent(event *parser.BuffEvent) {
	// Пока что просто логируем баффы, можно добавить статистику по баффам
	// В будущем можно добавить отслеживание времени действия баффов, их эффективности и т.д.
}

// updateDPSStats пересчитывает статистику DPS за текущий бой
func (c *Calculator) updateDPSStats() {
	if c.session.CurrentCombat == nil || !c.session.CurrentCombat.IsActive {
		// Если нет активного боя, DPS = 0
		c.session.DPSStats.CurrentDPS = 0
		return
	}

	// Рассчитываем урон за текущий бой
	damageInCombat := 0
	now := time.Now()
	combatStartTime := c.session.CurrentCombat.StartTime

	for _, event := range c.session.RecentEvents {
		// События должны быть в рамках текущего боя
		if event.Timestamp.After(combatStartTime) {
			if damageEvent, ok := event.Event.(*parser.DamageEvent); ok {
				if damageEvent.IsDealt {
					damageInCombat += damageEvent.Amount
				}
			}
		}
	}

	// DPS = урон за бой / длительность боя
	combatDuration := now.Sub(combatStartTime).Seconds()
	if combatDuration > 0 {
		c.session.DPSStats.CurrentDPS = float64(damageInCombat) / combatDuration
	} else {
		c.session.DPSStats.CurrentDPS = 0
	}

	// Обновляем максимум
	if c.session.DPSStats.CurrentDPS > c.session.DPSStats.MaxDPS {
		c.session.DPSStats.MaxDPS = c.session.DPSStats.CurrentDPS
	}

	c.session.DPSStats.TotalDamage = c.session.Stats.TotalDamage
	c.session.DPSStats.Duration = time.Since(c.session.StartTime)
}

// updateHPSStats пересчитывает статистику HPS за текущий бой
func (c *Calculator) updateHPSStats() {
	if c.session.CurrentCombat == nil || !c.session.CurrentCombat.IsActive {
		// Если нет активного боя, HPS = 0
		c.session.HPSStats.CurrentHPS = 0
		return
	}

	// Рассчитываем исцеление за текущий бой
	healingInCombat := 0
	now := time.Now()

	for _, event := range c.session.RecentEvents {
		// События должны быть в рамках текущего боя
		if event.Timestamp.After(c.session.CurrentCombat.StartTime) &&
			(c.session.CurrentCombat.IsActive || event.Timestamp.Before(c.session.CurrentCombat.EndTime)) {
			if healEvent, ok := event.Event.(*parser.HealEvent); ok && !healEvent.IsDealt {
				healingInCombat += healEvent.Amount
			}
		}
	}

	// HPS = исцеление за бой / длительность боя
	combatDuration := now.Sub(c.session.CurrentCombat.StartTime).Seconds()
	if combatDuration > 0 {
		c.session.HPSStats.CurrentHPS = float64(healingInCombat) / combatDuration
	} else {
		c.session.HPSStats.CurrentHPS = 0
	}

	// Обновляем максимум
	if c.session.HPSStats.CurrentHPS > c.session.HPSStats.MaxHPS {
		c.session.HPSStats.MaxHPS = c.session.HPSStats.CurrentHPS
	}

	c.session.HPSStats.TotalHealing = c.session.Stats.TotalHealing
	c.session.HPSStats.Duration = time.Since(c.session.StartTime)
}

// checkCombatStatus проверяет статус боя и при необходимости начинает новый
func (c *Calculator) checkCombatStatus(now time.Time) {
	// Если нет активного боя, начинаем новый
	if c.session.CurrentCombat == nil || !c.session.CurrentCombat.IsActive {
		c.startNewCombat(now)
		return
	}

	// Проверяем, прошло ли 10 секунд без активности
	if now.Sub(c.session.LastActivity) >= 10*time.Second {
		// Завершаем текущий бой
		c.endCurrentCombat(now)
		// Не начинаем новый бой автоматически - ждем следующего события
	}
}

// startNewCombat начинает новый бой
func (c *Calculator) startNewCombat(now time.Time) {
	c.session.CurrentCombat = &Combat{
		ID:        generateSessionID(),
		StartTime: now,
		IsActive:  true,
	}
	fmt.Printf("Started new combat: %s\n", c.session.CurrentCombat.ID)
}

// endCurrentCombat завершает текущий бой
func (c *Calculator) endCurrentCombat(now time.Time) {
	if c.session.CurrentCombat != nil && c.session.CurrentCombat.IsActive {
		c.session.CurrentCombat.EndTime = now
		c.session.CurrentCombat.IsActive = false
		c.session.CurrentCombat.Duration = now.Sub(c.session.CurrentCombat.StartTime)
		fmt.Printf("Ended combat: %s, Duration: %v\n", c.session.CurrentCombat.ID, c.session.CurrentCombat.Duration)
	}
}

// addRecentEvent добавляет событие в список недавних событий
func (c *Calculator) addRecentEvent(event interface{}) {
	now := time.Now()
	c.session.RecentEvents = append(c.session.RecentEvents, CombatEvent{
		Timestamp: now,
		Event:     event,
	})

	// Удаляем события старше 2 минут (оставляем буфер)
	twoMinutesAgo := now.Add(-2 * time.Minute)
	var recentEvents []CombatEvent
	for _, evt := range c.session.RecentEvents {
		if evt.Timestamp.After(twoMinutesAgo) {
			recentEvents = append(recentEvents, evt)
		}
	}
	c.session.RecentEvents = recentEvents
}

// startNewSession начинает новую сессию боя
func (c *Calculator) startNewSession() {
	c.session = &CombatSession{
		ID:        generateSessionID(),
		StartTime: time.Now(),
		IsActive:  true,
		Stats: CombatStats{
			TotalDamage:      0,
			TotalHealing:     0,
			TotalHits:        0,
			TotalHealingHits: 0,
			CritHits:         0,
			CritHealing:      0,
			TotalKills:       0,
		},
		DPSStats: DPSStats{
			CurrentDPS:  0,
			MaxDPS:      0,
			TotalDamage: 0,
			Duration:    0,
		},
		HPSStats: HPSStats{
			CurrentHPS:   0,
			MaxHPS:       0,
			TotalHealing: 0,
			Duration:     0,
		},
		Abilities:    make(map[string]*AbilityStats),
		Targets:      make(map[string]*TargetStats),
		RecentEvents: make([]CombatEvent, 0),
		LastActivity: time.Now(),
	}
}

// GetSession возвращает текущую сессию
func (c *Calculator) GetSession() *CombatSession {
	return c.session
}

// ResetSession сбрасывает текущую сессию
func (c *Calculator) ResetSession() {
	c.startNewSession()
}

// EndSession завершает текущую сессию
func (c *Calculator) EndSession() {
	c.session.EndTime = time.Now()
	c.session.IsActive = false
	c.session.Stats.Duration = c.session.EndTime.Sub(c.session.StartTime)
}

// Вспомогательные функции
func generateSessionID() string {
	return time.Now().Format("20060102150405")
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
