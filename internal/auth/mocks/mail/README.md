# Mail sender mock

This mock is used to test the mail sender layer.

## Regenerate mocks

Run from the root of the project.

```shell
mockery \ 
    --dir=internal/auth/commands/handlers \ 
    --name=userEmailVerificationSender \ 
    --filename=common.go \ 
    --output=internal/auth/mocks/mail \ 
    --outpkg=mocks_mail \ 
    --structname=UserEmailVerificationSender \ 
    --filename=verification_sender.go 
```