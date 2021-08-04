package booking

import (
	"database/sql"

	"github.com/pkg/errors"

	"github.com/antonzhukov/spacetrouble/internal/entity"
)

type Postgres struct {
	client *sql.DB
}

func NewPostgres(client *sql.DB) *Postgres {
	return &Postgres{client: client}
}

func (p *Postgres) Add(b *entity.Booking) error {
	q := `INSERT INTO bookings(
	first_name, last_name, gender, birthday, launchpad_id, destination_id, launch_date
) VALUES ($1, $2, $3, $4, $5, $6, $7);`
	stmt, err := p.client.Prepare(q)
	if err != nil {
		return errors.Wrap(err, "Prepare failed")
	}

	_, err = stmt.Exec(b.FirstName, b.LastName, b.Gender, b.Birthday, b.LaunchpadID, b.DestinationID, b.LaunchDate)
	if err != nil {
		return errors.Wrap(err, "Exec failed")
	}
	return nil
}

func (p *Postgres) GetAll() ([]*entity.Booking, error) {
	q := `SELECT id, first_name, last_name, gender, birthday, launchpad_id, destination_id, launch_date FROM bookings`
	it, err := p.client.Query(q)
	if err != nil {
		return nil, errors.Wrap(err, "Query failed")
	}
	var res []*entity.Booking
	for it.Next() {
		b := &entity.Booking{}
		err = it.Scan(
			&b.ID,
			&b.FirstName,
			&b.LastName,
			&b.Gender,
			&b.Birthday,
			&b.LaunchpadID,
			&b.DestinationID,
			&b.LaunchDate,
		)
		if err != nil {
			return nil, errors.Wrap(err, "Scan failed")
		}

		res = append(res, b)
	}

	return res, nil
}
