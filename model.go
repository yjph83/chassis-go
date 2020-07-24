package chassis

import (
	"time"

	"github.com/jinzhu/gorm"
)

//Page page
type Page struct {
	List   interface{} `json:"list,omitempty"`
	Total  uint        `json:"total,omitempty"`
	Offset uint        `json:"offset,omitempty"`
	Index  uint        `json:"page_index,omitempty"`
	Size   uint        `json:"page_size,omitempty"`
	Pages  uint        `json:"pages,omitempty"`
}

//Pagination 新建分页查询
type Pagination struct {
	Offset    uint        `json:"offset,omitempty"`
	Limit     uint        `json:"limit,omitempty"`
	Condition interface{} `json:"condition,omitempty"`
}

//NewPage new page
func NewPage(data interface{}, index, size, count uint) *Page {
	var pages uint
	if count%size == 0 {
		pages = count / size
	} else {
		pages = count/size + 1
	}
	return &Page{
		List:   data,
		Total:  count,
		Size:   size,
		Offset: index * size,
		Index:  index,
		Pages:  pages,
	}
}

//NewPagination pagination query
func NewPagination(db *gorm.DB, model interface{}, pageIndex, pageSize uint) *Page {
	var count uint
	db.Count(&count)
	if count > 0 && count > pageIndex*pageSize {
		db.Limit(pageSize).
			Offset(pageIndex * pageSize).
			Find(model)
		return NewPage(model, pageIndex, pageSize, count)
	}
	return nil
}

//paginationQuery pagination query。 pageIndex 必须从 1 开始， 如果从0 开始则必须 将 (pageIndex - 1) 换成 pageIndex
func PaginationQuery(db *gorm.DB, model interface{}, pageIndex, pageSize uint) (*Page, error) {
	var count uint
	db.Count(&count)
	var err error
	if count > 0 && count > (pageIndex - 1)*pageSize{
		err = db.Limit(pageSize).Offset((pageIndex - 1) * pageSize).Find(model).Error
		return NewPageInfo(model, pageIndex, pageSize, count), err
	}
	return nil, err
}

//NewPage new Chassis.page
func NewPageInfo(data interface{}, index, size, count uint) *Page {
	var pages uint
	if count%size == 0 {
		pages = count / size
	} else {
		pages = count/size + 1
	}
	return &Page{
		List:   data,
		Total:  count,
		Size:   size,
		Offset: (index - 1) * size,
		Index:  index,
		Pages:  pages,
	}
}

//SampleBaseDO model with id pk
type SampleBaseDO struct {
	ID uint `gorm:"primary_key" json:"id"`
}

//Model gorm model
type BaseDO struct {
	ID        uint       `gorm:"primary_key" json:"id"`            // primary key
	CreatedAt time.Time  `json:"created_at,omitempty"`             // created time
	UpdatedAt time.Time  `json:"updated_at,omitempty"`             //updated time
	DeletedAt *time.Time `sql:"index" json:"deleted_at,omitempty"` //deleted time
}

//ComplexBaseDO gorm model composed Model add Addition
type ComplexBaseDO struct {
	BaseDO
	Version  uint   `json:"version"` //version opt lock
	Addition string `json:"addition,omitempty"`
}
