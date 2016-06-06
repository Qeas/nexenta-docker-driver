package nedapi

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
)


type Client struct {
	SVIP              string
	Endpoint          string
	DefaultAPIPort    int
	DefaultVolSize    int64 //bytes
	DefaultAccountID  int64
	DefaultTenantName string
	Config            *Config
}


type Config struct {
	IOProtocol	string // NFS, iSCSI, NBD, S3
	EndPoint	string // server:/export, IQN, devname, 
	TenantName	string
	AccessKey	string
	SecretKey	string
	MountPoint	string
}

func ReadParseConfig(fname string) (Config, error) {
	content, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatal("Error processing config file: ", err)
	}
	var conf Config
	err = json.Unmarshal(content, &conf)
	if err != nil {
		log.Fatal("Error parsing config file: ", err)
	}
	return conf, nil
}


func ClientAlloc(configFile string) (c *Client, err error) {
	conf, err := ReadParseConfig(configFile)
	if err != nil {
		log.Fatal("Error initializing client from Config file: ", configFile, "(", err, ")")
	}

	//TODO:
	//DefaultApiPort
	//DefaultAccountID
	NexentaClient := &Client{
		SVIP:		conf.IOProtocol,
		Endpoint:	conf.EndPoint,
		DefaultAPIPort:	8888,
		DefaultAccountID:	9999,
		DefaultTenantName: conf.TenantName,
		Config:	&conf,
	}

	return NexentaClient, nil
}

func (c *Client) Request(method string, params interface{}, id int) (response []byte, err error) {
	log.Debug("Issue request to SolidFire Endpoint...")
	if c.Endpoint == "" {
		log.Error("Endpoint is not set, unable to issue requests")
		err = errors.New("Unable to issue json-rpc requests without specifying Endpoint")
		return nil, err
	}
	data, err := json.Marshal(map[string]interface{}{
		"method": method,
		"id":     id,
		"params": params,
	})

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	log.Debugf("POST request to: %+v", c.Endpoint)
	Http := &http.Client{Transport: tr}
	resp, err := Http.Post(c.Endpoint,
		"json-rpc",
		strings.NewReader(string(data)))
	if err != nil {
		log.Errorf("Error encountered posting request: %v", err)
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}

	var prettyJson bytes.Buffer
	_ = json.Indent(&prettyJson, body, "", "  ")
	log.WithField("", prettyJson.String()).Debug("request:", id, " method:", method, " params:", params)

	errresp := APIError{}
	json.Unmarshal([]byte(body), &errresp)
	if errresp.Error.Code != 0 {
		err = errors.New("Received error response from API request")
		return body, err
	}
	return body, nil
}
