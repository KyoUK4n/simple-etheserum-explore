#!/usr/bin/env zsh

# generate apis
# goctl api go --api etherscan.api --dir . --style go_zero

# generate dao
#goctl model mysql ddl --database etherscan -d ./ -s ./ddl/*.sql --style go_zero

# generate abi bind
#abigen --abi abi/ERC20.json --pkg=v1 --type ERC20 --out=pkg/contracts/v1/erc20.go
#abigen --v2 --abi abi/ERC20.json --pkg=v2 --type ERC20 --out=pkg/contracts/v2/erc20.go

cd web
npm install
npm run build

cd ../

go mod tidy && go mod vendor
go run .