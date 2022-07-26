package main

import (
	"aurora-tracker/noaa"
	"aurora-tracker/utils"
	"os"
)

var (
	threshold = os.Getenv("KP_INDEX_THRESHOLD")
)

func main() {
	kpIndexThreshold := 7
	if threshold != "" {		
		kpIndexThreshold = utils.ConvertToInt(threshold)
	}
	noaa.CheckForAuroraProbability(kpIndexThreshold)
}
