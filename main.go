package main

import (
	"bufio"
	"estiam/dictionary"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
)

func main() {
	dict, err := dictionary.New("dictionary/dictionary.json")
	if err != nil {
		fmt.Println("An error occured while initializing the dictionary : ", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Which action do you want to perform ? :")
		fmt.Println("1. Add a word")
		fmt.Println("2. Define a word")
		fmt.Println("3. Remove a word")
		fmt.Println("4. List all words")
		fmt.Println("5. Update a word")
		fmt.Println("6. Exit")

		fmt.Print("Enter the number corresponding to your choice : ")

		var choice int
		_, err := fmt.Scan(&choice)
		if err != nil {
			fmt.Println("Invalid input. Please enter a number.")
			continue
		}
		switch choice {
		case 1:
			actionAdd(dict, reader)
		case 2:
			actionDefine(dict, reader)
		case 3:
			actionRemove(dict, reader)
		case 4:
			actionList(dict)
		case 5:
			actionUpdate(dict, reader)
		case 6:
			fmt.Println("Exiting the program ...")
			return
		default:
			fmt.Println("Invalid command. Please try again.")
		}
	}
}

func actionAdd(d *dictionary.Dictionary, reader *bufio.Reader) {
	fmt.Print("Enter a word to add into the dictionary : ")
	word, _ := reader.ReadString('\n')
	word = word[:len(word)-1]

	fmt.Print("Enter a definition for '", word, "' : ")
	definition, _ := reader.ReadString('\n')
	definition = definition[:len(definition)-1]

	err := d.Add(word, definition)
	if err != nil {
		fmt.Println("Error while adding the word '", word, "' : ", err)
		return
	}
	fmt.Println("Word '", word, "' added successfully !")

}

func actionDefine(d *dictionary.Dictionary, reader *bufio.Reader) {
	fmt.Print("Enter the word to define : ")
	word, _ := reader.ReadString('\n')
	word = word[:len(word)-1]

	entry, err := d.Get(word)
	if err != nil {
		fmt.Println("Word '", word, "' not found in the dictionary.")
	} else {
		fmt.Printf("Definition : %s\n", entry)
	}
}

func actionRemove(d *dictionary.Dictionary, reader *bufio.Reader) {
	fmt.Print("Enter a word to remove : ")
	word, _ := reader.ReadString('\n')
	word = word[:len(word)-1]

	err := d.Remove(word)
	if err != nil {
		fmt.Println("Error while removing the word '", word, "' : ", err)
		return
	}
	fmt.Println("Word '", word, "' removed successfully !")
}

func actionUpdate(d *dictionary.Dictionary, reader *bufio.Reader) {
	fmt.Print("Enter the word to update : ")
	word, _ := reader.ReadString('\n')
	word = word[:len(word)-1]

	_, err := d.Get(word)
	if err != nil {
		fmt.Println("Word '", word, "' not found in the dictionary.")
		return
	}

	fmt.Print("Enter the new definition for '", word, "' : ")
	newDefinition, _ := reader.ReadString('\n')
	newDefinition = newDefinition[:len(newDefinition)-1]

	err = d.Add(word, newDefinition)
	if err != nil {
		fmt.Println("Error while updating word '", word, "' : ", err)
		return
	}
	fmt.Println("Word '", word, "' updated successfully!")
}

func actionList(d *dictionary.Dictionary) {
	words, entries, err := d.List()

	if err != nil {
		fmt.Println("Error while listing the dictionary : ", err)
		return
	}

	if len(words) == 0 {
		fmt.Println("The dictionary is empty.")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Word", "Definition"})

	for _, word := range words {
		table.Append([]string{word, entries[word].Definition})
	}

	fmt.Println("Dictionary content :")
	fmt.Println("\n")

	table.Render()
	fmt.Println("\n")
}
