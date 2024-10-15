package i2pkeys

import (
	"bytes"
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"os"
	"strings"
)

var (
	i2pB64enc *base64.Encoding = base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-~")
	i2pB32enc *base32.Encoding = base32.NewEncoding("abcdefghijklmnopqrstuvwxyz234567")
)

// If you set this to true, Addr will return a base64 String()
var StringIsBase64 bool

// The public and private keys associated with an I2P destination. I2P hides the
// details of exactly what this is, so treat them as blobs, but generally: One
// pair of DSA keys, one pair of ElGamal keys, and sometimes (almost never) also
// a certificate. String() returns you the full content of I2PKeys and Addr()
// returns the public keys.
type I2PKeys struct {
	Address I2PAddr // only the public key
	Both    string  // both public and private keys
}

// Creates I2PKeys from an I2PAddr and a public/private keypair string (as
// generated by String().)
func NewKeys(addr I2PAddr, both string) I2PKeys {
	log.WithField("addr", addr).Debug("Creating new I2PKeys")
	return I2PKeys{addr, both}
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) (bool, error) {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		log.WithField("filename", filename).Debug("File does not exist")
		return false, nil
	} else if err != nil {
		log.WithError(err).WithField("filename", filename).Error("Error checking file existence")
		return false, fmt.Errorf("error checking file existence: %w", err)
	}
	exists := !info.IsDir()
	if exists {
		log.WithField("filename", filename).Debug("File exists")
	} else {
		log.WithField("filename", filename).Debug("File is a directory")
	}
	return !info.IsDir(), nil
}

// load keys from non standard format
func LoadKeysIncompat(r io.Reader) (k I2PKeys, err error) {
	log.Debug("Loading keys from reader")
	var buff bytes.Buffer
	_, err = io.Copy(&buff, r)
	if err == nil {
		parts := strings.Split(buff.String(), "\n")
		k = I2PKeys{I2PAddr(parts[0]), parts[1]}
		log.WithField("keys", k).Debug("Loaded keys")
	}
	log.WithError(err).Error("Error copying from reader, did not load keys")
	return
}

// load keys from non-standard format by specifying a text file.
// If the file does not exist, generate keys, otherwise, fail
// closed.
func LoadKeys(r string) (I2PKeys, error) {
	log.WithField("filename", r).Debug("Loading keys from file")
	exists, err := fileExists(r)
	if err != nil {
		log.WithError(err).Error("Error checking if file exists")
		return I2PKeys{}, err
	}
	if !exists {
		log.WithError(err).Error("File does not exist")
		return I2PKeys{}, fmt.Errorf("file does not exist: %s", r)
	}
	fi, err := os.Open(r)
	if err != nil {
		log.WithError(err).WithField("filename", r).Error("Error opening file")
		return I2PKeys{}, fmt.Errorf("error opening file: %w", err)
	}
	defer fi.Close()
	log.WithField("filename", r).Debug("File opened successfully")
	return LoadKeysIncompat(fi)
}

// store keys in non standard format
func StoreKeysIncompat(k I2PKeys, w io.Writer) error {
	log.Debug("Storing keys")
	_, err := io.WriteString(w, k.Address.Base64()+"\n"+k.Both)
	if err != nil {
		log.WithError(err).Error("Error writing keys")
		return fmt.Errorf("error writing keys: %w", err)
	}
	log.WithField("keys", k).Debug("Keys stored successfully")
	return nil
}
func StoreKeys(k I2PKeys, r string) error {
	log.WithField("filename", r).Debug("Storing keys to file")
	if _, err := os.Stat(r); err != nil {
		if os.IsNotExist(err) {
			log.WithField("filename", r).Debug("File does not exist, creating new file")
			fi, err := os.Create(r)
			if err != nil {
				log.WithError(err).Error("Error creating file")
				return err
			}
			defer fi.Close()
			return StoreKeysIncompat(k, fi)
		}
	}
	fi, err := os.Open(r)
	if err != nil {
		log.WithError(err).Error("Error opening file")
		return err
	}
	defer fi.Close()
	return StoreKeysIncompat(k, fi)
}

func (k I2PKeys) Network() string {
	return k.Address.Network()
}

// Returns the public keys of the I2PKeys.
func (k I2PKeys) Addr() I2PAddr {
	return k.Address
}

func (k I2PKeys) Public() crypto.PublicKey {
	return k.Address
}

func (k I2PKeys) Private() []byte {
	log.Debug("Extracting private key")
	src := strings.Split(k.String(), k.Addr().String())[0]
	var dest []byte
	_, err := i2pB64enc.Decode(dest, []byte(src))
	if err != nil {
		log.WithError(err).Error("Error decoding private key")
		panic(err)
	}
	return dest
}

type SecretKey interface {
	Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) (signature []byte, err error)
}

func (k I2PKeys) SecretKey() SecretKey {
	var pk ed25519.PrivateKey = k.Private()
	return pk
}

