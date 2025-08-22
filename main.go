package main

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"unicode"
)

func main() {
	WordCount("go is fun go")
	AreAnagrams("Hello", "World")
	FirstUnique("aabdbcc")
	RemoveDuplicates([]int{4, 4, 4, 4})
	RemoveElement([]int{1, 2, 3}, 5)
	IsPalindrome("А роза упала на лапу Азора")
	drawChessboard(8)
}

func WordCount(s string) map[string]int {
	words := strings.Fields(s)
	result := make(map[string]int)
	for _, word := range words {
		result[word]++
	}
	fmt.Println(result)
	return result
}

func AreAnagrams(s1, s2 string) bool {
	s1ToLower := strings.ToLower(s1)
	s1Trimmed := strings.Trim(s1ToLower, " ")
	s2ToLower := strings.ToLower(s2)
	s2Trimmed := strings.Trim(s2ToLower, " ")
	firstMap := make(map[rune]int)
	secondMap := make(map[rune]int)
	for _, word := range s1Trimmed {
		firstMap[word]++
	}
	for _, r := range s2Trimmed {
		secondMap[r]++
	}
	if reflect.DeepEqual(firstMap, secondMap) {
		fmt.Println(true)
		return true
	}
	fmt.Println(false)
	return false
}

func FirstUnique(s string) rune {
	count := make(map[rune]int)
	for _, r := range s {
		count[r]++
	}
	for _, r := range s {
		if count[r] == 1 {
			fmt.Println(string(r))
			return r
		}
	}
	fmt.Println(0)
	return rune(0)
}

func RemoveDuplicates(nums []int) []int {
	var result []int
	for _, num := range nums {
		if slices.Contains(result, num) {
			continue
		}
		result = append(result, num)
	}
	fmt.Println(result)
	return result
}

func RemoveElement(nums []int, index int) ([]int, error) {
	if index > len(nums)-1 || index < 0 {
		fmt.Println(nil)
		return nil, errors.New(fmt.Sprintf("Index out of range: %d", index))
	}
	var left = nums[:index]
	var right = nums[index+1:]
	var result []int
	result = append(left, right...)
	fmt.Println(result)
	return result, nil
}

func IsPalindrome(s string) bool {
	var cleaned strings.Builder
	for _, char := range s {
		if unicode.IsLetter(char) {
			cleaned.WriteRune(unicode.ToLower(char))
		}
	}

	cleanedStr := cleaned.String()
	runes := []rune(cleanedStr)
	length := len(runes)
	for i := 0; i < length/2; i++ {
		if runes[i] != runes[length-i-1] {
			fmt.Println(false)
			return false
		}
	}
	fmt.Println(true)
	return true
}

func drawChessboard(size int) {
	board := make([][]rune, size)
	for i := 0; i < size; i++ {
		board[i] = make([]rune, size)
		for j := 0; j < size; j++ {
			if (i+j)%2 == 0 {
				board[i][j] = ' '
			} else {
				board[i][j] = '#'
			}
		}
	}
	// Отрисовка доски
	for _, row := range board {
		for _, cell := range row {
			fmt.Printf("%c", cell)
		}
		fmt.Println()
	}
}
