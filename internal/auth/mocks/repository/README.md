# Repository mock

This mock is used to test the repository layer.

## Regenerate mocks

Run from the root of the project.

```shell
mockery \ 
    --dir=internal/auth/repository \ 
    --name=TxQuerier \ 
    --filename=with_tx.go \ 
    --output=internal/auth/mocks/repository \ 
    --outpkg=mocks_repository 
```