func (k I2PKeys) PrivateKey() crypto.PrivateKey {
	var pk ed25519.PrivateKey = k.Private()
	_, err := pk.Sign(rand.Reader, []byte("nonsense"), crypto.Hash(0))
	if err != nil {
		log.WithError(err).Warn("Error in private key signature")
		//TODO: Elgamal, P256, P384, P512, GOST? keys?
	}
	return pk
}

func (k I2PKeys) Ed25519PrivateKey() *ed25519.PrivateKey {
	return k.SecretKey().(*ed25519.PrivateKey)
}

/*func (k I2PKeys) ElgamalPrivateKey() *ed25519.PrivateKey {
	return k.SecretKey().(*ed25519.PrivateKey)
}*/

//func (k I2PKeys) Decrypt(rand io.Reader, msg []byte, opts crypto.DecrypterOpts) (plaintext []byte, err error) {
//return k.SecretKey().(*ed25519.PrivateKey).Decrypt(rand, msg, opts)
//}

func (k I2PKeys) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) (signature []byte, err error) {
	return k.SecretKey().(*ed25519.PrivateKey).Sign(rand, digest, opts)
}

// Returns the keys (both public and private), in I2Ps base64 format. Use this
// when you create sessions.
func (k I2PKeys) String() string {
	return k.Both
}

func (k I2PKeys) HostnameEntry(hostname string, opts crypto.SignerOpts) (string, error) {
	sig, err := k.Sign(rand.Reader, []byte(hostname), opts)
	if err != nil {
		log.WithError(err).Error("Error signing hostname")
		return "", fmt.Errorf("error signing hostname: %w", err)
	}
	return string(sig), nil
}

// I2PAddr represents an I2P destination, almost equivalent to an IP address.
// This is the humongously huge base64 representation of such an address, which
// really is just a pair of public keys and also maybe a certificate. (I2P hides
// the details of exactly what it is. Read the I2P specifications for more info.)
type I2PAddr string

// an i2p destination hash, the .b32.i2p address if you will
type I2PDestHash [32]byte

// create a desthash from a string b32.i2p address
func DestHashFromString(str string) (dhash I2PDestHash, err error) {
	log.WithField("address", str).Debug("Creating desthash from string")
	if strings.HasSuffix(str, ".b32.i2p") && len(str) == 60 {
		// valid
		_, err = i2pB32enc.Decode(dhash[:], []byte(str[:52]+"===="))
		if err != nil {
			log.WithError(err).Error("Error decoding base32 address")
		}
	} else {
		// invalid
		err = errors.New("invalid desthash format")
		log.WithError(err).Error("Invalid desthash format")
	}
	return
}

// create a desthash from a []byte array
func DestHashFromBytes(str []byte) (dhash I2PDestHash, err error) {
	log.Debug("Creating DestHash from bytes")
	if len(str) == 32 {
		// valid
		//_, err = i2pB32enc.Decode(dhash[:], []byte(str[:52]+"===="))
		log.WithField("str", str).Debug("Copying str to desthash")
		copy(dhash[:], str)
	} else {
		// invalid
		err = errors.New("invalid desthash format")
		log.WithField("str", str).Error("Invalid desthash format")
	}
	return
}

// get string representation of i2p dest hash(base32 version)
func (h I2PDestHash) String() string {
	b32addr := make([]byte, 56)
	i2pB32enc.Encode(b32addr, h[:])
	return string(b32addr[:52]) + ".b32.i2p"
}

// get base64 representation of i2p dest sha256 hash(the 44-character one)
func (h I2PDestHash) Hash() string {
	hash := sha256.New()
	hash.Write(h[:])
	digest := hash.Sum(nil)
	buf := make([]byte, 44)
	i2pB64enc.Encode(buf, digest)
	return string(buf)
}

// Returns "I2P"
func (h I2PDestHash) Network() string {
	return "I2P"
}

// Returns the base64 representation of the I2PAddr
func (a I2PAddr) Base64() string {
	return string(a)
}

// Returns the I2P destination (base32-encoded)
func (a I2PAddr) String() string {
	if StringIsBase64 {
		return a.Base64()
	}
	return string(a.Base32())
}

// Returns "I2P"
func (a I2PAddr) Network() string {
	return "I2P"
}

// Creates a new I2P address from a base64-encoded string. Checks if the address
// addr is in correct format. (If you know for sure it is, use I2PAddr(addr).)
func NewI2PAddrFromString(addr string) (I2PAddr, error) {
	log.WithField("addr", addr).Debug("Creating new I2PAddr from string")
	if strings.HasSuffix(addr, ".i2p") {
		if strings.HasSuffix(addr, ".b32.i2p") {
			// do a lookup of the b32
			log.Warn("Cannot convert .b32.i2p to full destination")
			return I2PAddr(""), errors.New("cannot convert .b32.i2p to full destination")
		}
		// strip off .i2p if it's there
		addr = addr[:len(addr)-4]
	}
	addr = strings.Trim(addr, "\t\n\r\f ")
	// very basic check
	if len(addr) > 4096 || len(addr) < 516 {
		log.Error("Invalid I2P address length")
		return I2PAddr(""), errors.New(addr + " is not an I2P address")
	}
	buf := make([]byte, i2pB64enc.DecodedLen(len(addr)))
	if _, err := i2pB64enc.Decode(buf, []byte(addr)); err != nil {
		log.Error("Address is not base64-encoded")
		return I2PAddr(""), errors.New("Address is not base64-encoded")
	}
	log.Debug("Successfully created I2PAddr from string")
	return I2PAddr(addr), nil
}

