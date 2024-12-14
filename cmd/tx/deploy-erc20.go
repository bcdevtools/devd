package tx

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/bcdevtools/devd/v2/cmd/utils"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/cobra"
)

func GetDeployErc20EvmTxCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy-erc20",
		Short: "Deploy an ERC20 contract",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			ethClient8545, _ := mustGetEthClient(cmd)

			_, ecdsaPrivateKey, _, from := mustSecretEvmAccount(cmd)

			nonce, err := ethClient8545.NonceAt(context.Background(), *from, nil)
			utils.ExitOnErr(err, "failed to get nonce of sender")

			chainId, err := ethClient8545.ChainID(context.Background())
			utils.ExitOnErr(err, "failed to get chain ID")

			deploymentBytes, err := hex.DecodeString("60806040523480156200001157600080fd5b506040518060400160405280600381526020017f45324500000000000000000000000000000000000000000000000000000000008152506040518060400160405280600381526020017f453245000000000000000000000000000000000000000000000000000000000081525081600390816200008f9190620004ac565b508060049081620000a19190620004ac565b505050620000c3336c01431e0fae6d7217caa0000000620000c960201b60201c565b620006ae565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036200013b576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016200013290620005f4565b60405180910390fd5b6200014f600083836200022d60201b60201c565b806002600082825462000163919062000645565b92505081905550806000808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254620001ba919062000645565b925050819055508173ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8360405162000221919062000691565b60405180910390a35050565b505050565b600081519050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60006002820490506001821680620002b457607f821691505b602082108103620002ca57620002c96200026c565b5b50919050565b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b600060088302620003347fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82620002f5565b620003408683620002f5565b95508019841693508086168417925050509392505050565b6000819050919050565b6000819050919050565b60006200038d62000387620003818462000358565b62000362565b62000358565b9050919050565b6000819050919050565b620003a9836200036c565b620003c1620003b88262000394565b84845462000302565b825550505050565b600090565b620003d8620003c9565b620003e58184846200039e565b505050565b5b818110156200040d5762000401600082620003ce565b600181019050620003eb565b5050565b601f8211156200045c576200042681620002d0565b6200043184620002e5565b8101602085101562000441578190505b620004596200045085620002e5565b830182620003ea565b50505b505050565b600082821c905092915050565b6000620004816000198460080262000461565b1980831691505092915050565b60006200049c83836200046e565b9150826002028217905092915050565b620004b78262000232565b67ffffffffffffffff811115620004d357620004d26200023d565b5b620004df82546200029b565b620004ec82828562000411565b600060209050601f8311600181146200052457600084156200050f578287015190505b6200051b85826200048e565b8655506200058b565b601f1984166200053486620002d0565b60005b828110156200055e5784890151825560018201915060208501945060208101905062000537565b868310156200057e57848901516200057a601f8916826200046e565b8355505b6001600288020188555050505b505050505050565b600082825260208201905092915050565b7f45524332303a206d696e7420746f20746865207a65726f206164647265737300600082015250565b6000620005dc601f8362000593565b9150620005e982620005a4565b602082019050919050565b600060208201905081810360008301526200060f81620005cd565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000620006528262000358565b91506200065f8362000358565b92508282019050808211156200067a576200067962000616565b5b92915050565b6200068b8162000358565b82525050565b6000602082019050620006a8600083018462000680565b92915050565b611a8f80620006be6000396000f3fe6080604052600436106100f75760003560e01c806379cc67901161008a578063a457c2d711610059578063a457c2d714610340578063a9059cbb1461037d578063d3f24791146103ba578063dd62ed3e146103e5576100fe565b806379cc67901461029e57806385bd4767146102c757806389b8217a146102f757806395d89b4114610315576100fe565b8063313ce567116100c6578063313ce567146101d057806339509351146101fb57806342966c681461023857806370a0823114610261576100fe565b806306fdde0314610100578063095ea7b31461012b57806318160ddd1461016857806323b872dd14610193576100fe565b366100fe57005b005b34801561010c57600080fd5b50610115610422565b604051610122919061112f565b60405180910390f35b34801561013757600080fd5b50610152600480360381019061014d91906111ea565b6104b4565b60405161015f9190611245565b60405180910390f35b34801561017457600080fd5b5061017d6104d2565b60405161018a919061126f565b60405180910390f35b34801561019f57600080fd5b506101ba60048036038101906101b5919061128a565b6104dc565b6040516101c79190611245565b60405180910390f35b3480156101dc57600080fd5b506101e56105dd565b6040516101f291906112f9565b60405180910390f35b34801561020757600080fd5b50610222600480360381019061021d91906111ea565b6105e6565b60405161022f9190611245565b60405180910390f35b34801561024457600080fd5b5061025f600480360381019061025a9190611314565b610692565b005b34801561026d57600080fd5b5061028860048036038101906102839190611341565b6106a6565b604051610295919061126f565b60405180910390f35b3480156102aa57600080fd5b506102c560048036038101906102c091906111ea565b6106ee565b005b6102e160048036038101906102dc9190611314565b610772565b6040516102ee9190611245565b60405180910390f35b6102ff610789565b60405161030c9190611245565b60405180910390f35b34801561032157600080fd5b5061032a6107fd565b604051610337919061112f565b60405180910390f35b34801561034c57600080fd5b50610367600480360381019061036291906111ea565b61088f565b6040516103749190611245565b60405180910390f35b34801561038957600080fd5b506103a4600480360381019061039f91906111ea565b610983565b6040516103b19190611245565b60405180910390f35b3480156103c657600080fd5b506103cf6109a1565b6040516103dc9190611245565b60405180910390f35b3480156103f157600080fd5b5061040c6004803603810190610407919061136e565b6109f2565b604051610419919061126f565b60405180910390f35b606060038054610431906113dd565b80601f016020809104026020016040519081016040528092919081815260200182805461045d906113dd565b80156104aa5780601f1061047f576101008083540402835291602001916104aa565b820191906000526020600020905b81548152906001019060200180831161048d57829003601f168201915b5050505050905090565b60006104c86104c1610a79565b8484610a81565b6001905092915050565b6000600254905090565b60006104e9848484610c4a565b6000600160008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000610534610a79565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050828110156105b4576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105ab90611480565b60405180910390fd5b6105d1856105c0610a79565b85846105cc91906114cf565b610a81565b60019150509392505050565b60006012905090565b60006106886105f3610a79565b848460016000610601610a79565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020546106839190611503565b610a81565b6001905092915050565b6106a361069d610a79565b82610ec7565b50565b60008060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b6000610701836106fc610a79565b6109f2565b905081811015610746576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161073d906115a9565b60405180910390fd5b61076383610752610a79565b848461075e91906114cf565b610a81565b61076d8383610ec7565b505050565b600081341461078057600080fd5b60019050919050565b6000600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166108fc61270f9081150290604051600060405180830381858888f193505050501580156107f5573d6000803e3d6000fd5b506001905090565b60606004805461080c906113dd565b80601f0160208091040260200160405190810160405280929190818152602001828054610838906113dd565b80156108855780601f1061085a57610100808354040283529160200191610885565b820191906000526020600020905b81548152906001019060200180831161086857829003601f168201915b5050505050905090565b6000806001600061089e610a79565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205490508281101561095b576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016109529061163b565b60405180910390fd5b610978610966610a79565b85858461097391906114cf565b610a81565b600191505092915050565b6000610997610990610a79565b8484610c4a565b6001905092915050565b60006109ab610a79565b600560006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506001905090565b6000600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b600033905090565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610af0576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610ae7906116cd565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610b5f576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610b569061175f565b60405180910390fd5b80600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92583604051610c3d919061126f565b60405180910390a3505050565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610cb9576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610cb0906117f1565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610d28576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d1f90611883565b60405180910390fd5b610d3383838361109a565b60008060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905081811015610db9576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610db090611915565b60405180910390fd5b8181610dc591906114cf565b6000808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550816000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254610e559190611503565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051610eb9919061126f565b60405180910390a350505050565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610f36576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610f2d906119a7565b60405180910390fd5b610f428260008361109a565b60008060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905081811015610fc8576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610fbf90611a39565b60405180910390fd5b8181610fd491906114cf565b6000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550816002600082825461102891906114cf565b92505081905550600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8460405161108d919061126f565b60405180910390a3505050565b505050565b600081519050919050565b600082825260208201905092915050565b60005b838110156110d95780820151818401526020810190506110be565b60008484015250505050565b6000601f19601f8301169050919050565b60006111018261109f565b61110b81856110aa565b935061111b8185602086016110bb565b611124816110e5565b840191505092915050565b6000602082019050818103600083015261114981846110f6565b905092915050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061118182611156565b9050919050565b61119181611176565b811461119c57600080fd5b50565b6000813590506111ae81611188565b92915050565b6000819050919050565b6111c7816111b4565b81146111d257600080fd5b50565b6000813590506111e4816111be565b92915050565b6000806040838503121561120157611200611151565b5b600061120f8582860161119f565b9250506020611220858286016111d5565b9150509250929050565b60008115159050919050565b61123f8161122a565b82525050565b600060208201905061125a6000830184611236565b92915050565b611269816111b4565b82525050565b60006020820190506112846000830184611260565b92915050565b6000806000606084860312156112a3576112a2611151565b5b60006112b18682870161119f565b93505060206112c28682870161119f565b92505060406112d3868287016111d5565b9150509250925092565b600060ff82169050919050565b6112f3816112dd565b82525050565b600060208201905061130e60008301846112ea565b92915050565b60006020828403121561132a57611329611151565b5b6000611338848285016111d5565b91505092915050565b60006020828403121561135757611356611151565b5b60006113658482850161119f565b91505092915050565b6000806040838503121561138557611384611151565b5b60006113938582860161119f565b92505060206113a48582860161119f565b9150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b600060028204905060018216806113f557607f821691505b602082108103611408576114076113ae565b5b50919050565b7f45524332303a207472616e7366657220616d6f756e742065786365656473206160008201527f6c6c6f77616e6365000000000000000000000000000000000000000000000000602082015250565b600061146a6028836110aa565b91506114758261140e565b604082019050919050565b600060208201905081810360008301526114998161145d565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006114da826111b4565b91506114e5836111b4565b92508282039050818111156114fd576114fc6114a0565b5b92915050565b600061150e826111b4565b9150611519836111b4565b9250828201905080821115611531576115306114a0565b5b92915050565b7f45524332303a206275726e20616d6f756e74206578636565647320616c6c6f7760008201527f616e636500000000000000000000000000000000000000000000000000000000602082015250565b60006115936024836110aa565b915061159e82611537565b604082019050919050565b600060208201905081810360008301526115c281611586565b9050919050565b7f45524332303a2064656372656173656420616c6c6f77616e63652062656c6f7760008201527f207a65726f000000000000000000000000000000000000000000000000000000602082015250565b60006116256025836110aa565b9150611630826115c9565b604082019050919050565b6000602082019050818103600083015261165481611618565b9050919050565b7f45524332303a20617070726f76652066726f6d20746865207a65726f2061646460008201527f7265737300000000000000000000000000000000000000000000000000000000602082015250565b60006116b76024836110aa565b91506116c28261165b565b604082019050919050565b600060208201905081810360008301526116e6816116aa565b9050919050565b7f45524332303a20617070726f766520746f20746865207a65726f20616464726560008201527f7373000000000000000000000000000000000000000000000000000000000000602082015250565b60006117496022836110aa565b9150611754826116ed565b604082019050919050565b600060208201905081810360008301526117788161173c565b9050919050565b7f45524332303a207472616e736665722066726f6d20746865207a65726f20616460008201527f6472657373000000000000000000000000000000000000000000000000000000602082015250565b60006117db6025836110aa565b91506117e68261177f565b604082019050919050565b6000602082019050818103600083015261180a816117ce565b9050919050565b7f45524332303a207472616e7366657220746f20746865207a65726f206164647260008201527f6573730000000000000000000000000000000000000000000000000000000000602082015250565b600061186d6023836110aa565b915061187882611811565b604082019050919050565b6000602082019050818103600083015261189c81611860565b9050919050565b7f45524332303a207472616e7366657220616d6f756e742065786365656473206260008201527f616c616e63650000000000000000000000000000000000000000000000000000602082015250565b60006118ff6026836110aa565b915061190a826118a3565b604082019050919050565b6000602082019050818103600083015261192e816118f2565b9050919050565b7f45524332303a206275726e2066726f6d20746865207a65726f2061646472657360008201527f7300000000000000000000000000000000000000000000000000000000000000602082015250565b60006119916021836110aa565b915061199c82611935565b604082019050919050565b600060208201905081810360008301526119c081611984565b9050919050565b7f45524332303a206275726e20616d6f756e7420657863656564732062616c616e60008201527f6365000000000000000000000000000000000000000000000000000000000000602082015250565b6000611a236022836110aa565b9150611a2e826119c7565b604082019050919050565b60006020820190508181036000830152611a5281611a16565b905091905056fea26469706673582212207360cfb0416db54607945978cde44d485fc600313ef741e92ad87f7161c82b0364736f6c63430008140033")

			txData := ethtypes.LegacyTx{
				Nonce:    nonce,
				GasPrice: big.NewInt(20_000_000_000),
				Gas:      2_000_000,
				To:       nil,
				Data:     deploymentBytes,
				Value:    common.Big0,
			}
			tx := ethtypes.NewTx(&txData)

			newContractAddress := crypto.CreateAddress(*from, nonce)

			fmt.Println("Deploying new ERC-20 contract using", *from)

			signedTx, err := ethtypes.SignTx(tx, ethtypes.LatestSignerForChainID(chainId), ecdsaPrivateKey)
			utils.ExitOnErr(err, "failed to sign tx")

			var buf bytes.Buffer
			err = signedTx.EncodeRLP(&buf)
			utils.ExitOnErr(err, "failed to encode tx")

			rawTxRLPHex := hex.EncodeToString(buf.Bytes())
			fmt.Printf("RawTx: 0x%s\n", rawTxRLPHex)

			err = ethClient8545.SendTransaction(context.Background(), signedTx)
			utils.ExitOnErr(err, "failed to send tx")

			fmt.Println("New contract deployed at")
			fmt.Println(newContractAddress)
		},
	}

	cmd.Flags().String(flagRpc, "", flagEvmRpcDesc)
	cmd.Flags().String(flagSecretKey, "", flagSecretKeyDesc)

	return cmd
}