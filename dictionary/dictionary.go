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
	addCh    chan addMethod
	removeCh chan removeMethod
	doneCh   chan struct{}
}

type addMethod struct {
	word       string
	definition string
	resultCh   chan error
}

type removeMethod struct {
	word     string
	resultCh chan error
}

func New(filePath string) (*Dictionary, error) {
	d := &Dictionary{
		filePath: filePath,
		addCh:    make(chan addMethod),
		removeCh: make(chan removeMethod),
		doneCh:   make(chan struct{}),
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if err := os.WriteFile(filePath, []byte("{}"), 0644); err != nil {
			return nil, err
		}
	}

	go d.concurrentAddRemoveHandler()

	return d, nil
}

func (d *Dictionary) Add(word string, definition string) error {
	resultCh := make(chan error)
	d.addCh <- addMethod{word: word, definition: definition, resultCh: resultCh}
	return <-resultCh
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
	resultCh := make(chan error)
	d.removeCh <- removeMethod{word: word, resultCh: resultCh}
	return <-resultCh
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

func (d *Dictionary) concurrentAddRemoveHandler() {
	for {
		select {
		case addMtd := <-d.addCh:
			select {
			case <-d.doneCh:
				return
			default:
			}

			entries, err := d.readFromFile()
			if err != nil {
				addMtd.resultCh <- err
				continue
			}
			entries[addMtd.word] = Entry{Definition: addMtd.definition}
			err = d.writeToFile(entries)
			if err != nil {
				addMtd.resultCh <- err
				continue
			}
			addMtd.resultCh <- nil

		case removeMtd := <-d.removeCh:
			select {
			case <-d.doneCh:
				return
			default:
			}

			entries, err := d.readFromFile()
			if err != nil {
				removeMtd.resultCh <- err
				continue
			}
			delete(entries, removeMtd.word)
			err = d.writeToFile(entries)
			if err != nil {
				removeMtd.resultCh <- err
				continue
			}
			removeMtd.resultCh <- nil
		}
	}
}

func (d *Dictionary) Stop() {
	close(d.doneCh)
}
