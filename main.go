package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// go run main.go -chunk=true -file=bigfile.txt
// go run main.go -rebuild=true -file=bigfile_chunk_000.gfs -ext=.txt
// go run main.go -delete=true

func main() {
	// flags
	rebuild := flag.Bool("rebuild", false, "rebuild a file from .gfs chunks")
	input := flag.String("file", "", "file or directory to chunk or rebuild")
	chunk := flag.Bool("chunk", false, "split file(s) into 64MB binary chunks")
	delete := flag.Bool("delete", false, "delete all .gfs chunks")

	flag.Parse()

	if *delete {
		fmt.Println("Deleting all .gfs chunks...")
		deleteChunks()
		return
	}

	if *chunk {
		info, err := os.Stat(*input)
		check(err)

		if info.IsDir() {
			fmt.Println("Chunking all files in directory:", *input)
			chunkAllInDir(*input)
		} else {
			fmt.Println("Chunking:", *input)
			chunkFile(*input)
		}
		return
	}

	if *rebuild {
		if !strings.HasSuffix(*input, ".meta") {
			fmt.Println("Rebuild must be done from a .meta file")
			return
		}

		fmt.Println("Rebuilding from meta:", *input)
		rebuildFromMeta(*input)
		return
	}

	fmt.Println("No valid flag. Use -chunk, -rebuild, or -delete.")
}

func chunkFile(filename string) {
	// chunk size
	const chunkSize = 64 * 1024 * 1024 // 64Mo
	data, err := os.ReadFile(filename)
	check(err)

	base := strings.TrimSuffix(filename, getExt(filename))
	ext := getExt(filename)

	outDir := "chunks"
	check(os.MkdirAll(outDir, 0755))

	var chunkList []string

	for i := 0; i*chunkSize < len(data); i++ {
		start := i * chunkSize
		end := min(start+chunkSize, len(data))
		chunk := data[start:end]

		outDir := "chunks"
		check(os.MkdirAll(outDir, 0755))

		chunkName := fmt.Sprintf("%s/%s_chunk_%03d.gfs", outDir, filepath.Base(base), i)
		check(os.WriteFile(chunkName, chunk, 0644))
		chunkList = append(chunkList, chunkName)
		fmt.Println("New chunk :", chunkName)
	}

	// Metadata file
	// interface{} pour accepter n'importe quel type
	meta := map[string]interface{}{
		"original": filename,
		"ext":      ext,
		"chunks":   chunkList,
	}
	metaData, err := json.MarshalIndent(meta, "", "  ")
	check(err)

	metaName := fmt.Sprintf("%s/%s.meta", outDir, filepath.Base(base))
	check(os.WriteFile(metaName, metaData, 0644))
	fmt.Println("Master/metadata :", metaName)
}

func rebuildFromMeta(metaFile string) {
	// Lire les métadonnées
	data, err := os.ReadFile(metaFile)
	check(err)

	var meta struct {
		Original string   `json:"original"`
		Ext      string   `json:"ext"`
		Chunks   []string `json:"chunks"`
	}
	check(json.Unmarshal(data, &meta))

	// Créer le fichier de sortie
	base := strings.TrimSuffix(metaFile, ".meta")
	output := "rebuild_" + base + meta.Ext

	check(os.MkdirAll(filepath.Dir(output), 0755))

	outFile, err := os.Create(output)
	check(err)
	defer outFile.Close()

	for _, chunk := range meta.Chunks {
		chunkData, err := os.ReadFile(chunk)
		check(err)
		_, err = outFile.Write(chunkData)
		check(err)
		fmt.Println("Chunk :", chunk)
	}

	fmt.Println("Rebuild :", output)
}

func deleteChunks() {
	count := 0

	err := filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		check(err)

		if !d.IsDir() && (strings.HasSuffix(path, ".gfs") || strings.HasSuffix(path, ".meta")) {
			err := os.Remove(path)
			check(err)
			fmt.Println("Deleted:", path)
			count++
		}
		return nil
	})

	check(err)

	if count == 0 {
		fmt.Println("No chunks or metadata files found.")
	} else {
		fmt.Printf("Deleted %d \n", count)
	}
}

func getExt(name string) string {
	if i := strings.LastIndex(name, "."); i != -1 {
		return name[i:]
	}
	return ""
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func chunkAllInDir(dir string) {
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		check(err)

		if d.IsDir() {
			// continue for subdirectories
			return nil
		}

		chunkFile(path)
		return nil
	})

	check(err)
}
