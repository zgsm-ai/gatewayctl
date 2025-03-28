package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/zgsm-ai/gatewayctl/internal/pkg/config"
	errv1 "github.com/zgsm-ai/gatewayctl/internal/pkg/error"
	log "github.com/zgsm-ai/gatewayctl/internal/pkg/logger"
)

type upstream struct {
	Type  string         `json:"type"`
	Nodes map[string]int `json:"nodes"`
}

type newRouterRequest struct {
	Uri      string   `json:"uri"`
	Upstream upstream `json:"upstream"`
}

func AddRouterToGateway(partUri, upstreamAddr string) (string, error) {
	conf := config.App.Gateway.Apisix
	adminUrl := conf.Url + "/" + partUri
	u, err := url.Parse(adminUrl)
	if err != nil {
		log.Error("url parse error: ", err)
		return "", err
	}

	uri := conf.Prefix + "/" + partUri
	reqBody := newRouterRequest{
		Uri: uri,
		Upstream: upstream{
			Type: "roundrobin",
			Nodes: map[string]int{
				upstreamAddr: 1,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		log.Error("json decode error: ", err)
		return "", err
	}

	req, err := http.NewRequest("PUT", u.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		log.Error("http.NewRequest error: ", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", conf.Token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("send request to gateway error: ", err)
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Error("json decode error: ", err)
		return "", err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		jsonData, err := json.Marshal(result)
		if err != nil {
			log.Error("json encode error: ", err)
			return "", err
		}
		return "", errv1.NewError(errv1.ErrRequest, string(jsonData))
	}

	log.Debugf("response body: %+v", result)

	return uri, nil
}

func RemoveRouterFromGateway(partUri string) error {
	conf := config.App.Gateway.Apisix
	adminUrl := conf.Url + "/" + partUri
	u, err := url.Parse(adminUrl)
	if err != nil {
		log.Error("url parse error: ", err)
		return err
	}

	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		log.Error("http.NewRequest error: ", err)
		return err
	}

	req.Header.Set("X-API-KEY", conf.Token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("send request to gateway error: ", err)
		return err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Error("json decode error: ", err)
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		jsonData, err := json.Marshal(result)
		if err != nil {
			log.Error("json encode error: ", err)
			return err
		}
		return errv1.NewError(errv1.ErrRequest, string(jsonData))
	}

	log.Debugf("response body: %+v", result)

	return nil
}
