package framework

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/crypto/ssh"
)

func ProcessTemplateWithParameters(srcFilePath, dstFilePath string, params ...string) string {
	By(fmt.Sprintf("Overriding the template from %s to %s", srcFilePath, dstFilePath))
	out := execute(Result{cmd: "oc", verb: "process", filePath: srcFilePath, params: params})
	filePath, err := writeJson(dstFilePath, out)
	Expect(err).ToNot(HaveOccurred())
	return filePath
}

func CreateResourceWithFilePathTestNamespace(filePath string) {
	By("Creating resource from the json file with the oc-create command")
	execute(Result{cmd: "oc", verb: "create", filePath: filePath, nameSpace: NamespaceTestDefault})
}

func DeleteResourceWithLabelTestNamespace(resourceType, resourceLabel string) {
	By(fmt.Sprintf("Deleting %s:%s from the json file with the oc-delete command", resourceType, resourceLabel))
	execute(Result{cmd: "oc", verb: "delete", resourceType: resourceType, resourceLabel: resourceLabel, nameSpace: NamespaceTestDefault})
}

func WaitUntilResourceReadyByNameTestNamespace(resourceType, resourceName, query, expectOut string) {
	By(fmt.Sprintf("Wait until %s with name %s ready", resourceType, resourceName))
	execute(Result{cmd: "oc", verb: "get", resourceType: resourceType, resourceName: resourceName, query: query, expectOut: expectOut, nameSpace: NamespaceTestDefault})
}

func WaitUntilResourceReadyByLabelTestNamespace(resourceType, label, query, expectOut string) {
	By(fmt.Sprintf("Wait until resource %s with label=%s ready", resourceType, label))
	execute(Result{cmd: "oc", verb: "get", resourceType: resourceType, resourceLabel: label, query: query, expectOut: expectOut, nameSpace: NamespaceTestDefault})
}

func writeJson(jsonFile string, json string) (string, error) {
	err := ioutil.WriteFile(jsonFile, []byte(json), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write the json file %s", jsonFile)
	}
	return jsonFile, nil
}

func RunOcDescribeCommand(resourceType, resourceName string) string {
	fmt.Printf("Getting 'oc describe' with: %s ", resourceName)
	return execute(Result{cmd: "oc", verb: "describe", resourceType: resourceType, resourceName: resourceName, nameSpace: NamespaceTestDefault})
}

// generatePrivateKey creates a RSA Private Key of specified byte size
func GeneratePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	// Private Key generation
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}

	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

// Will returns in the format "ssh-rsa ..."
func GeneratePublicKey(privatekey *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(privatekey)
	if err != nil {
		return nil, err
	}

	publicKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	return publicKeyBytes, nil
}

// encodePrivateKeyToPEM encodes Private Key from RSA to PEM format
func EncodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	// Get ASN.1 DER format
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)

	// pem.Block
	privateBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	// Private key in PEM format
	privatePEM := pem.EncodeToMemory(&privateBlock)

	return privatePEM
}
