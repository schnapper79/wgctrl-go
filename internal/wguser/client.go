package wguser

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"golang.zx2c4.com/wireguard/wgctrl/internal/wginternal"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

var _ wginternal.Client = &Client{}

// A Client provides access to userspace WireGuard device information.
type Client struct {
	dial func(device string) (net.Conn, error)
	find func() ([]string, error)
}

// New creates a new Client.
func New() (*Client, error) {
	return &Client{
		// Operating system-specific functions which can identify and connect
		// to userspace WireGuard devices. These functions can also be
		// overridden for tests.
		dial: dial,
		find: find,
	}, nil
}

// Close implements wginternal.Client.
func (c *Client) Close() error { return nil }

// Devices implements wginternal.Client.
func (c *Client) Devices() ([]*wgtypes.Device, error) {
	devices, err := c.find()
	if err != nil {
		return nil, err
	}

	var wgds []*wgtypes.Device
	for _, d := range devices {
		wgd, err := c.getDevice(d)
		if err != nil {
			return nil, err
		}

		wgds = append(wgds, wgd)
	}

	return wgds, nil
}

// Device implements wginternal.Client.
func (c *Client) Device(name string) (*wgtypes.Device, error) {
	devices, err := c.find()
	if err != nil {
		return nil, err
	}
	fmt.Printf("checking for [%v] , found [%v] devices\n",name,len(devices))

	for _, d := range devices {
		dn:=deviceName(d)
		fmt.Printf("checking  [%v] vs [%v]",name,dn)
		if name != dn {
			fmt.Println(" --> not equal\n")
			continue
		}
		fmt.Println(" --> yay, equal\n")
		return c.getDevice(d)
	}
	fmt.Println(" -->return with error\n")
	return nil, os.ErrNotExist
}

// ConfigureDevice implements wginternal.Client.
func (c *Client) ConfigureDevice(name string, cfg wgtypes.Config) error {
	devices, err := c.find()
	if err != nil {
		return err
	}
	fmt.Printf("checking for [%v] to configure , found [%v] devices\n",name,len(devices))

	for _, d := range devices {
		dn:=deviceName(d)
		fmt.Printf("checking  [%v] vs [%v]",name,dn)
		if name != dn {
			fmt.Println(" --> not equal\n")
			continue
		}
		fmt.Println(" --> yay, equal\n")
		return c.configureDevice(d, cfg)
	}
	fmt.Println(" -->return with error\n")
	return os.ErrNotExist
}

// deviceName infers a device name from an absolute file path with extension.
func deviceName(sock string) string {
	return strings.TrimSuffix(filepath.Base(sock), filepath.Ext(sock))
}

func panicf(format string, a ...interface{}) {
	panic(fmt.Sprintf(format, a...))
}
