## Summary
jwtserver is a minimal jwt REST server in an alpine based container. When started, the container will run a `netcat` based web server which invokes a script to sign a pre-existing claim. The sign script will create `iat`, `exp`, `auth_time` and `jti` fields as expected. Please read on for further info.

### local tools needed
- docker
- openssl (to sign locally outside container)
- jq (to sign locally outside container and to examine jwt)

### How to use
To build and test, use the `Makefile` provided.

`make local_test` will build and test the local endpoint. By default, the server will listen on `localhost:9090` and sign this [example claim](claims/example_claim.json)

The server will be left running and is accessible as shown below

```
curl -s localhost:9090/example_claim.json

eyJraWQiOiIvcHJpdmF0ZUtleS5wZW0iLCJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyAgImlzcyI6ICJUb2tlbiBTZXJ2aWNlIiwgICJzdWIiOiAiMTIzNDU2Nzg5IiwgICJuYW1lIjogIkpvaG4gRG9lIiwgICJhZG1pbiI6IHRydWUsICAiaWF0IjogMTYyMjY3NjYxNywgICJhdXRoX3RpbWUiOiAxNjIyNjc2NjE3LCAgImV4cCI6IDE2MjI2NzcyMTcsICAianRpIjogIjExMDcyYjM0LWI2NWUtNDFjYi05NjM0LTU1MmU4YjZjNTVjYiJ9.jKl_pSYruk_3NHPeQti5iqKxZKBsKIJ8VIc8LtBXOf2rGv0fjgJDs50B9BnAxu-SOp7SJP0B63jqIJYjoR0fqmjs2l0Jl1efF193Luo-EdtmSUc0vctLiP2h_EUGBlgfkitrIyDWkccKcRZapr-eJZgopLvrl9kFdRZyJmSkgUGEC0gvhTCTmeyPg5TUw2ytpm6OHqFKdfPJwWc85ZI-kjj6BpWdgeZkIrviAPsfbOEFmEcX6GL6OHaE70MnGS8PikIuUCSbcEGbkrMTh_K9w4YMBZtTIHIx2ERy0RS54fdJn_FbdjMXuLZeEZ_ejYtgF1Ozu9A2BJQoMkw7HbhuHw
```

### Adding custom claims
Copy additional claim files into the claims directory or map in a different directory on container run.

Here is how the run command will look if you are mapping in another directory for claims.

```
docker run --rm -itd \
-v$(pwd)/custom_claims:/claims \
-v$(pwd)/privateKey.pem:/privateKey.pem \
-p9090:9090 jwtserver
```

You can then get a different set of claims signed by using that file name in path as shown below.

```
curl -s localhost:9090/claim_with_email.json
```

### Changing token duration
There is a default duration set in the [Dockerfile](Dockerfile) as 600 seconds.
```
ENV DURATION 600
```
To override this at container start, to 1200 seconds, add `-eDURATION=1200` to the run command. This server will then sign all tokens with a duration of 1200 seconds instead of the default 600.

In addition to this, if you want a token to have a different duration, it can be done at token generation time as follows.

```
curl -s localhost:9090/example_claim.json?duration=1200
```

### Examine the jwt
Let us decode the second part of the `jwt` using [this script](../../scripts/jwt_decode_claims.sh)

```
curl -s localhost:9090/example_claim.json | ../../scripts/jwt_decode_claims.sh | jq
```

```json
{
  "iss": "Token Service",
  "sub": "123456789",
  "name": "John Doe",
  "admin": true,
  "iat": 1622677451,
  "auth_time": 1622677451,
  "exp": 1622678051,
  "jti": "3c8ccc13-28ad-43a2-8a3e-8fd1f625cf74"
}
```
This shows that in addition to contents of [example claim.json](claims/example_claim.json), the following fields are added:

- iat
- auth_time
- exp
- jti

### Run server
```
docker run -ePORT=9091 -eDURATION=600 -p9091:9091 -it -v$(pwd)/privateKey.pem:/privateKey.pem jwtserver
```

`PORT` and `DURATION` can be changed when launching the server.

In order to change the claims files, mount a directory with different claims files in json format as `/claims` and you can access those by `curl -s localhost:9090/your_claim_file.json`

## References
- See [sign_jwt_rsa.sh](sign_jwt_rsa.sh) for signing code including date and `jti` unique id generation.
