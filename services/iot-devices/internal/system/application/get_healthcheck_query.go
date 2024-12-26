package system_application

const GetHealthcheckQueryName = "GetHealthcheckQuery"

type GetHealthcheckQuery struct{}

func NewGetHealthcheckQuery() *GetHealthcheckQuery {
	return &GetHealthcheckQuery{}
}

func (hq GetHealthcheckQuery) Type() string {
	return GetHealthcheckQueryName
}

func (hq GetHealthcheckQuery) Data() map[string]interface{} {
	return map[string]interface{}{}
}
