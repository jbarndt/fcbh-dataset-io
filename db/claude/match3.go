package main

import (
	"fmt"
)

type Country struct {
	ISO3 string
	ISO1 string
	// Other fields can be added here
}

func matchAndUpdate(list1, list2 []Country) ([]Country, []string, []string) {
	// Create maps for faster lookup
	map2 := make(map[string]Country)
	for _, country := range list2 {
		map2[country.ISO3] = country
	}

	var updatedList1 []Country
	var missingInList2 []string
	seenInList1 := make(map[string]bool)

	// Check list1 against list2
	for _, country1 := range list1 {
		seenInList1[country1.ISO3] = true
		if country2, found := map2[country1.ISO3]; found {
			// Update ISO1 if found in list2
			country1.ISO1 = country2.ISO1
			updatedList1 = append(updatedList1, country1)
		} else {
			missingInList2 = append(missingInList2, country1.ISO3)
		}
	}

	// Check for items in list2 missing from list1
	var missingInList1 []string
	for iso3 := range map2 {
		if !seenInList1[iso3] {
			missingInList1 = append(missingInList1, iso3)
		}
	}

	return updatedList1, missingInList2, missingInList1
}

func main() {
	list1 := []Country{
		{ISO3: "USA", ISO1: ""},
		{ISO3: "CAN", ISO1: ""},
		{ISO3: "MEX", ISO1: ""},
		{ISO3: "BRA", ISO1: ""},
		{ISO3: "USA", ISO1: ""}, // Duplicate entry
	}

	list2 := []Country{
		{ISO3: "USA", ISO1: "US"},
		{ISO3: "CAN", ISO1: "CA"},
		{ISO3: "FRA", ISO1: "FR"},
		{ISO3: "DEU", ISO1: "DE"},
		{ISO3: "USA", ISO1: "US"}, // Duplicate entry
	}

	updatedList1, missingInList2, missingInList1 := matchAndUpdate(list1, list2)

	fmt.Println("Updated List 1:")
	for _, country := range updatedList1 {
		fmt.Printf("ISO3: %s, ISO1: %s\n", country.ISO3, country.ISO1)
	}

	fmt.Println("\nMissing in List 2:")
	for _, iso3 := range missingInList2 {
		fmt.Println(iso3)
	}

	fmt.Println("\nMissing in List 1:")
	for _, iso3 := range missingInList1 {
		fmt.Println(iso3)
	}
}
