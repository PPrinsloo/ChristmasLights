package main

func countAll(lights [][]bool) int {
	count, countCh := countLights(lights)

	for range lights {
		count += <-countCh
	}
	return count
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

func countRowLightsOn(lights []bool) int {
	count := 0
	for _, light := range lights {
		if light {
			count++
		}
	}
	return count
}
