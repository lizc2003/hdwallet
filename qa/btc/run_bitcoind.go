package btc

import (
	"fmt"
	"github.com/lizc2003/hdwallet/btc"
	"github.com/lizc2003/hdwallet/qa"
	"github.com/lizc2003/hdwallet/wallet"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	RPCPortRegtest    = 18443
	BitcoinBinPathEnv = "BITCOIN_BIN_PATH"
)

// RunOptions .
type RunOptions struct {
	NewTmpDir bool
	RPCPort   uint
	Args      []string
}

func RunBitcoind(optionsPtr *RunOptions) (*btc.BtcClient, func(), error) {
	var options RunOptions
	if optionsPtr == nil {
		options = RunOptions{}
	} else {
		options = *optionsPtr
	}

	if options.RPCPort == 0 {
		options.RPCPort = RPCPortRegtest
	}

	if qa.IsProcessRunning(options.RPCPort, "bitcoin") {
		return nil, nil, fmt.Errorf("bitcoind already running on port %d", options.RPCPort)
	}

	var killHooks []killHook

	var dataDir string
	if options.NewTmpDir {
		for _, arg := range options.Args {
			if strings.Contains(arg, "-datadir=") {
				return nil, nil, fmt.Errorf("already provide param datadir >> %v", arg)
			}
		}

		tmpDir := strings.TrimRight(os.TempDir(), "/")
		dataDir = tmpDir + "/btccli_bitcoind_datatmp_" + time.Now().Format(time.RFC3339) + "/"
		err := os.MkdirAll(dataDir, 0777)
		if err != nil {
			return nil, nil, fmt.Errorf("cannot create tmp dir: %v, err: %v", dataDir, err)
		}

		options.Args = append(options.Args, "-datadir="+dataDir)

		killHooks = append(killHooks, func() error {
			return os.RemoveAll(dataDir)
		})
	}

	//bitcoin/share/rpcauth$ python3 rpcauth.py rpcusr 233
	//String to be appended to bitcoin.conf:
	//rpcauth=rpcusr:656f9dabc62f0eb697c801369617dc60$422d7fca742d4a59460f941dc9247c782558367edcbf1cd790b2b7ff5624fc1b
	//Your password: 233
	args := []string{
		"-regtest",
		"-txindex",
		"-fallbackfee=0.0002",
		"-rpcauth=rpcusr:656f9dabc62f0eb697c801369617dc60$422d7fca742d4a59460f941dc9247c782558367edcbf1cd790b2b7ff5624fc1b",
		fmt.Sprintf("-rpcport=%d", options.RPCPort),
	}
	args = append(args, options.Args...)

	cmd := exec.Command(getCmdBitcoind(), args...)
	fmt.Println(cmd.Args)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		for _, hook := range killHooks {
			hook()
		}
		return nil, nil, err
	}

	closeChan := make(chan struct{})
	go func() {
		fmt.Println("Wait for message to kill bitcoind")
		<-closeChan
		fmt.Println("Received message, killing bitcoind ...")

		if e := cmd.Process.Kill(); e != nil {
			fmt.Println("kill bitcoind error:", e)
		}
		for _, hook := range killHooks {
			hook()
		}

		fmt.Println("bitcoind exited.")
		closeChan <- struct{}{}
	}()

	fmt.Println("waiting for bitcoind starting, 5 seconds")
	time.Sleep(5 * time.Second)

	killFunc := func() {
		closeChan <- struct{}{}
		time.Sleep(1 * time.Second)
		<-closeChan
	}

	host := fmt.Sprintf("http://127.0.0.1:%d", options.RPCPort)
	cli, err := btc.NewBtcClient(host, "rpcusr", "233", wallet.BtcChainRegtest)
	if err != nil {
		killFunc()
		return nil, nil, err
	}

	return cli, killFunc, nil
}

type killHook func() error

func getCmdBitcoind() string {
	p := os.Getenv(BitcoinBinPathEnv)
	if p == "" {
		panic("run bitcoind need set env: BITCOIN_BIN_PATH")
	}
	cmdBitcoind := p + "/bitcoind"
	fmt.Println("bitcoin bin path:", cmdBitcoind)
	return cmdBitcoind
}
