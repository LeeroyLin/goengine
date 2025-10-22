package pool

type IdPool[T interface{}] struct {
	chanData  chan T
	getter    func(idx int64) T
	onDestroy func(v T)
	max       int64
}

// NewIdPool 新建 Id池子
//
// parameters:
// @max: 最大缓存个数
// @getter: 获得实例的方法，传递生成下标
func NewIdPool[T interface{}](max int64, getter func(idx int64) T) *IdPool[T] {
	p := &IdPool[T]{
		max:      max,
		getter:   getter,
		chanData: make(chan T, max),
	}

	p.init(max)

	return p
}

func NewIdPoolWithOnDestroy[T interface{}](max int64, getter func(idx int64) T, onDestroy func(v T)) *IdPool[T] {
	p := NewIdPool(max, getter)
	p.onDestroy = onDestroy

	return p
}

func (p *IdPool[T]) init(num int64) {
	for i := int64(0); i < num; i++ {
		p.chanData <- p.getter(i)
	}
}

func (p *IdPool[T]) Get() T {
	select {
	case v := <-p.chanData:
		return v
	default:
		return p.getter(0)
	}
}

func (p *IdPool[T]) Set(v T) {
	select {
	case p.chanData <- v:
		return
	default:
		if p.onDestroy != nil {
			p.onDestroy(v)
		}
	}
}
