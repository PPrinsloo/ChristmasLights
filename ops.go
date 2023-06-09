package main

func turnOn(light *bool) {
	*light = true
}

func turnOff(light *bool) {
	*light = false
}

func toggle(light *bool) {
	*light = !*light
}
