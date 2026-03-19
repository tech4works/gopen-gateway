package vo

type Cacheable interface {
	Entry() (string, error)
}
