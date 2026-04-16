package main

import (
	"fmt"
)

//реализация стека на го + решение задачи

// стек для символов
type Stack struct {
	items []rune
}

func NewStack() *Stack {
	return &Stack{
		items: make([]rune, 0),
	}
}

// добавление в стек
func (s *Stack) Push(item rune) {
	s.items = append(s.items, item)
}

// удаление
func (s *Stack) Pop() (rune, error) {
	if s.IsEmpty() {
		return 0, fmt.Errorf("стек пуст")
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item, nil
}

func (s *Stack) Peek() (rune, error) {
	if s.IsEmpty() {
		return 0, fmt.Errorf("стек пуст")
	}
	return s.items[len(s.items)-1], nil
}

func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}

func (s *Stack) Size() int {
	return len(s.items)
}

func (s *Stack) Clear() {
	s.items = make([]rune, 0)
}

// проверка скобочной последовательности
func isValid(s string) bool {
	stack := NewStack()

	// Словарь соответствия закрывающих скобок открывающим
	pairs := map[rune]rune{
		')': '(',
		']': '[',
		'}': '{',
	}

	// Множество открывающих скобок
	openBrackets := map[rune]bool{
		'(': true,
		'[': true,
		'{': true,
	}

	for _, char := range s {
		//если скобка открывающася, то добавляем в стек
		if openBrackets[char] {
			stack.Push(char)
		} else if closing, ok := pairs[char]; ok { // ЗНАЧИТ ЗАКРЫВАЮЩАЯСЯ
			if stack.IsEmpty() {
				return false
			}

			top, _ := stack.Pop()
			if top != closing {
				return false
			}
		}
	}

	// если в конце ничего не осталось, значит у каждой скобки своя пара
	return stack.IsEmpty()
}

// варивативная часть
func isValidOnlyParentheses(s string) bool {
	stack := NewStack()

	for _, char := range s {
		if char == '(' {
			stack.Push(char)
		} else if char == ')' {
			if stack.IsEmpty() {
				return false
			}
			stack.Pop()
		}
	}

	return stack.IsEmpty()
}

func main() {
	testCases := []struct {
		expression string
		expected   bool
	}{
		{"()[]{}", true},
		{"([)]", false},
		{"{[]}", true},
		{"(", false},
		{")", false},
		{"", true},
	}

	for _, tc := range testCases {
		result := isValid(tc.expression)
		status := "ok"
		if result != tc.expected {
			status = "ne ok"
		}
		fmt.Printf("%s %s → %v (ожидалось: %v)\n",
			status, tc.expression, result, tc.expected)
	}

	//вариативная часть проверка
	fmt.Println("\n=== Проверка только круглых скобок ===")
	testCasesParentheses := []struct {
		expression string
		expected   bool
	}{
		{"()", true},
		{"(", false},
		{")", false},
		{"", true},
		{"(())", true},
	}

	for _, tc := range testCasesParentheses {
		result := isValidOnlyParentheses(tc.expression)
		status := "ok"
		if result != tc.expected {
			status = "ne ok"
		}
		fmt.Printf("%s %s → %v (ожидалось: %v)\n",
			status, tc.expression, result, tc.expected)
	}

	stack := NewStack()

	stack.Push('A')
	stack.Push('B')
	stack.Push('C')

	//size
	stack.Size()

	if top, err := stack.Peek(); err == nil {
		fmt.Printf("Верхний элемент (без удаления): %c\n", top)
	}

	// извлечение всех элементов
	fmt.Println("\nИзвлекаем элементы:")
	for !stack.IsEmpty() {
		item, _ := stack.Pop()
		fmt.Printf("Извлечен: %c\n", item)
	}

	fmt.Printf("\nСтек пуст: %v\n", stack.IsEmpty())
}
