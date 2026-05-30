package lab

import "strings"

const (
	imageFRR   = "lab-frr"
	spineASN   = 65000
	leafBaseAS = 65101
)

var Nodes = []string{"spine1", "spine2", "leaf1", "leaf2", "leaf3", "leaf4"}

func nodeASN(name string) int {
	if strings.HasPrefix(name, "spine") {
		return spineASN
	}
	return leafBaseAS + int(name[4]-'1')
}

func nodeLo(name string) string {
	if strings.HasPrefix(name, "spine") {
		return "10.255.0." + itoa(int(name[5]-'0'))
	}
	return "10.255.1." + itoa(int(name[4]-'0'))
}

func nodeSSHPort(name string) string {
	switch name {
	case "spine1":
		return "2201"
	case "spine2":
		return "2202"
	case "leaf1":
		return "2211"
	case "leaf2":
		return "2212"
	case "leaf3":
		return "2213"
	case "leaf4":
		return "2214"
	}
	return ""
}

type linkDef struct {
	a, b         string
	aEth, bEth   int
	aIP, bIP     string
	aName, bName string
}

var topoLinks = []linkDef{
	{"spine1", "leaf1", 1, 1, "10.0.1.1", "10.0.1.2", "", ""},
	{"spine1", "leaf2", 2, 1, "10.0.1.5", "10.0.1.6", "", ""},
	{"spine1", "leaf3", 3, 1, "10.0.1.9", "10.0.1.10", "", ""},
	{"spine1", "leaf4", 4, 1, "10.0.1.13", "10.0.1.14", "", ""},
	{"spine2", "leaf1", 1, 2, "10.0.2.1", "10.0.2.2", "", ""},
	{"spine2", "leaf2", 2, 2, "10.0.2.5", "10.0.2.6", "", ""},
	{"spine2", "leaf3", 3, 2, "10.0.2.9", "10.0.2.10", "", ""},
	{"spine2", "leaf4", 4, 2, "10.0.2.13", "10.0.2.14", "", ""},
}

func init() {
	for i := range topoLinks {
		topoLinks[i].aName = "eth" + itoa(topoLinks[i].aEth)
		topoLinks[i].bName = "eth" + itoa(topoLinks[i].bEth)
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}

func peerLinks(name string) []linkDef {
	var links []linkDef
	for _, l := range topoLinks {
		if l.a == name || l.b == name {
			links = append(links, l)
		}
	}
	return links
}

func peerAddr(link linkDef, myName string) string {
	if link.a == myName {
		return link.bIP
	}
	return link.aIP
}

func peerAS(link linkDef, myName string) int {
	if link.a == myName {
		return nodeASN(link.b)
	}
	return nodeASN(link.a)
}