func FiveHundredAs() I2PAddr {
	log.Debug("Generating I2PAddr with 500 'A's")
	s := ""
	for x := 0; x < 517; x++ {
		s += "A"
	}
	r, _ := NewI2PAddrFromString(s)
	return r
}

// Creates a new I2P address from a byte array. The inverse of ToBytes().
func NewI2PAddrFromBytes(addr []byte) (I2PAddr, error) {
	log.Debug("Creating I2PAddr from bytes")
	if len(addr) > 4096 || len(addr) < 384 {
		log.Error("Invalid I2P address length")
		return I2PAddr(""), errors.New("Not an I2P address")
	}
	buf := make([]byte, i2pB64enc.EncodedLen(len(addr)))
	i2pB64enc.Encode(buf, addr)
	return I2PAddr(string(buf)), nil
}

// Turns an I2P address to a byte array. The inverse of NewI2PAddrFromBytes().
func (addr I2PAddr) ToBytes() ([]byte, error) {
	return i2pB64enc.DecodeString(string(addr))
}

func (addr I2PAddr) Bytes() []byte {
	b, _ := addr.ToBytes()
	return b
}

// Returns the *.b32.i2p address of the I2P address. It is supposed to be a
// somewhat human-manageable 64 character long pseudo-domain name equivalent of
// the 516+ characters long default base64-address (the I2PAddr format). It is
// not possible to turn the base32-address back into a usable I2PAddr without
// performing a Lookup(). Lookup only works if you are using the I2PAddr from
// which the b32 address was generated.
func (addr I2PAddr) Base32() (str string) {
	return addr.DestHash().String()
}

func (addr I2PAddr) DestHash() (h I2PDestHash) {
	hash := sha256.New()
	b, _ := addr.ToBytes()
	hash.Write(b)
	digest := hash.Sum(nil)
	copy(h[:], digest)
	return
}

// Makes any string into a *.b32.i2p human-readable I2P address. This makes no
// sense, unless "anything" is an I2P destination of some sort.
func Base32(anything string) string {
	return I2PAddr(anything).Base32()
}

/*
HELLO VERSION MIN=3.1 MAX=3.1
DEST GENERATE SIGNATURE_TYPE=7
*/
func NewDestination() (*I2PKeys, error) {
	removeNewlines := func(s string) string {
		return strings.ReplaceAll(strings.ReplaceAll(s, "\r\n", ""), "\n", "")
	}
	//
	log.Debug("Creating new destination via SAM")
	conn, err := net.Dial("tcp", "127.0.0.1:7656")
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	_, err = conn.Write([]byte("HELLO VERSION MIN=3.1 MAX=3.1\n"))
	if err != nil {
		log.WithError(err).Error("Error writing to SAM bridge")
		return nil, err
	}
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		log.WithError(err).Error("Error reading from SAM bridge")
		return nil, err
	}
	if n < 1 {
		log.Error("No data received from SAM bridge")
		return nil, fmt.Errorf("no data received")
	}

	response := string(buf[:n])
	log.WithField("response", response).Debug("Received response from SAM bridge")

	if strings.Contains(string(buf[:n]), "RESULT=OK") {
		_, err = conn.Write([]byte("DEST GENERATE SIGNATURE_TYPE=7\n"))
		if err != nil {
			log.WithError(err).Error("Error writing DEST GENERATE to SAM bridge")
			return nil, err
		}
		n, err = conn.Read(buf)
		if err != nil {
			log.WithError(err).Error("Error reading destination from SAM bridge")
			return nil, err
		}
		if n < 1 {
			log.Error("No destination data received from SAM bridge")
			return nil, fmt.Errorf("no destination data received")
		}
		pub := strings.Split(strings.Split(string(buf[:n]), "PRIV=")[0], "PUB=")[1]
		_priv := strings.Split(string(buf[:n]), "PRIV=")[1]

		priv := removeNewlines(_priv) //There is an extraneous newline in the private key, so we'll remove it.

		log.WithFields(logrus.Fields{
			"_priv(pre-newline removal)": _priv,
			"priv":                       priv,
		}).Info("Removed newline")

		log.Debug("Successfully created new destination")

		return &I2PKeys{
			Address: I2PAddr(pub),
			Both:    pub + priv,
		}, nil

	}
	log.Error("No RESULT=OK received from SAM bridge")
	return nil, fmt.Errorf("no result received")
}
