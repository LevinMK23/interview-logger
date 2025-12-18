package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type InterviewLogger struct {
	app         fyne.App
	window      fyne.Window
	inputField  *widget.Entry
	outputField *widget.RichText
	logger      *Logger
	outputText  string
}

func NewInterviewLogger() *InterviewLogger {
	myApp := app.NewWithID("interview.logger")
	myWindow := myApp.NewWindow("Interview Logger")
	myWindow.Resize(fyne.NewSize(800, 600))

	// Получаем путь к домашней директории
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Ошибка получения домашней директории: %v\n", err)
		homeDir = "." // Fallback на текущую директорию
	}

	// Создаем путь к файлу лога в домашней директории
	logPath := filepath.Join(homeDir, "interview.log")
	logger, err := NewLogger(logPath)
	if err != nil {
		fmt.Printf("Ошибка создания логгера: %v\n", err)
	}

	il := &InterviewLogger{
		app:    myApp,
		window: myWindow,
		logger: logger,
	}

	il.setupUI()
	return il
}

func (il *InterviewLogger) setupUI() {
	// Редактируемое текстовое поле (однострочное для работы Enter)
	il.inputField = widget.NewEntry()
	il.inputField.SetPlaceHolder("Введите текст и нажмите Enter для добавления...")
	il.inputField.OnSubmitted = il.handleEnter

	// Нередактируемое текстовое поле
	il.outputField = widget.NewRichText()
	il.outputField.Wrapping = fyne.TextWrapWord
	il.outputText = ""

	// Устанавливаем стиль для черного текста
	il.outputField.Segments = []widget.RichTextSegment{}

	// Кнопка для очистки
	clearButton := widget.NewButton("Очистить", func() {
		il.inputField.SetText("")
		il.outputText = ""
		il.outputField.Segments = []widget.RichTextSegment{}
		il.outputField.Refresh()
	})

	// Кнопка для сохранения лога
	saveButton := widget.NewButton("Сохранить лог", func() {
		il.showSaveDialog()
	})

	// Размещение элементов
	outputLabel := widget.NewLabel("Лог интервью:")

	inputContainer := container.NewBorder(nil, nil, nil, nil, il.inputField)
	outputContainer := container.NewBorder(outputLabel, nil, nil, nil, il.outputField)

	buttons := container.NewHBox(clearButton, saveButton)
	content := container.NewBorder(nil, buttons, nil, nil,
		container.NewVSplit(inputContainer, outputContainer))

	il.window.SetContent(content)
}

func (il *InterviewLogger) handleEnter(text string) {
	if text == "" {
		return
	}

	// Получаем текущее время (только время без даты)
	timestamp := time.Now().Format("15:04:05")

	// Форматируем запись с временной меткой
	logEntry := fmt.Sprintf("[%s] %s", timestamp, text)

	// Добавляем в нередактируемое поле
	if il.outputText != "" {
		il.outputText += "\n"
	}
	il.outputText += logEntry

	// Обновляем RichText с новым текстом
	il.outputField.Segments = []widget.RichTextSegment{
		&widget.TextSegment{
			Text: il.outputText,
		},
	}
	il.outputField.Refresh()

	// Логируем в файл
	if il.logger != nil {
		il.logger.Log(timestamp, text)
	}

	// Очищаем редактируемое поле
	il.inputField.SetText("")
}

func (il *InterviewLogger) showSaveDialog() {
	// Создаем модальное окно
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Введите фамилию и имя...")

	var dialog *widget.PopUp

	// Функция сохранения
	saveFunc := func() {
		name := nameEntry.Text
		if name != "" {
			// Переименовываем файл лога с именем
			if il.logger != nil {
				il.logger.Close()
				err := il.logger.RenameToName(name)
				if err != nil {
					fmt.Printf("Ошибка переименования файла: %v\n", err)
				}
			}
			if dialog != nil {
				dialog.Hide()
			}
			il.app.Quit()
		}
	}

	// Кнопки диалога
	saveDialogButton := widget.NewButton("Сохранить", saveFunc)
	cancelButton := widget.NewButton("Отмена", func() {
		if dialog != nil {
			dialog.Hide()
		}
	})

	// Обработка Enter в поле ввода
	nameEntry.OnSubmitted = func(text string) {
		if text != "" {
			saveFunc()
		}
	}

	// Создаем контент диалога
	content := container.NewVBox(
		widget.NewLabel("Введите фамилию и имя:"),
		nameEntry,
		container.NewHBox(saveDialogButton, cancelButton),
	)

	// Создаем модальное окно
	dialog = widget.NewModalPopUp(
		container.NewBorder(
			nil, nil, nil, nil,
			container.NewPadded(content),
		),
		il.window.Canvas(),
	)

	dialog.Resize(fyne.NewSize(400, 150))
	dialog.Show()
}

func (il *InterviewLogger) Run() {
	il.window.ShowAndRun()
}

func main() {
	logger := NewInterviewLogger()
	logger.Run()
}
