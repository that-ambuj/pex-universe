package model

import (
	"database/sql"
	"pex-universe/types"
	"time"
)

type (
	Address struct {
		Id             uint64
		CreatedAt      *time.Time
		UpdatedAt      *time.Time
		Verified       bool
		FirstName      string
		LastName       string
		Company        types.NullString `swaggertype:"string"`
		StreetAddress1 string           `db:"street_address"`
		StreetAddress2 types.NullString `swaggertype:"string"`
		City           string
		Zip            string
		Phone          string
		Ext            string
		Email          string
		// Foreign Keys
		CountryId sql.NullInt64 `db:"country_id" json:"-"`
		StateId   sql.NullInt64 `db:"state_id" json:"-"`
		UserId    sql.NullInt64 `db:"user_id" json:"-"`
		State     *State
		Country   *Country
	}

	State struct {
		Id              uint64 `json:"-"`
		Name            string
		FullName        string
		Tax             float32
		Info            string
		Locked          bool `json:"-"`
		TaxEnabled      bool `json:"-"`
		FixedRate       bool `json:"-"`
		ShippingTaxable bool `json:"-"`
	}

	Country struct {
		Id          uint64 `json:"-"`
		Name        string
		DisplayName string
		PpCode      string
		Position    int64
		Locked      bool `json:"-"`
	}
)

func FindAddressesByUserId(db *sql.DB, userId uint64) ([]*Address, error) {
	addrs := []*Address{}

	query := `
SELECT a.id,
	a.created_at,
	a.updated_at,
	a.verified,
	a.first_name,
	a.last_name,
	a.company,
	a.street_address,
	a.street_address2,
	a.city,
	a.zip,
	a.phone,
	a.ext,
	a.email,

	s.name,
	s.full_name,
	s.tax,
	s.info,

	c.name,
	c.display_name,
	c.pp_code,
	c.locked
FROM addresses a
         JOIN states s ON a.state_id = s.id
         JOIN countries c ON a.country_id = c.id
WHERE a.user_id = ?;
`

	rows, err := db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		addr := &Address{}

		addr.State = &State{}
		addr.Country = &Country{}

		err = rows.Scan(
			&addr.Id,
			&addr.CreatedAt,
			&addr.UpdatedAt,
			&addr.Verified,
			&addr.FirstName,
			&addr.LastName,
			&addr.Company,
			&addr.StreetAddress1,
			&addr.StreetAddress2,
			&addr.City,
			&addr.Zip,
			&addr.Phone,
			&addr.Ext,
			&addr.Email,

			&addr.State.Name,
			&addr.State.FullName,
			&addr.State.Tax,
			&addr.State.Info,

			&addr.Country.Name,
			&addr.Country.DisplayName,
			&addr.Country.PpCode,
			&addr.Country.Locked,
		)

		if err != nil {
			return nil, err
		}

		addrs = append(addrs, addr)
	}

	return addrs, nil
}
