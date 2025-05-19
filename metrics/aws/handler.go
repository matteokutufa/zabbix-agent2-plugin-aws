// handler.go
package aws

import (
	"encoding/json"
	"fmt"
	"sync"

	"git.zabbix.com/ap/plugin-support/plugin"
	"git.zabbix.com/ap/plugin-support/zbxerr"

	"github.com/matteokutufa/zabbix-agent2-plugin-aws/aws"
)

// SessionPool gestisce un pool di sessioni AWS
type SessionPool struct {
	mu       sync.Mutex
	sessions map[string]*aws.Client
}

// clientPool è il pool globale di client AWS
var clientPool = SessionPool{
	sessions: make(map[string]*aws.Client),
}

// getClient restituisce un client AWS per l'account specificato
func getClient(accountID string) (*aws.Client, error) {
	clientPool.mu.Lock()
	defer clientPool.mu.Unlock()

	// Verifica se esiste già un client per questo account
	if client, exists := clientPool.sessions[accountID]; exists {
		return client, nil
	}

	// Carica la configurazione degli account
	accounts, err := aws.LoadAccounts(AccountFile())
	if err != nil {
		return nil, fmt.Errorf("failed to load accounts: %v", err)
	}

	// Verifica che l'account esista
	account, exists := accounts[accountID]
	if !exists {
		return nil, fmt.Errorf("account ID %s not found in configuration", accountID)
	}

	// Crea un nuovo client AWS
	client, err := aws.NewClient(account)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS client: %v", err)
	}

	// Memorizza il client nel pool
	clientPool.sessions[accountID] = client

	return client, nil
}

// validateParams verifica che il numero di parametri sia corretto
func validateParams(params []string, minParams int) error {
	if len(params) < minParams {
		return zbxerr.ErrorTooFewParameters
	}

	return nil
}

// formatZabbixOutput formatta l'output per Zabbix
func formatZabbixOutput(data interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %v", err)
	}

	return string(jsonData), nil
}

// Ping esegue un controllo di salute del plugin
func Ping(ctx plugin.ContextProvider, params []string, _ bool) (interface{}, error) {
	return 1, nil
}
