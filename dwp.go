package main

import (
	"bufio"
	"crypto/rand"
	"encoding/binary"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Define command-line flags
	rolls := flag.Int("r", 10, "number of Diceware numbers to generate")
	dictFile := flag.String("d", "", "path to Diceware dictionary file")
	showPassphrase := flag.Bool("p", false, "output complete passphrase")

	// Parse command-line flags
	flag.Parse()

	// Check for invalid input
	if *rolls < 1 {
		fmt.Fprintf(os.Stderr, "Error: Number of rolls must be at least 1\n")
		printUsage()
		os.Exit(1)
	}

	// Load dictionary if specified
	var dict map[int]string
	var err error
	if *dictFile != "" {
		dict, err = loadDictionary(*dictFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading dictionary: %v\n", err)
			os.Exit(1)
		}
	}

	// Number of dice rolls per Diceware number (always 5)
	numDice := 5

	var passphraseWords []string

	// Generate and print Diceware numbers and words
	for i := 0; i < *rolls; i++ {
		dicewareNumber, err := generateDicewareNumber(numDice)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating Diceware number: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Diceware number %d: %05d", i+1, dicewareNumber)
		if dict != nil {
			if word, ok := dict[dicewareNumber]; ok {
				fmt.Printf(" - %s", word)
				passphraseWords = append(passphraseWords, word)
			} else {
				fmt.Printf(" - (word not found in dictionary for number %05d)", dicewareNumber)
			}
		}
		fmt.Println()
	}

	// Output complete passphrase if requested
	if *showPassphrase && len(passphraseWords) > 0 {
		fmt.Printf("\nComplete passphrase: %s\n", strings.Join(passphraseWords, " "))
	}
}

func generateDicewareNumber(numDice int) (int, error) {
	result := 0
	for i := 0; i < numDice; i++ {
		roll, err := secureRandInt(6)
		if err != nil {
			return 0, fmt.Errorf("failed to generate random number: %v", err)
		}
		roll++ // Add 1 to get a number between 1 and 6
		result = result*10 + int(roll)
	}
	return result, nil
}

func secureRandInt(max int32) (int32, error) {
	var n int32
	err := binary.Read(rand.Reader, binary.BigEndian, &n)
	if err != nil {
		return 0, err
	}
	n = n & 0x7FFFFFFF
	return n % max, nil
}

func loadDictionary(filename string) (map[int]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	dict := make(map[int]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "\t")
		if len(parts) >= 2 {
			var number int
			if _, err := fmt.Sscanf(parts[0], "%d", &number); err == nil {
				dict[number] = parts[1]
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return dict, nil
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [-r rolls] [-d dictionary] [-p]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  -r rolls       number of Diceware numbers to generate (default 10)\n")
	fmt.Fprintf(os.Stderr, "  -d dictionary  path to Diceware dictionary file\n")
	fmt.Fprintf(os.Stderr, "  -p             output complete passphrase\n")
	flag.PrintDefaults()
}



