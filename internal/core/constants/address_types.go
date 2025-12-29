package constants

type AddressTypeID int

const (
	ProvincialAddressTypeID AddressTypeID = iota + 1
	ResidentialAddressTypeID
)

type AddressType struct {
	ID   AddressTypeID
	Name string
}

var AddressTypes = map[AddressTypeID]AddressType{
	ProvincialAddressTypeID:  {ID: ProvincialAddressTypeID, Name: "Provincial"},
	ResidentialAddressTypeID: {ID: ResidentialAddressTypeID, Name: "Residential"},
}
