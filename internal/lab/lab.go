package lab

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/LinkStateDev/linkstate-cli/internal/ui"
)

func Up() error {
	d := newDockerClient()

	allUp := true
	for _, name := range Nodes {
		if !d.containerRunning(name) {
			allUp = false
			break
		}
	}
	if allUp {
		fmt.Println("── Lab already up")
		Status()
		printHints()
		return nil
	}

	fmt.Println("── SSH key")
	if err := ensureSSHKey(); err != nil {
		return fmt.Errorf("key: %w", err)
	}
	fmt.Println("   ", sshPubPath())
	pub, err := readPubKey()
	if err != nil {
		return fmt.Errorf("read key: %w", err)
	}

	fmt.Println("── FRR image")
	if err := ensureFRRImage(); err != nil {
		return fmt.Errorf("image: %w", err)
	}

	fmt.Println("── Creating containers")
	for _, name := range Nodes {
		fmt.Printf("   %s (:%s)\n", name, nodeSSHPort(name))
		if err := d.createContainer(name, name, nodeSSHPort(name), pub); err != nil {
			return fmt.Errorf("container %s: %w", name, err)
		}
	}

	parallel(Nodes, func(name string) {
		d.waitContainer(name, 30*time.Second)
	})
	time.Sleep(2 * time.Second)

	fmt.Println("── Veth pairs (8)")
	var vwg sync.WaitGroup
	for i, l := range topoLinks {
		vwg.Add(1)
		go func(idx int, l linkDef) {
			defer vwg.Done()
			if err := vethRun(idx, l.a, l.b, l.aName, l.bName); err != nil {
				fmt.Fprintf(os.Stderr, "link %s:%s↔%s:%s: %v\n", l.a, l.aName, l.b, l.bName, err)
			}
		}(i, l)
	}
	vwg.Wait()

	fmt.Println("── Assigning IPs")
	parallel(topoLinks, func(l linkDef) {
		d.execOk(l.a, "ip", "addr", "add", l.aIP+"/30", "dev", l.aName)
		d.execOk(l.b, "ip", "addr", "add", l.bIP+"/30", "dev", l.bName)
	})
	parallel(Nodes, func(name string) {
		lo := nodeLo(name)
		d.execOk(name, "ip", "addr", "add", lo+"/32", "dev", "lo")
		d.execOk(name, "sysctl", "-w", "net.ipv4.ip_forward=1")
	})

	fmt.Println("── eBGP config")
	parallel(Nodes, func(name string) {
		asn := nodeASN(name)
		lo := nodeLo(name)

		d.execOk(name, "vtysh", "-c", "conf t",
			"-c", "route-map PERMIT-ALL permit 1",
			"-c", "exit")

		for _, l := range peerLinks(name) {
			peerIP := peerAddr(l, name)
			d.execOk(name, "vtysh", "-c", "conf t",
				"-c", fmt.Sprintf("router bgp %d", asn),
				"-c", fmt.Sprintf("neighbor %s remote-as %d", peerIP, peerAS(l, name)),
				"-c", fmt.Sprintf("neighbor %s timers 3 10", peerIP),
				"-c", fmt.Sprintf("neighbor %s route-map PERMIT-ALL in", peerIP),
				"-c", fmt.Sprintf("neighbor %s route-map PERMIT-ALL out", peerIP),
			)
		}
		d.execOk(name, "vtysh", "-c", "conf t",
			"-c", fmt.Sprintf("router bgp %d", asn),
			"-c", fmt.Sprintf("network %s/32", lo))
	})

	fmt.Println("── BGP convergence")
	waitBGP(d)

	fmt.Println("\n── Lab ready ──")
	Status()
	printHints()
	return nil
}

func Down() {
	d := newDockerClient()
	fmt.Println("── Tearing down")
	parallel(Nodes, func(name string) {
		if d.containerExists(name) {
			d.rmContainer(name)
		}
	})
	fmt.Println("Done.")
}

func Status() {
	d := newDockerClient()
	rows, err := d.listLabContainers()
	if err != nil {
		fmt.Fprintf(os.Stderr, "list: %v\n", err)
		return
	}
	if len(rows) == 0 {
		fmt.Println("no lab containers")
		return
	}

	fmt.Println()
	for _, c := range rows {
		name := strings.TrimPrefix(c.Names[0], "/")
		state := "DOWN"
		stateStyle := ui.Error
		if c.State == "running" {
			state = "UP"
			stateStyle = ui.Success
		}

		lo := nodeLo(name)
		port := nodeSSHPort(name)
		asn := nodeASN(name)

		var peers []string
		for _, l := range peerLinks(name) {
			if l.a == name {
				peers = append(peers, l.b)
			} else {
				peers = append(peers, l.a)
			}
		}

		fmt.Printf("   %s  %s  lo %-12s  ssh :%-4s  AS %-5d  ↔ %s\n",
			ui.Bold.Render(name),
			stateStyle.Render(state),
			lo,
			port,
			asn,
			strings.Join(peers, " "),
		)
	}
	fmt.Println()
}

func printHints() {
	fmt.Printf("  %-12s %s  — vtysh console\n", "connect", ui.Code.Render("lst lab connect <name>"))
	fmt.Printf("  %-12s %s  — SSH into device\n", "ssh", ui.Code.Render("lst lab ssh <name>"))
	fmt.Println()
}

func Connect(name string) error {
	if !newDockerClient().containerRunning(name) {
		return fmt.Errorf("%s is not running", name)
	}
	cmd := exec.Command("docker", "exec", "-it", name, "vtysh")
	cmd.Env = append(os.Environ(), "DOCKER_CLI_HINTS=off")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func SSH(name string) error {
	d := newDockerClient()
	if !d.containerRunning(name) {
		return fmt.Errorf("%s is not running", name)
	}
	cmd := exec.Command("ssh",
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"-i", sshKeyPath(),
		"-p", nodeSSHPort(name),
		"root@127.0.0.1",
	)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func waitBGP(d *dockerClient) {
	expected := map[string]int{}
	for _, name := range Nodes {
		if strings.HasPrefix(name, "spine") {
			expected[name] = 4
		} else {
			expected[name] = 2
		}
	}

	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		ok := true
		for _, name := range Nodes {
			out, err := d.exec(name, "vtysh", "-c", "show bgp summary")
			if err != nil {
				time.Sleep(time.Second)
				ok = false
				break
			}
			if strings.Count(out, "10.0.") < expected[name] {
				ok = false
				break
			}
		}
		if ok {
			return
		}
		time.Sleep(time.Second)
	}
}

func parallel[T any](items []T, fn func(T)) {
	var wg sync.WaitGroup
	for _, item := range items {
		wg.Add(1)
		go func(v T) { defer wg.Done(); fn(v) }(item)
	}
	wg.Wait()
}
