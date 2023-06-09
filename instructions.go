package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

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
