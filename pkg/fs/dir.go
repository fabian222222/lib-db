package fs

import (
	"os"
)

func DoesDirExist(name string) bool {
	if name == "" {
		return false
	}
	info, err := os.Stat("./../../databases/" + name)

	if err != nil {
		return false 
	}
	return info.IsDir()
}