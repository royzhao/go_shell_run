package main

import (
	"log"
)

type DataClient struct {
	Peer *client
}

type Container struct {
	Image string  `json:"image"`
	Id    string  `json:"id"`
	Port  string  `json:"port"`
	Cpu   float64 `json:"cpu"`
	Mem   float64 `json:"mem"`
}
type SendData struct {
	Host       string      `json:"host"`
	Cpu        float64     `json:"cpu"`
	Mem        float64     `json:"mem"`
	Containers []Container `json:"containers"`
}

func NewPeerClient(endpoint string) (*DataClient, error) {
	c, err := newClient(endpoint)
	if err != nil {
		return nil, err
	}
	return &DataClient{
		Peer: c,
	}, nil
}

func (d *DataClient) Send(data SendData) error {
	_, _, err := d.Peer.do("POST", "/api/machine/stat/", data, true, nil)
	if err != nil {
		return err
	}
	log.Println(data)
	return nil
}
