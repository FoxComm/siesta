/* Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License. */

package siesta

import (
	crand "crypto/rand"
	"math/rand"
	"net"
	"reflect"
	"testing"
	"runtime"
)

const TCPListenerAddress = "localhost:4455"

func assert(t *testing.T, actual interface{}, expected interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		_, fn, line, _ := runtime.Caller(1)
		t.Errorf("Expected %v, actual %v\n@%s:%d", expected, actual, fn, line)
	}
}

func checkErr(t *testing.T, err error) {
	if err != nil {
		_, fn, line, _ := runtime.Caller(1)
		t.Errorf("%s\n @%s:%d", err, fn, line)
	}
}

func randomBytes(n int) []byte {
	b := make([]byte, n)
	crand.Read(b)
	return b
}

func randomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZйцукенгшщзхъфывапролджэжячсмитьбюЙЦУКЕНГШЩЗХЪФЫВАПРОЛДЖЭЯЧСМИТЬБЮ0123456789!@#$%^&*()")

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)[:n]
}

func testRequest(t *testing.T, request Request, expected []byte) {
	sizing := NewSizingEncoder()
	request.Write(sizing)
	bytes := make([]byte, sizing.Size())
	encoder := NewBinaryEncoder(bytes)
	request.Write(encoder)

	assert(t, bytes, expected)
}

func decode(t *testing.T, response Response, bytes []byte) {
	decoder := NewBinaryDecoder(bytes)
	err := response.Read(decoder)
	checkErr(t, err)
}

func awaitForTCPRequestAndReturn(t *testing.T, bufferSize int, resultChannel chan []byte) {
	netName := "tcp"
	addr, _ := net.ResolveTCPAddr(netName, "localhost:4455")
	listener, err := net.ListenTCP(netName, addr)
	if err != nil {
		t.Errorf("Unable to start tcp request listener: %s", err)
	}

	buffer := make([]byte, bufferSize)
	conn, err := listener.AcceptTCP()
	if err != nil {
		t.Error(err)
	}
	_, err = conn.Read(buffer)
	if err != nil {
		t.Error(err)
	}

	resultChannel <- buffer
}
