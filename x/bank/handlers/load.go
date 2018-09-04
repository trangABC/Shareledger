package handlers

import (
	sdk "bitbucket.org/shareringvn/cosmos-sdk/types"
	"fmt"

	"github.com/sharering/shareledger/x/auth"
	"github.com/sharering/shareledger/x/bank/messages"
)

//--------------------------------
// Handler for the message

func HandleMsgLoad(am auth.AccountMapper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		loadMsg, ok := msg.(messages.MsgLoad)

		// IMPORTANT
		// TODO: require a list of limited accounts which are priviledged to load coins
		signer := auth.GetSigner(ctx)
		fmt.Println("Signer: ", signer)

		if !ok {
			return sdk.NewError(2, 1, "MsgLoad is malformed").Result()
		}

		// Credit the account
		var resT sdk.Result
		if resT = handleTo(ctx, am, loadMsg.Account, loadMsg.Amount); !resT.IsOK() {
			return resT
		}

		return sdk.Result{
			Log:  resT.Log,
			Data: resT.Data,
			Tags: loadMsg.Tags(),
		}

	}
}
