package main

import "time"

type Database struct {
	localDir string
	Name     string
	Groups   []Group
}

type Group []struct {
	ID   string
	Name string
}

type Entry struct {
	Hash             string
	GroupID          string
	IBAN             string
	Buchungstag      time.Time
	Valutadatum      time.Time
	Buchungstext     string
	Verwendungszweck string
	Amount           float64
	OtherName        string
	OtherIBAN        string
}
