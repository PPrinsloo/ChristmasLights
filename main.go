package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// square struct
type rect struct {
	rowStart, columnStart, rowEnd, columnEnd int
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

	// print result
	fmt.Printf("Total lights turned on %d\n\n", count)

	elapsed := time.Since(start)
	fmt.Printf("Processing took %s\n", elapsed)
}

func processList(instList instructionList, lights [][]bool) {
	for _, inst := range instList {
		doInstruction(inst, &lights)
	}
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
