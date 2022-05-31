package keeper

import (
	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/ibc-go/v3/modules/apps/nft-transfer/types"
)

// GetClassTrace retreives the full identifiers trace and base classId from the store.
func (k Keeper) GetClassTrace(ctx sdk.Context, denomTraceHash tmbytes.HexBytes) (types.ClassTrace, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ClassTraceKey)
	bz := store.Get(denomTraceHash)
	if bz == nil {
		return types.ClassTrace{}, false
	}

	denomTrace := k.MustUnmarshalClassTrace(bz)
	return denomTrace, true
}

// ClassPathFromHash returns the full class path prefix from an ibc classId with a hash
// component.
func (k Keeper) ClassPathFromHash(ctx sdk.Context, classID string) (string, error) {
	// trim the class prefix, by default "ibc/"
	hexHash := classID[len(types.ClassPrefix+"/"):]

	hash, err := types.ParseHexHash(hexHash)
	if err != nil {
		return "", sdkerrors.Wrap(types.ErrInvalidClassID, err.Error())
	}

	classTrace, found := k.GetClassTrace(ctx, hash)
	if !found {
		return "", sdkerrors.Wrap(types.ErrTraceNotFound, hexHash)
	}
	return classTrace.GetFullClassPath(), nil
}

// HasClassTrace checks if a the key with the given denomination trace hash exists on the store.
func (k Keeper) HasClassTrace(ctx sdk.Context, denomTraceHash tmbytes.HexBytes) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ClassTraceKey)
	return store.Has(denomTraceHash)
}

// SetClassTrace sets a new {trace hash -> class trace} pair to the store.
func (k Keeper) SetClassTrace(ctx sdk.Context, denomTrace types.ClassTrace) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ClassTraceKey)
	bz := k.MustMarshalClassTrace(denomTrace)
	store.Set(denomTrace.Hash(), bz)
}