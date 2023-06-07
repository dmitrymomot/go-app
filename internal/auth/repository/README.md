# Repository

This is template for a repository. It is intended to be used as a starting point for new repositories.

## Usage
 
TODO: Add usage instructions

## Generate mocks

Example for auth service repository, replace `auth` with your service name.

```shell
mockery \
    --dir=internal/auth/repository \
    --name=TxQuerier --filename=with_tx.go \
    --output=internal/auth/mocks/repository \
    --outpkg=mocks_repository
```