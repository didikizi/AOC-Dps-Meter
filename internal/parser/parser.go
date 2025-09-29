package parser

import (
	"bufio"
	"encoding/json"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Parser парсит логи Ashes of Creation
type Parser struct {
	// Регулярные выражения для парсинга событий
	damageDealtRegex    *regexp.Regexp
	damageReceivedRegex *regexp.Regexp
	healReceivedRegex   *regexp.Regexp
	killRegex           *regexp.Regexp
	buffReceivedRegex   *regexp.Regexp
	buffAppliedRegex    *regexp.Regexp
	buffRemovedRegex    *regexp.Regexp
}

// NewParser создает новый парсер
func NewParser() *Parser {
	return &Parser{
		// Урон нанесенный: "83 damage(Crit) dealt to Wilderherd Berserker - Weapon_Wand_Projectile_1"
		damageDealtRegex: regexp.MustCompile(`(\d+(?:,\d+)*) damage(\(Crit\))?(\(Lethal\))? dealt to (.+) - (.+)`),
		// Урон полученный: "80 damage(Crit) received from Wilderherd Berserker - Axe Strike"
		damageReceivedRegex: regexp.MustCompile(`(\d+(?:,\d+)*) damage(\(Crit\))?(\(Lethal\))? received from (.+) - (.+)`),
		// Исцеление полученное: "103 healing(Crit) received from Your - Cleric_SoothingGlow"
		healReceivedRegex: regexp.MustCompile(`(\d+(?:,\d+)*) healing(\(Crit\))? received from (.+) - (.+)`),
		// Убийство: "95 damage(Crit)(Lethal) dealt to Wilderherd Berserker - Weapon_Wand_Projectile_1 [&Kill][KILL]Killed Wilderherd Berserker"
		killRegex: regexp.MustCompile(`(\d+(?:,\d+)*) damage(\(Crit\))?(\(Lethal\))? dealt to (.+) - (.+) \[&Kill\]\[KILL\]Killed (.+)`),
		// Бафф получен: "Received  [Divine Power]"
		buffReceivedRegex: regexp.MustCompile(`Received\s+\[(.+)\]`),
		// Бафф применен: "Applied [Volatile] to [Wilderherd Berserker]"
		buffAppliedRegex: regexp.MustCompile(`Applied \[(.+)\] to \[(.+)\]`),
		// Бафф снят: "Removed [Volatile] from [Wilderherd Berserker]"
		buffRemovedRegex: regexp.MustCompile(`Removed \[(.+)\] from \[(.+)\]`),
	}
}

// ParseLine парсит одну строку лога
func (p *Parser) ParseLine(line string) (interface{}, error) {
	// Парсим JSON
	var event CombatEvent
	if err := json.Unmarshal([]byte(line), &event); err != nil {
		return nil, err
	}

	// Проверяем, что это событие боя
	if event.Category != "LogAoC_CombatLog" {
		return nil, nil
	}

	// Парсим время
	timestamp, err := time.Parse("2006-01-02T15:04:05.000Z", event.Timestamp)
	if err != nil {
		return nil, err
	}

	// Парсим сообщение в зависимости от типа события
	if damageEvent := p.parseDamageEvent(event, timestamp); damageEvent != nil {
		return damageEvent, nil
	}

	if healEvent := p.parseHealEvent(event, timestamp); healEvent != nil {
		return healEvent, nil
	}

	if killEvent := p.parseKillEvent(event, timestamp); killEvent != nil {
		return killEvent, nil
	}

	if buffEvent := p.parseBuffEvent(event, timestamp); buffEvent != nil {
		return buffEvent, nil
	}

	return nil, nil
}

// parseDamageEvent парсит событие урона
func (p *Parser) parseDamageEvent(event CombatEvent, timestamp time.Time) *DamageEvent {
	// Проверяем урон нанесенный
	if matches := p.damageDealtRegex.FindStringSubmatch(event.Message); matches != nil {
		amount, _ := strconv.Atoi(strings.ReplaceAll(matches[1], ",", ""))
		return &DamageEvent{
			Timestamp: timestamp,
			Amount:    amount,
			IsCrit:    matches[2] == "(Crit)",
			IsLethal:  matches[3] == "(Lethal)",
			Target:    matches[4],
			Source:    "You", // Предполагаем, что игрок наносит урон
			Ability:   matches[5],
			IsDealt:   true,
		}
	}

	// Проверяем урон полученный
	if matches := p.damageReceivedRegex.FindStringSubmatch(event.Message); matches != nil {
		amount, _ := strconv.Atoi(strings.ReplaceAll(matches[1], ",", ""))
		return &DamageEvent{
			Timestamp: timestamp,
			Amount:    amount,
			IsCrit:    matches[2] == "(Crit)",
			IsLethal:  matches[3] == "(Lethal)",
			Target:    "You", // Предполагаем, что игрок получает урон
			Source:    matches[4],
			Ability:   matches[5],
			IsDealt:   false,
		}
	}

	return nil
}

// parseHealEvent парсит событие исцеления
func (p *Parser) parseHealEvent(event CombatEvent, timestamp time.Time) *HealEvent {
	// Проверяем исцеление полученное
	if matches := p.healReceivedRegex.FindStringSubmatch(event.Message); matches != nil {
		amount, _ := strconv.Atoi(strings.ReplaceAll(matches[1], ",", ""))
		return &HealEvent{
			Timestamp: timestamp,
			Amount:    amount,
			IsCrit:    matches[2] == "(Crit)",
			Target:    "You", // Предполагаем, что игрок получает исцеление
			Source:    matches[3],
			Ability:   matches[4],
			IsDealt:   false,
		}
	}

	return nil
}

// parseKillEvent парсит событие убийства
func (p *Parser) parseKillEvent(event CombatEvent, timestamp time.Time) *KillEvent {
	if matches := p.killRegex.FindStringSubmatch(event.Message); matches != nil {
		amount, _ := strconv.Atoi(strings.ReplaceAll(matches[1], ",", ""))
		return &KillEvent{
			Timestamp: timestamp,
			Target:    matches[4],
			Source:    "You", // Предполагаем, что игрок убивает
			Ability:   matches[5],
			Damage:    amount,
			IsCrit:    matches[2] == "(Crit)",
		}
	}

	return nil
}

// parseBuffEvent парсит событие баффа/дебаффа
func (p *Parser) parseBuffEvent(event CombatEvent, timestamp time.Time) *BuffEvent {
	// Проверяем получение баффа
	if matches := p.buffReceivedRegex.FindStringSubmatch(event.Message); matches != nil {
		return &BuffEvent{
			Timestamp: timestamp,
			Type:      "Received",
			BuffName:  matches[1],
			Target:    "You", // Предполагаем, что игрок получает бафф
			Source:    "Unknown",
		}
	}

	// Проверяем применение баффа
	if matches := p.buffAppliedRegex.FindStringSubmatch(event.Message); matches != nil {
		return &BuffEvent{
			Timestamp: timestamp,
			Type:      "Applied",
			BuffName:  matches[1],
			Target:    matches[2],
			Source:    "You", // Предполагаем, что игрок применяет бафф
		}
	}

	// Проверяем снятие баффа
	if matches := p.buffRemovedRegex.FindStringSubmatch(event.Message); matches != nil {
		return &BuffEvent{
			Timestamp: timestamp,
			Type:      "Removed",
			BuffName:  matches[1],
			Target:    matches[2],
			Source:    "Unknown",
		}
	}

	return nil
}

// ParseFile парсит весь файл лога
func (p *Parser) ParseFile(filename string) ([]interface{}, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var events []interface{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if event, err := p.ParseLine(line); err == nil && event != nil {
			events = append(events, event)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

// ParseFileFromLine парсит файл начиная с определенной строки
func (p *Parser) ParseFileFromLine(filename string, startLine int) ([]interface{}, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var events []interface{}
	scanner := bufio.NewScanner(file)
	currentLine := 0

	for scanner.Scan() {
		currentLine++
		if currentLine <= startLine {
			continue
		}

		line := scanner.Text()
		if event, err := p.ParseLine(line); err == nil && event != nil {
			events = append(events, event)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
