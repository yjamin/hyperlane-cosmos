package types

// NewGenesisState creates a new genesis state with default values.
func NewGenesisState() *GenesisState {
	return &GenesisState{
		Params:        Params{},
		Tokens:        []HypToken{},
		RemoteRouters: []GenesisRemoteRouterWrapper{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
func (gs *GenesisState) Validate() error {
	// nothing to validate

	return nil
}
