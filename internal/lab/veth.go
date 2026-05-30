//go:build linux

package lab

import (
	"fmt"
	"os/exec"
	"strconv"

	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

func VethLink(args []string) error {
	if len(args) != 5 {
		return fmt.Errorf("usage: veth-link <idx> <ctr1> <ctr2> <if1> <if2>")
	}
	d := newDockerClient()
	if err := createVethPair(d, args[0], args[1], args[2], args[3], args[4]); err != nil {
		return fmt.Errorf("veth-link: %w", err)
	}
	fmt.Printf("linked %s:%s ↔ %s:%s\n", args[1], args[3], args[2], args[4])
	return nil
}

func createVethPair(d *dockerClient, idx, ctrA, ctrB, ifA, ifB string) error {
	pidA, err := d.containerPID(ctrA)
	if err != nil {
		return fmt.Errorf("pid %s: %w", ctrA, err)
	}
	pidB, err := d.containerPID(ctrB)
	if err != nil {
		return fmt.Errorf("pid %s: %w", ctrB, err)
	}

	tmpA := "vl-" + idx + "-a"
	tmpB := "vl-" + idx + "-b"

	if err := netlink.LinkAdd(&netlink.Veth{
		LinkAttrs: netlink.LinkAttrs{Name: tmpA},
		PeerName:  tmpB,
	}); err != nil {
		return fmt.Errorf("add: %w", err)
	}
	if err := moveToNS(tmpA, pidA, ifA); err != nil {
		return fmt.Errorf("%s: %w", ctrA, err)
	}
	if err := moveToNS(tmpB, pidB, ifB); err != nil {
		return fmt.Errorf("%s: %w", ctrB, err)
	}
	return nil
}

func moveToNS(name string, pid int, newName string) error {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return err
	}
	ns := fmt.Sprintf("/proc/%d/ns/net", pid)
	fd, err := unix.Open(ns, unix.O_RDONLY, 0)
	if err != nil {
		return fmt.Errorf("open %s: %w", ns, err)
	}
	defer unix.Close(fd)

	if err := netlink.LinkSetNsFd(link, fd); err != nil {
		return err
	}
	return nsDo(ns, "ip", "link", "set", name, "name", newName, "up")
}

func nsDo(ns string, arg ...string) error {
	args := append([]string{"--net=" + ns}, arg...)
	out, err := exec.Command("nsenter", args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("nsenter %v: %s", arg, out)
	}
	return nil
}

func vethRun(idx int, ctrA, ctrB, ifA, ifB string) error {
	return createVethPair(newDockerClient(), strconv.Itoa(idx), ctrA, ctrB, ifA, ifB)
}
