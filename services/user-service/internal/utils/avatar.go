package utils

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
)

// SelectRandomAvatar - выбирает случайную аватарку из набора
func SelectRandomAvatar() (string, error) {
	avatarDir := "assets/avatars"

	files, err := os.ReadDir(avatarDir)
	if err != nil {
		return "", err
	}

	var avatarFiles []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".png" {
			avatarFiles = append(avatarFiles, file.Name())
		}
	}

	if len(avatarFiles) == 0 {
		return "", errors.New("no avatar files found")
	}

	randomIndex := rand.Intn(len(avatarFiles))

	return avatarFiles[randomIndex], nil
}

// GetAvatarPath - возвращает путь к аватарке
func GetAvatarPath(avatarFileName string) string {
	return fmt.Sprintf("assets/avatars/%s", avatarFileName)
}
