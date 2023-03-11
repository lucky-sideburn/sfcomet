package main

import (
	"context"
	"flag"
	"log"

	vault "github.com/hashicorp/vault-client-go"
)

func main() {
	ctx := context.Background()
	log.Println("Starting Safecomet agent")
	vaultTokenPtr := flag.String("token", "", "Vault Token")
	vaultAddressPtr := flag.String("address", "http://127.0.0.1:8200", "Vault Address")
	vaultEnableTls := flag.Bool("tls-verify", true, "Verify TLS (example: -tls-verify=false)")
	tls := vault.TLSConfiguration{}
	tls.ServerCertificate.FromFile = "/tmp/vault-ca.pem"
	flag.Parse()
	client, err := vault.New()

	if *vaultEnableTls {
		client, err = vault.New(
			vault.WithAddress(*vaultAddressPtr),
			vault.WithTLS(tls),
		)
	} else {
		client, err = vault.New(
			vault.WithAddress(*vaultAddressPtr),
		)
	}

	if err != nil {
		log.Fatalf("Error: Unable to initialize Vault client: %v", err)
	} else {
		log.Println("Vault client created correctly")
	}

	if *vaultTokenPtr == "" {
		log.Fatalf("Error: Vault token has not bee set")
	}
	client.SetToken(*vaultTokenPtr)

	s, err := client.Secrets.KVv2Read(ctx, "my-secret")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("secret retrieved:", s.Data)

	// log.Println("Access granted!")
}
