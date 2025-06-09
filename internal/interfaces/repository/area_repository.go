package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Isshinfunada/weather-bot/internal/entity"
)

// インターフェース
type AreaRepository interface {
	FindHierarchyByClass20ID(ctx context.Context, class20ID string) (*entity.HierarchyArea, error)
	FindAreasByname(ctx context.Context, namePattern string) ([]*entity.AreaClass20, error)
	// 追加メソッド
	FindOfficeByName(ctx context.Context, name string) (*entity.AreaOffice, error)
	FindAreaClass10ByOfficeID(ctx context.Context, officeID string) ([]*entity.AreaClass10, error)
	FindAreaClass15ByClass10ID(ctx context.Context, class10ID string) ([]*entity.AreaClass15, error)
	FindAreaClass20ByClass15ID(ctx context.Context, class15ID string) ([]*entity.AreaClass20, error)
}

// 実装構造体
type areaRepository struct {
	db *sql.DB
}

// コンストラクタ
func NewAreaRepository(db *sql.DB) AreaRepository {
	return &areaRepository{db: db}
}

func (r *areaRepository) FindHierarchyByClass20ID(ctx context.Context, class20ID string) (*entity.HierarchyArea, error) {
	query := `
        SELECT
            c20.id, c20.name, c20.en_name, c20.parent_id,
            c15.id, c15.name, c15.en_name, c15.parent_id,
            c10.id, c10.name, c10.en_name, c10.parent_id,
            o.id, o.name, o.en_name, o.parent_id,
            ct.id, ct.name, ct.en_name, ct.office_name
        FROM area_class20 c20
        JOIN area_class15 c15 ON c20.parent_id = c15.id
        JOIN area_class10 c10 ON c15.parent_id = c10.id
        JOIN area_offices o ON c10.parent_id = o.id
        JOIN area_centers ct ON o.parent_id = ct.id
        WHERE c20.id = $1
	`
	row := r.db.QueryRowContext(ctx, query, class20ID)

	var (
		c20 entity.AreaClass20
		c15 entity.AreaClass15
		c10 entity.AreaClass10
		o   entity.AreaOffice
		ct  entity.AreaCenter
	)
	err := row.Scan(
		&c20.ID, &c20.Name, &c20.EnName, &c20.ParentID,
		&c15.ID, &c15.Name, &c15.EnName, &c15.ParentID,
		&c10.ID, &c10.Name, &c10.EnName, &c10.ParentID,
		&o.ID, &o.Name, &o.EnName, &o.ParentID,
		&ct.ID, &ct.Name, &ct.EnName, &ct.OfficeName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query error: %w", err)
	}

	return &entity.HierarchyArea{
		Class20: &c20,
		Class15: &c15,
		Class10: &c10,
		Office:  &o,
		Center:  &ct,
	}, nil
}

func (r *areaRepository) FindAreasByname(ctx context.Context, namePattern string) ([]*entity.AreaClass20, error) {
	query := `
		SELECT id, name, en_name, parent_id
		FROM area_class20
		WHERE name LIKE $1
	`

	rows, err := r.db.QueryContext(ctx, query, namePattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var areas []*entity.AreaClass20
	for rows.Next() {
		var area entity.AreaClass20
		if err := rows.Scan(&area.ID, &area.Name, &area.EnName, &area.ParentID); err != nil {
			return nil, err
		}
		areas = append(areas, &area)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return areas, nil
}

// 都道府県名から都道府県情報を検索
func (r *areaRepository) FindOfficeByName(ctx context.Context, name string) (*entity.AreaOffice, error) {
	query := `
		SELECT id, name, en_name, parent_id
		FROM area_offices
		WHERE name LIKE $1
		LIMIT 1
	`

	row := r.db.QueryRowContext(ctx, query, "%"+name+"%")

	var office entity.AreaOffice
	err := row.Scan(&office.ID, &office.Name, &office.EnName, &office.ParentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query error: %w", err)
	}

	return &office, nil
}

// 都道府県IDから一次細分区域一覧を取得
func (r *areaRepository) FindAreaClass10ByOfficeID(ctx context.Context, officeID string) ([]*entity.AreaClass10, error) {
	query := `
		SELECT id, name, en_name, parent_id
		FROM area_class10
		WHERE parent_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, officeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var areas []*entity.AreaClass10
	for rows.Next() {
		var area entity.AreaClass10
		if err := rows.Scan(&area.ID, &area.Name, &area.EnName, &area.ParentID); err != nil {
			return nil, err
		}
		areas = append(areas, &area)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return areas, nil
}

// 一次細分区域IDから市町村等一覧を取得
func (r *areaRepository) FindAreaClass15ByClass10ID(ctx context.Context, class10ID string) ([]*entity.AreaClass15, error) {
	query := `
		SELECT id, name, en_name, parent_id
		FROM area_class15
		WHERE parent_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, class10ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var areas []*entity.AreaClass15
	for rows.Next() {
		var area entity.AreaClass15
		if err := rows.Scan(&area.ID, &area.Name, &area.EnName, &area.ParentID); err != nil {
			return nil, err
		}
		areas = append(areas, &area)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return areas, nil
}

// 市町村等IDから詳細地域一覧を取得
func (r *areaRepository) FindAreaClass20ByClass15ID(ctx context.Context, class15ID string) ([]*entity.AreaClass20, error) {
	query := `
		SELECT id, name, en_name, parent_id
		FROM area_class20
		WHERE parent_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, class15ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var areas []*entity.AreaClass20
	for rows.Next() {
		var area entity.AreaClass20
		if err := rows.Scan(&area.ID, &area.Name, &area.EnName, &area.ParentID); err != nil {
			return nil, err
		}
		areas = append(areas, &area)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return areas, nil
}
