package transport

import (
	"testing"

	"github.com/MadBase/MadNet/interfaces"
	"github.com/sirupsen/logrus"
)

func TestTransportfail(t *testing.T) {
	logger := logrus.New()
	nodePrivKey1, err := newTransportPrivateKey()
	if err != nil {
		t.Fatal(err)
	}
	nodePrivKey1Hex := serializeTransportPrivateKey(nodePrivKey1)

	nodePrivKey2, err := newTransportPrivateKey()
	if err != nil {
		t.Fatal(err)
	}
	nodePrivKey2Hex := serializeTransportPrivateKey(nodePrivKey2)

	transport1, err := NewP2PTransport(logger, testCID, nodePrivKey1Hex, t1Port, t1Host)
	if err != nil {
		t.Fatal(err)
	}
	defer transport1.Close()
	nodeAddr1 := transport1.NodeAddr()

	transport2, err := NewP2PTransport(logger, testCIDFail, nodePrivKey2Hex, t2Port, t2Host)
	if err != nil {
		t.Fatal(err)
	}
	defer transport2.Close()

	complete1 := make(chan struct{})
	complete2 := make(chan struct{})
	complete3 := make(chan struct{})

	go dialer2(t, transport2, nodeAddr1, complete1)
	go accept2(transport1, complete2)
	go accept2(transport2, complete3)

	<-complete1
	<-complete2
	<-complete3
}

func dialer2(t *testing.T, transport interfaces.P2PTransport, addr interfaces.NodeAddr, complete chan struct{}) {
	defer close(complete)
	defer transport.Close()
	conn, err := transport.Dial(addr, 1)
	if err != nil {
		return
	}
	if conn != nil {
		t.Error("Got a connection that should not have occurred in dialer")
	}
}

func accept2(transport interfaces.P2PTransport, complete chan struct{}) {
	defer close(complete)
	_, err := transport.Accept()
	if err != nil {
		return
	}
}
