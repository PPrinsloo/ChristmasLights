package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// square struct
type rect struct {
	rowStart, columnStart, rowEnd, columnEnd int
}

// instruction struct
type instruction struct {
	operation string
	rectangle rect
}

// instructionList
type instructionList []instruction

func (instList *instructionList) addInstruction(instruction2 instruction) {
	*instList = append(*instList, instruction2)
}

type LightOperation func(*bool)

const size = 1000

func main() {
	start := time.Now()
	// setup up 1000 * 1000 lights
	lights := make([][]bool, size)
	for i := range lights {
		lights[i] = make([]bool, size)
	}

	instList := fillInstructionList()

	processList(instList, lights)
	count := countAll(lights)

	fmt.Printf("Total lights turned on %d\n\n", count)
	elapsed := time.Since(start)
	fmt.Printf("Processing took %s\n", elapsed)
}

func countAll(lights [][]bool) int {
	count, countCh := countLights(lights)

	for range lights {
		count += <-countCh
	}
	return count
}

func processList(instList instructionList, lights [][]bool) {
	for _, inst := range instList {
		doInstruction(inst, &lights)
	}
}

func fillInstructionList() instructionList {
	instList := make(instructionList, 0)

	// open text file
	file, err := os.Open("lights.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	line := ""

	scanner := bufio.NewScanner(file)
	// read line by line
	for scanner.Scan() {
		cleanAndStoreInstuction(line, scanner, &instList)
	}

	return instList
}

func countLights(lights [][]bool) (int, chan int) {
	count := 0
	countCh := make(chan int)

	for _, row := range lights {
		go func(row []bool) {
			countCh <- countRowLightsOn(row)
		}(row)
	}
	return count, countCh
}

func doInstruction(inst instruction, lights *[][]bool) {
	switch inst.operation {
	case "on":
		operateOnLights(inst.rectangle, lights, turnOn)
	case "off":
		operateOnLights(inst.rectangle, lights, turnOff)
	case "toggle":
		operateOnLights(inst.rectangle, lights, toggle)
	}
}

func turnOn(light *bool) {
	*light = true
}

func turnOff(light *bool) {
	*light = false
}

func toggle(light *bool) {
	*light = !*light
}

func operateOnLights(rectangle rect, lights *[][]bool, operation LightOperation) {

	numCPU := runtime.NumCPU()
	chunkSize := (rectangle.rowEnd - rectangle.rowStart + 1 + numCPU - 1) / numCPU

	var wg sync.WaitGroup
	wg.Add(numCPU)

	for p := 0; p < numCPU; p++ {
		go operateOnChunk(p, rectangle, lights, operation, &wg, chunkSize)
	}

	wg.Wait()
}

func operateOnChunk(p int, rectangle rect, lights *[][]bool, operation LightOperation, wg *sync.WaitGroup, chunkSize int) {
	defer wg.Done()
	rowStart := rectangle.rowStart + p*chunkSize
	rowEnd := min(rectangle.rowStart+(p+1)*chunkSize, rectangle.rowEnd+1)
	for i := rowStart; i < rowEnd; i++ {
		for j := rectangle.columnStart; j <= rectangle.columnEnd; j++ {
			operation(&(*lights)[i][j])
		}
	}
}

func countRowLightsOn(lights []bool) int {
	count := 0
	for _, light := range lights {
		if light {
			count++
		}
	}
	return count
}

func cleanAndStoreInstuction(line string, scanner *bufio.Scanner, instList *instructionList) {
	// do instruction
	line = scanner.Text()
	//handle error
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	operation := ""

	// clean instruction
	if strings.Contains(line, "turn on") {
		operation = "on"
		line = strings.Replace(line, "turn on ", "", 1)
	} else if strings.Contains(line, "turn off") {
		operation = "off"
		line = strings.Replace(line, "turn off ", "", 1)
	} else if strings.Contains(line, "toggle") {
		operation = "toggle"
		line = strings.Replace(line, "toggle ", "", 1)
	}

	// store instruction
	instList.addInstruction(instruction{operation, getRect(line)})
}

func getRect(line string) rect {
	parts := strings.Split(line, " through ")
	start := strings.Split(parts[0], ",")
	end := strings.Split(parts[1], ",")

	rowStart, _ := strconv.Atoi(start[0])
	columnStart, _ := strconv.Atoi(start[1])
	rowEnd, _ := strconv.Atoi(end[0])
	columnEnd, _ := strconv.Atoi(end[1])

	// create rectangle
	rectangle := rect{rowStart, columnStart, rowEnd, columnEnd}

	return rectangle
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
