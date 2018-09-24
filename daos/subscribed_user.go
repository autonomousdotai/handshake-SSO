package daos

import (
	"github.com/ninjadotorg/handshake-dispatcher/models"
)

// SubscribedUserDAO : DAO
type SubscribedUserDAO struct{}

// CountUsersByProduct : product
func (s SubscribedUserDAO) CountUsersByProduct(product string) (int, error) {
	users := []models.SubscribedUser{}
	var count int
	err := models.Database().Where("subscribed_user.product = ?", product).Find(&users).Count(&count).Error
	return count, err
}
