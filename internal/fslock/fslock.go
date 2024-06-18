package fslock

import (
	"errors"
	"os"
	"path/filepath"
	"syscall"
)

type Lock struct {
	filepath string
	filename string
	fd       int
}

func New(filepath, filename string) *Lock {
	if filepath == "" {
		// If filepath is empty, default to system's temporary directory
		filepath = os.TempDir()
	}

	return &Lock{
		filepath: filepath,
		filename: filename,
	}
}

func (l *Lock) Lock() error {
	if err := l.open(); err != nil {
		return err
	}
	err := syscall.Flock(l.fd, syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		syscall.Close(l.fd)
	}
	if err == syscall.EWOULDBLOCK {
		return errors.New("Application instance already running")
	}
	return err
}

func (l *Lock) open() error {
	fullPath := filepath.Join(l.filepath, l.filename)

	fd, err := syscall.Open(fullPath, syscall.O_CREAT|syscall.O_RDONLY, 0600)
	if err != nil {
		return err
	}
	l.fd = fd
	return nil
}

func (l *Lock) Unlock() error {
	if err := syscall.Close(l.fd); err != nil {
		return err
	}

	// Attempt to delete the lock file
	fullPath := filepath.Join(l.filepath, l.filename)
	if err := os.Remove(fullPath); err != nil {
		return err
	}

	return nil
}
