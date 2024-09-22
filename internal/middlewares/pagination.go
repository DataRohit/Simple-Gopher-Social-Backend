package middlewares

import (
	"context"
	"gopher-social-backend-server/pkg/constants"
	"gopher-social-backend-server/pkg/utils"
	"net/http"
)

func PaginationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limit, err := utils.ParseLimitOffsetQueryParam(r, string(constants.LimitKey), constants.DefaultLimit, constants.MinLimit, constants.MaxLimit)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		offset, err := utils.ParseLimitOffsetQueryParam(r, string(constants.OffsetKey), constants.DefaultOffset, constants.DefaultOffset, -1)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), constants.LimitKey, limit)
		ctx = context.WithValue(ctx, constants.OffsetKey, offset)
		r = r.WithContext(ctx)

		limitFromContext := r.Context().Value(constants.LimitKey)
		offsetFromContext := r.Context().Value(constants.OffsetKey)

		if limitFromContext == nil || offsetFromContext == nil {
			utils.WriteError(w, http.StatusInternalServerError, "failed to set limit and offset in context")
			return
		}

		if limitFromContext.(int) < 0 || limitFromContext.(int) > constants.MaxLimit {
			utils.WriteError(w, http.StatusBadRequest, "invalid limit: must be between 1 and 20")
			return
		}

		if offsetFromContext.(int) < 0 {
			utils.WriteError(w, http.StatusBadRequest, "invalid offset: must be greater than or equal to 0")
			return
		}

		next.ServeHTTP(w, r)
	})
}
