# Copyright 2018 MongoDB Inc. All rights reserved. Use of this source code is
# governed by the Apache-2.0 license that can be found in the LICENSE file.


#!/bin/bash

cd /root/go/src/github.com/mongodb/redalert
go get ./...
go test -v ./...
