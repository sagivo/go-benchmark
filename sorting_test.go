package main

import (
	"bufio"
	"log"
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
	if len(l) < 2 {
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

func load() ([]int, error) {
	lines, err := readLines("arr.txt")
	if err != nil {
		return nil, err
	}

	l := make([]int, len(lines))
	for i, line := range lines {
		if i <= len(l) {
			l[i], err = strconv.Atoi(line)
			if err != nil {
				return nil, err
			}
		}
	}
	return l, nil
}

func BenchmarkMergesort(b *testing.B) {
	l, err := load()
	if err != nil {
		log.Fatalln("Error while loading input file: ", err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		mergeSort(l)
	}
}

func BenchmarkMergesortAsync(b *testing.B) {
	l, err := load()
	if err != nil {
		log.Fatalln("Error while loading input file: ", err)
	}
	c := make(chan []int)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		mergeSortAsync(l, c)
		l = <-c
	}
}

func BenchmarkQuicksort(b *testing.B) {
	l, err := load()
	if err != nil {
		log.Fatalln("Error while loading input file: ", err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		sort.Ints(l)
	}
}
