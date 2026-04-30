package hcutil

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// ContainsString checks if a string is contained within a string slice
func ContainsString(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

// ContainsInt checks if an int is contained within an int slice
func ContainsInt(slice []int, item int) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}

	return false
}

// RemoveString removes a string from a string slice if it exists
func RemoveString(stringSlice []string, itemToRemove string) []string {
	for idx, item := range stringSlice {
		if item == itemToRemove {
			ret := make([]string, 0)
			ret = append(ret, stringSlice[:idx]...)
			return append(ret, stringSlice[idx+1:]...)
		}
	}

	return stringSlice
}

// HasCommonElements determines if there are any common elements in both
// of the provided arrays.
func HasCommonElements(arr1 []string, arr2 []string) bool {
	// Create a map to store the elements from arr1
	elementsMap := make(map[string]bool)

	// Add elements from arr1 to the map
	for _, elem := range arr1 {
		elementsMap[elem] = true
	}

	// Check if any element from arr2 exists in the map
	for _, elem := range arr2 {
		if elementsMap[elem] {
			return true
		}
	}

	// No common elements found
	return false
}

// HasSameElements checks if two string slices have the same elements, regardless of order
func HasSameElements(arr1 []string, arr2 []string) bool {
	if len(arr1) != len(arr2) {
		return false
	}

	// Create a map to count occurrences of each element in arr1
	countMap := make(map[string]int)

	for _, elem := range arr1 {
		countMap[elem]++
	}

	// Decrease the count for each element in arr2
	for _, elem := range arr2 {
		if count, exists := countMap[elem]; exists && count > 0 {
			countMap[elem]--
		} else {
			return false // Element not found or count is zero
		}
	}

	return true // All elements matched
}

// IsHostedOnAWS determines if the currently running instance is on AWS nor not
func IsHostedOnAWS() bool {
	if os.Getenv("LAMBDA_TASK_ROOT") != "" || os.Getenv("AWS_EXECUTION_ENV") != "" {
		return true
	}

	return false
}

// GenerateGUID generates a unique guid used as an id
func GenerateGUID() string {
	u := uuid.NewV4()
	return u.String()
}

func GenerateUniqueShortID() string {
	guid := uuid.NewV4()
	guidParts := strings.Split(guid.String(), "-")

	return strings.ToUpper(guidParts[4])
}

// GetValueFromStringMap retrieves the value from the map if the key exists. This does a
// case-insensitive key value search
func GetValueFromStringMap(key string, source map[string]string) string {
	lowerKey := strings.ToLower(key)

	for k, v := range source {
		if strings.ToLower(k) == lowerKey {
			return v
		}
	}

	return ""
}

// GetCurrentTimestamp gets the current time in the globally desired time format
func GetCurrentTimestamp() string {
	return time.Now().Format(time.RFC3339)
}

// CapitalizeFirstLetter capitalizes the first letter of the given the string
func CapitalizeFirstLetter(input string) string {
	if len(input) == 0 {
		return input
	}
	var result []string

	parts := strings.Split(input, " ")

	for _, part := range parts {
		firstCharUpper := strings.ToTitle(part[:1])
		result = append(result, firstCharUpper+part[1:])
	}

	return strings.Join(result, " ")
}

// StructToMap converts a struct to a map while maintaining the json alias as keys
func StructToMap(obj interface{}) (newMap map[string]interface{}, err error) {
	data, err := json.Marshal(obj) // Convert to a json string

	if err != nil {
		return
	}

	err = json.Unmarshal(data, &newMap) // Convert to a map
	return
}

func ObjectToJSON(object interface{}) string {
	jsonData, err := json.Marshal(object)
	if err != nil {
		log.Printf("failed to marshal object to JSON. Please verify object is valid. Object: %+v, Error: %+v", object, err)
		return ""
	}

	return string(jsonData)
}

func StringToTitle(input string) string {
	input = strings.ToLower(input)
	return cases.Title(language.Und, cases.NoLower).String(input)
}

