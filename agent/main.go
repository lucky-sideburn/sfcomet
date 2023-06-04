package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	vault "github.com/hashicorp/vault-client-go"
)

type Comet struct {
	path             []string
	dynamic_sentinel bool
	fencing          string
}

func Keys(m map[string]interface{}) (keys []string) {
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func writeOutputToFile(output string, outputFile string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	fmt.Fprintln(w, output)
	return w.Flush()
}

func buildCometList(role string, ctx context.Context, client vault.Client) {

	var currentComet Comet

	log.Println("Add comet configuration for role", role)

	secret, err := client.Secrets.KVv2Read(
		ctx,
		"roles/"+role,
		vault.WithMountPath("sfcomet"),
	)

	if err != nil {
		log.Fatal(err)
	}

	secretData := secret.Data["data"]
	secretMap := secretData.(map[string]interface{})
	currentComet.fencing = secretMap["fencing"].(string)

	for _, secretMapInterface := range secretMap["path"].([]interface{}) {
		currentComet.path = append(currentComet.path, secretMapInterface.(string))
	}

	cometList = append(cometList, currentComet)
}

var cometList []Comet

func main() {
	ctx := context.Background()
	log.Println("Starting sfagent agent")
	vaultTokenPtr := flag.String("token", "{{ vault_token }}", "Vault Token")
	vaultAddressPtr := flag.String("address", "{{ vault_url }}", "Vault Address")
	vaultEnableTls := flag.Bool("use-tls", true, "Use TLS (example: -use-tls=false)")
	vaultTlsCaPath := flag.String("ca-file", "{{ agent_vault_ca_path }}", "CA cert path (example: -ca-file=./rootCA.crt)")
	roles := flag.String("roles", "{{ agent_active_roles }}", "Define the roles of the node  (example: -roles=default or -roles=default,database)")

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
		log.Println("Creating Vault client with TLS disabled")
		client, err = vault.New(
			vault.WithAddress(*vaultAddressPtr),
		)
	}

	if err != nil {
		log.Fatalf("Error: Unable to initialize Vault client: %v", err)
	} else {
		log.Println("Vault client created correctly")
	}

	if *vaultTokenPtr != "" {
		client.SetToken(*vaultTokenPtr)
	} else {
		log.Fatalf("Error: Vault token has not been set. This is for now the only way to authenticate the agent. To do use approle")
	}

	listRoles := strings.Split(*roles, ",")

	for _, role := range listRoles {
		log.Println("Found role", role)
		buildCometList(role, ctx, *client)
	}

	log.Println(cometList)

	for _, comet := range cometList {
		for _, fileSentinel := range comet.path {
			log.Println("Creating sentinel file:", fileSentinel)
			writeOutputToFile("foobar", fileSentinel)
			var p [8]byte
			log.Println(rand.Read(p[:]))
		}
	}

}
