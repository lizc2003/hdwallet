package eth

import (
	"fmt"
	"github.com/lizc2003/hdwallet/eth"
	"github.com/lizc2003/hdwallet/qa"
	"os"
	"os/exec"
	"time"
)

// ganache-cli --account="0xc356cfe48d1ddcd2320b62553fe8739978d0478e2e940d6c30190e7637f51c76,200"

const ganacheCliPort = 8545

func RunGanache(accountPrivateKeys ...string) (*eth.EthClient, func(), error) {
	if qa.IsProcessRunning(ganacheCliPort, "ganache-cli") {
		return nil, nil, fmt.Errorf("ganache already running on port %d", ganacheCliPort)
	}

	args := []string{}
	if len(accountPrivateKeys) > 0 {
		for _, a := range accountPrivateKeys {
			args = append(args, fmt.Sprintf("--account=0x%s,100000000000000000000", a)) //100 Ether
		}
	}

	cmd := exec.Command("ganache-cli", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println(cmd.Args)
	err := cmd.Start()
	if err != nil {
		return nil, nil, err
	}

	closeChan := make(chan struct{})

	go func() {
		fmt.Println("Wait for message to kill ganache-cli")
		<-closeChan
		fmt.Println("Received message,killing ganache-cli ...")

		if e := cmd.Process.Kill(); e != nil {
			fmt.Println("kill ganache-cli error:", e)
		}

		fmt.Println("ganache-cli exited.")
		closeChan <- struct{}{}
	}()

	fmt.Println("waiting for ganache-cli starting, 3 seconds")
	time.Sleep(3 * time.Second)

	killFunc := func() {
		closeChan <- struct{}{}
		time.Sleep(1 * time.Second)
		<-closeChan
	}

	cli, err := eth.NewEthClient(fmt.Sprintf("http://127.0.0.1:%d", ganacheCliPort))
	if err != nil {
		killFunc()
		return nil, nil, err
	}

	return cli, killFunc, nil
}
