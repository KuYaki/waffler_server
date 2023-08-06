package handler

import (
	"github.com/KuYaki/waffler_server/internal/infrastructure/errors"
	"github.com/KuYaki/waffler_server/internal/infrastructure/midlleware"
	"github.com/KuYaki/waffler_server/internal/infrastructure/tools/cryptography"
	"net/http"
	"strconv"
)

func ExtractUser(r *http.Request) (cryptography.UserFromClaims, error) {
	ctx := r.Context()
	u, ok := ctx.Value(midlleware.UserRequest{}).(cryptography.UserClaims)
	if !ok {
		return cryptography.UserFromClaims{}, errors.TokenExtractUserError
	}
	userID, err := strconv.Atoi(u.ID)
	if err != nil {
		return cryptography.UserFromClaims{}, errors.TokenExtractUserError
	}

	return cryptography.UserFromClaims{
		ID: userID,
	}, nil
}
