package main

import (
	"os"
	"fmt"
	"bufio"
	"context"
	"strings"
	"math/big"
	"encoding/json"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func init() {
	if len(os.Args) < 2 {
		panic("failed: len(os.Args) < 2")
	}
	var (
		userLoadStr = ""
		userLoadExist = false
	)
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		switch {
		case strings.HasPrefix(arg, "-loaduser:"):
			userLoadStr = strings.Replace(arg, "-loaduser:", "", 1)
			userLoadExist = true
		}
	}
	if !userLoadExist {
		panic("failed: !userLoadExist")
	}
	if ClientETH == nil {
		panic("failed: connect to ETH")
	}
	if Instance == nil {
		panic("failed: instance is nil")
	}
	User = loadUser(userLoadStr)
	if User == nil {
		panic("failed: load user")
	}
}

func main() {
	var (
		message string
		splited []string
	)
	for {
		message = inputString("> ")
		splited = strings.Split(message, " ")
		switch splited[0] {
		case "/exit":
			os.Exit(0)
		case "/user":
			if len(splited) < 2 {
				fmt.Println("failed: len(user) < 2\n")
				continue
			}
			switch splited[1] {
			case "address":
				userAddress()
			case "purse":
				userPurse()
			case "balance":
				userBalance()
			default:
				fmt.Println("command undefined\n")
			}
		case "/chain":
			if len(splited) < 3 {
				fmt.Println("failed: len(chain) < 3\n")
				continue
			}
			switch splited[1] {
			case "get":
				chainGet(splited[2], splited[2:])
			case "create":
				switch splited[2] {
				case "estate":
					// chain create estate address info squere usefulSquere
					chainCreateEstate(splited[2:])
				case "present":
					// chain create present id_estate address
					chainCreatePresent(splited[2:])
				default:
					fmt.Println("command undefined\n")
				}
			case "cancel":
				switch splited[2] {
				case "present":
					// chain cancel present id_present
					chainCancelPresent(splited[2], splited[2:])
				default:
					fmt.Println("command undefined\n")
				}
			case "confirm":
				switch splited[2] {
				case "present":
					// chain confirm present id_present
					chainConfirmPresent(splited[2:])
				default:
					fmt.Println("command undefined\n")
				}
			default:
				fmt.Println("command undefined\n")
			}
		default:
			fmt.Println("command undefined\n")
		}
	}
}

func chainCreateEstate(splited []string) {
	if len(splited) != 5 {
		fmt.Println("failed: len(splited) != 5\n")
		return
	}
	var (
		squere = new(big.Int)
		usefulSquere = new(big.Int)
		ok bool
	)
	squere, ok = squere.SetString(splited[3], 10)
	if !ok {
		fmt.Println("failed: conv(str1) to num\n")
		return
	}
	usefulSquere, ok = usefulSquere.SetString(splited[4], 10)
	if !ok {
		fmt.Println("failed: conv(str2) to num\n")
		return
	}
	var address common.Address
	if splited[1] == "my" {
		address = User.AddressEth
	} else {
		address = common.HexToAddress(splited[1])
	}
	tx, err := Instance.CreateEstate(
		resetAuth(User), 
		address, 
		splited[2],
		squere,
		usefulSquere,
	)
	if err != nil {
		fmt.Println(err, "\n")
		return
	}
	fmt.Println("Tx:", tx.Hash().Hex(), "\n")
}

func chainCreatePresent(splited []string) {
	if len(splited) != 3 {
		fmt.Println("failed: len(splited) != 3\n")
		return
	}
	var (
		estateId = new(big.Int)
		ok bool
	)
	estateId, ok = estateId.SetString(splited[1], 10)
	if !ok {
		fmt.Println("failed: conv(str1) to num\n")
		return
	}
	tx, err := Instance.CreatePresent(
		resetAuth(User), 
		estateId, 
		common.HexToAddress(splited[2]),
	)
	if err != nil {
		fmt.Println(err, "\n")
		return
	}
	fmt.Println("Tx:", tx.Hash().Hex(), "\n")
}

