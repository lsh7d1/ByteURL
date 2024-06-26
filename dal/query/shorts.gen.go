// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"byteurl/dal/model"
)

func newShort(db *gorm.DB, opts ...gen.DOOption) short {
	_short := short{}

	_short.shortDo.UseDB(db, opts...)
	_short.shortDo.UseModel(&model.Short{})

	tableName := _short.shortDo.TableName()
	_short.ALL = field.NewAsterisk(tableName)
	_short.ID = field.NewInt64(tableName, "id")
	_short.Lurl = field.NewString(tableName, "lurl")
	_short.Md5 = field.NewString(tableName, "md5")
	_short.Surl = field.NewString(tableName, "surl")
	_short.CreateAt = field.NewTime(tableName, "create_at")
	_short.CreateBy = field.NewString(tableName, "create_by")
	_short.IsDel = field.NewInt32(tableName, "is_del")

	_short.fillFieldMap()

	return _short
}

// short 长短链映射表
type short struct {
	shortDo shortDo

	ALL      field.Asterisk
	ID       field.Int64  // 主键
	Lurl     field.String // 长链接
	Md5      field.String // 长链接MD5
	Surl     field.String // 短链接
	CreateAt field.Time   // 创建时间
	CreateBy field.String // 创建者
	IsDel    field.Int32  // 是否删除：0正常1删除

	fieldMap map[string]field.Expr
}

func (s short) Table(newTableName string) *short {
	s.shortDo.UseTable(newTableName)
	return s.updateTableName(newTableName)
}

func (s short) As(alias string) *short {
	s.shortDo.DO = *(s.shortDo.As(alias).(*gen.DO))
	return s.updateTableName(alias)
}

func (s *short) updateTableName(table string) *short {
	s.ALL = field.NewAsterisk(table)
	s.ID = field.NewInt64(table, "id")
	s.Lurl = field.NewString(table, "lurl")
	s.Md5 = field.NewString(table, "md5")
	s.Surl = field.NewString(table, "surl")
	s.CreateAt = field.NewTime(table, "create_at")
	s.CreateBy = field.NewString(table, "create_by")
	s.IsDel = field.NewInt32(table, "is_del")

	s.fillFieldMap()

	return s
}

func (s *short) WithContext(ctx context.Context) IShortDo { return s.shortDo.WithContext(ctx) }

func (s short) TableName() string { return s.shortDo.TableName() }

func (s short) Alias() string { return s.shortDo.Alias() }

func (s short) Columns(cols ...field.Expr) gen.Columns { return s.shortDo.Columns(cols...) }

func (s *short) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := s.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (s *short) fillFieldMap() {
	s.fieldMap = make(map[string]field.Expr, 7)
	s.fieldMap["id"] = s.ID
	s.fieldMap["lurl"] = s.Lurl
	s.fieldMap["md5"] = s.Md5
	s.fieldMap["surl"] = s.Surl
	s.fieldMap["create_at"] = s.CreateAt
	s.fieldMap["create_by"] = s.CreateBy
	s.fieldMap["is_del"] = s.IsDel
}

func (s short) clone(db *gorm.DB) short {
	s.shortDo.ReplaceConnPool(db.Statement.ConnPool)
	return s
}

func (s short) replaceDB(db *gorm.DB) short {
	s.shortDo.ReplaceDB(db)
	return s
}

type shortDo struct{ gen.DO }

