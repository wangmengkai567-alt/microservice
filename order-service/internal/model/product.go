package model

// ProductInfo 从商品服务获取的商品信息（非数据库模型）
type ProductInfo struct {
	ID          uint
	Name        string
	Description string
	Price       float64
	Stock       int
}
