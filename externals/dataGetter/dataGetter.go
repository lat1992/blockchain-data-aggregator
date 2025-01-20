package dataGetter

type DataGetterService interface {
}

type DataGetter struct {
}

func NewDataGetter() *DataGetter {
	return &DataGetter{}
}

func (g *DataGetter) GetData(url string, params map[string]string) (map[string]interface{}, error) {
	return nil, nil
}
