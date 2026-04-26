package replacer

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"
)

var src string = `
        Date: %date%
        Year: %date.Y%
        Month: %date.M%
        Day: %date.D%
        Time: %time%
        Hour: %time.H%
        Minute: %time.M%
        Second: %time.S%
        Ms: %time.Z%
        Ns: %time.ns%

        %undefined%

        %hello%, %name%
        %greet Michael%
        %list one two three four%

        %undeflist one two ...%
    `

var src1 string = `
	Date: $date$
	Day: $date.D$
	S.Ms: $time.S$.$time.Z$

	Line: $line$

	Code: $code$
`

func hello() string {
	return "Hello"
}

func greet(name string) string {
	return "Hello " + name
}

func list(vals ...string) string {
	return "Some: " + strings.Join(vals, ", ")
}

func TestNewReplacer(t *testing.T) {
	repl := New()
	repl.Add("name", "John")
	repl.Add("hello", hello)
	repl.Add("greet", greet)
	repl.Add("list", list)

	{
		fmt.Println("\n\nReplace:")
		dst := repl.Replace(src)
		fmt.Println(dst)
	}

	{
		fmt.Println("\n\nReplaceC")
		dst := repl.ReplaceC(src)
		fmt.Println(dst)
	}

	{
		fmt.Println("\n\nReplaceE")
		dst, e := repl.ReplaceE(src)
		if e != nil {
			fmt.Println("Errors: ")
			for _, er := range e {
				fmt.Println("\t", er.Error())
			}
		}
		fmt.Println(dst)
	}

	{
		fmt.Println("\n\nReplaceCE")
		dst, e := repl.ReplaceCE(src)
		if e != nil {
			fmt.Println("Errors: ")
			for _, er := range e {
				fmt.Println("\t", er.Error())
			}
		}
		fmt.Println(dst)
	}

}

func TestReplace2(t *testing.T) {
	r := New()
	rf := NewFix(time.Now())
	e := r.Add("line", "123")
	if e != nil {
		fmt.Println("\t", e.Error())
	}
	r.Add("code", "NewCode")
	rf.Add("code", "NewCode")
	r.MChar = '$'
	rf.MChar = '$'

	dst := r.Replace(src1)
	dstf := rf.Replace(src1)
	fmt.Println(dstf)
	fmt.Println(dst)
	fmt.Printf("T: %s\n", time.Now().Format("15:04:05.000"))
	wms := rand.Intn(2000)
	fmt.Printf("Sleep: %d ms\n", wms)
	time.Sleep(time.Duration(wms) * time.Millisecond)
	dst = r.Replace(src1)
	dstf = rf.Replace(src1)
	fmt.Println(dstf)
	fmt.Println(dst)
	fmt.Printf("T: %s\n", time.Now().Format("15:04:05.000"))
}
