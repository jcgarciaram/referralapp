package referralapp_api

import (
    "fmt"
    "time"
    "strconv"
    "math/rand"
    
    qrcode "github.com/skip2/go-qrcode"
    "github.com/jcgarciaram/general-api/apiutils"
)

type ReferralCode struct {
    ReferralCodePk              int                 `json:"referral_code_pk"`
    ReferralCodeId              string              `json:"referral_code_id"`
    StoreId                     string              `json:"store_id"`
    GeneratedPhone              string              `json:"generated_phone"`
    ExpirationDate              apiutils.CustomTime `json:"expiration_date"`
    UsedCount                   int                 `json:"used_count"`
    VerificationCode            int                 `json:"verification_code"`
    VerificationCodeDisplay     string              `json:"verification_code_display"`
}


func (rc *ReferralCode) generateVerificationCode() {
    s1 := rand.NewSource(time.Now().UnixNano())
    r1 := rand.New(s1)

	
	n1 := r1.Intn(100)
	n2 := r1.Intn(100)
	n3 := r1.Intn(100)

	display := fmt.Sprintf("%02d %02d %02d", n1, n2, n3)
	actualInt := fmt.Sprintf("%02d%02d%02d", n1, n2, n3)
    
    rc.VerificationCode, _ = strconv.Atoi(actualInt)
    rc.VerificationCodeDisplay = display
}

func (rc *ReferralCode) generateQRCode() ([]byte, error) {
    
    url := "1 " + rc.StoreId + " " + rc.ReferralCodeId

    return qrcode.Encode(url, qrcode.Medium, 256)
}