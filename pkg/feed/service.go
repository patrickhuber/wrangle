package feed

type Service interface {
	List(request *ListRequest) (*ListResponse, error)
	Update(request *UpdateRequest) (*UpdateResponse, error)
	Generate(request *GenerateRequest) (*GenerateResponse, error)
}
