package types

// Distribution module event types
var (
	EventTypeSetWithdrawAddress = "set_withdraw_address"
	EventTypeRewards            = "rewards"
	EventTypeNonDenomRewards    = "non_denom_rewards"
	EventTypeCommission         = "commission"
	EventTypeWithdrawCommission = "withdraw_commission"
	EventTypeProposerReward     = "proposer_reward"

	AttributeKeyWithdrawAddress = "withdraw_address"
	AttributeKeyValidator       = "validator"

	AttributeValueCategory = ModuleName
)
