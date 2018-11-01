package pos

import (
	"fmt"

	sdk "bitbucket.org/shareringvn/cosmos-sdk/types"
	wire "bitbucket.org/shareringvn/cosmos-sdk/wire"

	"github.com/sharering/shareledger/constants"
	keep "github.com/sharering/shareledger/x/pos/keeper"
	posTypes "github.com/sharering/shareledger/x/pos/type"

	abci "github.com/tendermint/abci/types"
)

// query endpoints supported by the staking Querier
const (
	QueryValidators          = "validators"
	QueryValidator           = "validator"
	QueryDelegator           = "delegator"
	QueryDelegation          = "delegation"
	QueryUnbondingDelegation = "unbondingDelegation"
	QueryDelegatorValidators = "delegatorValidators"
	QueryDelegatorValidator  = "delegatorValidator"
	QueryPool                = "pool"
	QueryParameters          = "parameters"
	QueryValidatorDistInfo   = "validatorDistInfo"
)

// creates a querier for staking REST endpoints

func NewQuerier(k keep.Keeper, cdc *wire.Codec) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryValidators:
			return queryValidators(ctx, cdc, k)
		case QueryValidator:
			return queryValidator(ctx, cdc, req, k)
		case QueryValidatorDistInfo:
			return queryValidatorDistInfo(ctx, cdc, req, k)
		case QueryDelegation:
			return queryDelegation(ctx, cdc, req, k)
		/*
			case QueryDelegator:
				return queryDelegator(ctx, cdc, req, k)
			case QueryDelegation:
				return queryDelegation(ctx, cdc, req, k)
			case QueryUnbondingDelegation:
				return queryUnbondingDelegation(ctx, cdc, req, k)
			case QueryDelegatorValidators:
				return queryDelegatorValidators(ctx, cdc, req, k)
			case QueryDelegatorValidator:
				return queryDelegatorValidator(ctx, cdc, req, k)
			case QueryPool:
				return queryPool(ctx, cdc, k)
			case QueryParameters:
				return queryParameters(ctx, cdc, k) */
		default:
			return nil, sdk.ErrUnknownRequest("unknown stake query endpoint")
		}
	}
}

// defines the params for the following queries:
// - 'custom/stake/delegator'
// - 'custom/stake/delegatorValidators'
type QueryDelegatorParams struct {
	DelegatorAddr sdk.Address
}

// defines the params for the following queries:
// - 'custom/stake/validator'
type QueryValidatorParams struct {
	ValidatorAddr sdk.Address
}

// defines the params for the following queries:
// - 'custom/stake/delegation'
// - 'custom/stake/unbondingDelegation'
// - 'custom/stake/delegatorValidator'
type QueryBondsParams struct {
	DelegatorAddr sdk.Address
	ValidatorAddr sdk.Address
}

type QueryValidatorDistParams struct {
	ValidatorAddr sdk.Address
}

type QueryDelegationParams struct {
	ValidatorAddr sdk.Address
	DelegatorAddr sdk.Address
}

func queryValidators(ctx sdk.Context, cdc *wire.Codec, k keep.Keeper) (res []byte, err sdk.Error) {
	stakeParams := k.GetParams(ctx)
	validators := k.GetValidators(ctx, stakeParams.MaxValidators)

	res, errRes := cdc.MarshalJSON(validators)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("could not marshal result to JSON: %s", errRes.Error()))
	}
	return res, nil
}

func queryValidator(ctx sdk.Context, cdc *wire.Codec, req abci.RequestQuery, k keep.Keeper) (res []byte, err sdk.Error) {
	var params QueryValidatorParams

	//errRes := cdc.UnmarshalJSON(req.Data, &params)
	errRes := cdc.UnmarshalBinary(req.Data, &params)
	if errRes != nil {
		return []byte{}, sdk.ErrUnknownAddress(fmt.Sprintf("incorrectly formatted request address: %s", err.Error()))
	}

	validator, found := k.GetValidator(ctx, params.ValidatorAddr)
	if !found {
		return []byte{}, posTypes.ErrNoValidatorFound(posTypes.DefaultCodespace)
	}

	res, errRes = cdc.MarshalJSON(validator)
	if errRes != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("could not marshal result to JSON: %s", errRes.Error()))
	}
	return res, nil
}

func queryValidatorDistInfo(
	ctx sdk.Context,
	cdc *wire.Codec,
	req abci.RequestQuery,
	k keep.Keeper,
) (
	res []byte, err sdk.Error,
) {
	var params QueryValidatorDistParams

	errRes := cdc.UnmarshalBinary(req.Data, &params)

	if errRes != nil {
		return []byte{},
			sdk.ErrUnknownAddress(fmt.Sprintf(constants.POS_INVALID_VALIDATOR_ADDRESS, err.Error()))
	}

	vdi, found := k.GetValidatorDistInfo(ctx, params.ValidatorAddr)
	if !found {
		return []byte{},
			posTypes.ErrNoValidatorFound(posTypes.DefaultCodespace)
	}

	res, errRes = cdc.MarshalJSON(vdi)
	if errRes != nil {
		return nil,
			sdk.ErrInternal(fmt.Sprintf(constants.POS_MARSHAL_ERROR, errRes.Error()))
	}

	return res, nil
}

