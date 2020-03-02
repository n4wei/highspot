package collection

import (
	mixtape_pkg "github.com/n4wei/highspot/mixtape"
	"github.com/n4wei/highspot/models"
	"github.com/n4wei/highspot/util"
)

// The purpose of this interface is to decouple the top level code
// in main from the object implementing the ApplyChanges logic.
type Collection interface {
	ApplyChanges(changes *models.Changes) error
}

// This function is really simple, but we could use the factory pattern
// for more complex instantiation needs
func New(mixtape *models.Mixtape, logger util.Logger) Collection {
	return mixtape_pkg.New(mixtape, logger)
}
