// Package ipset provides interaction with ipset (ipset.netfilter.org)
//
// It runs ipset in interactive mode (ipset -) instead of running the binary on every call.
//
// Example of usage:
//
//		package main
//		import (
//			"github.com/42wim/ipsetd/ipset"
//			"fmt"
//		)
//
//		func main() {
//	 		ipset := NewIPset("/usr/sbin/ipset")
//	 		fmt.Print(ipset.Cmd("version"))
//	 		fmt.Print(ipset.Cmd("create abc hash:ip"))
//	 		fmt.Print(ipset.Cmd("add abc 1.2.3.4"))
//  	}
//
package ipset
