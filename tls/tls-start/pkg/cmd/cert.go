package cmd

import (
	"fmt"
	"os"

	"github.com/jbrcoleman/golang-devops/tls/tls-start/pkg/cert"
	"github.com/spf13/cobra"
)

var certKeyPath string
var certPath string
var certName string

func init() {
	createCmd.AddCommand(certCreateCmd)
	certCreateCmd.Flags().StringVarP(&certKeyPath, "key-out", "k", "server.key", "destination path for cert key")
	certCreateCmd.Flags().StringVarP(&certPath, "cert-out", "o", "server.cert", "destination path for cert")
	certCreateCmd.Flags().StringVarP(&certName, "name", "n", "", "name of cert in config file")
	certCreateCmd.Flags().StringVar(&caKey, "ca-key", "ca.key", "ca key to sign cert")
	certCreateCmd.Flags().StringVar(&caCert, "ca-cert", "ca.cert", "ca cert")
	certCreateCmd.MarkFlagRequired("ca-key")
	certCreateCmd.MarkFlagRequired("ca-cert")
	certCreateCmd.MarkFlagRequired("name")
}

var certCreateCmd = &cobra.Command{
	Use:   "cert",
	Short: "cert commands",
	Long:  `commands to create the certificates`,
	Run: func(cmd *cobra.Command, args []string) {
		caKeyBytes, err := os.ReadFile(caKey)
		if err != nil {
			fmt.Printf("CA key error: %s\n", err)
			return
		}
		caCertBytes, err := os.ReadFile(caCert)
		if err != nil {
			fmt.Printf("CA cert error: %s\n", err)
			return
		}
		certConfig, ok := config.Cert[certName]
		if !ok {
			fmt.Println("Could not find certificate name in config\n", config.Cert)
			return
		}
		err = cert.CreateCert(certConfig, caKeyBytes, caCertBytes, certKeyPath, certPath)
		if err != nil {
			fmt.Printf("Create cert error: %s\n", err)
			return
		}
		fmt.Printf("cert created. Key: %s, cert: %s", certKeyPath, certPath)
	},
}
