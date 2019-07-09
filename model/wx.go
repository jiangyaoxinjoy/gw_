package model

import (
	"fmt"
	"gw/utils"
	"strconv"
	"time"

	"github.com/go-xorm/xorm"

	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

type SetupDevice struct {
	DeviceId    string `json:"deviceId" binding:"required"`
	TokenString string `json:"token" binding:"required"`
	Encrypt     string `json:"encrypt" binding:"required"`
	Lng         string `json:"lng" binding:"required"`
	Lat         string `json:"lat" binding:"required"`
	Address     string `json:"address" binding:"required"`
}

type AlertCount struct {
	Water int64 `json:"water"`
	Open  int64 `json:"open"`
	Down  int64 `json:"down"`
	Loss  int64 `json:"loss"`
}

func (wx *Model) WxAllAlertCount(selectCompanyId int, companyId int) (AlertCount, error) {
	var (
		queryCompanId int
		count         AlertCount
		water         int64
		open          int64
		down          int64
		loss          int64
	)
	db, _ := utils.Connect()
	queryCompanId = selectCompanyId
	if companyId != 1 {
		queryCompanId = companyId
	}
	if queryCompanId == 0 {
		water, _ = db.Where("1=1").And("state = ?", "10").Count(new(GwDevice))
		open, _ = db.Where("1=1").And("state = ?", "20").Count(new(GwDevice))
		down, _ = db.Where("1=1").And("state = ?", "30").Count(new(GwDevice))
		loss, _ = db.Where("1=1").And("state = ?", "40").Count(new(GwDevice))
	} else {
		water, _ = db.Where("company_id = ?", queryCompanId).And("state = ?", "10").Count(new(GwDevice))
		open, _ = db.Where("company_id = ?", queryCompanId).And("state = ?", "20").Count(new(GwDevice))
		down, _ = db.Where("company_id = ?", queryCompanId).And("state = ?", "30").Count(new(GwDevice))
		loss, _ = db.Where("company_id = ?", queryCompanId).And("state = ?", "70").Count(new(GwDevice))
	}
	count.Water = water
	count.Open = open
	count.Down = down
	count.Loss = loss

	// if err := xSession.Find(&devices); err != nil {
	// 	return devices, err
	// }
	return count, nil
}

func (wx *Model) WxSetupDevice(query SetupDevice) error {
	var (
		device GwDevice
	)
	res, _ := AesDecrypt(query.Encrypt, []byte("gwechatguanwei99gwechatguanwei99"))
	fmt.Println(string(res))
	if query.DeviceId != string(res) {
		return fmt.Errorf("设备号错误")
	}
	// device.DeviceId = query.DeviceId
	device.Lat = query.Lat
	device.Lng = query.Lng
	device.Address = query.Address
	device.Status = 1
	device.Setuptime = strconv.FormatInt(time.Now().Unix(), 10)
	db, _ := utils.Connect()
	db.ShowSQL(true)
	count, err := db.Where("device_id = ?", query.DeviceId).Cols("lat", "lng", "status", "setuptime", "address").Update(&device)
	fmt.Println(count)
	if err != nil {
		return err
	}
	return nil
}

func AesEncrypt(origData []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "AES", err
	}
	//blockSize := block.BlockSize()
	//fmt.Println(blockSize)
	origData = PKCS7Padding(origData, 16)
	iv, _ := base64.StdEncoding.DecodeString("AAAAAAAAAAAAAAAAAAAAAA==")
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return base64.StdEncoding.EncodeToString(crypted), nil
}

func AesDecrypt(encryptedData string, key []byte) ([]byte, error) {
	crypted, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//blockSize := block.BlockSize()
	iv, _ := base64.StdEncoding.DecodeString("AAAAAAAAAAAAAAAAAAAAAA==")
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}

func PKCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	unPadding := int(plantText[length-1])
	if unPadding < 1 || unPadding > 32 {
		unPadding = 0
	}
	return plantText[:(length - unPadding)]
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func (wx *Model) WxStateAlert(selectCompanyId int, companyId int, state string) ([]GwDevice, error) {
	var (
		queryCompanId int
		xSession      *xorm.Session
		device        []GwDevice
	)
	db, _ := utils.Connect()
	queryCompanId = selectCompanyId
	if companyId != 1 {
		queryCompanId = companyId
	}
	if queryCompanId == 0 {
		xSession = db.Where("1=1").And("state = ?", state)
	} else {
		xSession = db.Where("company_id = ?", queryCompanId).And("state = ?", state)
	}
	err := xSession.Find(&device)
	if err != nil {
		return device, err
	}
	return device, nil
}
