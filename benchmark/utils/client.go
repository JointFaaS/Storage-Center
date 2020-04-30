package utils

import (
	"encoding/json"
	"net/http"
	"net/url"
	"bytes"
)

type GetResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type SetRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type SetResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Client struct {
	addr string
}

func (c *Client) Get(n string) (*string, error)  {
	resp, err := http.PostForm("http://" + c.addr + "/get", url.Values{"key": []string{n}})
	if err != nil {
		return nil, err
	}
	getResponse := GetResponse{}
	err = json.NewDecoder(resp.Body).Decode(&getResponse)
	if err != nil {
		return nil, err
	}
	return &getResponse.Value, nil
}

func (c *Client) Set(n string, v string) (error)  {
	body := SetRequest{
		Key: n,
		Value: v,
	}
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(body)
	
	resp, err := http.Post("http://" + c.addr + "/set", "application/json" , buf)
    if err != nil {
        return err
	}
	setResponse := SetResponse{}
	err = json.NewDecoder(resp.Body).Decode(&setResponse)
	if err != nil {
		return err
	}
	return nil
}

func NewClient(addr string) (*Client) {
	return &Client{
		addr: addr,
	}
}