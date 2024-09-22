package middlewares

import (
	"context"
	"gopher-social-backend-server/pkg/constants"
	"gopher-social-backend-server/pkg/utils"
	"net/http"
)

func OrderingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		orderby := utils.ParseOrderByQueryParam(r, string(constants.OrderByKey), constants.DefaultOrderBy)
		desc := utils.ParseDescQueryParam(r, string(constants.DescKey), constants.DefaultDesc)

		if utils.IsSQLInjection(orderby) {
			utils.WriteError(w, http.StatusBadRequest, "invalid orderby: SQL injection detected")
			return
		}

		ctx := context.WithValue(r.Context(), constants.OrderByKey, orderby)
		ctx = context.WithValue(ctx, constants.DescKey, desc)
		r = r.WithContext(ctx)

		if desc != "true" && desc != "false" {
			utils.WriteError(w, http.StatusBadRequest, "invalid desc: must be true or false")
			return
		}

		next.ServeHTTP(w, r)
	})
}
