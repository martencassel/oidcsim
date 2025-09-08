package delegation

import (
	"context"
	"database/sql"
	"strings"

	delegationapp "github.com/martencassel/oidcsim/internal/application/delegation"
	"github.com/martencassel/oidcsim/internal/domain/delegation"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) FindByUserAndClient(ctx context.Context, userID, clientID string) (*delegation.Delegation, error) {
	const q = `
        SELECT id, user_id, client_id, scopes, created_at
        FROM delegations
        WHERE user_id = $1 AND client_id = $2
        LIMIT 1`
	var d delegation.Delegation
	var scopes string
	err := r.db.QueryRowContext(ctx, q, userID, clientID).
		Scan(&d.ID, &d.UserID, &d.ClientID, &scopes, &d.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	d.Scopes = splitScopes(scopes)
	return &d, nil
}

func (r *PostgresRepo) Save(ctx context.Context, d delegation.Delegation) error {
	const q = `
        INSERT INTO delegations (id, user_id, client_id, scopes, created_at)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (user_id, client_id) DO UPDATE
        SET scopes = $4, created_at = $5`
	_, err := r.db.ExecContext(ctx, q,
		d.ID, d.UserID, d.ClientID, joinScopes(d.Scopes), d.CreatedAt)
	return err
}

func (r *PostgresRepo) Delete(ctx context.Context, userID, clientID string) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM delegations WHERE user_id = $1 AND client_id = $2`,
		userID, clientID)
	return err
}

func splitScopes(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, " ")
}

func joinScopes(scopes []string) string {
	return strings.Join(scopes, " ")
}

// FindByID retrieves a delegation by its ID.
func (r *PostgresRepo) FindByID(ctx context.Context, id string) (*delegation.Delegation, error) {
	const q = `
		SELECT id, user_id, client_id, scopes, created_at
		FROM delegations
		WHERE id = $1
		LIMIT 1`
	var d delegation.Delegation
	var scopes string
	err := r.db.QueryRowContext(ctx, q, id).
		Scan(&d.ID, &d.UserID, &d.ClientID, &scopes, &d.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	d.Scopes = splitScopes(scopes)
	return &d, nil
}

var _ delegationapp.Repository = (*PostgresRepo)(nil)
