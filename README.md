# Gitleaks Rules

> custom rules

## Usage

create a github action for your repo in `.github/workflows/.gitleaks.yml`

### Use Default Rules

```
name: gitleaks

on: [push,pull_request]

jobs:
  gitleaks:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
      with:
        fetch-depth: '2'
    - name: wget
      uses: wei/wget@v1
      with:
        args: -O .gitleaks.toml https://raw.githubusercontent.com/ycjcl868/gitleaks/master/.gitleaks.toml
    - name: gitleaks-action
      uses: zricethezav/gitleaks-action@master
```

### Using your own configuration

create a `.gitleaks.toml` in the root of your repo directory.

```
name: gitleaks

on: [push,pull_request]

jobs:
  gitleaks:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
      with:
        fetch-depth: '2'
    - name: gitleaks-action
      uses: zricethezav/gitleaks-action@master
```
