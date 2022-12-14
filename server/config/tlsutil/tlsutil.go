/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package tlsutil

import (
	"crypto/tls"
	"errors"
	"os"

	"github.com/apache/servicecomb-kie/pkg/cipherutil"
	"github.com/apache/servicecomb-kie/server/config"
	"github.com/go-chassis/foundation/stringutil"
	"github.com/go-chassis/foundation/tlsutil"
	"github.com/go-chassis/openlog"
)

var ErrRootCAMissing = errors.New("rootCAFile is empty in config file")

func Config(c *config.TLS) (*tls.Config, error) {
	var password string
	if c.CertPwdFile != "" {
		pwdBytes, err := os.ReadFile(c.CertPwdFile)
		if err != nil {
			openlog.Error("read cert password file failed: " + err.Error())
			return nil, err
		}
		password = cipherutil.TryDecrypt(stringutil.Bytes2str(pwdBytes))
	}
	if c.RootCA == "" {
		openlog.Error(ErrRootCAMissing.Error())
		return nil, ErrRootCAMissing
	}
	opts := append(tlsutil.DefaultClientTLSOptions(),
		tlsutil.WithVerifyPeer(c.VerifyPeer),
		tlsutil.WithVerifyHostName(false),
		tlsutil.WithKeyPass(password),
		tlsutil.WithCA(c.RootCA),
		tlsutil.WithCert(c.CertFile),
		tlsutil.WithKey(c.KeyFile),
	)
	return tlsutil.GetClientTLSConfig(opts...)
}
