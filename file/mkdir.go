package file

import (
	"os"
	"path/filepath"
)

func Mkdir(dirs []string) error {
	path := ""
	for _, v := range dirs {
		path = filepath.Join(path, v)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if err := os.Mkdir(path, 0755); err != nil {
				return err
			}
		}
	}
	return nil
}
