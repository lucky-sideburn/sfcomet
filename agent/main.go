package main

import (
	"context"
	"flag"
	"log"

	vault "github.com/hashicorp/vault-client-go"
)

func main() {
	ctx := context.Background()
	log.Println("Starting sfagent agent")
	vaultTokenPtr := flag.String("token", "", "Vault Token")
	vaultAddressPtr := flag.String("address", "http://127.0.0.1:8200", "Vault Address")
	vaultEnableTls := flag.Bool("use-tls", true, "Use TLS (example: -use-tls=false)")
	vaultTlsCaPath := flag.String("ca-file", "./rootCA.crt", "CA cert path (example: -ca-file=./rootCA.crt)")
	flag.Parse()

	tls := vault.TLSConfiguration{}
	tls.ServerCertificate.FromFile = *vaultTlsCaPath

	client, err := vault.New()

	if *vaultEnableTls {
		log.Println("Creating client with TLS enabled")
		log.Println("Use this CA cert:", *vaultTlsCaPath)

		client, err = vault.New(
			vault.WithAddress(*vaultAddressPtr),
			vault.WithTLS(tls),
		)
	} else {
		log.Println("Creating client with TLS disabled")
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

	s, err := client.Secrets.KVv2Read(ctx, "comets/default_nodes")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("secret retrieved:", s.Data)

	// log.Println("Access granted!")
}
