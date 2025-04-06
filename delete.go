// delete.go
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func removeBlockFiles() error {
	// filepath ou regexp, au choix
	files, err := filepath.Glob("block_*.txt")
	if err != nil {
		return fmt.Errorf("Error on filde search: %v", err)
	}

	// suprrime ficheir trouvés
	// %v pour afficher l'erreur, %s pour afficher le nom du fichier
	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			fmt.Printf("Error on file delete %s: %v\n", file, err)
		} else {
			fmt.Printf("File Delete: %s\n", file)
		}
	}
	return nil
}

func main() {
	err := removeBlockFiles()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Blocks are deleted successfully.")
	}
}
