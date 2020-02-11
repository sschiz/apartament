package postgres

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres driver
	"github.com/sschiz/apartament/internal/apartment"
	"github.com/sschiz/apartament/models"
	"strings"
)

type apt struct {
	Rooms        int
	Area         float64
	Floor        int
	Rent         float64
	City         string
	District     string
	Address      string
	Corpus       string
	Floors       int
	Name         string
	MinAptNumber int `db:"min_apartment_number"`
	MaxAptNumber int `db:"max_apartment_number"`
}

func (a apt) toApartment() *models.Apartment {
	return &models.Apartment{
		Floor: a.Floor,
		Rooms: a.Rooms,
		Area:  a.Area,
		Rent:  a.Rent,
		House: &models.House{
			City:     a.City,
			District: a.District,
			Address:  a.Address,
			Corpus:   a.Corpus,
			Floors:   a.Floors,
		},
		ApartmentComplex: &models.ApartmentComplex{
			Name:       a.Name,
			Apartments: [2]int{a.MaxAptNumber, a.MaxAptNumber},
		},
	}
}

type ApartmentRepository struct {
	db *sqlx.DB
}

// NewApartmentRepository creates ApartmentRepository
func NewApartmentRepository(db *sqlx.DB) *ApartmentRepository {
	return &ApartmentRepository{db: db}
}

// Create creates new apartment using transactions
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
			"SELECT * FROM houses WHERE city = :city AND address = :address AND corpus = :corpus"+
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

// Get returns apartments from DB by incoming at
func (a ApartmentRepository) Get(ctx context.Context, at *models.Apartment, opts ...apartment.Option) ([]*models.Apartment, error) {
	if at.House == nil || len(at.House.City) == 0 {
		return nil, apartment.ErrWrongHouse
	}
	if len(at.House.District) == 0 && len(at.House.Address) == 0 {
		return nil, apartment.ErrWrongHouse
	}

	options := apartment.Options{
		Limit:      apartment.DefaultLimit,
		Offset:     apartment.DefaultOffset,
		OrderField: apartment.DefaultOrderField,
	}

	for _, o := range opts {
		o.Apply(&options)
	}

	var optsPart string
	if options.OrderField != apartment.DefaultOrderField {
		optsPart += fmt.Sprintf("ORDER BY %s ", options.OrderField)
	}
	if options.Limit != apartment.DefaultLimit {
		optsPart += fmt.Sprintf("LIMIT %d ", options.Limit)
	}
	if options.Offset != apartment.DefaultOffset {
		optsPart += fmt.Sprintf("OFFSET %d", options.Offset)
	}

	query := "SELECT rooms," +
		"area," +
		"rent::numeric," +
		"city," +
		"floor," +
		"district," +
		"address," +
		"corpus," +
		"floors," +
		"name," +
		"min_apartment_number," +
		"max_apartment_number " +
		"FROM apartments " +
		"INNER JOIN houses h on apartments.house_id = h.id " +
		"INNER JOIN apartment_complexes ac on h.ac_name = ac.name WHERE"

	where := generateWhere(at)

	query += " " + where + " " + optsPart

	rows, err := a.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []*models.Apartment
	for rows.Next() {
		apt := &apt{}
		err := rows.StructScan(apt)
		if err != nil {
			return nil, err
		}

		res = append(res, apt.toApartment())
	}

	return res, nil
}

func generateWhere(apt *models.Apartment) string {
	conditions := make([]string, 0, 10)
	if apt.Floor != 0 {
		conditions = append(conditions, fmt.Sprint("floor = ", apt.Floor))
	}

	if apt.Rooms != 0 {
		conditions = append(conditions, fmt.Sprint("rooms = ", apt.Rooms))
	}

	if apt.Rent != 0 {
		conditions = append(conditions, fmt.Sprint("rent::numeric = ", apt.Rent))
	}

	if apt.Area != 0 {
		conditions = append(conditions, fmt.Sprint("area = ", apt.Area))
	}

	if apt.ApartmentComplex != nil && len(apt.ApartmentComplex.Name) != 0 {
		conditions = append(conditions, "name = '"+apt.ApartmentComplex.Name+"'")
	}

	if len(apt.House.City) != 0 {
		conditions = append(conditions, "city = '"+apt.House.City+"'")
	}

	if len(apt.House.District) != 0 {
		conditions = append(conditions, "district = '"+apt.House.District+"'")
	}

	if len(apt.House.Address) != 0 {
		conditions = append(conditions, "address = '"+apt.House.Address+"'")
	}

	if len(apt.House.Corpus) != 0 {
		conditions = append(conditions, "corpus = '"+apt.House.Corpus+"'")
	}

	if apt.House.Floors != 0 {
		conditions = append(conditions, fmt.Sprint("floors = ", apt.House.Floors))
	}

	return strings.Join(conditions, " AND ")
}
