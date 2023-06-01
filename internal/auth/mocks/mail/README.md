```shell
mockery \                               
    --dir=internal/auth/commands/handlers \
    --name=userEmailVerificationSender --filename=common.go \
    --output=internal/auth/mocks/mail --outpkg=mocks_mail \
    --structname=UserEmailVerificationSender --filename=verification_sender.go
```