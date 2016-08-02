package errors

import (
	"fmt"
	"runtime"
	"strings"
)

// Frame represents a program counter inside a stack frame.
type Frame uintptr

// pc returns the program counter for this frame;
// multiple frames may have the same PC value.
func (f Frame) pc() uintptr {
	return uintptr(f) - 1
}

func (f Frame) Func() *runtime.Func {
	return runtime.FuncForPC(f.pc())
}

func (f Frame) FileLine() (string, int) {
	fn := f.Func()
	if fn == nil {
		return "<unknown>", 0
	}
	return fn.FileLine(f.pc())
}

func (f Frame) Name() string {
	fn := f.Func()
	if fn == nil {
		return "<unknown>"
	}
	return fn.Name()
}

// Format formats the frame according to the fmt.Formatter interface.
//    %v    github.com/subchen/gstack/errors/stack.go:83
//    %+v   errors.Callers() (github.com/subchen/gstack/errors/stack.go:83)
func (f Frame) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		fn := f.Func()
		if fn == nil {
			fmt.Fprint(s, "<unknown source>")
			return
		}

		file, line := fn.FileLine(f.pc())
		file, name := trimGOPATH(file, fn.Name())
		if s.Flag('+') {
			fmt.Fprintf(s, "%s() (%s:%d)", name, file, line)
		} else {
			fmt.Fprintf(s, "%s:%d", file, line)
		}
	}
}

// Stack represents a stack of program counters.
type Stack []Frame

// Format formats the stack according to the fmt.Formatter interface.
//    %v    at github.com/subchen/gstack/errors/stack.go:83
//          at example/main.go:9
//          at example/main.go:16
//    %+v   at errors.Callers() (github.com/subchen/gstack/errors/stack.go:83)
//          at main.test() (example/main.go:9)
//          at main.main() (example/main.go:16)
func (st Stack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			for _, f := range st {
				fmt.Fprintf(s, "\tat %+v\n", f)
			}
		} else {
			for _, f := range st {
				fmt.Fprintf(s, "\tat %v\n", f)
			}
		}
	}
}

// Callers returns a stack for caller
func Callers(skip int) Stack {
	pcs := make([]uintptr, 32) // max depth is 32
	n := runtime.Callers(skip, pcs)

	stack := make([]Frame, n)
	for i := 0; i < n; i++ {
		stack[i] = Frame(pcs[i])
	}
	return stack
}

func trimGOPATH(file, name string) (string, string) {
	// Here we want to get the source file path relative to the compile time
	// GOPATH. As of Go 1.6.x there is no direct way to know the compiled
	// GOPATH at runtime, but we can infer the number of path segments in the
	// GOPATH. We note that fn.Name() returns the function name qualified by
	// the import path, which does not include the GOPATH. Thus we can trim
	// segments from the beginning of the file path until the number of path
	// separators remaining is one more than the number of path separators in
	// the function name. For example, given:
	//
	//    GOPATH     /home/user
	//    file       /home/user/src/pkg/sub/file.go
	//    fn.Name()  pkg/module.Type.Method
	//
	// We want to produce:
	//
	//    file       pkg/sub/file.go
	//    fn.name    module.Type.Method
	//
	// From this we can easily see that fn.Name() has one less path separator
	// than our desired output. We count separators from the end of the file
	// path until it finds two more than in the function name and then move
	// one character forward to preserve the initial path segment without a
	// leading separator.
	const sep = "/"
	goal := strings.Count(name, sep) + 2
	i := len(file)
	for n := 0; n < goal; n++ {
		i = strings.LastIndex(file[:i], sep)
		if i == -1 {
			// not enough separators found, set i so that the slice expression
			// below leaves file unmodified
			i = -len(sep)
			break
		}
	}
	// get back to 0 or trim the leading separator
	file = file[i+len(sep):]
	if pos := strings.LastIndex(name, "/"); pos != -1 {
		name = name[pos+1 : len(name)]
	}
	return file, name
}
