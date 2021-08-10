package cmds

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thetatoken/theta-eth-rpc-adaptor/common"
	"github.com/thetatoken/theta/cmd/thetacli/cmd/utils"
	tcommon "github.com/thetatoken/theta/common"
	"github.com/thetatoken/theta/wallet"
	ks "github.com/thetatoken/theta/wallet/softwallet/keystore"
	wtypes "github.com/thetatoken/theta/wallet/types"
)

var cfgPath string

const testAmount = 10
const testFile = "/testAddresses"
const testAccountPassword = "123"

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "theta-eth-rpc-adaptor",
	Short: "Theta ETH RPC Adaptor",
	Long:  `Theta ETH RPC Adaptor`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	checkWallets()
	RootCmd.PersistentFlags().StringVar(&cfgPath, "config", getDefaultConfigPath(), fmt.Sprintf("config path (default is %s)", getDefaultConfigPath()))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AddConfigPath(cfgPath)
	// Search config (without extension).
	viper.SetConfigName("config")

	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func getDefaultConfigPath() string {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return path.Join(home, ".theta-eth-rpc-adaptor")
}

func checkWallets() {
	keysDirPath := path.Join(getDefaultConfigPath(), "keys")
	log.Infof("Using keyDirPath: %v\n", keysDirPath)
	_, err := os.Stat(keysDirPath + "/testAddresses")
	if os.IsNotExist(err) { //firstTime
		err = createAccounts(keysDirPath, &common.TestWallets)
	}
	if err != nil {
		log.Errorf("failed to create test accounts, %v", err)
		return
	}
	if addresesBytes, err := ioutil.ReadFile(keysDirPath + testFile); err == nil {
		err = getAccounts(keysDirPath+"/keys", addresesBytes, &common.TestWallets)
		if err != nil {
			log.Errorf("failed to get test accounts, %v", err)
		}
	}
}

func createAccounts(keyPath string, ab *common.AddressBook) error {
	log.Infof("creating key at %s\n", keyPath)
	accountbytes := make([]byte, 20*testAmount)
	for i := 0; i < testAmount; i++ {
		wallet, err := wallet.OpenWallet(keyPath, wtypes.WalletTypeSoft, true)
		if err != nil {
			utils.Error("Failed to open wallet: %v\n", err)
		}
		address, err := wallet.NewKey(testAccountPassword)
		if err != nil {
			utils.Error("Failed to generate new key: %v\n", err)
		}
		copy(accountbytes[i*tcommon.AddressLength:(i+1)*tcommon.AddressLength], address[0:tcommon.AddressLength])
		fmt.Printf("Successfully created test key: %v\n", address.Hex())
	}
	ioutil.WriteFile(keyPath+testFile, accountbytes, 0777)
	return nil
}

func getAccounts(keyPath string, accountbytes []byte, ab *common.AddressBook) error {
	keystore, err := ks.NewKeystoreEncrypted(keyPath, ks.StandardScryptN, ks.StandardScryptP)
	if err != nil {
		log.Fatalf("Failed to open key store: %v", err)
	}
	addresses := make([]tcommon.Address, testAmount)
	for i := 0; i < testAmount; i++ {
		copy(addresses[i][0:tcommon.AddressLength], accountbytes[i*tcommon.AddressLength:(i+1)*tcommon.AddressLength])
		log.Infof("opening test wallet %s \n", addresses[i].Hex())
		nodeKey, err := keystore.GetKey(addresses[i], testAccountPassword)
		if err != nil {
			log.Errorf("Failed to open wallet, err is %v\n", err)
			return err
		}
		(*ab)[strings.ToLower(addresses[i].Hex())] = nodeKey.PrivateKey
	}
	return nil
}
