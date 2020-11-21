package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

func main() {
	// build isd:region map
	regions, err := os.Open("regionmap.csv")
	if err != nil {
		panic("Failed to open region map")
	}
	defer regions.Close()

	var regionmap = make(map[string]int)

	reader := csv.NewReader(regions)
	reader.Comma = '\t'

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err == csv.ErrFieldCount || len(record) < 2 {
			continue
		}
		region, err := strconv.Atoi(record[1])
		if err != nil {
			continue
		}
		regionmap[record[0]] = region
	}

	// build all state audition result isd count map
	result, err := os.Open("result.csv")

	if err != nil {
		panic("Failed to open audition result")
	}
	defer result.Close()

	reader = csv.NewReader(result)
	reader.Comma = '\t'
	isdResultmap := make(map[string]int)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err == csv.ErrFieldCount || len(record) < 5 {
			continue
		}
		isdResultmap[record[4]]++
	}
	//fmt.Println(regionmap)

	/*
		A - Regions 1, *6/7, 16, 22
		B - Regions 5, 24, 30, 31
		C - Regions 2, 3, 4, 25
		D - Regions 8, 20, 21, 26
		E - Regions 12, 18, 29, 32
		F - Regions 9, 10, 27, 33
		G - Regions 11, 14, 15, 28
		H - Regions 13, 17, 19, 23
	*/
	areaMap := [][]int{
		{1, 6, 7, 16, 22},
		{5, 24, 30, 31},
		{2, 3, 4, 25},
		{8, 20, 21, 26},
		{12, 18, 29, 32},
		{9, 10, 27, 33},
		{11, 14, 15, 28},
		{13, 17, 19, 23},
	}

	// build area-region map
	regionAreaMap := make(map[int]int)
	for a, rs := range areaMap {
		for _, r := range rs {
			regionAreaMap[r] = a
		}
	}

	// area final count that made to all-state
	var areaAllstate = make(map[int]int)
	var regionAllstate = make(map[int]int)

	for isd, count := range isdResultmap {
		region := regionmap[isd]
		if region == 0 {
			fmt.Println("Not found isd in region", isd)
		}
		regionAllstate[region] += count
		area := regionAreaMap[region]
		areaAllstate[area] += count
	}

	// print result
	fmt.Println("\nRegion All State")
	for k, val := range regionAllstate {
		fmt.Printf("Region %d: %d\n", k, val)
	}
	fmt.Println("\nArea All State")
	for k, val := range areaAllstate {
		fmt.Printf("Area %d: %d\n", k+1, val)
	}
}
