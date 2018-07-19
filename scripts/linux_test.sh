# Copyright 2018 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved. Use of this source code is
# governed by the Apache-2.0 license that can be found in the LICENSE file.


#!/bin/bash

cd /root/go/src/github.com/chasinglogic/redalert
go get ./...
go test -v ./...
