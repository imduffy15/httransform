package ca

import (
	"crypto/x509"
	"testing"

	"github.com/stretchr/testify/suite"
)

var testCaCACert = []byte(`-----BEGIN CERTIFICATE-----
MIICWzCCAcSgAwIBAgIJAJ34yk7oiKv5MA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV
BAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBX
aWRnaXRzIFB0eSBMdGQwHhcNMTgxMjAyMTQyNTAyWhcNMjgxMTI5MTQyNTAyWjBF
MQswCQYDVQQGEwJBVTETMBEGA1UECAwKU29tZS1TdGF0ZTEhMB8GA1UECgwYSW50
ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKB
gQDL7Hzfmx7xfFWTRm26t/lLsCZwOri6VIzp2dYM5Hp0dV4XUZ+q60nEbHwN3Usr
GKAK/Rsr9Caam3A18Upn2ly69Tyr29kVK+PlsOgSSCUnAYcqT166/j205n3CGNLL
OPtQKfAT/iH3dPBObd8N4FR9FlXiYIiAp1opCbyu2mlHiwIDAQABo1MwUTAdBgNV
HQ4EFgQUOJ+uGtIhHxXHPNESBNI4YbwAl+wwHwYDVR0jBBgwFoAUOJ+uGtIhHxXH
PNESBNI4YbwAl+wwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOBgQCW
s7P0wJ8ON8ieEJe4pAfACpL6IyhZ5YK/C/hip+czxdvZHc5zngVwHP2vsIcHKBTr
8qXoHgh2gaXqwn8kRVNnZzWrxgSe8IR3oJ2yTbLAxqDS42SPfRLAUpy9sK/tEEGM
rMk/LWMzH/S6bLcsAm0GfVIrUNfg0eF0ZVIjxINBVA==
-----END CERTIFICATE-----`)

var testCaPrivateKey = []byte(`-----BEGIN PRIVATE KEY-----
MIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBAMvsfN+bHvF8VZNG
bbq3+UuwJnA6uLpUjOnZ1gzkenR1XhdRn6rrScRsfA3dSysYoAr9Gyv0JpqbcDXx
SmfaXLr1PKvb2RUr4+Ww6BJIJScBhypPXrr+PbTmfcIY0ss4+1Ap8BP+Ifd08E5t
3w3gVH0WVeJgiICnWikJvK7aaUeLAgMBAAECgYAk+/kR3OJZzcD/evB/wsoV7haq
mBvUv2znJLjrkayb3oV4GTeqGg5A76P4J8BwSoEMPSdma1ttAu/w+JgUCchzVPwU
34Sr80mYawOmGVGJsDnrrYA2w51Nj42e71pmRc9IqNLwFEhW5Uy7eASf3THJMWDl
F2M6xAVYr+X0eKLf4QJBAO8lVIIMnzIReSZukWBPp6GKmXOuEkWeBOfnYC2HOVZq
1M/E6naOP2MBk9CWG4o9ysjcZ1hosi3/txxrc8VmBAkCQQDaS651dpQ3TRE//raZ
s79ZBEdMCMlgXB6CPrZpvLz/3ZPcLih4MJ59oVkeFHCNct7ccQcQu4XHMGNBIRBh
kpvzAkEAlS/AjHC7T0y/O052upJ2jLweBqBtHaj6foFE6qIVDugOYp8BdXw/5s+x
GsrJ22+49Z0pi2mk3jVMUhpmWprNoQJBANdAT0v2XFpXfQ38bTQMYT82j9Myytdg
npjRm++Rs1AdvoIbZb52OqIoqoaVoxJnVchLD6t5LYXnecesAcok1e8CQEKB7ycJ
6yVwnBE3Ua9CHcGmrre6HmEWdPy1Zyb5DQC6duX46zEBzti9oWx0DJIQRZifeCvw
4J45NsSQjuuAAWs=
-----END PRIVATE KEY-----`)

type CATestSuite struct {
	suite.Suite

	ca *CA
}

func (suite *CATestSuite) SetupTest() {
	ca, err := NewCA(testCaCACert, testCaPrivateKey, 1000, 100)
	if err != nil {
		panic(err)
	}
	suite.ca = &ca
}

func (suite *CATestSuite) TearDownTest() {
	suite.ca.Close()
}

func (suite *CATestSuite) TestDoubleGet() {
	conf1, err := suite.ca.Get("hostname")
	suite.Nil(err)

	conf2, err := suite.ca.Get("hostname")
	suite.Nil(err)

	suite.Equal(conf1.Get().Certificates[0].PrivateKey, conf2.Get().Certificates[0].PrivateKey)
	suite.Equal(conf1.Get().Certificates[0].Certificate[0], conf2.Get().Certificates[0].Certificate[0])
}

func (suite *CATestSuite) TestSigner() {
	conf, err := suite.ca.Get("hostname")
	suite.Nil(err)

	cert := conf.Get().Certificates[0]
	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	suite.Nil(err)

	pool := x509.NewCertPool()
	pool.AddCert(suite.ca.ca.Leaf)

	suite.Nil(cert.Leaf.VerifyHostname("hostname"))
	suite.NotNil(cert.Leaf.VerifyHostname("hostname2"))
	suite.Nil(cert.Leaf.CheckSignatureFrom(suite.ca.ca.Leaf))

	_, err = cert.Leaf.Verify(x509.VerifyOptions{
		DNSName: "hostname",
		Roots:   pool,
	})
	suite.Nil(err)
}

func TestCA(t *testing.T) {
	suite.Run(t, &CATestSuite{})
}
