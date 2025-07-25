## es client
`krypton-cli es` is a client module to aid with integration and load testing of the `enroll` service.
Note that [krypton-cli](../README.md) is the client program for all `krypton` services. `es` appears as
a module in `krypton-cli`

### How to run
Refer to [krypton-cli](../README.md) for build/install steps. If you install, `krypton-cli` will be
available as a standard command.

### Supported commands
The following commands are supported initially. Make sure you run `make build` to create `krypton-cli` locally.
To see a list of supported commands, use cli with `es` arg
```
$ krypton-cli es
Supported commands are:
- enroll
- enroll_and_wait
- create_enroll_token
- get_enroll_token
- get_certificate
- get_device_token
- renew_enroll
- unenroll
```
### Getting help on any command
To get help, do `-help` on any command
```
$ krypton-cli es enroll -help
Usage of krypton-cli:
  -batch_size uint
        number of routines to run. parallel mode only. (default 100)
  -count uint
        how many times to run command (default 1)
  -dsts_server string
        dsts server (default "http://localhost:7001/api/v1")
  -jwt_token string
        provide a jwt token string
  -parallel
        should the runs be in parallel
  -retry_count uint
        how many times to retry failures (default 3)
  -server string
        enroll server (default "http://localhost:7979/api/v1")
  -stdin
        use stdin for command input
  -token_server string
        test token server (default "http://localhost:9090/api/v1/token")
  -token_type string
        token type (default "azuread")
  -verbose
        verbose logs (default false)
```

#### Chaining commands
Command results are in `json` format. This helps combining multiple commands with minimal external tools
to create useful workflows. Some examples are as follows.

- Scenario 1
  - `get_device_token` command returns a json with `device_id` and `device_token` as elements.
  - `renew_enroll` command takes a `device_id` and `device_token` to renew an existing enroll..
  - We can chain these commands without additional tools
```
krypton-cli es get_device_token | krypton-cli es renew_enroll -stdin
{"id":"9033aa7a-6095-4eac-99e9-7ed40cba1daa","bearer":"Bearer eyJhbGci...."}
```
- Scenario 2
  - `get_device_token` command returns a json with `device_id` and `device_token` as elements.
  - `unenroll` command takes a `device_id` and `device_token` to remove an existing enroll.
  - We can chain these commands without additional tools
```
krypton-cli es get_device_token | krypton-cli es unenroll -stdin

{"id":"f947eee2-cafc-4c18-8de3-71123a54336c","bearer":"Bearer eyJhbGci...."}
```
- Scenario 4
  - use `-parallel` flag on any command to increase throughput
  - example from Scenario 2 will show significant increase in throughput when `get_device_token` gets `parallel` flag
```
time krypton-cli es get_device_token -count 10 | krypton-cli es unenroll -stdin > /dev/null

real    0m24.183s
user    0m3.588s
sys     0m0.316s
time krypton-cli es get_device_token -count 10 -parallel | krypton-cli es unenroll -stdin > /dev/null

real    0m6.176s
user    0m3.722s
sys     0m0.316s
```

#### Commands needing app tokens
Note that `get_enroll_token` will need an app token instead of a regular user or device token.
An usage example is as follows.

1. Create an enrollment token with a user token
```
$ CLI_PROFILE_NAME=local go run main.go es create_enroll_token
{
  "tenant_id":"8ab17067-0057-483a-bd0c-1f52561d0f29",
  "enroll_token":"XyEFByGUQfOeZDdFvvgBloBEPmlpNkZyurvFjQqrKcfwhShAYXUrCuNzcJpjJUrm",
  "issued_at":-62135596800,
  "expires_at":1694510625
}
```
2. Fetch this enroll token for the tenant. Note this will need an app token and will fail
```
$ CLI_PROFILE_NAME=local go run main.go es get_enroll_token -tenant_id 8ab17067-0057-483a-bd0c-1f52561d0f29

2023/09/05 09:32:54 Please provide -tenant_id, -app_token or specify -stdin for standard input
2023/09/05 09:32:54 Use -help to see available args
2023/09/05 09:32:54 Missing arguments
```

Correct sequence between 1 and 2 is to do the following
```
CLI_PROFILE_NAME=local krypton-cli auth app_token -server http://localhost:7001/api/v1
eyJhbGciOiJSUzUxMiIsImtpZCI6ImZhYTQwOTgyNzI1M2VhODJmOTE3YmVlMGQ4NmZmYzZkODA2ZGNhMTA0YjMwMDAzMzE1NTVkYTQyZjZkNjAxYjUiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJIUCBEZXZpY2UgVG9rZW4gU2VydmljZSIsInN1YiI6ImRlN2U1OTVmLTlhY2EtNDMzNC05ZjQ3LTIzNTJkMDBhY2FjZSIsImV4cCI6MTY5MzkxMDMwNSwibmJmIjoxNjkzOTA2NzA1LCJpYXQiOjE2OTM5MDY3MDUsImp0aSI6Ijg2MzI3ZDYyLTRhMDctNGZkOS1iOTU5LTVkOWZhZjY4YzU1NyIsInR5cCI6ImFwcCJ9.JJ1F92CeY7rYXw-LxIYEfB2RmYPQiQtewlCICU_ms_n2NkmlZ3AGTVbbWzXmj3F5N1eoSK_DnRN0-tqNJ5VjTKb4zv21zbg2R3ws0Lz9Xjvn7dhZGllB9O-2vfktabfe61LuNf4pzKbcarcYJ1qoAQot2QqF9JtIPMcgw6gGKr1vCScsgwB2dpbCmAG3q-jvcEAniG-kyzzKQ_koM08uGM4-nPjYAouPETpSgnClQZKeRp4b0GOTVnq3zzOLy_5Ce8jjNIEUqE1PvGPDkMBeKdOeIb0bDXv6ddUwvoZv5d72PqUr5TWFsqIWTVkRhDCCIPqGaqwCNsQ5V0mcrWC8g7IQJs--YWr22GR2ksIunQXeqymJL4P_GtbUaRPdKoqG4-8BpHhy2ly1diy75kyh8IXgQxamQDJSgWgsYvCrwOc-P02Sz1byirHl_DztCUbO2XrAIz6N6DtNmdyVk51w_19_6bECN35k3C6h3TKelytLcGB6CiuG9zEw7HqXhsaxMwfyY-82uMlJkiUgB6AgosahJrv8ZAzfd_UvBIZDhXlAznzxw7zd-KPLkTf71PPmtTcq5IqBezqpG0lSTguTHIvcqv7CfWuYQq7GE50nVvtXG21yyz1ZjYXJS9rt7sbjuLzZ-6yi6kpxShrpRKjJRfzN9E2qU6RAOjd4JmvwkR8
2023/09/05 09:38:25 Elapsed: 13.891487ms, Processed: 1
```

Note that this is cached in cli
```
$ ls -al ~/.krypton-cli/app_token
-rw-r--r-- 1 user user 1067 Sep  5 09:38 /home/user/.krypton-cli/app_token
```

Now, the get call will work
```
$ CLI_PROFILE_NAME=local go run main.go es get_enroll_token -tenant_id 8ab17067-0057-483a-bd0c-1f52561d0f29
{"tenant_id":"8ab17067-0057-483a-bd0c-1f52561d0f29","enroll_token":"XyEFByGUQfOeZDdFvvgBloBEPmlpNkZyurvFjQqrKcfwhShAYXUrCuNzcJpjJUrm","issued_at":1693905825,"expires_at":1694510625}

2023/09/05 09:40:09 Elapsed: 4.56568ms, Processed: 1
```
