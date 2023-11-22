package dictionary

import (
	"encoding/json"
	"errors"
	"os"
)

type Entry struct {
	Definition string `json:"definition"`
}

func (e Entry) String() string {
	return e.Definition
}

type Dictionary struct {
	filePath string
}

func New(filePath string) (*Dictionary, error) {
	d := &Dictionary{filePath: filePath}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if err := os.WriteFile(filePath, []byte("{}"), 0644); err != nil {
			return nil, err
		}
	}
	return d, nil
}

func (d *Dictionary) Add(word string, definition string) error {
	entries, err := d.readFromFile()
	if err != nil {
		return err
	}

	if _, exists := entries[word]; exists {
		return errors.New("Word already exists in the dictionary")
	}

	entries[word] = Entry{Definition: definition}
	return d.writeToFile(entries)
}

func (d *Dictionary) Update(word string, newDefinition string) error {
	entries, err := d.readFromFile()
	if err != nil {
		return err
	}

	_, exists := entries[word]
	if !exists {
		return errors.New("Word not found in the dictionary")
	}

	entries[word] = Entry{Definition: newDefinition}
	return d.writeToFile(entries)
}

func (d *Dictionary) Get(word string) (Entry, error) {
	entries, err := d.readFromFile()
	if err != nil {
		return Entry{}, err
	}

	entry, isEntryExists := entries[word]
	if !isEntryExists {
		return Entry{}, errors.New("Sorry, the word you are looking for is not available in the dictionary")
	}
	return entry, nil
}

func (d *Dictionary) Remove(word string) error {
	entries, err := d.readFromFile()
	if err != nil {
		return err
	}

	_, exists := entries[word]
	if !exists {
		return errors.New("Word not found in the dictionary")
	}

	delete(entries, word)
	return d.writeToFile(entries)
}

func (d *Dictionary) List() ([]string, map[string]Entry, error) {
	entries, err := d.readFromFile()
	if err != nil {
		return nil, nil, err
	}

	words := make([]string, 0, len(entries))
	for word := range entries {
		words = append(words, word)
	}

	return words, entries, nil
}
func (d *Dictionary) readFromFile() (map[string]Entry, error) {
	data, err := os.ReadFile(d.filePath)
	if err != nil {
		return nil, err
	}

	var entries map[string]Entry
	err = json.Unmarshal(data, &entries)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func (d *Dictionary) writeToFile(entries map[string]Entry) error {
	data, err := json.Marshal(entries)
	if err != nil {
		return err
	}

	err = os.WriteFile(d.filePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
