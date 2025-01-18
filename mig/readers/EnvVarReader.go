package readers

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type EnvVarReader struct {
	FilePath string
}

func NewEnvVarReader() *EnvVarReader {
	return &EnvVarReader{
		FilePath: os.Getenv("ENVVER_PATH"),
	}
}
func (srv *EnvVarReader) Read() {
	file, err := os.Open(srv.FilePath + "\\envvars.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	properties := make(map[string]string)

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		position := strings.Index(line, ":")
		value := line[position+1 : len(line)]
		name := line[0:position]
		//fmt.Println(position)
		//fmt.Println(fmt.Sprintf("%s:%s", name, value))
		count++
		properties[name] = strings.TrimSpace(value)

	}

	//fmt.Println(len(properties), count)
	/*var config Config
	if _, err := toml.DecodeFile(configfile, &config); err != nil {
		log.Fatal(err)
	}*/
	//log.Print(config.Index)

}

type Conf struct {
	Key string
	Val string
}
