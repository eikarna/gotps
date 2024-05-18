package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func parseItemData(data string) map[string]string {
	parsedData := make(map[string]string)

	// Split string by "\"
	elements := strings.Split(data, "\\")

	// Define field names
	fields := []string{
		"add_item", "item_id", "editable_type", "item_category", "action_type", "hit_sound_type",
		"name", "texture", "texture_hash", "item_kind", "val1", "texture_x", "texture_y",
		"spread_type", "is_stripey_wallpaper", "collision_type", "break_hits", "drop_chance",
		"clothing_type", "rarity", "max_amount", "extra_file", "extra_file_hash", "audio_volume",
		"pet_name", "pet_prefix", "pet_suffix", "pet_ability", "seed_base", "seed_overlay",
		"tree_base", "tree_leaves", "seed_color", "seed_overlay_color", "grow_time", "val2",
		"is_rayman", "extra_options", "texture2", "extra_options2", "data_position_80",
		"punch_options", "data_version_12", "int_version_13", "int_version_14", "data_version_15",
		"str_version_15", "str_version_16", "int_version_17",
	}

	// Assign values to parsedData
	for i, field := range fields {
		if i < len(elements) {
			parsedData[field] = elements[i]
		}
	}

	return parsedData
}

func cleanLine(line string) string {
	// Trim whitespace and non-printable characters from the end of the line
	return strings.TrimRightFunc(line, unicode.IsSpace)
}

func printItemData(data map[string]string) {
	fmt.Println("Item Details:")
	for key, value := range data {
		fmt.Printf("%s: %s\n", key, value)
	}
	fmt.Println("====================================")
}

func main() {
	// Check if search string is provided as command line argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <search_string>")
		return
	}

	// Extract search string from command line argument
	searchString := strings.Join(os.Args[1:], " ")

	// Open the file
	file, err := os.Open("items.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Flag to indicate whether any item matching the search string is found
	found := false

	// Read each line and parse it
	for scanner.Scan() {
		line := scanner.Text()
		// Ignore lines starting with "//" (comments)
		if !strings.HasPrefix(line, "//") {
			line = cleanLine(line) // Clean up the line
			parsedData := parseItemData(line)
			// Check if any field contains the search string
			for _, value := range parsedData {
				if strings.Contains(value, searchString) {
					printItemData(parsedData)
					found = true
					break
				}
			}
		}
	}

	// Check if any item matching the search string is found
	if !found {
		fmt.Println("No item found matching the search string:", searchString)
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning file:", err)
		return
	}
}
