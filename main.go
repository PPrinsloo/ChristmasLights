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

/**
Deploy one million lights in a 1000x1000 grid
Lights in your grid are numbered from 0 to 999 in each direction;
	the lights at each corner are at 0,0, 0,999, 999,999, and 999,0.
The instructions include whether to turn on, turn off, or toggle various inclusive ranges given as coordinate pairs.
Each coordinate pair represents opposite corners of a rectangle, inclusive;
	a coordinate pair like 0,0 through 2,2 therefore refers to 9 lights in a 3x3 square.
The lights all start turned off.
For example:
turn on 0,0 through 999,999 would turn on (or leave on) every light.
toggle 0,0 through 999,0 would toggle the first line of 1000 lights, turning off the ones that were on, and turning on the ones that were off.
turn off 499,499 through 500,500 would turn off (or leave off) the middle four lights.

More examples:
turn off 660,55 through 986,197
turn off 341,304 through 638,850
turn off 199,133 through 461,193
toggle 322,558 through 977,958
toggle 537,781 through 687,941
turn on 226,196 through 599,390
turn on 240,129 through 703,297
turn on 317,329 through 451,798

Kata Goal: After following the instructions, how many lights are lit?
*/

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

	printLightsGrid(lights)
}

func countAll(lights [][]bool) int {
	count, countCh := countLights(lights)

	for range lights {
		count += <-countCh
	}
	return count
}

func processList(instList instructionList, lights [][]bool) {
	start := time.Now()
	for _, inst := range instList {
		doInstruction(inst, &lights)
	}
	elapsed := time.Since(start)
	fmt.Printf("Processing instructions took %s\n", elapsed)
}

func fillInstructionList() instructionList {
	start := time.Now()
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

	elapsed := time.Since(start)
	fmt.Printf("Filling instructions set took %s\n", elapsed)

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
		turnOn(inst.rectangle, lights)
	case "off":
		turnOff(inst.rectangle, lights)
	case "toggle":
		toggle(inst.rectangle, lights)
	}
}

func turnOn(rectangle rect, lights *[][]bool) {
	start := time.Now()

	numCPU := runtime.NumCPU()
	chunkSize := (rectangle.rowEnd - rectangle.rowStart + 1 + numCPU - 1) / numCPU

	var wg sync.WaitGroup
	wg.Add(numCPU)

	for p := 0; p < numCPU; p++ {
		go func(p int) {
			defer wg.Done()
			rowStart := rectangle.rowStart + p*chunkSize
			rowEnd := min(rectangle.rowStart+(p+1)*chunkSize, rectangle.rowEnd+1)
			for i := rowStart; i < rowEnd; i++ {
				for j := rectangle.columnStart; j <= rectangle.columnEnd; j++ {
					(*lights)[i][j] = true
				}
			}
		}(p)
	}

	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("Turn on took %s\n", elapsed)
}

func turnOff(rectangle rect, lights *[][]bool) {
	start := time.Now()

	numCPU := runtime.NumCPU()
	chunkSize := (rectangle.rowEnd - rectangle.rowStart + 1 + numCPU - 1) / numCPU

	var wg sync.WaitGroup
	wg.Add(numCPU)

	for p := 0; p < numCPU; p++ {
		go func(p int) {
			defer wg.Done()
			rowStart := rectangle.rowStart + p*chunkSize
			rowEnd := min(rectangle.rowStart+(p+1)*chunkSize, rectangle.rowEnd+1)
			for i := rowStart; i < rowEnd; i++ {
				for j := rectangle.columnStart; j <= rectangle.columnEnd; j++ {
					(*lights)[i][j] = false
				}
			}
		}(p)
	}

	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("Turn off took %s\n", elapsed)
}

func toggle(rectangle rect, lights *[][]bool) {
	start := time.Now()

	numCPU := runtime.NumCPU()
	chunkSize := (rectangle.rowEnd - rectangle.rowStart + 1 + numCPU - 1) / numCPU

	var wg sync.WaitGroup
	wg.Add(numCPU)

	for p := 0; p < numCPU; p++ {
		go func(p int) {
			defer wg.Done()
			rowStart := rectangle.rowStart + p*chunkSize
			rowEnd := min(rectangle.rowStart+(p+1)*chunkSize, rectangle.rowEnd+1)
			for i := rowStart; i < rowEnd; i++ {
				for j := rectangle.columnStart; j <= rectangle.columnEnd; j++ {
					(*lights)[i][j] = !(*lights)[i][j]
				}
			}
		}(p)
	}

	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("Toggle took %s\n", elapsed)
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

func printLightsGrid(lights [][]bool) {
	for _, row := range lights {
		for _, light := range row {
			if light {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