func chainCancelPresent(category string, splited []string) {
	if len(splited) != 2 {
		fmt.Println("failed: len(splited) != 2\n")
		return
	}
	var (
		tx *types.Transaction
		err error
		num = new(big.Int)
		ok bool
	)
	num, ok = num.SetString(splited[1], 10)
	if !ok {
		fmt.Println("failed: conv(str1) to num\n")
		return
	}
	switch category {
	case "present":
		tx, err = Instance.CancelPresent(
			resetAuth(User), 
			num,
		)
	case "sale":
		tx, err = Instance.CancelSale(
			resetAuth(User), 
			num,
		)
	case "rent":
		tx, err = Instance.CancelRent(
			resetAuth(User), 
			num,
		)
	}
	if err != nil {
		fmt.Println(err, "\n")
		return
	}
	fmt.Println("Tx:", tx.Hash().Hex(), "\n")
}

func chainConfirmPresent(splited []string) {
	if len(splited) != 2 {
		fmt.Println("failed: len(splited) != 2\n")
		return
	}
	var (
		presentNumber = new(big.Int)
		ok bool
	)
	presentNumber, ok = presentNumber.SetString(splited[1], 10)
	if !ok {
		fmt.Println("failed: conv(str1) to num\n")
		return
	}
	tx, err := Instance.ConfirmPresent(
		resetAuth(User), 
		presentNumber,
	)
	if err != nil {
		fmt.Println(err, "\n")
		return
	}
	fmt.Println("Tx:", tx.Hash().Hex(), "\n")
}

func chainGet(category string, splited []string) {
	if len(splited) != 2 {
		fmt.Println("failed: len(splited) != 2\n")
		return
	}
	var (
		inc = big.NewInt(1)
		err error
		num *big.Int
		jsonData []byte
	)
	switch category {
	case "estates":
		num, err = Instance.GetEstatesNumber(&bind.CallOpts{From: User.AddressEth})
	case "presents":
		num, err = Instance.GetPresentsNumber(&bind.CallOpts{From: User.AddressEth})
	default:
		fmt.Println("undefined category\n")
		return
	}
	if err != nil {
		fmt.Println(err, "\n")
		return
	}
	for index := big.NewInt(0); index.Cmp(num) == -1; index.Add(index, inc) {
		switch category {
		case "estates":
			data := getEstates(index)
			if data == nil {
				fmt.Println("data is nil\n")
				return
			}
			if splited[1] == "my" && User.AddressHex != data.Owner.Hex() {
				continue
			}
			if splited[1] != "all" && splited[1] != "my" && 
				strings.ToLower(splited[1]) != strings.ToLower(data.Owner.Hex()) {
				continue
			}
			jsonData, err = json.MarshalIndent(data, "", "\t")
		case "presents":
			data := getPresents(index)
			if data == nil {
				fmt.Println("data is nil\n")
				return
			}
			if data.Finished {
				continue
			}
			if splited[1] == "my" && 
				(User.AddressHex != data.AddressFrom.Hex() && User.AddressHex != data.AddressTo.Hex()){
				continue
			}
			if splited[1] != "all" && splited[1] != "my" && 
				(strings.ToLower(splited[1]) != strings.ToLower(data.AddressFrom.Hex()) && 
				strings.ToLower(splited[1]) != strings.ToLower(data.AddressTo.Hex())) {
				continue
			}
			jsonData, err = json.MarshalIndent(data, "", "\t")
		default:
			fmt.Println("undefined category\n")
			return
		}
		if err != nil {
			fmt.Println(err, "\n")
			return
		}
		fmt.Println(string(jsonData))
	}
	fmt.Println()
}

func userAddress() {
	fmt.Println("Address:", User.AddressHex, "\n")
}

func userPurse() {
	fmt.Println("Purse:", User.Purse, "\n")
}

func userBalance() {
	balance, err := ClientETH.BalanceAt(context.Background(), User.AddressEth, nil)
	if err != nil {
		fmt.Println(err, "\n")
		return
	}
	fmt.Println("Balance:", balance, "\n")
}

func inputString(begin string) string {
	fmt.Print(begin)
	msg, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	return strings.Replace(msg, "\n", "", 1)
}
