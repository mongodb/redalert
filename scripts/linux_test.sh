#!/bin/bash

source ~/.venv/bin/activate
cd /vagrant
pip install -r requirements.txt
pip install -r requirements.dev.txt
pytest -m $PLATFORM
pytest -m agnostic
