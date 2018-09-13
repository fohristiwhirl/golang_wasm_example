# golang_wasm_example

To compile:

* Set env variable GOOS to `js`
* Set env variable GOARCH to `wasm`
* Build as `main.wasm`

In Windows Powershell one can do:

```
$env:GOOS = 'js'
$env:GOARCH = 'wasm'
go build -o main.wasm main.go
```

To run:

* Make sure Apache (or whatever) is set up to specify the Content-Type of .wasm files as `application/wasm`
