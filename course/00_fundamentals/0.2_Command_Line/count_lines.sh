#!/bin/bash

#проверка: передан ли аргумент с путем к папкею если нет -используем текущую папку
TARGET_DIR="${1:-.}"

echo "Сканирование директории: $TARGET_DIR"

total_lines=$(find "$TARGET_DIR" -type f -name "*.go" -exec cat {} + | grep -v "^$"  |  wc -l)

echo "Общее количество строк в .go файлах: $total_lines"
