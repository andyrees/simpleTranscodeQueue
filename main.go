package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

func listDir(dirPath string) ([]string, error) {
	rawFiles, err := ioutil.ReadDir(dirPath)
	var dirFiles []string

	if err != nil {
		return dirFiles, err
	}

	for _, f := range rawFiles {
		switch {
		case strings.HasSuffix(f.Name(), "mp4"):
			dirFiles = append(dirFiles, f.Name())
		case strings.HasSuffix(f.Name(), "mov"):
			dirFiles = append(dirFiles, f.Name())
		case strings.HasSuffix(f.Name(), "mp3"):
			dirFiles = append(dirFiles, f.Name())
		case strings.HasSuffix(f.Name(), "mpg"):
			dirFiles = append(dirFiles, f.Name())
		case strings.HasSuffix(f.Name(), "avi"):
			dirFiles = append(dirFiles, f.Name())
		case strings.HasSuffix(f.Name(), "m4v"):
			dirFiles = append(dirFiles, f.Name())
		case strings.HasSuffix(f.Name(), "flv"):
			dirFiles = append(dirFiles, f.Name())
		case strings.HasSuffix(f.Name(), "wav"):
			dirFiles = append(dirFiles, f.Name())
		}
	}
	return dirFiles, nil
}

func transcodeFile(filepath string) error {
	fname := path.Base(filepath)
	fdir := path.Dir(filepath)

	outfilename := path.Join(fdir, fmt.Sprintf("%s.ts", strings.Split(fname, ".")[0]))

	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		log.Printf("ffmpeg Error: %v", err.Error())
	}

	cmd := exec.Command(
		ffmpegPath,
		"-threads",
		"auto",
		"-i",
		filepath,
		"-q:v",
		"1",
		"-q:a",
		"1",
		"-f",
		"mpegts",
		outfilename,
	)
	err2 := cmd.Start()
	if err2 != nil {
		return err2
	}
	err2 = cmd.Wait()
	return err2
}

func worker(id int, jobs <-chan string, results chan<- string) {
	for j := range jobs {
		log.Println("worker", id, "processing job", path.Base(j))
		time.Sleep(time.Second)
		err := transcodeFile(j)
		if err != nil {
			results <- fmt.Sprintf("%s is in Error", path.Base(j))
		}
		results <- fmt.Sprintf("%s is Done", path.Base(j))
	}
}

func fcheck(filepath string) (string, error) {
	_, err := os.Stat(filepath)
	if err != nil {
		err := errors.New("NO SUCH FILE")
		return "NO SUCH FILE", err
	}
	return "FILE EXISTS", nil
}

var (
	sourceDirectory = flag.String("sourceDirectory", os.Getenv("HOME"), "source dircetory of media files to transcode")
)

func main() {
	flag.Parse()

	jobs := make(chan string, 100)
	results := make(chan string, 100)

	// Start 3 workers
	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	dirFilez, err := listDir(*sourceDirectory)
	if err != nil {
		log.Fatalln(err)
	}

	for _, j := range dirFilez {
		sourceFile := path.Join(*sourceDirectory, j)
		_, err := fcheck(sourceFile)
		if err != nil {
			fmt.Println("File does not exist: ", path.Base(sourceFile))
		}
		log.Println("sending: ", path.Base(sourceFile))
		jobs <- sourceFile
	}
	close(jobs)

	for a := 1; a <= len(dirFilez); a++ {
		log.Printf("%s\n", <-results)
	}
	log.Println("All DONE")
}
