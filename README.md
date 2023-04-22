# Virtual SSD (VSSD)

A way to host a server that can hold data from files, json data, plain text, etc in a key:value style. Not for production use.

## How to use

⚠️ please have golang installed on your machine before using this ⚠️

1. Clone the repository `git clone https://github.com/cookieforpres/vssd.git`
2. Run the server `go mod tidy && go build`
3. Now run `./vssd -h` for more information on how to use the tcp server or http server
