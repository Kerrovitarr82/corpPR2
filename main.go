package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

type FileAnalysis struct {
	FileName  string
	WordCount int
	CharCount int
}

func analyzeFile(filePath string, wg *sync.WaitGroup, mu *sync.Mutex, results *[]FileAnalysis) {
	defer wg.Done()

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Ошибка открытия файла %s: %v\n", filePath, err)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var wordCount, charCount int

	for {
		line, err := reader.ReadString('\n')
		wordCount += len(strings.Fields(line))
		charCount += len([]rune(line))
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Ошибка чтения файла %s: %v\n", filePath, err)
			return
		}
	}

	mu.Lock()
	*results = append(*results, FileAnalysis{
		FileName:  filePath,
		WordCount: wordCount,
		CharCount: charCount,
	})
	mu.Unlock()
}

func main() {
	files := []string{"text.txt", "text.txt", "text.txt", "text.txt", "text.txt"} // список файлов

	var wg sync.WaitGroup
	var mu sync.Mutex
	var results []FileAnalysis

	for _, file := range files {
		wg.Add(1)
		go analyzeFile(file, &wg, &mu, &results)
	}

	wg.Wait()

	fmt.Println("Результаты анализа:")
	totalWords, totalChars := 0, 0
	for i, res := range results {
		fmt.Printf("%d. %s: %d слов, %d символов\n", i+1, res.FileName, res.WordCount, res.CharCount)
		totalWords += res.WordCount
		totalChars += res.CharCount
	}
	fmt.Printf("\nИтог: %d слов, %d символов\n", totalWords, totalChars)
}
