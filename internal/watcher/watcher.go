package watcher

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"aocdpsmetr/internal/parser"

	"github.com/fsnotify/fsnotify"
)

// Watcher отслеживает изменения в файле лога
type Watcher struct {
	filename string
	parser   *parser.Parser
	callback func([]interface{})
	watcher  *fsnotify.Watcher
	ctx      context.Context
	cancel   context.CancelFunc
	lastLine int
}

// NewWatcher создает новый watcher
func NewWatcher(filename string, callback func([]interface{})) *Watcher {
	ctx, cancel := context.WithCancel(context.Background())

	return &Watcher{
		filename: filename,
		parser:   parser.NewParser(),
		callback: callback,
		ctx:      ctx,
		cancel:   cancel,
		lastLine: 0,
	}
}

// Start начинает мониторинг файла
func (w *Watcher) Start() error {
	// Создаем watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}
	w.watcher = watcher

	// Добавляем файл для мониторинга
	if err := w.watcher.Add(w.filename); err != nil {
		return fmt.Errorf("failed to add file to watcher: %w", err)
	}

	// Читаем существующий файл для получения начальной позиции
	if err := w.readExistingFile(); err != nil {
		return fmt.Errorf("failed to read existing file: %w", err)
	}

	// Обрабатываем весь существующий файл при старте
	if err := w.processExistingFile(); err != nil {
		fmt.Printf("Warning: failed to process existing file: %v\n", err)
	}

	// Запускаем цикл мониторинга
	go w.watchLoop()

	return nil
}

// Stop останавливает мониторинг
func (w *Watcher) Stop() {
	if w.cancel != nil {
		w.cancel()
	}
	if w.watcher != nil {
		w.watcher.Close()
	}
}

// readExistingFile читает существующий файл для получения начальной позиции
func (w *Watcher) readExistingFile() error {
	file, err := os.Open(w.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0

	for scanner.Scan() {
		lineCount++
	}

	w.lastLine = lineCount
	fmt.Println("Last line:", w.lastLine)
	return scanner.Err()
}

// processExistingFile обрабатывает весь существующий файл
func (w *Watcher) processExistingFile() error {
	file, err := os.Open(w.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var events []interface{}

	for scanner.Scan() {
		line := scanner.Text()
		if event, err := w.parser.ParseLine(line); err == nil && event != nil {
			events = append(events, event)
		}
	}

	if len(events) > 0 && w.callback != nil {
		fmt.Printf("Processing %d existing events\n", len(events))
		w.callback(events)
	}

	return scanner.Err()
}

// watchLoop основной цикл мониторинга
func (w *Watcher) watchLoop() {
	ticker := time.NewTicker(100 * time.Millisecond) // Проверяем каждые 100мс
	defer ticker.Stop()

	for {
		select {
		case <-w.ctx.Done():
			return
		case event := <-w.watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				w.processFileUpdate()
			}
		case err := <-w.watcher.Errors:
			if err != nil {
				fmt.Printf("Watcher error: %v\n", err)
			}
		case <-ticker.C:
			// Периодически проверяем файл на случай пропущенных событий
			w.processFileUpdate()
		}
	}
}

// processFileUpdate обрабатывает обновление файла
func (w *Watcher) processFileUpdate() {
	file, err := os.Open(w.filename)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentLine := 0
	var newEvents []interface{}

	// Пропускаем уже прочитанные строки
	for scanner.Scan() {
		currentLine++
		if currentLine <= w.lastLine {
			continue
		}

		line := scanner.Text()
		if event, err := w.parser.ParseLine(line); err == nil && event != nil {
			newEvents = append(newEvents, event)
			fmt.Printf("Parsed event: %T\n", event)
		} else if err != nil {
			fmt.Printf("Parse error: %v\n", err)
		}
	}

	// Обновляем позицию
	w.lastLine = currentLine

	// Вызываем callback с новыми событиями
	if len(newEvents) > 0 && w.callback != nil {
		w.callback(newEvents)
	}
}