func queryDelegation(
	ctx sdk.Context,
	cdc *wire.Codec,
	req abci.RequestQuery,
	k keep.Keeper,
) (
	res []byte, err sdk.Error,
) {
	var params QueryDelegationParams

	errRes := cdc.UnmarshalBinary(req.Data, &params)

	if errRes != nil {
		return []byte{},
			sdk.ErrUnknownAddress(fmt.Sprintf(constants.POS_INVALID_PARAMS, err.Error()))
	}

	delegation, found := k.GetDelegation(ctx, params.DelegatorAddr, params.ValidatorAddr)
	if !found {
		return []byte{},
			posTypes.ErrNoDelegationFound(posTypes.DefaultCodespace)
	}

	res, errRes = cdc.MarshalJSON(delegation)

	if errRes != nil {
		return nil,
			sdk.ErrInternal(fmt.Sprintf(constants.POS_MARSHAL_ERROR, errRes.Error()))
	}

	return res, nil
}

/*
func queryDelegator(ctx sdk.Context, cdc *wire.Codec, req abci.RequestQuery, k keep.Keeper) (res []byte, err sdk.Error) {
	var params QueryDelegatorParams
	errRes := cdc.UnmarshalJSON(req.Data, &params)
	if errRes != nil {
		return []byte{}, sdk.ErrUnknownAddress(fmt.Sprintf("incorrectly formatted request address: %s", errRes.Error()))
	}
	delegations := k.GetAllDelegatorDelegations(ctx, params.DelegatorAddr)
	unbondingDelegations := k.GetAllUnbondingDelegations(ctx, params.DelegatorAddr)
	redelegations := k.GetAllRedelegations(ctx, params.DelegatorAddr)

	summary := types.DelegationSummary{
		Delegations:          delegations,
		UnbondingDelegations: unbondingDelegations,
		Redelegations:        redelegations,
	}

	res, errRes = codec.MarshalJSONIndent(cdc, summary)
	if errRes != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("could not marshal result to JSON: %s", errRes.Error()))
	}
	return res, nil
}

func queryDelegatorValidators(ctx sdk.Context, cdc *codec.Codec, req abci.RequestQuery, k keep.Keeper) (res []byte, err sdk.Error) {
	var params QueryDelegatorParams

	stakeParams := k.GetParams(ctx)

	errRes := cdc.UnmarshalJSON(req.Data, &params)
	if errRes != nil {
		return []byte{}, sdk.ErrUnknownAddress(fmt.Sprintf("incorrectly formatted request address: %s", errRes.Error()))
	}

	validators := k.GetDelegatorValidators(ctx, params.DelegatorAddr, stakeParams.MaxValidators)

	res, errRes = codec.MarshalJSONIndent(cdc, validators)
	if errRes != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("could not marshal result to JSON: %s", errRes.Error()))
	}
	return res, nil
}

func queryDelegatorValidator(ctx sdk.Context, cdc *codec.Codec, req abci.RequestQuery, k keep.Keeper) (res []byte, err sdk.Error) {
	var params QueryBondsParams

	errRes := cdc.UnmarshalJSON(req.Data, &params)
	if errRes != nil {
		return []byte{}, sdk.ErrUnknownRequest(fmt.Sprintf("incorrectly formatted request address: %s", errRes.Error()))
	}

	validator, err := k.GetDelegatorValidator(ctx, params.DelegatorAddr, params.ValidatorAddr)
	if err != nil {
		return
	}

	res, errRes = codec.MarshalJSONIndent(cdc, validator)
	if errRes != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("could not marshal result to JSON: %s", errRes.Error()))
	}
	return res, nil
}

func queryDelegation(ctx sdk.Context, cdc *codec.Codec, req abci.RequestQuery, k keep.Keeper) (res []byte, err sdk.Error) {
	var params QueryBondsParams

	errRes := cdc.UnmarshalJSON(req.Data, &params)
	if errRes != nil {
		return []byte{}, sdk.ErrUnknownRequest(fmt.Sprintf("incorrectly formatted request address: %s", errRes.Error()))
	}

	delegation, found := k.GetDelegation(ctx, params.DelegatorAddr, params.ValidatorAddr)
	if !found {
		return []byte{}, types.ErrNoDelegation(types.DefaultCodespace)
	}

	res, errRes = codec.MarshalJSONIndent(cdc, delegation)
	if errRes != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("could not marshal result to JSON: %s", errRes.Error()))
	}
	return res, nil
}

func queryUnbondingDelegation(ctx sdk.Context, cdc *codec.Codec, req abci.RequestQuery, k keep.Keeper) (res []byte, err sdk.Error) {
	var params QueryBondsParams

	errRes := cdc.UnmarshalJSON(req.Data, &params)
	if errRes != nil {
		return []byte{}, sdk.ErrUnknownRequest(fmt.Sprintf("incorrectly formatted request address: %s", errRes.Error()))
	}

	unbond, found := k.GetUnbondingDelegation(ctx, params.DelegatorAddr, params.ValidatorAddr)
	if !found {
		return []byte{}, types.ErrNoUnbondingDelegation(types.DefaultCodespace)
	}

	res, errRes = codec.MarshalJSONIndent(cdc, unbond)
	if errRes != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("could not marshal result to JSON: %s", errRes.Error()))
	}
	return res, nil
}

func queryPool(ctx sdk.Context, cdc *codec.Codec, k keep.Keeper) (res []byte, err sdk.Error) {
	pool := k.GetPool(ctx)

	res, errRes := codec.MarshalJSONIndent(cdc, pool)
	if errRes != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("could not marshal result to JSON: %s", errRes.Error()))
	}
	return res, nil
}

func queryParameters(ctx sdk.Context, cdc *codec.Codec, k keep.Keeper) (res []byte, err sdk.Error) {
	params := k.GetParams(ctx)

	res, errRes := codec.MarshalJSONIndent(cdc, params)
	if errRes != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("could not marshal result to JSON: %s", errRes.Error()))
	}
	return res, nil
}
*/
