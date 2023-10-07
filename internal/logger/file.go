package logger

import (
	"bufio"
	"fmt"
	"os"

	"github.com/AndrivA89/key-value-store/internal/entity"
)

type FileLogger struct {
	events       chan<- entity.Event
	errors       <-chan error
	lastSequence uint64
	file         *os.File
}

func (l *FileLogger) WritePut(key, value string) {
	l.events <- entity.Event{
		EventType: entity.EventPut,
		Key:       key,
		Value:     value,
	}
}

func (l *FileLogger) WriteDelete(key string) {
	l.events <- entity.Event{
		EventType: entity.EventDelete,
		Key:       key,
	}
}

func (l *FileLogger) Err() <-chan error {
	return l.errors
}

func NewFileLogger(fileName string) (Logger, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("cannot open transaction log file: %w", err)
	}

	return &FileLogger{file: file}, nil
}

func (l *FileLogger) Run() {
	events := make(chan entity.Event, 16)
	l.events = events

	errors := make(chan error, 1)
	l.errors = errors

	go func() {
		for e := range events {
			l.lastSequence++

			_, err := fmt.Fprintf(
				l.file,
				"%d\t%d\t%s\t%s\n",
				l.lastSequence,
				e.EventType,
				e.Key,
				e.Value,
			)

			if err != nil {
				errors <- err
				return
			}
		}
	}()
}

func (l *FileLogger) ReadEvents() (<-chan entity.Event, <-chan error) {
	scanner := bufio.NewScanner(l.file)
	outEvent := make(chan entity.Event)
	outError := make(chan error, 1)

	go func() {
		var e entity.Event

		defer close(outEvent)
		defer close(outError)

		for scanner.Scan() {
			line := scanner.Text()

			if _, err := fmt.Sscanf(line, "%d\t%d\t%s\t%s\n",
				&e.Sequence, &e.EventType, &e.Key, &e.Value); err != nil {
				outError <- fmt.Errorf("input parse error: %w", err)
				return
			}

			if l.lastSequence >= e.Sequence {
				outError <- fmt.Errorf("transaction numbers out of sequence")
				return
			}

			l.lastSequence = e.Sequence
			outEvent <- e
		}

		if err := scanner.Err(); err != nil {
			outError <- fmt.Errorf("transaction log read failure: %w", err)
			return
		}
	}()

	return outEvent, outError
}
