package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// handle input/output
func Readfile(filename string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("file no readable")
		os.Exit(1)
	}
	return string(data)
}

// convert string of data in byte slice
func Writefile(filename string, data string) {
	err := os.WriteFile(filename, []byte(data), 0644)
	if err != nil {
		fmt.Println("file no writeable")
		os.Exit(1)
	}
}

// define a map for punctuation marks
var punctuationMarks = map[string]bool{
	",": true,
	".": true,
	"!": true,
	"?": true,
	":": true,
	";": true,
}

func ModifySlice(s string) []string {
	slice := strings.Fields(s)
	// Define a map for vowels marks
	vowels := map[string]bool{
		"a": true,
		"e": true,
		"i": true,
		"o": true,
		"u": true,
		"h": true,
		"A": true,
		"E": true,
		"I": true,
		"O": true,
		"U": true,
		"H": true,
	}
	// if for loop encounters 'a' and the next word starts with a vowel it replaces 'a' with "an"
	for i := 0; i < len(slice); i++ {
		word := slice[i]
		if (word == "a" || word == "A") && i < len(slice)-1 && vowels[string(slice[i+1][0])] {

			slice[i] += "n" // lisätään vain n jossa a on vain ensimmäinen kirjain
		}
	}
	return slice
}
func Modifytext(slice []string) string {

	flag := false
	for i := 0; i < len(slice); i++ {
		word := []byte(slice[i])
		if slice[i] == "(hex)" && i > 0 {
			decimalValue := ConvertHexToDecimal(slice[i-1])
			slice[i-1] = decimalValue
			slice[i] = ""
		} else if slice[i] == "(bin)" && i > 0 {
			decimalValue := ConvertBinToDecimal(slice[i-1])
			slice[i-1] = decimalValue
			slice[i] = ""
		} else if slice[i] == "(up)" && i > 0 {
			slice[i-1] = strings.ToUpper(slice[i-1])
			slice[i] = ""
		} else if slice[i] == "(low)" && i > 0 {
			slice[i-1] = strings.ToLower(slice[i-1])
			slice[i] = ""
		} else if slice[i] == "(cap)" && i > 0 {
			slice[i-1] = Capitilize(slice[i-1])
			slice[i] = ""
			slice = append(slice[:i], slice[i+1:]...)
			i--

		} else if (slice[i] == "(up,") || (slice[i] == "(low,") || (slice[i] == "(cap,") {
			num := ExtractNumber(slice[i+1])
			if num > 0 && i >= num {
				for j := 1; j <= num; j++ {

					if slice[i] == "(up," {
						slice[i-j] = strings.ToUpper(slice[i-j])
					} else if slice[i] == "(low," {
						slice[i-j] = strings.ToLower(slice[i-j])
					} else if slice[i] == "(cap," {
						slice[i-j] = Capitilize(slice[i-j])
					}
				}
			}
			slice = append(slice[:i], slice[i+2:]...)
			i--
			//slice[i], slice[i+1] = "", ""

			// add punctuation to the prev element
		} else if punctuationMarks[slice[i]] && i > 0 {
			slice[i-1] += slice[i]
			slice = append(slice[:i], slice[i+1:]...)
			i--

			// check if the word has more than one character & if the first char is punctuation mark

		} else if len(word) > 1 && punctuationMarks[string(word[0])] {
			j := 0
			for j < len(word) && punctuationMarks[string(word[j])] {
				j++
			}
			if i > 0 {
				slice[i-1] += string(word[:j])
			}
			slice[i] = " " + string(word[j:])

			// use flag to distinguish start / end single quote
		} else if slice[i] == "'" {
			if !flag {
				slice[i+1] = "'" + slice[i+1]
				slice[i] = ""
				// slice = append(slice[:i], slice[i+1:]...)
				// i--
				flag = true
			} else {
				slice[i-1] += "'"
				slice[i] = ""
				// slice = append(slice[:i], slice[i+1:]...)
				// i--
				flag = false
			}
		}
	}
	return strings.Join(slice, " ")
}

func Capitilize(s string) string {
	s = strings.ToLower(s)
	letters := []byte(s)
	for i, c := range letters {
		if c >= 'a' && c <= 'z' {
			c -= 32
			letters[i] = c
			break
		}
	}
	// return strings.ToUpper(strings(s[0])) + strings.ToLower(strings(s[1:]))
	return string(letters)
}

// ExtractNumber extracts a number from a string "3)"
func ExtractNumber(s string) int {
	end := strings.Index(s, ")")
	number, err := strconv.Atoi(s[:end])
	if err != nil {
		fmt.Println("error in Extract the number: ", err)
		os.Exit(1)
	}
	return number
}

// parseint converts strings to values.. Formatint converts values to strings
func ConvertHexToDecimal(hexStr string) string {
	num, err := strconv.ParseInt(hexStr, 16, 64)
	if err != nil {
		fmt.Println("problem convert from hex to decimal:", err)
		os.Exit(1)
	}
	return strconv.FormatInt(num, 10)
}

func ConvertBinToDecimal(binStr string) string {
	num, err := strconv.ParseInt(binStr, 2, 64)
	if err != nil {
		fmt.Println("problem convert from binary to decimal:", err)
		os.Exit(1)
	}
	return strconv.FormatInt(num, 10)
}

func main() {
	if len(os.Args) != 3 {
		log.Fatal("Usage: go run main.go <inputfile> <outputfile>")
	}
	text := Readfile(os.Args[1])

	modifedSlice := ModifySlice(text)

	modifiedText := Modifytext(modifedSlice)

	modifiedText = strings.Replace(modifiedText, "  ", " ", -1)

	Writefile(os.Args[2], modifiedText)

}
