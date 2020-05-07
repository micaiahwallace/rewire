package rewire

import "crypto/rsa"

type KeyRequest struct {
}

type KeyResp struct {
}

func CreateKeyRequest(pub *rsa.PublicKey) (*KeyRequest, *KeyResp) {
	req := &KeyRequest{}
	res := &KeyResp{}
	return req, res
}
