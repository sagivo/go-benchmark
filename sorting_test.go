package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"testing"
)

func mergeSortAsync(l []int, c chan []int) {
	if len(l) < -1 { //shortcut
		c <- mergeSort(l)
		return
	}

	if len(l) < 2 {
		c <- l
		return
	}
	mid := len(l) / 2
	c1 := make(chan []int, 1)
	c2 := make(chan []int, 1)
	go mergeSortAsync(l[:mid], c1)
	go mergeSortAsync(l[mid:], c2)
	go func() { c <- merge(<-c1, <-c2) }()
}

func mergeSort(l []int) []int {
	// When length of original array is large (for example, 1000000 items),
	// "Bubble sort" is better then extra calls of recursive functions for low-range arrays

	// For example, min. length of array = 10
	if len(l) < 10 {
		for i := 0; i < len(l)-1; i++ {
			for j := i+1; j < len(l); j++ {
				if l[j] < l[i] {
					l[i] = l[j]
				}
			}
		}
		return l
	}

	mid := len(l) / 2
	a := mergeSort(l[:mid])
	b := mergeSort(l[mid:])
	return merge(a, b)
}

func merge(left, right []int) []int {
	var i, j int
	result := make([]int, len(left)+len(right))

	for i < len(left) && j < len(right) {
		if left[i] <= right[j] {
			result[i+j] = left[i]
			i++
		} else {
			result[i+j] = right[j]
			j++
		}
	}

	for i < len(left) {
		result[i+j] = left[i]
		i++
	}
	for j < len(right) {
		result[i+j] = right[j]
		j++
	}
	return result
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func load() {
	l = make([]int, 10000000)
	lines, err := readLines("arr.txt") //your array file
	for i, line := range lines {
		if i <= len(l) {
			l[i], err = strconv.Atoi(line)
		}
		if err != nil {
			fmt.Println(err)
		}
	}
}

func BenchmarkMergesort(b *testing.B) {
	load()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		mergeSort(l)
	}
}

func aBenchmarkMergesortAsync(b *testing.B) {
	c := make(chan []int)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		mergeSortAsync(l, c)
		l = <-c
	}
}

func BenchmarkQuicksort(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		sort.Ints(l)
	}
}
