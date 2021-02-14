package domain

const (
	tagSuffix = "99999999999999999999999"
	// TagPreorderTrytes is the trytes for pre-order tag
	TagPreorderTrytes = "ZBYB" + tagSuffix
	// TagRegisterTrytes is the trytes for register tag
	TagRegisterTrytes = "ACQB" + tagSuffix
	// TagRenewTrytes is the trytes for renew tag
	TagRenewTrytes = "ACXB" + tagSuffix
	// TagUpdateTrytes is the trytes for update tag
	TagUpdateTrytes = "DCNB" + tagSuffix
	// TagTransferTrytes is the trytes for transfer tag
	TagTransferTrytes = "CCPB" + tagSuffix
	// TagRevokeTrytes is the trytes for revoke tag
	TagRevokeTrytes = "ACEC" + tagSuffix
)

// PreorderName sends the transaction for preordering a name
func (d *Domain) PreorderName(name string) error {
	return nil
}

// RegisterName sends the transaction for registering a name
func (d *Domain) RegisterName()

// RenewName sends the transaction for renewing a name
func (d *Domain) RenewName()

// UpdateName sends the transaction for updating a name
func (d *Domain) UpdateName()

// TransferName sends the transaction for transfering a name
func (d *Domain) TransferName()

// RevokeName sends the transaction for recoking a name
func (d *Domain) RevokeName()
