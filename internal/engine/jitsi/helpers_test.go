package jitsi

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"testing"

	"github.com/zarazaex69/j"
)

func encodeForTest(t *testing.T, data []byte) string {
	t.Helper()
	return base64.StdEncoding.EncodeToString(data)
}

func makeBridgeMessage(class string, fields map[string]any) j.BridgeMessage {
	return j.BridgeMessage{
		Class:  class,
		Fields: fields,
	}
}

func makeBridgeMessageFrom(from string, fields map[string]any) j.BridgeMessage {
	return j.BridgeMessage{
		Class:  "EndpointMessage",
		From:   from,
		Fields: fields,
	}
}

func makeBridgeFrame(t *testing.T, payload []byte) string {
	t.Helper()
	return makeBridgeFrameForEpoch(t, 0x10203040, 0, payload)
}

func makeBridgeFrameForEpoch(t *testing.T, senderEpoch, receiverEpoch uint32, payload []byte) string {
	t.Helper()
	framed := append([]byte{}, bridgeMagic[:]...)
	var hdr [8]byte
	binary.BigEndian.PutUint32(hdr[0:4], senderEpoch)
	binary.BigEndian.PutUint32(hdr[4:8], receiverEpoch)
	framed = append(framed, hdr[:]...)
	framed = append(framed, payload...)
	return base64.StdEncoding.EncodeToString(framed)
}

func decodeBridgeFrameForTest(t *testing.T, frame []byte) (uint32, uint32, []byte) {
	t.Helper()
	const epochHeaderLen = 8
	if len(frame) < len(bridgeMagic)+epochHeaderLen {
		t.Fatalf("bridge frame length = %d, want at least %d", len(frame), len(bridgeMagic)+epochHeaderLen)
	}
	if !bytes.Equal(frame[:len(bridgeMagic)], bridgeMagic[:]) {
		t.Fatalf("bridge frame magic = %x, want %x", frame[:len(bridgeMagic)], bridgeMagic)
	}
	off := len(bridgeMagic)
	return binary.BigEndian.Uint32(frame[off : off+4]),
		binary.BigEndian.Uint32(frame[off+4 : off+epochHeaderLen]),
		frame[off+epochHeaderLen:]
}
