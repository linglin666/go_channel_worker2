package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func worker(id int, pathes <-chan string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for path := range pathes {
		fmt.Printf("worker %d starts on path %s\n", id, path)
		time.Sleep(time.Second * 2)
		results <- fmt.Sprintf("Worker %d finished path %s", id, path)
	}
}

func main() {
	const numWorkers = 2

	dirPath := "./test/"
	files, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}
	fmt.Println("Files in directory:", dirPath, files)

	pathes := make(chan string, len(files))
	fmt.Println("pathes channel cap:", cap(pathes), pathes)

	results := make(chan string, len(files))

	var waitGroup sync.WaitGroup

	for wk := 0; wk < numWorkers; wk++ {
		waitGroup.Add(1)
		go worker(wk, pathes, results, &waitGroup)
	}

	// send all the file paths to the pathes channel
	for _, file := range files {
		pathes <- filepath.Join(dirPath + file.Name())
	}
	close(pathes)

	// wait for all workers to finish
	waitGroup.Wait()
	close(results)

	// print results
	for res := range results {
		fmt.Println(res)
	}
}
