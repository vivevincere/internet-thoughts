package wordcloud

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type Word_Cloud struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
}

func WordCloud(words string) [30][]string {
	err := os.Remove("python/wordcloud.txt")
	if err != nil {
		fmt.Println("remove file failed")
	}
	err = ioutil.WriteFile("python/wordcloud.txt", []byte(words), 0644)
	if err != nil {
		fmt.Println("write to file failed")
	}
	cmd := exec.Command("python/.env/bin/python3", "python/wordcloud.py")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	var toRet [30][]string
	values := strings.Split(string(out), "\n")
	for i := 0; i < 30; i++ {
		wordCount := strings.Split(values[i], " ")
		toRet[i] = wordCount
	}
	return toRet
}
