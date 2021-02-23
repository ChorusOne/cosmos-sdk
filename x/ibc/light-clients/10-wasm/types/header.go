package types

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/ibc/core/exported"
)

var _ exported.Header = (*Header)(nil)

func (m *Header) ClientType() string {
	return m.Type
}

func (m *Header) GetHeight() exported.Height {
	return m.Height
}

func (m *Header) ValidateBasic() error {
	if m.Data == nil || len(m.Data) == 0 {
		return fmt.Errorf("data cannot be empty")
	}

	if m.CodeId == nil || len(m.CodeId) == 0 {
		return fmt.Errorf("codeid cannot be empty")
	}

	return nil
}
