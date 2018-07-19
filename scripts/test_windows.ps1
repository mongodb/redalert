# Copyright 2018 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved. Use of this source code is
# governed by the Apache-2.0 license that can be found in the LICENSE file.


python -m venv %APPDATA%\test_venv
%APPDATA%\test_venv\Scripts\activate
cd C:\vagrant
python -m pip install -r requirements.txt
python -m pip install -r requirements.dev.txt
python -m pytest -m windows
python -m pytest -m agnostic
