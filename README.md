# API
This is the API the Oomph authenticator will use in the near future to authenticate users
that wish to use the Oomph proxy. This will be replacing the authentication server, which would end up
being a hassle to maintain. The reasons this API is public is

1. To make an open-source project that [I](https://www.github.com/ethaniccc) am able to put on my portfolio.
2. To ensure that Oomph's authentication is not secure only through obscurity.

## Exploit/Bug Finds
If you find a bug or exploit in the API, please feel free to DM me on Discord `@ethaniccc`. I'm a broke college student at the moment, but I'd be willing to compensate you if your exploit is legitimate and effective :)

## Preventing MiTM Exploits on the Authenticator
Because this project is open-source, some bad actors may attempt to clone this repository in attempt to host their own authentication servers and use Oomph without permission. To prevent MiTM exploits, we verify the API's SSL certificate on the authenticator.

## TODO
Some documentation on how the authenticator is supposed to use the API.