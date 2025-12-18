package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Logger struct {
	file     *os.File
	filename string
	mutex    sync.Mutex
}

func NewLogger(filename string) (*Logger, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия файла лога: %w", err)
	}

	return &Logger{
		file:     file,
		filename: filename,
	}, nil
}

func (l *Logger) Log(timestamp, text string) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	_, err := fmt.Fprintf(l.file, "[%s] %s\n", timestamp, text)
	if err != nil {
		return fmt.Errorf("ошибка записи в лог: %w", err)
	}

	// Синхронизируем запись на диск
	return l.file.Sync()
}

// RenameToName переименовывает файл лога на основе имени
func (l *Logger) RenameToName(name string) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	// Закрываем текущий файл
	if l.file != nil {
		l.file.Close()
	}

	// Нормализуем имя для использования в имени файла
	normalizedName := normalizeFileName(name)
	if normalizedName == "" {
		normalizedName = "interview"
	}

	// Формируем новое имя файла
	dir := filepath.Dir(l.filename)
	newFilename := filepath.Join(dir, normalizedName+".log")

	// Переименовываем файл
	err := os.Rename(l.filename, newFilename)
	if err != nil {
		return fmt.Errorf("ошибка переименования файла: %w", err)
	}

	return nil
}

// normalizeFileName нормализует имя для использования в имени файла
func normalizeFileName(name string) string {
	// Заменяем пробелы на подчеркивания
	name = strings.ReplaceAll(name, " ", "_")

	// Удаляем недопустимые символы для имен файлов
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range invalidChars {
		name = strings.ReplaceAll(name, char, "")
	}

	// Убираем лишние подчеркивания
	name = strings.Trim(name, "_")

	return name
}

func (l *Logger) Close() error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.file != nil {
		return l.file.Close()
	}
	return nil
}
