package trx

import (
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"google.golang.org/grpc"
	"time"
)

type TrxClient struct {
	RpcClient *client.GrpcClient
}

func NewTrxClient(grpcHost string, apiKey string) (*TrxClient, error) {
	cli := client.NewGrpcClient(grpcHost)
	if apiKey != "" {
		cli.SetAPIKey(apiKey)
	}
	cli.SetTimeout(30 * time.Second)
	err := cli.Start(grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &TrxClient{RpcClient: cli}, nil
}

func (this *TrxClient) GetBlockHeight() (int64, error) {
	b, err := this.RpcClient.GetNowBlock()
	if err != nil {
		return 0, err
	}
	return b.BlockHeader.RawData.Number, nil
}

/*
func (this *TrxClient) GetAccount(addr string) (int64, error) {
	a, err := this.RpcClient.GetAccount(addr)
	if err != nil {
		return 0, err
	}

	return b.BlockHeader.RawData.Number, nil
}
*/
/*
const (
	UrlVisible = "?visible=true"
)

type TrxClient struct {
	httpClient *http.Client
	baseUrl    string
	headers    http.Header
}

func NewTrxClient(baseUrl string, apiKey string) (*TrxClient, error) {
	timeout := 30 * time.Second
	hc := &http.Client{
		Timeout: timeout + 2*time.Second,
		Transport: &http.Transport{
			MaxIdleConns:          200,
			MaxIdleConnsPerHost:   100,
			IdleConnTimeout:       60 * time.Second,
			DisableCompression:    true,
			ResponseHeaderTimeout: timeout,
			DialContext: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).DialContext,
		},
	}

	contentType := "application/json"
	headers := make(http.Header, 2)
	headers.Set("accept", contentType)
	headers.Set("content-type", contentType)
	if apiKey != "" {
		headers.Set("TRON-PRO-API-KEY", apiKey)
	}

	return &TrxClient{httpClient: hc, baseUrl: baseUrl, headers: headers}, nil
}

func (this *TrxClient) Call(method string, req interface{}, result interface{}) error {
	var err error
	var httpReq *http.Request

	if req == nil {
		httpReq, err = http.NewRequest("GET", this.baseUrl+method+UrlVisible, nil)
	} else {
		b, err := json.Marshal(req)
		if err != nil {
			return err
		}
		//fmt.Println(string(b))

		httpReq, err = http.NewRequest(http.MethodPost, this.baseUrl+method, bytes.NewReader(b))
	}
	if err != nil {
		return err
	}

	httpReq.Header = this.headers.Clone()
	resp, err := this.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("trx request fail: %d, %s", resp.StatusCode, string(body))
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}
	return nil
}

func (this *TrxClient) GetNowBlock() (*core.Block, error) {
	var b core.Block
	err := this.Call("wallet/getnowblock", nil, &b)
	if err != nil {
		return nil, err
	}
	return &b, nil
}
*/
