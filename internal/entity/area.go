package entity

type AreaCenter struct {
	ID         string
	Name       string
	EnName     string
	OfficeName string
}

type AreaOffice struct {
	ID       string
	Name     string
	EnName   string
	ParentID string // area_centers.id
}

type AreaClass10 struct {
	ID       string
	Name     string
	EnName   string
	ParentID string // area_offices.id
}

type AreaClass15 struct {
	ID       string
	Name     string
	EnName   string
	ParentID string // area_class10.id
}

type AreaClass20 struct {
	ID       string
	Name     string
	EnName   string
	ParentID string // area_class15.id
}

// 階層をまとめて表現するための複合エンティティ
type HierarchyArea struct {
	Class20 *AreaClass20
	Class15 *AreaClass15
	Class10 *AreaClass10
	Office  *AreaOffice
	Center  *AreaCenter
}
