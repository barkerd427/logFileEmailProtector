package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

var wg sync.WaitGroup

func removeEmail(fp string, fi os.FileInfo, err error) error {
	if err != nil {
		log.Println(err)
		return nil
	}
	if fi.IsDir() {
		return nil
	}
	matched, err := filepath.Match("*.log", fi.Name())
	if err != nil {
		log.Println(err)
		return err
	}
	if matched {
		log.Println(fp)
		b, err := ioutil.ReadFile(fp)
		if err != nil {
			panic(err)
		}

		r := regexp.MustCompile(`[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.?[a-zA-Z0-9-.]+`)

		wg.Add(1)
		go func() {
			defer wg.Done()
			ioutil.WriteFile(fp+"protected", r.ReplaceAll(b, []byte("*protected*")), 0664)
		}()
		//os.Stdout.Write(r.ReplaceAll(b, []byte("*protected*")))
	}
	return nil
}

func main() {
	filepath.Walk(os.Args[1], removeEmail)
	wg.Wait()
}
