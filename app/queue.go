package app

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"time"
)

var c chunk

func Start(storePath string) {
	files, err := ioutil.ReadDir(storePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(files), func(i, j int) { files[i], files[j] = files[j], files[i] })
		for _, file := range files {
			fName := file.Name()
			fmt.Println(fName)
			if filepath.Ext(fName) != ".mp3" {
				continue
			}
			f, err := os.Open(path.Join(storePath, fName))
			defer f.Close()
			if err != nil {
				fmt.Println(err)
				return
			}

			go c.Load(f)
			<-done
		}
	}
}
