package ipset

import (
	"bufio"
	"github.com/kr/pty"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// IPset struct
type IPset struct {
	stdin  *bufio.Writer
	stdout *bufio.Reader
	pty    *os.File
}

// NewIPset starts ipset specified with path in interactive mode (ipset - ) and returns a new IPset.
func NewIPset(path string) *IPset {
	cmd := exec.Command(path, "-")
	f, _ := pty.Start(cmd)
	ipset := &IPset{pty: f, stdin: bufio.NewWriter(f), stdout: bufio.NewReader(f)}
	buf := make([]byte, 1000)
	ipset.stdout.Read(buf)
	return ipset
}

func NewIPsetExtra(path string, args ...string) *IPset {
	args = append(args, "-")
	cmd := exec.Command(path, args...)
	f, _ := pty.Start(cmd)
	ipset := &IPset{pty: f, stdin: bufio.NewWriter(f), stdout: bufio.NewReader(f)}
	buf := make([]byte, 1000)
	ipset.stdout.Read(buf)
	return ipset
}

// Cmd executes the ipset command and returns the output.
func (ipset *IPset) Cmd(cmd string) string {
	ipset.stdin.WriteString(cmd + "\n")
	ipset.stdin.Flush()
	res := ipset.read()
	res = strings.TrimPrefix(res, cmd+"\r\n")
	res = strings.Replace(res, "ipset> ", "", -1)
	res = strings.Replace(res, "\r", "", -1)
	res = strings.TrimPrefix(res, "\n")
	return (res)
}

func (ipset *IPset) read() string {
	loadStr := ""
	re := regexp.MustCompile("ipset> ")
	buf := make([]byte, 10000)
	for {
		n, _ := ipset.stdout.Read(buf)
		loadStr += string(buf[:n])
		if re.MatchString(string(buf[:n])) {
			break
		}
	}
	return loadStr
}
