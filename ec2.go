package main

// import (
// 	"io/ioutil"

// 	"github.com/aws/aws-sdk-go/aws"
// 	"github.com/aws/aws-sdk-go/aws/session"
// 	"github.com/aws/aws-sdk-go/service/ec2"
// )

// type importKeypairFromPublicMaterialInput struct {
// 	keyName           string
// 	publicKeyPath     string
// 	awsSessionOptions *session.Options
// }

// // func (k KeyPairStruct) getImportKeypairInput() ec2.ImportKeyPairInput {

// // 	return
// func importKeypairFromPublicMaterial(awsSessionOptions *session.Options, keypariPublicMaterialPath string) *ec2.ImportKeyPairOutput {
// 	publicKeyMaterial, err := ioutil.ReadFile(k.publicKeyPath)
// 	if err != nil {
// 		exitErrorf("Something went wrong when reading public key material: %v", err)
// 	}

// 	importKeypairInput := ec2.ImportKeyPairInput{KeyName: aws.String(k.keyName), PublicKeyMaterial: publicKeyMaterial}
// 	session := newSessionFromSessionOptions(*awsSessionOptions)
// }

// //KeyPairExists
// func (k KeyPairStruct) KeyPairExists() (result bool) {
// 	result = false
// 	return result
// }

// func (k KeyPairStruct) ImportKeypair() {
// 	if k.KeyPairExists() == false {
// 		importKeypairFromPublicMaterial(k.publicKeyPath)
// 	}
// }
