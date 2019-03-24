package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"time"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/ssh"
	"github.com/docker/machine/libmachine/state"
	api "github.com/vscale/go-vscale"
)

const (
	defaultLocation = "spb0"
	defaultMadeFrom = "ubuntu_16.04_64_001_docker"
	defaultRplan    = "small"
	defaultSwapFile = 0
)

// Driver ...
type Driver struct {
	*drivers.BaseDriver

	AccessToken string
	Location    string
	MadeFrom    string
	Rplan       string
	ScaletID    int64
	ScaletName  string
	SSHKeyID    int64
	SwapFile    int
}

// NewDriver ...
func NewDriver(hostName string, storePath string) *Driver {
	return &Driver{
		Location: defaultLocation,
		MadeFrom: defaultMadeFrom,
		Rplan:    defaultRplan,
		SwapFile: defaultSwapFile,
		BaseDriver: &drivers.BaseDriver{
			MachineName: hostName,
			StorePath:   storePath,
		},
	}
}

func (d *Driver) getClient() *api.WebClient {
	return api.NewClient(d.AccessToken)
}

func (d *Driver) createSSHKey() error {
	// Generate new SSH Key pair
	err := ssh.GenerateSSHKey(d.GetSSHKeyPath())
	if err != nil {
		return err
	}

	// Read public SSH Key
	publicKey, err := ioutil.ReadFile(d.GetSSHKeyPath() + ".pub")
	if err != nil {
		return err
	}

	// Upload public SSH Key to Vscale
	key, _, err := d.getClient().SSHKey.Create(string(publicKey), d.MachineName)

	d.SSHKeyID = key.ID

	return err
}

func (d *Driver) createScalet() error {
	client := d.getClient()

	scalet, _, err := client.Scalet.CreateWithoutPassword(
		d.MadeFrom,
		d.Rplan,
		d.MachineName,
		d.Location,
		true,
		[]int64{d.SSHKeyID},
		true,
	)
	if err != nil {
		return err
	}

	d.ScaletID = scalet.CTID

	for {
		scalet, _, err := client.Scalet.Get(d.ScaletID)
		if err != nil {
			return err
		}

		if scalet.Active == true && scalet.PublicAddresses.Address != "" {
			d.IPAddress = scalet.PublicAddresses.Address

			break
		}

		time.Sleep(1 * time.Second)
	}

	return nil
}

func (d *Driver) createSwapFile() error {
	if d.SwapFile == 0 {
		return nil
	}

	for {
		sshConnection := drivers.WaitForSSH(d)

		if sshConnection == nil {
			_, err := drivers.RunSSHCommandFromDriver(d, fmt.Sprintf(`
					touch /var/swap.img && \
					chmod 600 /var/swap.img && \
					dd if=/dev/zero of=/var/swap.img bs=1MB count=%d && \
					mkswap /var/swap.img && swapon /var/swap.img && \
					echo '/var/swap.img    none    swap    sw    0    0' >> /etc/fstab
				`, d.SwapFile))

			if err != nil {
				return err
			}

			break
		}

		time.Sleep(3 * time.Second)
	}

	return nil
}

// DriverName ...
func (d *Driver) DriverName() string {
	return "vscale"
}

// GetCreateFlags ...
func (d *Driver) GetCreateFlags() []mcnflag.Flag {
	return []mcnflag.Flag{
		mcnflag.StringFlag{
			EnvVar: "VSCALE_ACCESS_TOKEN",
			Name:   "vscale-access-token",
			Usage:  "Access token",
		},
		mcnflag.StringFlag{
			EnvVar: "VSCALE_LOCATION",
			Name:   "vscale-location",
			Usage:  "Location",
			Value:  defaultLocation,
		},
		mcnflag.StringFlag{
			EnvVar: "VSCALE_MADE_FROM",
			Name:   "vscale-made-from",
			Usage:  "Made from",
			Value:  defaultMadeFrom,
		},
		mcnflag.StringFlag{
			EnvVar: "VSCALE_RPLAN",
			Name:   "vscale-rplan",
			Usage:  "Rplan",
			Value:  defaultRplan,
		},
		mcnflag.IntFlag{
			EnvVar: "VSCALE_SWAP_FILE",
			Name:   "vscale-swap-file",
			Usage:  "Swap file",
			Value:  defaultSwapFile,
		},
	}
}

// SetConfigFromFlags ...
func (d *Driver) SetConfigFromFlags(flags drivers.DriverOptions) error {
	d.AccessToken = flags.String("vscale-access-token")
	d.Location = flags.String("vscale-location")
	d.MadeFrom = flags.String("vscale-made-from")
	d.Rplan = flags.String("vscale-rplan")
	d.SwapFile = flags.Int("vscale-swap-file")

	d.SetSwarmConfigFromFlags(flags)

	if d.AccessToken == "" {
		return fmt.Errorf("vscale driver requres the --vscale-access-token option")
	}

	return nil
}

// GetSSHHostname ...
func (d *Driver) GetSSHHostname() (string, error) {
	return d.GetIP()
}

// GetURL ...
func (d *Driver) GetURL() (string, error) {
	ip, err := d.GetIP()

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("tcp://%s", net.JoinHostPort(ip, "2376")), nil
}

// GetState ...
func (d *Driver) GetState() (state.State, error) {
	scalet, _, err := d.getClient().Scalet.Get(d.ScaletID)

	if err != nil {
		return state.Error, err
	}

	switch scalet.Status {
	case "defined":
		return state.Starting, nil
	case "started":
		return state.Running, nil
	case "stopped":
		return state.Stopped, nil
	}

	return state.None, nil
}

// PreCreateCheck ...
func (d *Driver) PreCreateCheck() error {
	if d.getClient() == nil {
		return fmt.Errorf("cannot create client")
	}

	return nil
}

// Create Scalet
func (d *Driver) Create() error {
	// Create SSH key
	if err := d.createSSHKey(); err != nil {
		return err
	}

	// Create Scalet
	if err := d.createScalet(); err != nil {
		return err
	}

	// Create Swap for Scalet
	err := d.createSwapFile()

	return err
}

// Start Scalet
func (d *Driver) Start() error {
	_, _, err := d.getClient().Scalet.Start(d.ScaletID, true)

	return err
}

// Stop Scalet
func (d *Driver) Stop() error {
	_, _, err := d.getClient().Scalet.Stop(d.ScaletID, true)

	return err
}

// Restart Scalet
func (d *Driver) Restart() error {
	_, _, err := d.getClient().Scalet.Restart(d.ScaletID, true)

	return err
}

// Remove Scalet
func (d *Driver) Remove() error {
	client := d.getClient()

	// Remove SSH Key
	if _, _, err := client.SSHKey.Remove(d.SSHKeyID); err != nil {
		return err
	}

	// Remove Scalet
	if _, _, err := client.Scalet.Remove(d.ScaletID, true); err != nil {
		return err
	}

	return nil
}

// Kill Scalet
func (d *Driver) Kill() error {
	// NOTE: Vscale not implemented kill option
	return d.Stop()
}
