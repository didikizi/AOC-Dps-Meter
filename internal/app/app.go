package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"aocdpsmetr/internal/metrics"
	"aocdpsmetr/internal/watcher"
)

// App struct
type App struct {
	ctx        context.Context
	calculator *metrics.Calculator
	watcher    *watcher.Watcher
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		calculator: metrics.NewCalculator(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	fmt.Println("App startup completed")
}

// DomReady is called after the front-end dom has been loaded
func (a *App) DomReady(ctx context.Context) {
	fmt.Println("DOM ready")
}

// BeforeClose is called when the app is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the app to continue,
// false will continue shutdown as normal.
func (a *App) BeforeClose(ctx context.Context) (prevent bool) {
	return false
}

// Shutdown is called at application shutdown
func (a *App) Shutdown(ctx context.Context) {
	if a.watcher != nil {
		a.watcher.Stop()
	}
	fmt.Println("App shutdown")
}

// findLogFile ищет файл логов в стандартных местах
func (a *App) findLogFile() string {
	// Возможные пути для файла логов AOC
	possiblePaths := []string{
		// Стандартный путь для Windows
		filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local", "AOC", "Saved", "Logs", "AOC.log"),
		// Альтернативные пути
		filepath.Join(os.Getenv("USERPROFILE"), "Documents", "AOC", "Logs", "AOC.log"),
		filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming", "AOC", "Logs", "AOC.log"),
		// Пользовательский путь (замените на свой путь к логам)
		// "C:\\Your\\Custom\\Path\\To\\AOC\\Logs\\AOC.log",
		// Локальный файл для тестирования
		"AOC.log",
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			fmt.Printf("Found log file: %s\n", path)
			return path
		}
	}

	fmt.Println("Log file not found in any standard location")
	return ""
}

func (a *App) StartMonitoring() string {
	fmt.Println("StartMonitoring called")
	if a.watcher != nil {
		fmt.Println("Already monitoring")
		return "Already monitoring"
	}

	// Ищем файл логов в стандартных местах
	logPath := a.findLogFile()
	if logPath == "" {
		return "Log file not found in standard locations"
	}
	fmt.Println("Creating watcher for:", logPath)

	// Проверяем существование файла
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		fmt.Println("Log file does not exist:", logPath)
		return "Log file not found: " + logPath
	}

	a.watcher = watcher.NewWatcher(logPath, func(events []interface{}) {
		fmt.Printf("Processing %d events\n", len(events))
		for _, event := range events {
			a.calculator.ProcessEvent(event)
		}
	})

	if err := a.watcher.Start(); err != nil {
		fmt.Println("Failed to start monitoring:", err)
		return "Failed to start monitoring: " + err.Error()
	}

	fmt.Println("Monitoring started successfully")
	return "Monitoring started"
}

func (a *App) StopMonitoring() string {
	fmt.Println("StopMonitoring called")
	if a.watcher == nil {
		fmt.Println("Watcher is nil, not monitoring")
		return "Not monitoring"
	}

	fmt.Println("Stopping watcher...")
	a.watcher.Stop()
	a.watcher = nil
	fmt.Println("Monitoring stopped successfully")
	return "Monitoring stopped"
}

func (a *App) ResetStats() string {
	a.calculator.ResetSession()
	return "Statistics reset"
}

// OpenDevTools opens the developer tools
func (a *App) OpenDevTools() string {
	fmt.Println("OpenDevTools called")
	return "DevTools opening attempted"
}

// GetLogPath returns the current log file path
func (a *App) GetLogPath() string {
	return a.findLogFile()
}

func (a *App) GetStats() map[string]interface{} {
	session := a.calculator.GetSession()

	critRate := 0.0
	if session.Stats.TotalHits > 0 {
		critRate = float64(session.Stats.CritHits) / float64(session.Stats.TotalHits) * 100
	}

	healingCritRate := 0.0
	if session.Stats.TotalHealingHits > 0 {
		healingCritRate = float64(session.Stats.CritHealing) / float64(session.Stats.TotalHealingHits) * 100
	}

	stats := map[string]interface{}{
		"maxDps":          session.DPSStats.MaxDPS,
		"dps":             session.DPSStats.CurrentDPS,
		"damage":          session.Stats.TotalDamage,
		"hits":            session.Stats.TotalHits,
		"crits":           session.Stats.CritHits,
		"maxHps":          session.HPSStats.MaxHPS,
		"hps":             session.HPSStats.CurrentHPS,
		"healing":         session.Stats.TotalHealing,
		"healingHits":     session.Stats.TotalHealingHits,
		"healingCrits":    session.Stats.CritHealing,
		"critRate":        critRate,
		"healingCritRate": healingCritRate,
		"kills":           session.Stats.TotalKills,
		"duration":        time.Since(session.StartTime).Seconds(),
		"isActive":        session.IsActive,
	}

	fmt.Printf("Returning stats: %+v\n", stats)
	return stats
}

func (a *App) GetAbilities() []map[string]interface{} {
	session := a.calculator.GetSession()
	abilities := make([]*metrics.AbilityStats, 0, len(session.Abilities))

	for _, ability := range session.Abilities {
		if ability.Damage > 0 {
			abilities = append(abilities, ability)
		}
	}

	// Сортируем по урону
	for i := 0; i < len(abilities); i++ {
		for j := i + 1; j < len(abilities); j++ {
			if abilities[i].Damage < abilities[j].Damage {
				abilities[i], abilities[j] = abilities[j], abilities[i]
			}
		}
	}

	result := make([]map[string]interface{}, 0, len(abilities))
	for _, ability := range abilities {
		critRate := 0.0
		if ability.Hits > 0 {
			critRate = float64(ability.Crits) / float64(ability.Hits) * 100
		}

		healingCritRate := 0.0
		if ability.HealingHits > 0 {
			healingCritRate = float64(ability.CritHealing) / float64(ability.HealingHits) * 100
		}

		result = append(result, map[string]interface{}{
			"name":            ability.Name,
			"damage":          ability.Damage,
			"healing":         ability.Healing,
			"hits":            ability.Hits,
			"crits":           ability.Crits,
			"critRate":        critRate,
			"healingHits":     ability.HealingHits,
			"healingCrits":    ability.CritHealing,
			"healingCritRate": healingCritRate,
		})
	}

	return result
}

func (a *App) GetTargets() []map[string]interface{} {
	session := a.calculator.GetSession()
	targets := make([]*metrics.TargetStats, 0, len(session.Targets))

	for _, target := range session.Targets {
		if target.Damage > 0 {
			targets = append(targets, target)
		}
	}

	// Сортируем по урону
	for i := 0; i < len(targets); i++ {
		for j := i + 1; j < len(targets); j++ {
			if targets[i].Damage < targets[j].Damage {
				targets[i], targets[j] = targets[j], targets[i]
			}
		}
	}

	result := make([]map[string]interface{}, 0, len(targets))
	for _, target := range targets {
		critRate := 0.0
		if target.Hits > 0 {
			critRate = float64(target.Crits) / float64(target.Hits) * 100
		}

		healingCritRate := 0.0
		if target.HealingHits > 0 {
			healingCritRate = float64(target.CritHealing) / float64(target.HealingHits) * 100
		}

		result = append(result, map[string]interface{}{
			"name":            target.Name,
			"damage":          target.Damage,
			"healing":         target.Healing,
			"hits":            target.Hits,
			"crits":           target.Crits,
			"critRate":        critRate,
			"kills":           target.Kills,
			"healingHits":     target.HealingHits,
			"healingCrits":    target.CritHealing,
			"healingCritRate": healingCritRate,
		})
	}

	return result
}