// StringToArray converts a string into an array by splitting on multiple aspects, including
// comma, semicolon, spaces, carriage return and line feeds
func StringToArray(input string) []string {
	dirty := strings.ReplaceAll(input, ",", " ")
	dirty = strings.ReplaceAll(dirty, ";", " ")
	dirty = strings.ReplaceAll(dirty, "\r", " ")
	dirty = strings.ReplaceAll(dirty, "\n", " ")
	dirty = strings.ReplaceAll(dirty, "\\", " ")

	parts := strings.Split(dirty, " ")
	var result []string

	for _, part := range parts {
		cleanPart := strings.Trim(part, " ,;\r\n")
		if len(cleanPart) == 5 {
			result = append(result, cleanPart)
		}
	}

	return result
}

// GetStartOfNextMonth returns the start of the following month in UTCTimeFormat
func GetStartOfNextMonth() string {
	now := time.Now().UTC()
	nextMonth := now.AddDate(0, 1, 0)
	startOfNextMonth := time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, time.UTC)
	return startOfNextMonth.Format(time.RFC3339)
}

func SortStringIntMapByKey(inputMap map[string]int) map[string]int {
	// Create a slice of keys
	var keys []string
	for key := range inputMap {
		keys = append(keys, key)
	}

	// Sort the keys
	sort.Strings(keys)

	// Create a new map with sorted keys
	sortedMap := make(map[string]int)
	for _, key := range keys {
		sortedMap[key] = inputMap[key]
	}

	return sortedMap
}

func SortStringStringMapByKey(inputMap map[string]string) map[string]string {
	// Create a slice of keys
	var keys []string
	for key := range inputMap {
		keys = append(keys, key)
	}

	// Sort the keys
	sort.Strings(keys)

	// Create a new map with sorted keys
	sortedMap := make(map[string]string)
	for _, key := range keys {
		sortedMap[key] = inputMap[key]
	}

	return sortedMap
}

// RemoveSpacesAndPunctuation remove all spaces, numbers and punctuation from a string
func RemoveSpacesAndPunctuation(input string) string {
	// Define a function to check if a character is a punctuation mark or space
	isPunctOrSpace := func(r rune) bool {
		return unicode.IsPunct(r) || unicode.IsSpace(r)
	}

	// Replace all punctuation marks and spaces with an empty string
	cleanedString := strings.Map(func(r rune) rune {
		if isPunctOrSpace(r) {
			return -1 // Remove the character
		}
		return r // Keep non-punctuation and non-space characters
	}, input)

	return cleanedString
}

func EscapeString(input string) string {
	input = strings.ReplaceAll(input, "\r", " ")
	input = strings.ReplaceAll(input, "\n", " ")
	input = strings.ReplaceAll(input, "\"", "\\\"")
	return input
}

// RemoveNonNumeric removes all non-numeric characters from a string. This is useful for phone numbers.
func RemoveNonNumeric(phone string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) {
			return r
		}
		return -1
	}, phone)
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int) int {
	rnd := rand.Intn(max - min + 1)

	return min + rnd
}

// SplitFullName splits a full name into first and last names
func SplitFullName(fullName string) (string, string) {
	parts := strings.Fields(fullName)
	if len(parts) == 0 {
		return "", ""
	}

	firstName := parts[0]
	lastName := strings.Join(parts[1:], " ")

	// Capitalize first letter of each name part
	firstName = cases.Title(language.Und, cases.NoLower).String(firstName)
	lastName = cases.Title(language.Und, cases.NoLower).String(lastName)

	return firstName, lastName
}

func UnescapeString(input string) string {
	input = strings.ReplaceAll(input, "\\n", "\n")
	input = strings.ReplaceAll(input, "\\r", "\r")
	input = strings.ReplaceAll(input, "\\\"", "\"")
	return input
}

// SmartSplit splits a string by a separator, ignoring separators that are within quotes. The returned values do not include the quotes.
func SmartSplit(input string, separator rune) []string {
	var result []string
	var current strings.Builder
	inQuotes := false

	for _, char := range input {
		switch char {
		case '"':
			inQuotes = !inQuotes
			//current.WriteRune(char)
		case separator:
			if inQuotes {
				current.WriteRune(char)
			} else {
				result = append(result, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(char)
		}
	}
	result = append(result, current.String())
	return result
}

// stringToInt converts a string to an int, returning 0 if the conversion fails
func StringToInt(input string) int {
	value, err := strconv.Atoi(input)
	if err != nil {
		return 0
	}
	return value
}
