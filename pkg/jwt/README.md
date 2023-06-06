# JWT wrapper

This package provides a wrapper around the [golang-jwt](github.com/golang-jwt/jwt/v5) package.

## Generate mocks

Example for auth service repository, replace `auth` with your service name.

```shell
mockery --dir=pkg/jwt \
 --name=Interactor --filename=interactor_mock.go \
 --output=pkg/jwt
```