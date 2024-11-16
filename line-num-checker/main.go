package main

import (
	"bufio"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
)

var wg sync.WaitGroup

func getLinesOfCode(path string, linesCount *int) {
	defer wg.Done()
	file, err := os.Open(path)
	if err != nil {

		panic(err)
	}

	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Fatalf("error reading line: %s", err)
		}
		if len(line) > 0 {
			// Avoid counting comments
			if len(line) > 2 && string(line[0]) == "/" && (string(line[1]) == "/" || string(line[1]) == "*") {
				break
			}

			*linesCount = *linesCount + 1
		}
	}
}

func dfs(rootPath string) struct {
	linesCount      int
	goRoutinesCount int
} {

	linesCount := 0
	goRoutinesCount := 0

	files, err := os.ReadDir(rootPath)
	if err != nil {
		log.Fatalln("Issue reading dir", rootPath)
		panic(err)
	}

	if len(files) == 0 {
		return struct {
			linesCount      int
			goRoutinesCount int
		}{
			linesCount:      0,
			goRoutinesCount: 0,
		}
	}

	for i := 0; i < len(files); i++ {
		var file = files[i]
		if file.Name() == "node_modules" {
			continue
		}
		if !file.IsDir() {
			wg.Add(1)
			go getLinesOfCode(rootPath+"/"+file.Name(), &linesCount)
		} else {
			var root = rootPath + "/" + file.Name()
			var data = dfs(root)
			linesCount = linesCount + data.linesCount
			goRoutinesCount = goRoutinesCount + data.goRoutinesCount
		}
	}
	goRoutinesCount = goRoutinesCount + runtime.NumGoroutine()
	wg.Wait()
	return struct {
		linesCount      int
		goRoutinesCount int
	}{
		linesCount:      linesCount,
		goRoutinesCount: goRoutinesCount,
	}
}

/*
func create(count int, root string) {
	if count < 1 {
		return
	}

	for i := 1; i <= count; i++ {
		var dirName = strconv.Itoa(i)
		//log.Println(dirName)
		log.Println("CREATE", root+dirName)
		err := os.MkdirAll(root+dirName, 0750)
		if err != nil {
			log.Fatalln(err)
		}

		for x := 1; x <= count; x++ {
			var xToString = strconv.Itoa(x)
			file, err := os.Create(root + dirName + "/" + dirName + "test" + xToString + ".js")
			if err != nil {
				panic(err)
			}
			var fileContent = "Hello from" + dirName + "\n"
			for y := 1; y <= i; y++ {
				fileContent = fileContent + "Hello from" + dirName + "\n"
			}
			file.WriteString(fileContent)
		}
		create(count-1, root+dirName+"/")
	}
} */

func main() {
	linesCount := 0
	goRoutinesCount := 0

	start := time.Now()
	var data = dfs(".")
	linesCount = data.linesCount
	goRoutinesCount = data.goRoutinesCount

	wg.Wait()
	elapsed := time.Since(start)

	log.Println("LINES COUNT", elapsed, linesCount, goRoutinesCount)

}
