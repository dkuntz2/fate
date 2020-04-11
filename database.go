package fate

import (
	// "database/sql"
	_ "github.com/jackc/pgx/stdlib"

	"fmt"
	"github.com/jackc/pgtype"
	"github.com/jmoiron/sqlx"
)

type DbService struct {
	*sqlx.DB
}

type rowScanner interface {
	StructScan(interface{}) error
}

var (
	singletonDb *DbService
)

func ProvideDb() *DbService {
	if singletonDb == nil {
		dbHost := EnvValue("DB_HOST")
		dbUser := EnvValue("DB_USER")
		dbName := EnvValue("DB_NAME")
		dbPass := EnvValue("DB_PASS")

		dbUrl := fmt.Sprintf(
			"host=%s user=%s dbname=%s password=%s",
			dbHost, dbUser, dbName, dbPass,
		)

		db, err := sqlx.Connect("pgx", dbUrl)
		if err != nil {
			panic(err)
		}

		singletonDb = &DbService{db}
	}

	return singletonDb
}

type dbCharacter struct {
	Pk           int               `db:"pk"`
	Name         string            `db:"name"`
	Player       string            `db:"player"`
	Refresh      int8              `db:"refresh"`
	FatePoints   int8              `db:"fate_points"`
	HighConcept  string            `db:"high_concept"`
	Trouble      string            `db:"trouble"`
	Aspects      *pgtype.TextArray `db:"aspects"`
	Stress       int8              `db:"stress"`
	Consequences *pgtype.TextArray `db:"consequences"`

	Careful  int8 `db:"careful"`
	Clever   int8 `db:"clever"`
	Flashy   int8 `db:"flashy"`
	Forceful int8 `db:"forceful"`
	Quick    int8 `db:"quick"`
	Sneaky   int8 `db:"sneaky"`
}

func (char dbCharacter) Character() *Character {
	aspects := []string{}
	if char.Aspects.Status == pgtype.Present {
		for _, aspect := range char.Aspects.Elements {
			aspects = append(aspects, aspect.String)
		}
	}

	consequences := []string{}
	if char.Consequences.Status == pgtype.Present {
		for _, conseq := range char.Consequences.Elements {
			consequences = append(consequences, conseq.String)
		}
	}

	return &Character{
		Name:         char.Name,
		Player:       char.Player,
		Refresh:      char.Refresh,
		FatePoints:   char.FatePoints,
		HighConcept:  char.HighConcept,
		Trouble:      char.Trouble,
		Aspects:      aspects,
		Stress:       char.Stress,
		Consequences: consequences,

		Approaches: &Approaches{
			Careful:  char.Careful,
			Clever:   char.Clever,
			Flashy:   char.Flashy,
			Forceful: char.Forceful,
			Quick:    char.Quick,
			Sneaky:   char.Sneaky,
		},
	}
}

func (char *Character) toDb() *dbCharacter {
	dbChar := &dbCharacter{
		Name:        char.Name,
		Player:      char.Player,
		Refresh:     char.Refresh,
		FatePoints:  char.FatePoints,
		HighConcept: char.HighConcept,
		Trouble:     char.Trouble,
		Stress:      char.Stress,

		Careful:  char.Approaches.Careful,
		Clever:   char.Approaches.Clever,
		Flashy:   char.Approaches.Flashy,
		Forceful: char.Approaches.Forceful,
		Quick:    char.Approaches.Quick,
		Sneaky:   char.Approaches.Sneaky,

		Aspects:      &pgtype.TextArray{},
		Consequences: &pgtype.TextArray{},
	}

	dbChar.Aspects.Set(char.Aspects)
	dbChar.Consequences.Set(char.Consequences)

	return dbChar
}

func (db *DbService) AllCharacters() ([]*Character, error) {
	rows, err := db.Queryx("SELECT * FROM character")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	chars := []*Character{}
	for rows.Next() {
		charDb := dbCharacter{}
		err = rows.StructScan(&charDb)
		if err != nil {
			return nil, err
		}

		chars = append(chars, charDb.Character())
	}

	return chars, nil
}

func (db *DbService) SaveCharacter(char *Character) error {
	charDb := char.toDb()
	_, err := db.NamedExec(
		`INSERT INTO character (
			player, name, refresh, fate_points, high_concept, trouble, aspects, stress, consequences, careful, clever, flashy, forceful, quick, sneaky
		) VALUES (
			:player, :name, :refresh, :fate_points, :high_concept, :trouble, :aspects, :stress, :consequences, :careful, :clever, :flashy, :forceful, :quick, :sneaky
		) ON CONFLICT (player) DO UPDATE SET
			name = excluded.name,
			refresh = excluded.refresh,
			fate_points = excluded.fate_points,
			high_concept = excluded.high_concept,
			trouble = excluded.trouble,
			aspects = excluded.aspects,
			stress = excluded.stress,
			consequences = excluded.consequences,
			careful = excluded.careful,
			clever = excluded.clever,
			flashy = excluded.flashy,
			forceful = excluded.forceful,
			quick = excluded.quick,
			sneaky = excluded.sneaky`,
		charDb,
	)

	return err
}

func (db *DbService) ResetFatePoints() error {
	_, err := db.Exec("UPDATE character SET fate_points = refresh")
	return err
}
