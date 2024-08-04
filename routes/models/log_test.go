package models

import (
	"fmt"
	"strings"
	"testing"
)

func TestName(t *testing.T) {
	/**

	safsd
	*/

	str := "Hello, World!"
	replaced := strings.Replace(str, "World", "Gopher", -1)
	fmt.Println("Replaced string:", replaced)

	out := []byte(str)
	print(out)
	
}
