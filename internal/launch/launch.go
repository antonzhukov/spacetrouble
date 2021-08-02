package launch

import "github.com/antonzhukov/spacetrouble/internal/entity"

type Provider interface {
	GetLaunches() []*entity.Launch
}
