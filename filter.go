package main

import "strings"

func filterTwitterURL(text string) string {
	var key string = "twitter.com/"
	if !strings.Contains(text, key) {
		return text
	}
	var spaceArr []string = strings.Split(text, " ")
	for i, spaceUnit := range spaceArr {
		if strings.Contains(spaceUnit, key) {
			spaceArr[i] = strings.Split(spaceUnit, "?")[0]
		}
	}
	return strings.Join(spaceArr, " ")
}
