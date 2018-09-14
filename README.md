# golang_wasm_example

This is, sadly, much slower than a pure JS version...

Compilation (for Windows):

```
$env:GOOS = 'js'
$env:GOARCH = 'wasm'
go build -o main.wasm main.go
```

To run:

* Make sure Apache (or whatever) is set up to specify the Content-Type of .wasm files as `application/wasm`
