/* This file lifted from https://github.com/grpc/grpc-go/blob/master/testdata/testdata.go */

/*
 * Copyright 2017 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package tls

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"path/filepath"
	"runtime"
)

// basepath is the root directory of this package.
var basepath string

func init() {
	_, currentFile, _, _ := runtime.Caller(0)
	basepath = filepath.Dir(currentFile)
}

// Path returns the absolute path the given relative file or directory path,
// relative to the google.golang.org/grpc/testdata directory in the user's GOPATH.
// If rel is already absolute, it is returned unmodified.
func Path(rel string) string {
	if filepath.IsAbs(rel) {
		return rel
	}

	return filepath.Join(basepath, rel)
}

func LocalhostTLSConfig() *tls.Config {
	cert, err := tls.LoadX509KeyPair(Path("localhost.crt"), Path("localhost.key"))
	if err != nil {
		panic(err)
	}

	ca, err := ioutil.ReadFile(Path("rootCA.crt"))
	if err != nil {
		panic(err)
	}

	pool := x509.NewCertPool()
	if ok := pool.AppendCertsFromPEM(ca); !ok {
		panic("Why is testdata/rootCA.crt not a valid PEM file?")
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    pool,
		RootCAs:      pool,
	}
}
