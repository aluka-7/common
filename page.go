package common

type Pagination struct {
	pageNo       int // 当前页码
	PageSize     int // 每页大小
	TotalPages   int // 总页数
	TotalRecords int // 总记录数
}

func DefaultPagination() *Pagination {
	return NewPagination(1, 20, 0)
}
func NewPagination(pageNo, pageSize, totalPages int) *Pagination {
	return &Pagination{
		pageNo:     pageNo,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}

// SetPageNumber 设置当前分页的页码信息，必须大于0，否则设置无效。默认是1
func (p *Pagination) SetPageNumber(pageNo int) {
	if p.pageNo <= 0 {
		return
	}
	p.pageNo = pageNo
}

// SetTotalRecord 设置当前分页信息的总记录数
func (p *Pagination) SetTotalRecord(totalRecords int) {
	p.TotalRecords = totalRecords
	p.computeTotalPages()
}

// SetPageSize 设置页大小，必须大于0，否则设置无效
func (p *Pagination) SetPageSize(pageSize int) {
	if pageSize <= 0 {
		return
	}
	p.PageSize = pageSize
	p.computeTotalPages()
}

//computeTotalPages 根据总记录数和页大小计算总页数，在每次设置完总记录数和页大小后都会自动 进行计算，即方法SetPageSize和SetTotalRecord被调用后会自动调用该方法进行总页数的计算。
func (p *Pagination) computeTotalPages() {
	p.TotalPages = p.TotalRecords / p.PageSize
	if (p.TotalRecords % p.PageSize) != 0 {
		p.TotalPages += 1
	}
	if p.TotalPages <= 0 {
		p.TotalPages = 1
	}
}

// Limit 获取mysql分页值（limit,offset）
func (p *Pagination) Limit() (limit int, offset int) {
	pageNo := 0
	if p.pageNo > 0 {
		pageNo = p.pageNo - 1
	}
	return p.PageSize, pageNo * p.PageSize
}

// IsFirst 如果当前页面是第一页，则返回true.
func (p *Pagination) IsFirst() bool {
	return p.pageNo == 1
}

// HasPrevious 如果存在相对于当前页面的上一页，则返回true.
func (p *Pagination) HasPrevious() bool {
	return p.pageNo > 1
}

func (p *Pagination) Previous() int {
	if !p.HasPrevious() {
		return p.pageNo
	}
	return p.pageNo - 1
}

// HasNext 如果存在相对于当前页面的下一页，则返回true.
func (p *Pagination) HasNext() bool {
	return p.TotalRecords > p.pageNo*p.PageSize
}

func (p *Pagination) Next() int {
	if !p.HasNext() {
		return p.pageNo
	}
	return p.pageNo + 1
}

// IsLast 如果当前页面是最后一页，则返回true.
func (p *Pagination) IsLast() bool {
	if p.TotalRecords == 0 {
		return true
	}
	return p.TotalRecords > (p.pageNo-1)*p.PageSize && !p.HasNext()
}

// Total returns number of total rows.
func (p *Pagination) Total() int {
	return p.TotalRecords
}

// TotalPage 返回总页数.
func (p *Pagination) TotalPage() int {
	if p.TotalRecords == 0 {
		return 1
	}
	if p.TotalRecords%p.PageSize == 0 {
		return p.TotalRecords / p.PageSize
	}
	return p.TotalRecords/p.PageSize + 1
}

type PaginationList struct {
	Pagination
	List interface{} // 记录集
}

// Paginator 在分页器中显示页面
type Paginator struct {
	num       int
	isCurrent bool
}

func (p *Paginator) Num() int {
	return p.num
}

func (p *Paginator) IsCurrent() bool {
	return p.isCurrent
}

func getMiddleIdx(numPages int) int {
	if numPages%2 == 0 {
		return numPages / 2
	}
	return numPages/2 + 1
}

// Pages 返回相对于当前页面的附近页码的列表。如果值为-1表示"..."，则不会显示更多页面。
func (p *Pagination) Pages() []*Paginator {
	if p.TotalPages == 0 {
		return []*Paginator{}
	} else if p.TotalPages == 1 && p.TotalPage() == 1 {
		// 仅显示当前页面。
		return []*Paginator{{1, true}}
	}

	// 总页数小于或等于.
	if p.TotalPage() <= p.TotalPages {
		pages := make([]*Paginator, p.TotalPage())
		for i := range pages {
			pages[i] = &Paginator{i + 1, i+1 == p.TotalPages}
		}
		return pages
	}

	numPages := p.TotalPages
	maxIdx := numPages - 1
	offsetIdx := 0
	hasMoreNext := false

	// 检查更多的上一页和下一页。
	previousNum := getMiddleIdx(p.TotalPages) - 1
	if previousNum > p.pageNo-1 {
		previousNum -= previousNum - (p.pageNo - 1)
	}
	nextNum := p.TotalPages - previousNum - 1
	if p.pageNo+nextNum > p.TotalPage() {
		delta := nextNum - (p.TotalPage() - p.pageNo)
		nextNum -= delta
		previousNum += delta
	}

	offsetVal := p.pageNo - previousNum
	if offsetVal > 1 {
		numPages++
		maxIdx++
		offsetIdx = 1
	}

	if p.pageNo+nextNum < p.TotalPage() {
		numPages++
		hasMoreNext = true
	}

	pages := make([]*Paginator, numPages)

	// 还有更多的前几页.
	if offsetIdx == 1 {
		pages[0] = &Paginator{-1, false}
	}
	// 还有更多的下一页.
	if hasMoreNext {
		pages[len(pages)-1] = &Paginator{-1, false}
	}

	// 检查前几页.
	for i := 0; i < previousNum; i++ {
		pages[offsetIdx+i] = &Paginator{i + offsetVal, false}
	}

	pages[offsetIdx+previousNum] = &Paginator{p.pageNo, true}

	// 检查下一页.
	for i := 1; i <= nextNum; i++ {
		pages[offsetIdx+previousNum+i] = &Paginator{p.pageNo + i, false}
	}
	return pages
}
