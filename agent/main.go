package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

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

func writeOutputToFile(output []byte, outputFile string) error {
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

	for true {
		for _, comet := range cometList {
			log.Printf("Fencing Mechanism for current cometList is %s", comet.fencing)
			fencingCodeVaultPath := "/sfcomet/data/fencing/" + comet.fencing

			secret, err := client.Read(
				ctx,
				fencingCodeVaultPath,
			)
			if err != nil {
				log.Fatal(err)
			}
			secretData := secret.Data["data"]
			secretMap := secretData.(map[string]interface{})
			log.Printf("Base64 Encoded string of %s is %s", comet.fencing, secretMap["base64code"].(string))

			decodedFencingCode, err := base64.StdEncoding.DecodeString(secretMap["base64code"].(string))
			decodedFencingCodeString := string(decodedFencingCode[:])

			if err != nil {
				panic(err)
			} else {
				log.Printf("Successfully decoded fencingCodeVaultPath")
			}

			for _, fileSentinel := range comet.path {
				if _, err := os.Stat(fileSentinel); err == nil {
					log.Printf("File sentinel %s already exists", fileSentinel)
					log.Printf("Reading checksum of %s from Vault", fileSentinel)
					vaultPath := "/sfcomet/data/sentinels/" + fileSentinel

					secret, err := client.Read(
						ctx,
						vaultPath,
					)
					if err != nil {
						log.Fatal(err)
					}

					secretData := secret.Data["data"]
					secretMap := secretData.(map[string]interface{})
					log.Printf("Checksum from Vault of %s is |%s|", fileSentinel, secretMap["checksum"].(string))
					checksumFromVault := secretMap["checksum"].(string)

					f, err := os.Open(fileSentinel)
					if err != nil {
						log.Fatal(err)
					}
					defer f.Close()

					h := sha256.New()
					if _, err := io.Copy(h, f); err != nil {
						log.Fatal(err)
					}

					localChecksum := hex.EncodeToString(h.Sum(nil))
					log.Printf("Local checksum of %s is |%s|", fileSentinel, localChecksum)

					if localChecksum == checksumFromVault {
						log.Printf("Local and remote checksums of %s are equals", fileSentinel)
					} else {
						log.Printf("Local and remote checksums of %s are NOT equals", fileSentinel)
						cmd := exec.Command(decodedFencingCodeString)

						out, err := cmd.Output()
						if err != nil {
							// if there was any error, print it here
							fmt.Println("could not run command: ", err)
						}
						log.Printf("%s", out)
					}

				} else {
					token := make([]byte, 4)
					rand.Read(token)
					log.Println("Creating sentinel file:", fileSentinel)
					randBuff := make([]byte, 1024)
					rand.Read(randBuff)
					ioutil.WriteFile(fileSentinel, randBuff, 0666)
					sum := sha256.Sum256(randBuff)
					checksum_string := hex.EncodeToString(sum[:])
					log.Printf("Checksum of %s is %s", fileSentinel, checksum_string)

					vaultPath := "/sfcomet/data/sentinels/" + fileSentinel

					_, err = client.Write(ctx, vaultPath, map[string]any{
						"data": map[string]any{
							"checksum": checksum_string,
						},
					})

					if err != nil {
						log.Fatal(err)
					}
					log.Println("secret written successfully")
				}
			}
		}
		time.Sleep(time.Second)
	}
}
