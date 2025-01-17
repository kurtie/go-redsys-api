package redsys

import (
	"crypto/hmac"
	"encoding/base64"
	"encoding/json"
	"net/url"
)

// Redsys Init this struct with your key to operate with the corresponding functions
type Redsys struct {
	Key string
}

// CreateMerchantParameters Return a string corresponding to a marshalled MerchantParametersRequest
func (r *Redsys) CreateMerchantParameters(data *MerchantParametersRequest) string {
	merchantMarshalledParams, _ := json.Marshal(data)
	return base64.URLEncoding.EncodeToString(merchantMarshalledParams)
}

// DecodeMerchantParameters Decode a response into a MerchantParametersResponse
func (r *Redsys) DecodeMerchantParameters(data string) MerchantParametersResponse {
	merchantParameters := MerchantParametersResponse{}
	decodedB64, _ := base64.URLEncoding.DecodeString(data)
	unscaped, err := url.QueryUnescape(string(decodedB64))
	if err != nil {
		unscaped = string(decodedB64)
	}
	json.Unmarshal([]byte(unscaped), &merchantParameters)
	return merchantParameters
}

// CreateMerchantSignature generates a merchant signature from MerchantParametersRequest
func (r *Redsys) CreateMerchantSignature(data *MerchantParametersRequest) string {
	stringMerchantParameters := r.CreateMerchantParameters(data)

	orderID := data.MerchantOrder

	encrypted := r.encrypt3DES(orderID)
	return r.mac256(stringMerchantParameters, encrypted)
}

// CreateMerchantSignatureNotif generates a signature for MerchantParametersResponse representing string
func (r *Redsys) CreateMerchantSignatureNotif(data string) string {
	merchantParametersResponse := r.DecodeMerchantParameters(data)

	orderID := merchantParametersResponse.Order
	encrypted := r.encrypt3DES(orderID)
	mac := r.mac256(data, encrypted)

	decodedMac, _ := base64.StdEncoding.DecodeString(mac)
	return base64.URLEncoding.EncodeToString(decodedMac)
}

// MerchantSignatureIsValid checks that two hmacs are equal
func (r *Redsys) MerchantSignatureIsValid(mac1 string, mac2 string) bool {
	return hmac.Equal([]byte(mac1), []byte(mac2))
}
