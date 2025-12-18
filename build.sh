#!/bin/bash

# Скрипт для сборки .app bundle для macOS

echo "Сборка приложения Interview Logger..."

# Собираем бинарник
go build -o "Interview Logger" .

# Создаем структуру .app bundle
mkdir -p "Interview Logger.app/Contents/MacOS"
mkdir -p "Interview Logger.app/Contents/Resources"

# Перемещаем бинарник
mv "Interview Logger" "Interview Logger.app/Contents/MacOS/Interview Logger"

# Делаем исполняемым
chmod +x "Interview Logger.app/Contents/MacOS/Interview Logger"

# Добавляем иконку, если она существует
if [ -f "interview.png" ]; then
    echo "Добавление иконки..."
    sips -s format icns interview.png --out "Interview Logger.app/Contents/Resources/icon.icns" 2>/dev/null
    if [ $? -eq 0 ]; then
        echo "Иконка добавлена успешно"
    else
        echo "Предупреждение: не удалось создать .icns, копирую PNG"
        cp interview.png "Interview Logger.app/Contents/Resources/icon.png"
    fi
fi

echo "Готово! Приложение создано: Interview Logger.app"
echo "Запустите его двойным кликом или через: open 'Interview Logger.app'"

