package engine

import (
	"bruit/bruit/clients/kraken"
	"errors"
	"log"
)

func IsFatalErr(err error) bool {
	if err == nil {
		return false
	}

	isFatal := errors.Is(err, kraken.ErrPairNotFound) ||
		errors.Is(err, kraken.ErrPubSocketNotInit) ||
		errors.Is(err, kraken.ErrBookSocketNotInit) ||
		errors.Is(err, kraken.ErrPrivSocketNotInit) ||
		errors.Is(err, kraken.ErrNotPubSocket) ||
		errors.Is(err, kraken.ErrNotBookSocket) ||
		errors.Is(err, kraken.ErrNotPrivSocket) ||
		errors.Is(err, kraken.ErrPubSocketNotConnected) ||
		errors.Is(err, kraken.ErrBookSocketNotConnected) ||
		errors.Is(err, kraken.ErrPrivSocketNotConnected)

	log.Println("error: ", err, "isFatal", isFatal)

	return isFatal
}
