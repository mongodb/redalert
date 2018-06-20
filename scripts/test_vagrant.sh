#!/bin/bash

export STATUS=$(vagrant status $PLATFORM)
if [[ $STATUS == *"running"* ]]; then
    vagrant provision --provision-with test $PLATFORM
else
    vagrant up --provision $PLATFORM
fi
