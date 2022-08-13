package mini

import (
	"bufio"
	"fmt"
	"github.com/metafates/mangal/style"
	"os"
	"strconv"
	"strings"
)

type input struct {
	value string
}

func (o input) asInt() (n int64, ok bool) {
	n, err := strconv.ParseInt(o.value, 10, 16)
	ok = err == nil
	return n, ok
}

func getInput(validator func(string) bool) (*input, error) {
	fmt.Print(style.Magenta("> "))
	reader := bufio.NewReader(os.Stdin)
	in, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	in = strings.TrimSpace(in)

	if !validator(in) {
		fmt.Println(style.Red("Invalid choice entered"))
		return getInput(validator)
	}

	return &input{value: in}, nil
}
