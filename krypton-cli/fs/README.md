## Overview
`krypton-cli fs` is a client program to aid with integration and load testing of the `fs` service.
Note that [krypton-cli](../README.md) is the client program for all `krypton` services.

### How to run
Refer to [krypton-cli](../README.md) for build/install steps. If you install, `krypton-cli` will be
available as a standard command.

### Supported commands
The following commands are supported initially. Make sure you run `make build` to create `krypton-cli` locally.
To see a list of supported commands, use cli with `fs` arg
```
$ krypton-cli fs
Supported commands are:
- create_file
- get_download_url
- get_file_details
- get_upload_url
```

#### Chaining commands
Command results are in `json` format. This helps combining multiple commands with minimal external tools
to create useful workflows. Some examples are as follows.

- Scenario 1
  - `create_file` command returns a json with `file_id` as one of the elements.
  - `get_file_details` command takes a `file_id` and prints details of the file.
  - If we could manipulate json in a pipeline, we can combine these commands
```
krypton-cli fs get_file_details -id=$(krypton-cli fs create_file | jq '.file.file_id')
{
  "file_id": 577,
  "tenant_id": "57836902-253f-41a6-8533-1e65c8c9892d",
  "device_id": "6d2ef5f2-5fc4-4fef-b2f3-b69c9b9854ac",
  "name": "fs_cli_upload_2553695053",
  "checksum": "dnRSeK2M7oR5V4FkCJXu9Q==",
  "size": 31,
  "status": "new",
  "created_at": "2023-04-19T03:39:45.774977Z",
  "updated_at": "2023-04-19T03:39:45.774977Z"
}
```
- Scenario 2
  - `create_file` command can run in batch mode making `n` files
  - `get_file_details` has an `stdin` flag which will allow json input `{"file_id":n}`
```
krypton-cli fs create_file -count 10 | krypton-cli fs get_file_details -stdin

{"file_id":210,"url":"","tenant_id":"0276024f-06d6-4619-a793-ab7e3ca44f52","device_id":"98ca710f-49d3-4e29-860b-d58e9ba31cf0","name":"fs_cli_upload_2536329754","checksum":"WcbdjUctL2jaLgW/b/lNug==","size":15,"status":"new","created_at":"2023-05-04T05:35:30.891696Z","updated_at":"2023-05-04T05:35:30.891696Z"}
...
```
- Scenario 3
  - `create_file` command is actually an internal pipe of `get_upload_url` and an `http PUT`
  - we can try that with `curl` as the http client
```
echo "hello" > /tmp/file
krypton-cli fs get_upload_url -filename /tmp/file | jq '.file.url' |
 xargs -n1 -I{} curl -XPUT -T"/tmp/file" -H"Content-Type: application/octet-stream" -H"Content-MD5:$(openssl dgst -md5 -binary /tmp/file | base64) {}
```
- Scenario 4
  - use `-parallel` flag on any command to increase throughput
  - example from Scenario 2 will show significant increase in throughput when `create_file` gets `parallel` flag
```
time go run ./main.go fs create_file -count 100 | go run ./main.go fs get_file_details -stdin > /dev/null

real    0m1.929s
user    0m0.531s
sys     0m0.140s
time go run ./main.go fs create_file -count 100 -parallel | go run ./main.go fs get_file_details -stdin > /dev/null

real    0m0.705s
user    0m0.515s
sys     0m0.145s
```

#### get upload url
`get_upload_url` creates a metadata entry in fs db and returns a `PUT` url that can be used with clients like `curl`
- file name can be specified using `-filename` flag. If not specified, `/tmp/1` is used as default.
- tenant id can be specified using `-tenant_id` flag. If not specified, a random uuid is used.
- device id can be specified using `-device_id` flag. If not specified, a random uuid is used.
```
krypton-cli fs get_upload_url -filename /tmp/file
{"file":{"file_id":640,"url":"http://fs_storage.local.dev:9000/mytestkrypton20221130/b9f02d55-b322-405d-b2aa-7c27edf2030e/ee614aba-f86e-4bbd-a619-163f9f27b805/640?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=minioadmin%2F20230419%2F%2Fs3%2Faws4_request&X-Amz-Date=20230419T040921Z&X-Amz-Expires=900&X-Amz-SignedHeaders=content-length%3Bcontent-md5%3Bhost&x-id=PutObject&X-Amz-Signature=630539ae6b843f0fb22df465195b7892ccd3399b39cbcea2a4c6763d2cbb291a"}}
```

#### create file
`create_file` gets an upload url (using a `get_upload_url` internally) and uploads a file with random contents.
- content length of uploaded files can be controlled using `-max_file_size` param
- tenant id can be specified using `-tenant_id` flag. If not specified, a random uuid is used.
- device id can be specified using `-device_id` flag. If not specified, a random uuid is used.

Example usage:
```
krypton-cli fs create_file
{"file":{"file_id":574,"url":"http://fs_storage.local.dev:9000/mytestkrypton20221130/5e71a30d-59f2-4e85-9bb9-10b78fff30b0/27d9a35a-69c4-4abe-9219-b1aa7ca8ac4f/574?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=minioadmin%2F20230419%2F%2Fs3%2Faws4_request&X-Amz-Date=20230419T025448Z&X-Amz-Expires=900&X-Amz-SignedHeaders=content-length%3Bcontent-md5%3Bhost&x-id=PutObject&X-Amz-Signature=e0359e8ddb4bc94b4645ec67fa61fb5a73edcebb85ebf339495956e39e25d753"}}
```

#### get file details
`get_file_details` returns info about a file. Files can be specified using a valid `id`.
- File id integer. If not specified, a default file id of `1` is used.

```
krypton-cli fs get_file_details -id 201
{"request_id":"a058568c-21e8-4b03-85e9-d71f9489e59f","response_time":"2023-04-19T04:14:33.505900361Z","file":{"file_id":201,"tenant_id":"1c5cb1bd-b1c3-4ded-996b-bfd37a44a6df","device_id":"e17d8e48-2f4e-4011-85ae-bf0308b2422d","name":"1","checksum":"McIu5m+I2u2bEthrjHRqWw==","size":1179,"status":"new","created_at":"2023-04-18T21:18:08.830087Z","updated_at":"2023-04-18T21:18:08.830087Z"}}
```
