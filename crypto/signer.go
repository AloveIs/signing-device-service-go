package crypto

// Signer defines a contract for different types of signing implementations.
type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
	GetAlgorithm() string
	PublicKey() string
}

// KeyMarhsaller implements a key pair that can be exported
// with its two
type KeyMarhsaller interface {
	// Encode the public and private key
	// Returns the public key, private key and an error
	Marshal() ([]byte, []byte, error)
	// TODO: make marshall return a named struct not to swap public and private keys
}

// Signer that can be marshalled
type MarshallableSigner interface {
	Signer
	KeyMarhsaller
}

type ECCSigner struct {
	ECCMarshaler
	*ECCKeyPair
}

func NewECDSASigner() (MarshallableSigner, error) {
	g := &ECCGenerator{}
	keys, err := g.Generate()
	return &ECCSigner{ECCKeyPair: keys}, err
}

type RSASigner struct {
	RSAMarshaler
	*RSAKeyPair
}

func NewRSASigner() (MarshallableSigner, error) {
	g := &RSAGenerator{}
	keys, err := g.Generate()
	return &RSASigner{RSAKeyPair: keys}, err
}

func (s *ECCSigner) PublicKey() string {
	public, _, err := s.Marshal()
	if err != nil {
		// TODO: handle this
		panic(err)
	}
	return string(public)
}

func (s *ECCSigner) Marshal() ([]byte, []byte, error) {
	return s.ECCMarshaler.Marshal(*s.ECCKeyPair)
}

func UnmarshalECDSASigner(privateKey []byte) (MarshallableSigner, error) {
	g := NewECCMarshaler()
	keys, err := g.Unmarshal(privateKey)
	return &ECCSigner{ECCKeyPair: keys}, err
}

func UnmarshalRSASigner(privateKey []byte) (MarshallableSigner, error) {
	g := NewRSAMarshaler()
	keys, err := g.Unmarshal(privateKey)
	return &RSASigner{RSAKeyPair: keys}, err
}

func (s *RSASigner) PublicKey() string {
	public, _, err := s.Marshal()
	if err != nil {
		// TODO: handle this
		panic(err)
	}
	return string(public)
}

func (s *RSASigner) Marshal() ([]byte, []byte, error) {
	return s.RSAMarshaler.Marshal(*s.RSAKeyPair)
}
