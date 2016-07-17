package ipset

import (
	"bufio"
	"errors"
	"github.com/kr/pty"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
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
func (ipset *IPset) Cmd(cmd string) (string, error) {
	if cmd == "\r\n" || cmd == "\n" || cmd == "" {
		return "", nil
	}
	cmd = strings.Replace(cmd, "\r\n", "\n", -1)
	cmd = strings.Replace(cmd, "\n\n", "\n", -1)
	ch := make(chan string)
	ipset.stdin.WriteString(cmd)
	ipset.stdin.Flush()
	go ipset.read(ch)
	select {
	case res := <-ch:
		{
			res = strings.Replace(res, "\r", "", -1)
			res = strings.TrimPrefix(res, cmd)
			res = strings.Replace(res, "ipset> ", "", -1)
			res = strings.TrimPrefix(res, "\n")
			return res, nil
		}
	case <-time.After(time.Second):
		return "", errors.New("timeout")
	}
}

func (ipset *IPset) read(ch chan string) {
	loadStr := ""
	re := regexp.MustCompile("ipset> ")
	buf := make([]byte, 10000)
	for {
		n, _ := ipset.stdout.Read(buf)
		loadStr += string(buf[:n])
		// check the 7 last bytes, should match "ipset> "
		if len(loadStr) > 7 {
			if re.MatchString(string(loadStr[len(loadStr)-7:])) {
				break
			}
		}
	}
	ch <- loadStr
}
