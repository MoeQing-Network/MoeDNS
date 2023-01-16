package utils

import (
	"encoding/csv"
	"io"
	"log"
	"net/netip"
	"os"
)

func FindPrefix(ip netip.Addr) bool {
	path := "cn-prefix.csv"
	csvfile, err := os.Open(path)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	defer csvfile.Close()

	// Parse the file
	r := csv.NewReader(csvfile)
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// fmt.Printf("%s %s \n", record[0], record[1])
		p, err := netip.ParsePrefix(record[0])
		if err != nil {
			log.Fatal(err)
		}
		if p.Contains(ip) {
			return true
		}
	}
	return false
}
