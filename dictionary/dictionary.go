package dictionary

import (
	"errors"
)

type Entry struct {
	Definition string
}

func (e Entry) String() string {
	return e.Definition
}

type Dictionary struct {
	entries map[string]Entry
}

func New() *Dictionary {
	return &Dictionary{
		entries: map[string]Entry{},
	}
}

func (d *Dictionary) Add(word string, definition string) {
	entry := Entry{Definition: definition}
	d.entries[word] = entry
}

func (d *Dictionary) Get(word string) (Entry, error) {
	entry, isEntryExists := d.entries[word]
	if !isEntryExists {
		return Entry{}, errors.New("Sorry the word you are looking for is not available in the dictionary")
	}
	return entry, nil
}

func (d *Dictionary) Remove(word string) {
	delete(d.entries, word)
}

func (d *Dictionary) List() ([]string, map[string]Entry) {
	words := make([]string, 0, len(d.entries))

	for word := range d.entries {
		words = append(words, word)
	}

	return words, d.entries
}
