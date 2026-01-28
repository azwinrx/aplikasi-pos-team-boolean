package adaptor

// Adaptor is a struct that holds all HTTP handlers for the application
type Adaptor struct {
	InventoriesAdaptor *InventoriesAdaptor
	StaffAdaptor       *StaffAdaptor
}

// NewAdaptor creates a new instance of Adaptor with all handlers
func NewAdaptor(inventoriesAdaptor *InventoriesAdaptor, staffAdaptor *StaffAdaptor) *Adaptor {
	return &Adaptor{
		InventoriesAdaptor: inventoriesAdaptor,
		StaffAdaptor:       staffAdaptor,
	}
}
