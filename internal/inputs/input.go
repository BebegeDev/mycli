package inputs

import (
	"bufio"
	"os"
	"strings"
)

// Ввод от пользователя
func Input() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	return input
}
