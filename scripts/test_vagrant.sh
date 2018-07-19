# Copyright 2018 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved. Use of this source code is
# governed by the Apache-2.0 license that can be found in the LICENSE file.


#!/bin/bash

export STATUS=$(vagrant status $PLATFORM)
if [[ $STATUS == *"running"* ]]; then
    vagrant provision --provision-with test $PLATFORM
else
    vagrant up --provision $PLATFORM
fi
