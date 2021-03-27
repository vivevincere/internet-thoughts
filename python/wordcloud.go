package wordcloud


import(
	"os/exec"
	"os"
	"io/ioutil"
	"fmt"
	"strings"
)




func wordCloud(words string) [6][]string{
	err := os.Remove("wordcloud.txt")
	if err != nil{
		fmt.Println("remove file failed")
	}
    err = ioutil.WriteFile("wordcloud.txt", []byte(words), 0644)
    if err != nil{
		fmt.Println("write to file failed")
	}
    cmd := exec.Command(".env/bin/python3","wordcloud.py")
    out, err := cmd.Output()
    if err != nil{
		fmt.Println(err)
	}
	var toRet [6][]string

	values := strings.Split(string(out), "\n")
    for i:= 0; i <6; i++{
    	wordCount := strings.Split(values[i]," ")
    	toRet[i] = wordCount
    }
    return toRet
}