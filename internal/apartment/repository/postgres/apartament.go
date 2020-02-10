package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres driver
	"github.com/sschiz/apartament/internal/apartment"
	"github.com/sschiz/apartament/models"
)

type ApartmentRepository struct {
	db *sqlx.DB
}

func NewApartmentRepository(db *sqlx.DB) *ApartmentRepository {
	return &ApartmentRepository{db: db}
}

func (a ApartmentRepository) Create(ctx context.Context, apartment *models.Apartment) error {
	tx, err := a.db.BeginTxx(ctx, nil)

	if err != nil {
		return err
	}

	defer tx.Rollback()

	acName, err := addAC(ctx, tx, apartment.ApartmentComplex)
	if err != nil {
		return err
	}

	id, err := addHouse(ctx, tx, apartment.House, acName)

	if err != nil {
		return err
	}

	res, err := tx.NamedExecContext(
		ctx,
		"INSERT INTO apartments (rooms, area, rent, house_id) "+
			"SELECT :rooms, :area, :rent, :house_id "+
			"WHERE NOT EXISTS("+
			"SELECT * FROM apartments WHERE rooms = :rooms AND area = :area AND rent = :rent AND house_id = :house_id"+
			")",
		map[string]interface{}{
			"rooms":    apartment.Rooms,
			"area":     apartment.Area,
			"rent":     apartment.Rent,
			"house_id": id,
		},
	)

	if err != nil {
		return err
	}

	ra, err := res.RowsAffected()

	if err != nil || ra == 0 {
		return err
	}

	return tx.Commit()
}

func (a ApartmentRepository) Get(ctx context.Context, apartment *models.Apartment, opts ...apartment.Option) ([]*models.Apartment, error) {
	panic("implement me")
}

func addAC(ctx context.Context, tx *sqlx.Tx, ac *models.ApartmentComplex) (name *string, err error) {
	if ac == nil {
		return
	}

	name = &ac.Name

	if ac.Apartments[0] == 0 {
		ac.Apartments[0] = 1
	}

	if ac.Apartments[1] == 0 {
		ac.Apartments[1] = 1
	}

	if ac.Apartments[0] > ac.Apartments[1] {
		ac.Apartments[0], ac.Apartments[1] = ac.Apartments[1], ac.Apartments[0]
	}

	_, err = tx.ExecContext(
		ctx,
		"INSERT INTO apartment_complexes (name, min_apartment_number, max_apartment_number) "+
			"SELECT $1, $2, $3 "+
			"WHERE NOT EXISTS(SELECT name "+
			"FROM apartment_complexes "+
			"WHERE name = $1)",
		ac.Name, ac.Apartments[0], ac.Apartments[1],
	)

	if err != nil {
		return nil, err
	}

	return
}

func addHouse(ctx context.Context, tx *sqlx.Tx, house *models.House, acName *string) (res int, err error) {
	if house == nil {
		return 0, apartment.ErrWrongHouse
	}

	_, err = tx.NamedExecContext(
		ctx,
		"INSERT INTO houses (city, district, address, corpus, ac_name) "+
			"SELECT :city, :district, :address, :corpus, :ac_name "+
			"WHERE NOT EXISTS("+
			"SELECT * FROM houses WHERE city = :city AND district = :district AND address = :address AND corpus = :corpus"+
			")",
		map[string]interface{}{
			"city":     house.City,
			"district": house.District,
			"address":  house.Address,
			"corpus":   house.Corpus,
			"ac_name":  acName,
		},
	)

	if err != nil {
		return
	}

	err = tx.QueryRowContext(
		ctx,
		"SELECT id FROM houses WHERE city = $1 AND district = $2 AND address = $3 AND corpus = $4",
		house.City, house.District, house.Address, house.Corpus,
	).Scan(&res)

	return
}
