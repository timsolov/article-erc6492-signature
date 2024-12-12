# Go package for validation ERC-6492 signature

This is a package for validation [ERC-6492](https://eips.ethereum.org/EIPS/eip-6492) signature.
That project was born as an example for my article "ERC-6492 signature verification in Go".
Despite this package was written for the learning purpose it is ready for production as well.

The package supports validation SCA and EOA signatures in order described in [ERC-6492](https://eips.ethereum.org/EIPS/eip-6492) spec:

1. ERC-6492
2. ERC-1271
3. ecrecover

Installation:

```shell
go get github.com/timsolov/article-erc6492-signature
```

The usage example you can find in test file.