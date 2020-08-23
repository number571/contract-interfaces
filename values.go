package main

import (
	"context"
	"math/big"
	"io/ioutil"
	"crypto/ecdsa"
	contract "./contracts"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type UserType struct {
	Purse string
	AddressHex string
	AddressEth common.Address
	PublicKey *ecdsa.PublicKey
	PrivateKey *ecdsa.PrivateKey
}

type Estate struct{
    Id *big.Int
    Owner common.Address
    Info string
    Squere *big.Int
    UsefulSquere *big.Int
    RenterAddress common.Address
    PresentStatus bool
    SaleStatus bool
    RentStatus bool
}

type Present struct {
	Id *big.Int
	EstateId *big.Int
	AddressFrom common.Address
	AddressTo common.Address
	Finished bool
}

var (
	User *UserType
	ClientETH     = connectToETH("http://127.0.0.1:7545") 
	Instance      = connectToContract(
		common.HexToAddress(readFile("contract.address")), 
		ClientETH,
	)
)

func loadUser(purse string) *UserType {
	priv, err := crypto.HexToECDSA(purse)
	if err != nil {
		return nil
	}
	pub, ok := priv.Public().(*ecdsa.PublicKey)
	if !ok {
		return nil
	}
	addressHex := crypto.PubkeyToAddress(*pub).Hex()
	addressEth  := common.HexToAddress(addressHex)
	return &UserType{
		Purse:      purse,
		AddressHex: addressHex,
		AddressEth: addressEth,
		PublicKey:  pub,
		PrivateKey: priv,
	}
}

func connectToContract(contractAddr common.Address, clientEth *ethclient.Client) *contract.Contract {
	instance, err := contract.NewContract(contractAddr, clientEth)
	if err != nil {
		return nil
	}
	return instance
}

func connectToETH(address string) *ethclient.Client {
	client, err := ethclient.Dial(address)
	if err != nil {
		return nil
	}
	return client
}

func readFile(filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return ""
	}
	return string(data)
}

func resetAuth(user *UserType) *bind.TransactOpts {
	nonce, err := ClientETH.PendingNonceAt(context.Background(), crypto.PubkeyToAddress(*user.PublicKey))
	if err != nil {
		return nil
	}

	gasPrice, err := ClientETH.SuggestGasPrice(context.Background())
	if err != nil {
		return nil
	}

	auth := bind.NewKeyedTransactor(user.PrivateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)

	auth.GasLimit = uint64(3000000)
	auth.GasPrice = gasPrice

	return auth
}

func getEstates(index *big.Int) *Estate {
	// (*big.Int, common.Address, string, *big.Int, *big.Int, common.Address, error)
	id, owner, info, squere, usefulsquere, renteraddress, err := Instance.GetEstates(&bind.CallOpts{From: User.AddressEth}, index)
	if err != nil {
		return nil
	}
	presentS, saleS, rentS, err := Instance.GetEstatesStatuses(&bind.CallOpts{From: User.AddressEth}, index)
	if err != nil {
		return nil
	}
	return &Estate{
		Id: id,
		Owner: owner,
		Info: info,
		Squere: squere,
		UsefulSquere: usefulsquere,
		RenterAddress: renteraddress,
		PresentStatus: presentS,
		SaleStatus: saleS,
		RentStatus: rentS,
	}
}

func getPresents(index *big.Int) *Present {
	// (*big.Int, common.Address, common.Address, bool, error)
	id, from, to, finished, err := Instance.GetPresents(&bind.CallOpts{From: User.AddressEth}, index)
	if err != nil {
		return nil
	}
	return &Present{
		Id: index,
		EstateId: id,
		AddressFrom: from,
		AddressTo: to,
		Finished: finished,
	}
}

type EstateStr struct {
	Id *big.Int
    Owner string
    Info string
    Squere *big.Int
    UsefulSquere *big.Int
    RenterAddress string
    PresentStatus bool
    SaleStatus bool
    RentStatus bool
}

func estatesToString(estate *Estate) *EstateStr {
	return &EstateStr{
		Id: estate.Id,
		Owner: estate.Owner.Hex(),
		Info: estate.Info,
		Squere: estate.Squere,
		UsefulSquere: estate.UsefulSquere,
		RenterAddress: estate.RenterAddress.Hex(),
		PresentStatus: estate.PresentStatus,
		SaleStatus: estate.SaleStatus,
		RentStatus: estate.RentStatus,
	}
}

type PresentStr struct {
	Id *big.Int
	EstateId *big.Int
	AddressFrom string
	AddressTo string
	Finished bool
}

func presentsToString(present *Present) *PresentStr {
	return &PresentStr{
		Id: present.Id,
		EstateId: present.EstateId,
		AddressFrom: present.AddressFrom.Hex(),
		AddressTo: present.AddressTo.Hex(),
		Finished: present.Finished,
	}
}
