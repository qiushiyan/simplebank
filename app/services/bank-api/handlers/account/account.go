package account

import (
	"context"
	"encoding/json"
	"net/http"
)

func ListAccounts(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return json.NewEncoder(w).Encode([]string{"account1", "account2"})
}
