//go:build !android
// +build !android

package anet

func GetNetDriver() NetDriver {
	return nil
}
