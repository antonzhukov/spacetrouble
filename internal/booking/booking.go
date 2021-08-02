package booking

import "github.com/antonzhukov/spacetrouble/internal/entity"

type Store interface {
	Add(b *entity.Booking) error
}
