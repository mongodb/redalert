python -m venv %APPDATA%\test_venv
%APPDATA%\test_venv\Scripts\activate
cd C:\vagrant
python -m pip install -r requirements.txt
python -m pip install -r requirements.dev.txt
python -m pytest -m windows
python -m pytest -m agnostic