type IShortDo interface {
	gen.SubQuery
	Debug() IShortDo
	WithContext(ctx context.Context) IShortDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IShortDo
	WriteDB() IShortDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IShortDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IShortDo
	Not(conds ...gen.Condition) IShortDo
	Or(conds ...gen.Condition) IShortDo
	Select(conds ...field.Expr) IShortDo
	Where(conds ...gen.Condition) IShortDo
	Order(conds ...field.Expr) IShortDo
	Distinct(cols ...field.Expr) IShortDo
	Omit(cols ...field.Expr) IShortDo
	Join(table schema.Tabler, on ...field.Expr) IShortDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IShortDo
	RightJoin(table schema.Tabler, on ...field.Expr) IShortDo
	Group(cols ...field.Expr) IShortDo
	Having(conds ...gen.Condition) IShortDo
	Limit(limit int) IShortDo
	Offset(offset int) IShortDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IShortDo
	Unscoped() IShortDo
	Create(values ...*model.Short) error
	CreateInBatches(values []*model.Short, batchSize int) error
	Save(values ...*model.Short) error
	First() (*model.Short, error)
	Take() (*model.Short, error)
	Last() (*model.Short, error)
	Find() ([]*model.Short, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Short, err error)
	FindInBatches(result *[]*model.Short, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.Short) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IShortDo
	Assign(attrs ...field.AssignExpr) IShortDo
	Joins(fields ...field.RelationField) IShortDo
	Preload(fields ...field.RelationField) IShortDo
	FirstOrInit() (*model.Short, error)
	FirstOrCreate() (*model.Short, error)
	FindByPage(offset int, limit int) (result []*model.Short, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IShortDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (s shortDo) Debug() IShortDo {
	return s.withDO(s.DO.Debug())
}

func (s shortDo) WithContext(ctx context.Context) IShortDo {
	return s.withDO(s.DO.WithContext(ctx))
}

func (s shortDo) ReadDB() IShortDo {
	return s.Clauses(dbresolver.Read)
}

func (s shortDo) WriteDB() IShortDo {
	return s.Clauses(dbresolver.Write)
}

func (s shortDo) Session(config *gorm.Session) IShortDo {
	return s.withDO(s.DO.Session(config))
}

func (s shortDo) Clauses(conds ...clause.Expression) IShortDo {
	return s.withDO(s.DO.Clauses(conds...))
}

func (s shortDo) Returning(value interface{}, columns ...string) IShortDo {
	return s.withDO(s.DO.Returning(value, columns...))
}

func (s shortDo) Not(conds ...gen.Condition) IShortDo {
	return s.withDO(s.DO.Not(conds...))
}

func (s shortDo) Or(conds ...gen.Condition) IShortDo {
	return s.withDO(s.DO.Or(conds...))
}

func (s shortDo) Select(conds ...field.Expr) IShortDo {
	return s.withDO(s.DO.Select(conds...))
}

func (s shortDo) Where(conds ...gen.Condition) IShortDo {
	return s.withDO(s.DO.Where(conds...))
}

func (s shortDo) Order(conds ...field.Expr) IShortDo {
	return s.withDO(s.DO.Order(conds...))
}

func (s shortDo) Distinct(cols ...field.Expr) IShortDo {
	return s.withDO(s.DO.Distinct(cols...))
}

func (s shortDo) Omit(cols ...field.Expr) IShortDo {
	return s.withDO(s.DO.Omit(cols...))
}

func (s shortDo) Join(table schema.Tabler, on ...field.Expr) IShortDo {
	return s.withDO(s.DO.Join(table, on...))
}

func (s shortDo) LeftJoin(table schema.Tabler, on ...field.Expr) IShortDo {
	return s.withDO(s.DO.LeftJoin(table, on...))
}

func (s shortDo) RightJoin(table schema.Tabler, on ...field.Expr) IShortDo {
	return s.withDO(s.DO.RightJoin(table, on...))
}

func (s shortDo) Group(cols ...field.Expr) IShortDo {
	return s.withDO(s.DO.Group(cols...))
}

func (s shortDo) Having(conds ...gen.Condition) IShortDo {
	return s.withDO(s.DO.Having(conds...))
}

func (s shortDo) Limit(limit int) IShortDo {
	return s.withDO(s.DO.Limit(limit))
}

func (s shortDo) Offset(offset int) IShortDo {
	return s.withDO(s.DO.Offset(offset))
}

func (s shortDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IShortDo {
	return s.withDO(s.DO.Scopes(funcs...))
}

func (s shortDo) Unscoped() IShortDo {
	return s.withDO(s.DO.Unscoped())
}

func (s shortDo) Create(values ...*model.Short) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Create(values)
}

func (s shortDo) CreateInBatches(values []*model.Short, batchSize int) error {
	return s.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (s shortDo) Save(values ...*model.Short) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Save(values)
}

func (s shortDo) First() (*model.Short, error) {
	if result, err := s.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Short), nil
	}
}

func (s shortDo) Take() (*model.Short, error) {
	if result, err := s.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Short), nil
	}
}

func (s shortDo) Last() (*model.Short, error) {
	if result, err := s.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Short), nil
	}
}

func (s shortDo) Find() ([]*model.Short, error) {
	result, err := s.DO.Find()
	return result.([]*model.Short), err
}

func (s shortDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Short, err error) {
	buf := make([]*model.Short, 0, batchSize)
	err = s.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (s shortDo) FindInBatches(result *[]*model.Short, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return s.DO.FindInBatches(result, batchSize, fc)
}

func (s shortDo) Attrs(attrs ...field.AssignExpr) IShortDo {
	return s.withDO(s.DO.Attrs(attrs...))
}

func (s shortDo) Assign(attrs ...field.AssignExpr) IShortDo {
	return s.withDO(s.DO.Assign(attrs...))
}

func (s shortDo) Joins(fields ...field.RelationField) IShortDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Joins(_f))
	}
	return &s
}

func (s shortDo) Preload(fields ...field.RelationField) IShortDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Preload(_f))
	}
	return &s
}

func (s shortDo) FirstOrInit() (*model.Short, error) {
	if result, err := s.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Short), nil
	}
}

func (s shortDo) FirstOrCreate() (*model.Short, error) {
	if result, err := s.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Short), nil
	}
}

func (s shortDo) FindByPage(offset int, limit int) (result []*model.Short, count int64, err error) {
	result, err = s.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = s.Offset(-1).Limit(-1).Count()
	return
}

func (s shortDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = s.Count()
	if err != nil {
		return
	}

	err = s.Offset(offset).Limit(limit).Scan(result)
	return
}

func (s shortDo) Scan(result interface{}) (err error) {
	return s.DO.Scan(result)
}

func (s shortDo) Delete(models ...*model.Short) (result gen.ResultInfo, err error) {
	return s.DO.Delete(models)
}

func (s *shortDo) withDO(do gen.Dao) *shortDo {
	s.DO = *do.(*gen.DO)
	return s
}
