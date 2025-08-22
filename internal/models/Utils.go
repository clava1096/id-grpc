package models

// TODO переместить это в слой домена или респонсов
type JsonResponse struct {
	Msg string `json:"msg"`
}

type JsonLinkResponse struct {
	Link string `json:"link"`
}
