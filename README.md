# go-knapsackcrypto
Knapsack cryptosystems

Implementation of knapsack crypto systems in Go.

## Schemes

### Basic Merkle-Hellman knapsack public-key encryption

References:
- Chapter 6 of "An Introductio  to Mathematical Cryptography" by Jeffrey Hoffstein, Jill Pipher and J.H Silverman.
- Chapter 8 of "Handbook of Applied Cryptography" by Alfred J. Menezes et al.

Remarks:
- Decryption will produce a bit slice of the same length as the super imcreasimg sequence of the private key. But the actual bit length of the encrypted plaintext message could be smaller. Therefore to correctly extract the plaintext message I make use of the fact that I know the byte length of the original message. I guess in real-life situations you would encrypt together with the plaintext message some code that allows the decrypting party to know how long the message actually is.
- The number of bits of the public key is 2 times the length of the super increasig sequence to the power of 2, then the bit size of the public key grows exponentially with the length of the super increasing sequence. The bit length of the plaintext message needs to be smaller than or equal to the length of the super increasing sequence of the private key, therefore this scheme does not seem to be very practical for encrypting long messages.

Example:

`go run examples/merklehellman/main.go`

### Multiple-iterated Merkle-Hellman knapsack encryption

TODO

### Chor-Rivest knapsack public-key encryption

TODO