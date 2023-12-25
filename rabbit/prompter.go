package rabbit

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func PromptForDialURL() string {
	fmt.Println("Enter queue url for example amqp://guest:guest@localhost:15672")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalln("Error occurred while reading input")
	}
	input = strings.TrimSuffix(input, "\n")

	fmt.Println(input)
	return input
}
