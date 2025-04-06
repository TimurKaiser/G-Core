package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

type Block struct {
	ID           int
	PhysicalFile string
	Size         int
}

type File struct {
	Name   string
	Size   int64
	Blocks []*Block
}

func main() {
	// Ouvrir le dataset
	file, err := os.Open("data1Go.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// taille du fichier
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	fileSize := fileInfo.Size()

	// taille de chaque bloc, 128 Mo
	blockSize := int64(128 * 1024 * 1024)
	numBlocks := (fileSize + blockSize - 1) / blockSize

	// node de depart
	myFile := File{
		Name:   fileInfo.Name(),
		Size:   fileSize,
		Blocks: []*Block{},
	}

	// diviser les fichiers en blocs
	buffer := make([]byte, blockSize)

	// créer les fichiers de blocs
	for i := int64(0); i < numBlocks; i++ {
		blockFileName := fmt.Sprintf("block_%d.txt", i+1)
		blockFile, err := os.Create(blockFileName)
		if err != nil {
			log.Fatal(err)
		}
		defer blockFile.Close()

		// offset pour le bloc actuel
		// seek pour se déplacer à la position du bloc
		offset := i * blockSize
		_, err = file.Seek(offset, io.SeekStart)
		if err != nil {
			log.Fatal(err)
		}

		// lire le bloc
		// n est le nombre d'octets lus
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		// ecrire le bloc dans le fichier
		blockFile.Write(buffer[:n])

		// créer un nouveau bloc et l'ajouter à la liste des blocs
		// ID du bloc est i+1 (commence à 1)
		newBlock := &Block{
			ID:           int(i + 1),
			PhysicalFile: blockFileName,
			Size:         n,
		}
		myFile.Blocks = append(myFile.Blocks, newBlock)

		// fmt.Printf("Block %d saved to %s (%d bytes)\n", newBlock.ID, newBlock.PhysicalFile, newBlock.Size)
	}

	// afficher la chaîne de blocs
	fmt.Printf("\nFile: %s (size: %d bytes)\n", myFile.Name, myFile.Size)
	fmt.Println("Chain list:")
	for _, b := range myFile.Blocks {
		fmt.Printf("  File -> Block %d -> %s (%d bytes)\n", b.ID, b.PhysicalFile, b.Size)
	}

}
