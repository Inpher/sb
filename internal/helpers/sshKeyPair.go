package helpers

// SSHKeyPair describes an SSH key pair
type SSHKeyPair struct {
	PublicKey          *PublicKey
	PrivateKeyFilepath string
}
