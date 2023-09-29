package models

import (
	"database/sql"
	"errors"
	"fmt"
)

type Gallery struct {
	ID     int
	UserID int
	Title  string
}

type GalleryService struct {
	DB *sql.DB
}

func (srv *GalleryService) Create(title string, userID int) (*Gallery, error) {
	gallery := Gallery{
		Title:  title,
		UserID: userID,
	}

	row := srv.DB.QueryRow(`
		INSERT INTO galleries (title, user_id)
		VALUES ($1, $2) RETURNING id;`, gallery.Title, gallery.UserID)
	err := row.Scan(&gallery.ID)
	if err != nil {
		return nil, fmt.Errorf("create gallery: %w", err)
	}

	return &gallery, nil
}

func (srv *GalleryService) ByID(id int) (*Gallery, error) {
	gallery := Gallery{
		ID: id,
	}
	row := srv.DB.QueryRow(`
		SELECT user_id, title
		FROM galleries
		WHERE id = $1;
	`, gallery.ID)
	err := row.Scan(&gallery.UserID, &gallery.Title)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("gallery by id: %w", err)
	}

	return &gallery, nil
}

func (srv *GalleryService) ByUserID(userID int) ([]Gallery, error) {
	rows, err := srv.DB.Query(`
		SELECT id, title
		FROM galleries
		WHERE user_id = $1;
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("galleries by user: %w", err)
	}

	var galleries []Gallery
	for rows.Next() {
		gallery := Gallery{
			UserID: userID,
		}
		err := rows.Scan(&gallery.ID, &gallery.Title)
		if err != nil {
			return nil, fmt.Errorf("galleries by user: %w", err)
		}
		galleries = append(galleries, gallery)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("galleries by user: %w", err)
	}

	return galleries, nil
}

func (srv *GalleryService) Update(gallery *Gallery) error {
	_, err := srv.DB.Exec(`
		UPDATE galleries
		SET title = $2
		WHERE id = $1;
	`, gallery.ID, gallery.Title)
	if err != nil {
		return fmt.Errorf("update gallery: %w", err)
	}

	return nil
}

func (srv *GalleryService) Delete(id int) error {
	_, err := srv.DB.Exec(`
		DELETE FROM galleries
		WHERE id = $1;
	`, id)
	if err != nil {
		return fmt.Errorf("delete gallery: %w", err)
	}
	return nil
}
