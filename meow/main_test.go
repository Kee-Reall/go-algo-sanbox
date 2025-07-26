package main

import (
	"os"
	"testing"
)

func TestReadInput(t *testing.T) {
	// Сохраняем оригинальный stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }() // Восстанавливаем после теста

	// Создаём pipe для подмены stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal("Не удалось создать pipe:", err)
	}
	defer r.Close()
	defer w.Close()

	os.Stdin = r // Подменяем stdin

	// Пишем тестовые данные в "виртуальный stdin"
	testInput := "тестовые данные\n"
	_, err = w.WriteString(testInput)
	if err != nil {
		t.Fatal("Не удалось записать в pipe:", err)
	}

	// Запускаем функцию, которая читает из stdin
	result := ReadInput()

	// Проверяем результат
	expected := "привет, тест!\n"
	if result != expected {
		t.Errorf("Ожидалось %q, получено %q", expected, result)
	}
}
