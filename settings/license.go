/*
license 的验证
注意：部分功能，依赖于setting的初始化完成
*/

package settings

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"myGin/pkg/utils"
)

const (
	LicenseProduct  = "bravotranscode"
	LicenseVersion  = "official"
	LicenseFileName = "license.txt"
)

type license struct {
	Customer     string         `json:"customer"`
	Product      string         `json:"product"`
	Features     licenseFeature `json:"features"`
	HardwareInfo string         `json:"hwinfo"`
	BuildTime    int64          `json:"build_time"`
	Version      string         `json:"version"`
	Expiration   int64          `json:"expiration"`
	LicenseId    int64          `json:"license_id"`
}

type licenseFeature struct {
	NormalConcurrentTasks int64 `json:"normal_concurrent_tasks"` // 并发数量
	NormalNodes           int64 `json:"normal_nodes"`            //转码节点个数
}

func NewTesterLicense() *license {
	return &license{
		Product:    LicenseProduct,
		Version:    LicenseVersion,
		Features:   licenseFeature{NormalConcurrentTasks: 1},
		Expiration: time.Now().AddDate(1, 0, 0).Unix(),
	}
}

func LoadLicense() error {
	l := &license{}
	if !l.Enable() {
		return fmt.Errorf("license not enable")
	}

	if err := l.Verify(); err != nil {
		return err
	}

	// 正式
	LicenseCfg = l
	return nil
}

func (l *license) Expire() bool {
	if l.Product != LicenseProduct {
		return true
	}
	if time.Now().Unix() >= l.Expiration {
		return true
	}
	return false
}

func (l *license) Verify() error {
	licenseFile := l.GetFilePath()
	if licenseFile == "" {
		return fmt.Errorf("license file empty")
	}

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*10)
	defer cancel()
	cmd := exec.CommandContext(ctx, VerifyLicense, licenseFile)
	var out, errOut bytes.Buffer
	cmd.Stderr = &errOut
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return errors.Wrapf(err, "do cmd: %s error: %s", cmd.String(), errOut.String())
	}

	if err := json.Unmarshal(out.Bytes(), l); err != nil {
		return errors.Wrapf(err, "json unmarshal body: %v failed", out.String())
	}

	if l.Expire() {
		return fmt.Errorf("license is expired")
	}
	return nil
}

func (l *license) GetFilePath() string {
	if AppCfg.LicenseFileDir == "" {
		return ""
	}
	return filepath.Join(AppCfg.LicenseFileDir, LicenseFileName)
}

func (l *license) Enable() bool {
	licenseFilePath := l.GetFilePath()
	if licenseFilePath == "" {
		return false
	}
	return utils.Exist(licenseFilePath)
}
