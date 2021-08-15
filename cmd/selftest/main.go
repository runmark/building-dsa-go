package selftest

import (
	"crypto/sha256"
	"fmt"
	"io"
	"strings"
)

func countLetters(r io.Reader) (res map[string]int, err error) {

	buf := make([]byte, 2048)
	out := map[string]int{}
	for {
		n, err := r.Read(buf)
		if err == io.EOF {
			return out, nil
		}
		if err != nil {
			return nil, err
		}

		for _, b := range buf[:n] {
			if (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') {
				out[string(b)]++
			}
		}

	}

}

type Currency int

const (
	RMB Currency = iota
	USD
	EUR
	GBP
)

func countDiffLetters(a, b [32]byte) int {

	total := 0
	for i := range a {
		if a[i] != b[i] {
			total += 1
		}
	}
	return total
}

type LogicProvider struct{}

func (lp LogicProvider) Process(data string) string {
	return "business logic"
}

type Logic interface {
	Process(data string) string
}

type Client struct {
	L Logic
}

func (c Client) Program() {
	data := "get data from somewhere"
	c.L.Process(data)
}

func main() {

	cli := Client{
		L: LogicProvider{},
	}

	cli.Program()

	m, _ := countLetters(strings.NewReader("hello"))
	fmt.Println(m)

	q := [...]int{1, 2, 3}
	fmt.Printf("%T\n", q)

	c := [...]string{USD: "$", RMB: "¥", GBP: "£", EUR: "€"}
	fmt.Printf("%T %v\n", c, c)
	fmt.Println(RMB, c[RMB])

	//a := [2]int{1, 2}
	//b := [...]int{1,2}
	//e := []int{1,2}

	sumx := sha256.Sum256([]byte("x"))
	sumX := sha256.Sum256([]byte("X"))

	fmt.Printf("%v\n", countDiffLetters(sumx, sumX))
	fmt.Printf("%v\n", countDiffLetters(sumx, sumx))

	s := make([]int, 0, 5)

	for i := 1; i < 4; i++ {
		s = append(s, i)
	}
	fmt.Println(s)

	var ma map[string]int
	ma = make(map[string]int)

	ma["hello"] = 3
	fmt.Println(ma)
}

