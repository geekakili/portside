package models

// Label holds information about a label
type Label struct {
	Name        string   `bow:"key" validate:"empty=false"` // Name of the label
	Description string   //Description for this label
	Images      []string //List of images classified by this label
}
