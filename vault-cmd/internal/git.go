package internal

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

func GetFileContentFromGit(r string, f string) []byte {
	fs := memfs.New()
	_, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		URL: r,
	})

	if err != nil {
		fmt.Printf("Can't fetch k8s information from git, please recheck your connection \n")
		fmt.Printf("%s \n", err.Error())
		os.Exit(2)
	}

	file, _ := fs.Open(f)
	rd, _ := ioutil.ReadAll(file)
	return rd
}
