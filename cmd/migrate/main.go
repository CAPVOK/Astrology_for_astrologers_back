package main

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"space/internal/dsn"
	"space/internal/model"
)

func main() {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Явно мигрировать только нужные таблицы
	err = db.AutoMigrate(&model.User{}, &model.Planet{}, &model.Constellation{}, &model.ConstellationPlanet{})
	if err != nil {
		panic("cant migrate db")
	}
}

/* CREATE TABLE
    constellations (
        constellation_id SERIAL PRIMARY KEY,
        name VARCHAR(50) NOT NULL,
        start_date TIMESTAMP NOT NULL,
        end_date TIMESTAMP NOT NULL,
        moderator_id INTEGER,
        user_id INTEGER NOT NULL,
        constellation_status VARCHAR(11) CHECK (
            constellation_status IN (
                'created',
                'inprogress',
                'completed',
                'deleted',
                'canceled'
            )
        ),
        creation_date TIMESTAMP NOT NULL,
        formation_date TIMESTAMP,
        confirmation_date TIMESTAMP,
        FOREIGN KEY (moderator_id) REFERENCES "users"(user_id),
        FOREIGN KEY (user_id) REFERENCES "users"(user_id)
    );

DROP TABLE constellations;

DELETE FROM constellations; */
